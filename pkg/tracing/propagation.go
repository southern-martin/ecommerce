package tracing

import (
	"context"

	"go.opentelemetry.io/otel"
)

// Ensure mapCarrier satisfies the propagation.TextMapCarrier interface at compile time.
// (The interface is from go.opentelemetry.io/otel/propagation but used implicitly via otel.GetTextMapPropagator.)

// mapCarrier adapts a map[string]string to the propagation.TextMapCarrier
// interface, enabling trace context injection and extraction for systems
// that use simple string maps for headers (e.g., NATS messages).
type mapCarrier map[string]string

func (mc mapCarrier) Get(key string) string {
	return mc[key]
}

func (mc mapCarrier) Set(key, value string) {
	mc[key] = value
}

func (mc mapCarrier) Keys() []string {
	keys := make([]string, 0, len(mc))
	for k := range mc {
		keys = append(keys, k)
	}
	return keys
}

// InjectContext injects the current trace context from ctx into the provided
// headers map. This is used to propagate trace context across NATS messages
// or any other transport that uses map[string]string for headers.
//
// If headers is nil, a new map is created and returned.
func InjectContext(ctx context.Context, headers map[string]string) map[string]string {
	if headers == nil {
		headers = make(map[string]string)
	}
	propagator := otel.GetTextMapPropagator()
	propagator.Inject(ctx, mapCarrier(headers))
	return headers
}

// ExtractContext extracts trace context from the provided headers map and
// returns a new context with the extracted span context attached.
//
// This is used to reconstruct trace context from incoming NATS messages
// or any other transport that uses map[string]string for headers.
func ExtractContext(ctx context.Context, headers map[string]string) context.Context {
	if headers == nil {
		return ctx
	}
	propagator := otel.GetTextMapPropagator()
	return propagator.Extract(ctx, mapCarrier(headers))
}
