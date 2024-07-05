package main

import (
	"context"

	"go.opentelemetry.io/otel/attribute"

	"github.com/wangweihong/gotoolbox/pkg/errors"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func initTracer(ctx context.Context, serviceName string) (*sdktrace.TracerProvider, error) {
	var secureOption otlptracegrpc.Option
	secureOption = otlptracegrpc.WithInsecure()

	exporter, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithEndpoint("192.168.134.218:4317"),
		secureOption,
	)
	if err != nil {
		return nil, errors.Wrap(err, "fail to create exporter")
	}

	// labels/tags/resources that are common to all traces.
	//resource := resource.NewWithAttributes(
	//	semconv.SchemaURL,
	//	semconv.ServiceNameKey.String(serviceName),
	//	attribute.String("some-attribute", "some-value"),
	//)

	resources, err := resource.New(ctx,
		resource.WithAttributes(
			attribute.String("service.name", serviceName),
			attribute.String("library.language", "go"),
		))

	if err != nil {
		return nil, errors.Wrap(err, "fail to set resource")

	}

	provider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resources),
		// set the sampling rate based on the parent span to 60%
		//sdktrace.WithSampler(trace.ParentBased(trace.TraceIDRatioBased(0.6))),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	otel.SetTracerProvider(provider)

	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{}, // W3C Trace Context format; https://www.w3.org/TR/trace-context/
		),
	)

	return provider, nil
}
