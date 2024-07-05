package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"runtime"

	"go.opentelemetry.io/otel/propagation"

	sdkmetric "go.opentelemetry.io/otel/sdk/metric"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

const serviceName = "AdderSvc"

func main() {
	os.Setenv("GODEBUG", "http2debug=2")
	ctx := context.Background()
	{
		tp, err := initTracer(ctx, serviceName)
		if err != nil {
			panic(err)
		}
		defer tp.Shutdown(ctx)

		mp, err := setupMetrics(ctx, serviceName)
		if err != nil {
			panic(err)
		}
		defer mp.Shutdown(ctx)
	}

	go serviceA(ctx, 8081)
	serviceB(ctx, 8082)
}

// curl -vkL http://127.0.0.1:8081/serviceA
func serviceA(ctx context.Context, port int) {
	mux := http.NewServeMux()
	mux.HandleFunc("/serviceA", serviceA_HttpHandler)
	handler := otelhttp.NewHandler(mux, "server.http.serviceA")
	serverPort := fmt.Sprintf(":%d", port)
	server := &http.Server{Addr: serverPort, Handler: handler}

	fmt.Println("serviceA listening on", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}

func getStackTrace() string {
	pc := make([]uintptr, 10)
	n := runtime.Callers(3, pc)
	frames := runtime.CallersFrames(pc[:n])

	var stackTrace string
	for {
		frame, more := frames.Next()
		stackTrace += fmt.Sprintf("%s\n\t%s:%d\n", frame.Function, frame.File, frame.Line)
		if !more {
			break
		}
	}
	return stackTrace
}

func serviceA_HttpHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer("myTracer").Start(r.Context(), "serviceA_HttpHandler")
	defer span.End()

	span.SetAttributes(attribute.String("parent_span_id", span.SpanContext().SpanID().String()))
	stackTrace := getStackTrace()
	span.SetAttributes(attribute.String("stack_trace", stackTrace))

	cli := &http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8082/serviceB", nil)
	if err != nil {
		panic(err)
	}

	req = req.WithContext(r.Context())

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))
	resp, err := cli.Do(req)
	if err != nil {
		panic(err)
	}

	w.Header().Add("SVC-RESPONSE", resp.Header.Get("SVC-RESPONSE"))
}

func serviceB(ctx context.Context, port int) {
	mux := http.NewServeMux()
	mux.HandleFunc("/serviceB", serviceB_HttpHandler)
	handler := otelhttp.NewHandler(mux, "server.http.serviceB")
	serverPort := fmt.Sprintf(":%d", port)
	server := &http.Server{Addr: serverPort, Handler: handler}

	fmt.Println("serviceB listening on", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}

func serviceB_HttpHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Extract the trace context from the incoming requests
	ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(r.Header))

	ctx, span := otel.Tracer("myTracer").Start(ctx, "serviceB_HttpHandler")
	defer span.End()

	answer := add(ctx, 42, 1813)
	w.Header().Add("SVC-RESPONSE", fmt.Sprint(answer))
	fmt.Fprintf(w, "hello from serviceB: Answer is: %d", answer)
}

func add(ctx context.Context, x, y int64) int64 {
	ctx, span := otel.Tracer("myTracer").Start(
		ctx,
		"add",
		// add labels/tags/resources(if any) that are specific to this scope.
		trace.WithAttributes(attribute.String("component", "addition")),
		trace.WithAttributes(attribute.String("someKey", "someValue")),
		trace.WithAttributes(attribute.Int("age", 89)),
	)
	defer span.End()

	counter, _ := sdkmetric.NewMeterProvider().
		Meter(
			"instrumentation/package/name",
			metric.WithInstrumentationVersion("0.0.1"),
		).
		Int64Counter(
			"add_counter",
			metric.WithDescription("how many times add function has been called."),
		)
	attrs := []attribute.KeyValue{
		attribute.String("component", "addition"),
		attribute.Int("age", 89),
	}
	counter.Add(
		ctx,
		1,
		// labels/tags
		metric.WithAttributes(attrs...),
	)

	log := NewLogrus(ctx).WithFields(logrus.Fields{
		"component": "addition",
		"age":       89,
	})
	log.Info("add_called")

	return x + y
}
