package tracing

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

// InitTracer sets up OpenTelemetry with an OTLP gRPC exporter. It configures
// a TracerProvider with service.name, service.version, and deployment.environment
// resource attributes, and sets the global TracerProvider and TextMapPropagator.
//
// It returns a cleanup function that should be called on application shutdown to
// flush any pending spans and release resources.
func InitTracer(serviceName, endpoint, environment string) (func(context.Context) error, error) {
	ctx := context.Background()

	opts := []otlptracegrpc.Option{}
	if endpoint != "" {
		opts = append(opts, otlptracegrpc.WithEndpoint(endpoint))
	}
	// For non-TLS endpoints (common in dev/staging), allow insecure connections.
	if environment != "production" {
		opts = append(opts, otlptracegrpc.WithInsecure())
	}

	exporter, err := otlptracegrpc.New(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP trace exporter: %w", err)
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion("1.0.0"),
			attribute.String("deployment.environment", environment),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Choose a sampler based on environment.
	var sampler sdktrace.Sampler
	switch environment {
	case "production":
		sampler = sdktrace.ParentBased(sdktrace.TraceIDRatioBased(0.1))
	case "staging":
		sampler = sdktrace.ParentBased(sdktrace.TraceIDRatioBased(0.5))
	default:
		sampler = sdktrace.AlwaysSample()
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sampler),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	cleanup := func(ctx context.Context) error {
		if err := tp.Shutdown(ctx); err != nil {
			return fmt.Errorf("failed to shutdown tracer provider: %w", err)
		}
		return nil
	}

	return cleanup, nil
}
