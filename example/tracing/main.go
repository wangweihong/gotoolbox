package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"

	"go.opentelemetry.io/otel/trace"

	"github.com/wangweihong/gotoolbox/pkg/callerutil"
	"github.com/wangweihong/gotoolbox/pkg/httpcli"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func doSomething(ctx context.Context) error {
	return fmt.Errorf("an error occurred in ServiceC")
}

func getStackTrace() []string {
	return callerutil.CallersDepth(10, 3).List()
}

// 设置OpenTelemetry将追踪数据导到本地控制台
func InitLocalTracer(serviceName string) (*sdktrace.TracerProvider, error) {
	exporter, err := stdouttrace.New(
		stdouttrace.WithWriter(os.Stdout),
		stdouttrace.WithPrettyPrint(),
	)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		)),
	)
	otel.SetTracerProvider(tp)

	return tp, nil
}

// 一个例子用于通过OpenTelemetry传递调用链，并在调用链中记录错误的调用栈
func main() {

	serverC := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.Tracer("serviceC").Start(r.Context(), "serviceC_HttpHandler")
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
		ctx, span := otel.Tracer("serviceB").Start(r.Context(), "serviceB_HttpHandler",
			// add labels/tags/resources(if any) that are specific to this scope.
			trace.WithAttributes(attribute.String("component", "addition")),
			trace.WithAttributes(attribute.String("someKey", "someValue")),
			trace.WithAttributes(attribute.Int("age", 89)))
		defer span.End()

		cli := &http.Client{
			Transport: otelhttp.NewTransport(http.DefaultTransport),
		}
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverC.URL, nil)
		if err != nil {
			stackTrace := getStackTrace()
			span.SetAttributes(attribute.Bool("error", true))
			span.SetAttributes(attribute.StringSlice("stack_trace", stackTrace))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resp, err := cli.Do(req)
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

	tracer, err := InitLocalTracer("serviceB")
	//shutdown, err := InitLocalTracer("serviceA")
	if err != nil {
		log.Fatalf("cannot initialize console exporter: %v", err)
	}
	defer func() {
		if err := tracer.Shutdown(context.Background()); err != nil {
			log.Fatalf("failed to shutdown TracerProvider: %v", err)
		}
	}()

	resp, err := httpcli.NewHttpRequestBuilder().
		WithEndpoint(serverB.URL).
		WithMethod("GET").
		Build().Invoke()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp.GetBody())

	//以下 数据是OpenTelemetry输出到终端的. 非HTTP回应
	/*
		generate by  InitLocalTracer
			{
			        "Name": "serviceC_handler",
			        "SpanContext": {
			                "TraceID": "1c619b2ac20ac5f3391478070974de8e",
			                "SpanID": "730a759393c34c4e",
			                "TraceFlags": "01",
			                "TraceState": "",
			                "Remote": false
			        },
			        "Parent": {
			                "TraceID": "00000000000000000000000000000000",
			                "SpanID": "0000000000000000",
			                "TraceFlags": "00",
			                "TraceState": "",
			                "Remote": false
			        },
			        "SpanKind": 1,
			        "StartTime": "2024-06-27T11:56:39.1636023+08:00",
			        "EndTime": "2024-06-27T11:56:39.1636023+08:00",
			        "Attributes": [
			                {
			                        "Key": "error",
			                        "Value": {
			                                "Type": "BOOL",
			                                "Value": true
			                        }
			                },
			                {
			                        "Key": "stack_trace",
			                        "Value": {
			                                "Type": "STRINGSLICE",
			                                "Value": [
			                                        "C:/goprogram/src/github.com/wangweihong/gotoolbox/example/tracing/main.go:99 main.main.func1",
			                                        "C:/Users/Administrator/go/go1.20.12/src/net/http/server.go:2122 net/http.HandlerFunc.ServeHTTP",
			                                        "C:/Users/Administrator/go/go1.20.12/src/net/http/server.go:2936 net/http.serverHandler.ServeHTTP",
			                                        "C:/Users/Administrator/go/go1.20.12/src/net/http/server.go:1995 net/http.(*conn).serve",
			                                        "C:/Users/Administrator/go/go1.20.12/src/runtime/asm_amd64.s:1598 runtime.goexit"
			                                ]
			                        }
			                }
			        ],
			        "Events": null,
			        "Links": null,
			        "Status": {
			                "Code": "Unset",
			                "Description": ""
			        },
			        "DroppedAttributes": 0,
			        "DroppedEvents": 0,
			        "DroppedLinks": 0,
			        "ChildSpanCount": 0,
			        "Resource": [
			                {
			                        "Key": "service.name",
			                        "Value": {
			                                "Type": "STRING",
			                                "Value": "serviceA"
			                        }
			                }
			        ],
			        "InstrumentationLibrary": {
			                "Name": "serviceC",
			                "Version": "",
			                "SchemaURL": ""
			        }
			}
			{
			        "Name": "serviceB_handler",
			        "SpanContext": {
			                "TraceID": "827b921267359273a2bdb3d7b9b1a6d8",
			                "SpanID": "32ffe2c3b9359389",
			                "TraceFlags": "01",
			                "TraceState": "",
			                "Remote": false
			        },
			        "Parent": {
			                "TraceID": "00000000000000000000000000000000",
			                "SpanID": "0000000000000000",
			                "TraceFlags": "00",
			                "TraceState": "",
			                "Remote": false
			        },
			        "SpanKind": 1,
			        "StartTime": "2024-06-27T11:56:39.1630614+08:00",
			        "EndTime": "2024-06-27T11:56:39.1636023+08:00",
			        "Attributes": [
			                {
			                        "Key": "error",
			                        "Value": {
			                                "Type": "BOOL",
			                                "Value": true
			                        }
			                },
			                {
			                        "Key": "stack_trace",
			                        "Value": {
			                                "Type": "STRINGSLICE",
			                                "Value": [
			                                        "C:/goprogram/src/github.com/wangweihong/gotoolbox/example/tracing/main.go:125 main.main.func2",
			                                        "C:/Users/Administrator/go/go1.20.12/src/net/http/server.go:2122 net/http.HandlerFunc.ServeHTTP",
			                                        "C:/Users/Administrator/go/go1.20.12/src/net/http/server.go:2936 net/http.serverHandler.ServeHTTP",
			                                        "C:/Users/Administrator/go/go1.20.12/src/net/http/server.go:1995 net/http.(*conn).serve",
			                                        "C:/Users/Administrator/go/go1.20.12/src/runtime/asm_amd64.s:1598 runtime.goexit"
			                                ]
			                        }
			                }
			        ],
			        "Events": null,
			        "Links": null,
			        "Status": {
			                "Code": "Unset",
			                "Description": ""
			        },
			        "DroppedAttributes": 0,
			        "DroppedEvents": 0,
			        "DroppedLinks": 0,
			        "ChildSpanCount": 0,
			        "Resource": [
			                {
			                        "Key": "service.name",
			                        "Value": {
			                                "Type": "STRING",
			                                "Value": "serviceA"
			                        }
			                }
			        ],
			        "InstrumentationLibrary": {
			                "Name": "serviceB",
			                "Version": "",
			                "SchemaURL": ""
			        }
			}
			{
			        "Name": "serviceA_handler",
			        "SpanContext": {
			                "TraceID": "02d205943cb14889e01d8800d8573e43",
			                "SpanID": "30ab7df644984441",
			                "TraceFlags": "01",
			                "TraceState": "",
			                "Remote": false
			        },
			        "Parent": {
			                "TraceID": "00000000000000000000000000000000",
			                "SpanID": "0000000000000000",
			                "TraceFlags": "00",
			                "TraceState": "",
			                "Remote": false
			        },
			        "SpanKind": 1,
			        "StartTime": "2024-06-27T11:56:39.162511+08:00",
			        "EndTime": "2024-06-27T11:56:39.1641204+08:00",
			        "Attributes": [
			                {
			                        "Key": "error",
			                        "Value": {
			                                "Type": "BOOL",
			                                "Value": true
			                        }
			                },
			                {
			                        "Key": "stack_trace",
			                        "Value": {
			                                "Type": "STRINGSLICE",
			                                "Value": [
			                                        "C:/goprogram/src/github.com/wangweihong/gotoolbox/example/tracing/main.go:151 main.main.func3",
			                                        "C:/Users/Administrator/go/go1.20.12/src/net/http/server.go:2122 net/http.HandlerFunc.ServeHTTP",
			                                        "C:/Users/Administrator/go/go1.20.12/src/net/http/server.go:2936 net/http.serverHandler.ServeHTTP",
			                                        "C:/Users/Administrator/go/go1.20.12/src/net/http/server.go:1995 net/http.(*conn).serve",
			                                        "C:/Users/Administrator/go/go1.20.12/src/runtime/asm_amd64.s:1598 runtime.goexit"
			                                ]
			                        }
			                }
			        ],
			        "Events": null,
			        "Links": null,
			        "Status": {
			                "Code": "Unset",
			                "Description": ""
			        },
			        "DroppedAttributes": 0,
			        "DroppedEvents": 0,
			        "DroppedLinks": 0,
			        "ChildSpanCount": 0,
			        "Resource": [
			                {
			                        "Key": "service.name",
			                        "Value": {
			                                "Type": "STRING",
			                                "Value": "serviceA"
			                        }
			                }
			        ],
			        "InstrumentationLibrary": {
			                "Name": "serviceA",
			                "Version": "",
			                "SchemaURL": ""
			        }
			}
	*/
}
