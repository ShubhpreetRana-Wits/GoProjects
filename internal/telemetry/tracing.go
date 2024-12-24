package telemetry

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace" // Alias this import
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

// Global Tracer instance
var Tracer trace.Tracer

// InitTelemetry initializes the telemetry system, including tracing.
func InitTelemetry(serviceName string) func() {
	// Set up an OTLP exporter
	client := otlptracegrpc.NewClient()
	exporter, err := otlptrace.New(context.Background(), client)
	if err != nil {
		log.Fatalf("Failed to create OTLP exporter: %v", err)
	}

	// Set up trace provider
	tp := sdktrace.NewTracerProvider( // Use sdktrace alias here
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		)),
	)
	otel.SetTracerProvider(tp)

	// Set global tracer
	Tracer = otel.Tracer(serviceName)

	// Shutdown function to clean up resources
	return func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Fatalf("Error shutting down tracer provider: %v", err)
		}
	}
}

// StartSpan creates a new span for tracing.
func StartSpan(ctx context.Context, tracer trace.Tracer, spanName string) (context.Context, trace.Span) {
	return tracer.Start(ctx, spanName)
}
