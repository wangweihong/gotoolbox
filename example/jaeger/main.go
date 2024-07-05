package main

import (
	"context"
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/codes"

	"github.com/wangweihong/gotoolbox/pkg/errors"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/credentials"

	"github.com/wangweihong/gotoolbox/pkg/callerutil"
	"github.com/wangweihong/gotoolbox/pkg/httpcli"
	"github.com/wangweihong/gotoolbox/pkg/json"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.14.0"
)

func doSomething(ctx context.Context) error {
	return errors.New("an error occurred in ServiceC")
}

func getStackTrace() []string {
	return callerutil.CallersDepth(10, 3).List()
}

func setupMetrics(ctx context.Context, serviceName string, tlsConfig *tls.Config) (*sdkmetric.MeterProvider, error) {
	connOpt := otlpmetricgrpc.WithInsecure()
	if tlsConfig != nil {
		connOpt = otlpmetricgrpc.WithTLSCredentials(credentials.NewTLS(tlsConfig))
	}

	exporter, err := otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithEndpoint("localhost:4317"),
		connOpt,
	)
	if err != nil {
		return nil, err
	}

	// labels/tags/resources that are common to all metrics.
	resource := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(serviceName),
		attribute.String("some-attribute", "some-value"),
	)

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(resource),
		sdkmetric.WithReader(
			// collects and exports metric data every 30 seconds.
			sdkmetric.NewPeriodicReader(exporter, sdkmetric.WithInterval(30*time.Second)),
		),
	)

	otel.SetMeterProvider(mp)

	return mp, nil
}

// 创建一个 logrus 钩子;
// （a） 将 TraceIds 和 spanIds 添加到日志中，
// （b） 将日志作为 span-events 添加到活动 span 中。
// 这个钩子从跟踪中查找任何 traceId 和 spanId，并将它们添加到日志事件中。
// 此外，它还会获取任何日志事件，并将它们作为 span 事件添加到跟踪中。
// 这使我们能够将日志与其相应的跟踪相关联，反之亦然。
// usage:
//
//	ctx, span := tracer.Start(ctx, "myFuncName")
//	l := NewLogrus(ctx)
//	l.Info("hello world")
func NewLogrus(ctx context.Context) *logrus.Entry {
	l := logrus.New()
	l.SetLevel(logrus.TraceLevel)
	l.AddHook(logrusTraceHook{})
	return l.WithContext(ctx)
}

// logrusTraceHook is a hook that;
// (a) adds TraceIds & spanIds to logs of all LogLevels
// (b) adds logs to the active span as events.
type logrusTraceHook struct{}

func (t logrusTraceHook) Levels() []logrus.Level { return logrus.AllLevels }

func (t logrusTraceHook) Fire(entry *logrus.Entry) error {
	ctx := entry.Context
	if ctx == nil {
		return nil
	}

	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return nil
	}

	{ // (a) adds TraceIds & spanIds to logs.
		sCtx := span.SpanContext()
		if sCtx.HasTraceID() {
			entry.Data["traceId"] = sCtx.TraceID().String()
		}
		if sCtx.HasSpanID() {
			entry.Data["spanId"] = sCtx.SpanID().String()
		}
	}

	{ // (b) adds logs to the active span as events.

		// code from: https://github.com/uptrace/opentelemetry-go-extra/tree/main/otellogrus
		// whose license(BSD 2-Clause) can be found at: https://github.com/uptrace/opentelemetry-go-extra/blob/v0.1.18/LICENSE
		attrs := make([]attribute.KeyValue, 0)
		logSeverityKey := attribute.Key("log.severity")
		logMessageKey := attribute.Key("log.message")
		attrs = append(attrs, logSeverityKey.String(entry.Level.String()))
		attrs = append(attrs, logMessageKey.String(entry.Message))

		span.AddEvent("log", trace.WithAttributes(attrs...))
		if entry.Level <= logrus.ErrorLevel {
			span.SetStatus(codes.Error, entry.Message)
		}
	}

	return nil
}

func setupGRPCTracing(ctx context.Context, serviceName string, tlsConfig *tls.Config) (*sdktrace.TracerProvider, error) {
	connOpt := otlptracegrpc.WithInsecure()
	if tlsConfig != nil {
		connOpt = otlptracegrpc.WithTLSCredentials(credentials.NewTLS(tlsConfig))
	}

	exporter, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithEndpoint("localhost:4317"),
		connOpt,
	)
	if err != nil {
		return nil, err
	}

	// labels/tags/resources that are common to all traces.
	resource := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(serviceName),
		attribute.String("some-attribute", "some-value"),
	)

	provider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource),
		// set the sampling rate based on the parent span to 60%
		// 配置批量跟踪数据, 并以60的速率采用 采样
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(0.6))),
	)

	otel.SetTracerProvider(provider)

	//设置传播器。通过传播机制，跟踪可以跨传输边界从一个服务传播/通信到另一个服务
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{}, // W3C Trace Context format; https://www.w3.org/TR/trace-context/
		),
	)

	return provider, nil
}

