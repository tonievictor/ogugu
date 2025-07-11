package telemetry

import (
	"context"
	"errors"
	"os"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

func NewOtlpExporter(ctx context.Context) (sdktrace.SpanExporter, error) {
	insecureopt := otlptracehttp.WithInsecure()

	val := os.Getenv("OTLP_ENDPOINT")
	if val == "" {
		return nil, errors.New("OTLP_ENDPOINT env variable not set")
	}

	endpointopt := otlptracehttp.WithEndpoint(val)
	return otlptracehttp.New(ctx, insecureopt, endpointopt)
}

func NewTraceProvider(exp sdktrace.SpanExporter) (*sdktrace.TracerProvider, error) {
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("ogugu"),
		),
	)
	if err != nil {
		return nil, err
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(r),
	), nil
}