func startHTTPTracing(serviceName string) (*sdktrace.TracerProvider, error) {
	headers := map[string]string{
		"content-type": "application/json",
	}

	//
	exporter, err := otlptrace.New(
		context.Background(),
		otlptracehttp.NewClient(
			otlptracehttp.WithEndpoint("localhost:4318"),
			otlptracehttp.WithHeaders(headers),
			otlptracehttp.WithInsecure(),
		),
	)
	if err != nil {
		return nil, errors.Errorf("creating new exporter: %w", err)
	}

	// labels/tags/resources that are common to all traces.
	resources := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(serviceName),
	)

	tracerprovider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(
			exporter,
			sdktrace.WithMaxExportBatchSize(sdktrace.DefaultMaxExportBatchSize),
			sdktrace.WithBatchTimeout(sdktrace.DefaultScheduleDelay*time.Millisecond),
			sdktrace.WithMaxExportBatchSize(sdktrace.DefaultMaxExportBatchSize),
		),
		sdktrace.WithResource(resources),
	)

	otel.SetTracerProvider(tracerprovider)

	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{}, // W3C Trace Context format; https://www.w3.org/TR/trace-context/
		),
	)
	return tracerprovider, nil
}

// 设置OpenTelemetry将追踪数据导到Jaeger
func InitJaegerTracer(serviceName string) (*sdktrace.TracerProvider, error) {
	headers := map[string]string{
		"content-type": "application/json",
	}

	exporter, err := otlptrace.New(
		context.Background(),
		otlptracehttp.NewClient(
			otlptracehttp.WithEndpoint("192.168.134.218:4318"),
			otlptracehttp.WithHeaders(headers),
			otlptracehttp.WithInsecure(),
		),
	)
	if err != nil {
		return nil, errors.Wrap(err, "create new exporter fail")
	}

	tracerprovider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(
			exporter,
			sdktrace.WithMaxExportBatchSize(sdktrace.DefaultMaxExportBatchSize),
			sdktrace.WithBatchTimeout(sdktrace.DefaultScheduleDelay*time.Millisecond),
			sdktrace.WithMaxExportBatchSize(sdktrace.DefaultMaxExportBatchSize)),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		)),
	)

	otel.SetTracerProvider(tracerprovider)

	return tracerprovider, nil
}

// 一个例子用于通过OpenTelemetry传递调用链，并在调用链中记录错误的调用栈
func main() {

	serverC := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tracer := otel.Tracer("serviceC")
		ctx, span := tracer.Start(context.Background(), "serviceC_handler")
		defer span.End()

		err := doSomething(ctx)
		if err != nil {
			stackTrace := callerutil.Stacks(-1)
			span.SetAttributes(attribute.Bool("error", true))
			span.SetAttributes(attribute.StringSlice("stack_trace", stackTrace))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte("ServiceC success"))
	}))
	defer serverC.Close()

	serverB := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tracer := otel.Tracer("serviceB")
		_, span := tracer.Start(context.Background(), "serviceB_handler")
		defer span.End()

		resp, err := http.Get(serverC.URL)
		if err != nil {
			stackTrace := getStackTrace()
			span.SetAttributes(attribute.Bool("error", true))
			span.SetAttributes(attribute.StringSlice("stack_trace", stackTrace))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		if resp.StatusCode != http.StatusOK {
			stackTrace := getStackTrace()
			span.SetAttributes(attribute.Bool("error", true))
			span.SetAttributes(attribute.StringSlice("stack_trace", stackTrace))
			http.Error(w, string(body), resp.StatusCode)
			return
		}
		w.Write(body)
	}))
	defer serverB.Close()

	serverA := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tracer := otel.Tracer("serviceA")
		_, span := tracer.Start(context.Background(), "serviceA_handler")
		defer span.End()

		resp, err := http.Get(serverB.URL)
		if err != nil {
			stackTrace := getStackTrace()
			span.SetAttributes(attribute.Bool("error", true))
			span.SetAttributes(attribute.StringSlice("stack_trace", stackTrace))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		if resp.StatusCode != http.StatusOK {
			stackTrace := getStackTrace()
			span.SetAttributes(attribute.Bool("error", true))
			span.SetAttributes(attribute.StringSlice("stack_trace", stackTrace))
			http.Error(w, string(body), resp.StatusCode)
			return
		}
		w.Write(body)
	}))
	defer serverA.Close()

	tracer, err := InitJaegerTracer("serviceA")
	if err != nil {
		log.Fatalf("cannot initialize console exporter: %v", err)
	}
	defer func() {
		if err := tracer.Shutdown(context.Background()); err != nil {
			log.Fatalf("failed to shutdown TracerProvider: %v", err)
		}
	}()

	resp, err := httpcli.NewHttpRequestBuilder().
		WithEndpoint(serverA.URL).
		WithMethod("GET").
		Build().Invoke()
	if err != nil {
		log.Fatal(err)
	}
	json.PrettyPrint([]byte(resp.GetBody()))

}
