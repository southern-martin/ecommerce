package circuitbreaker

import (
	"context"
	"strings"
	"sync"

	"google.golang.org/grpc"
)

// BreakerRegistry manages a collection of circuit breakers keyed by service
// name. Each gRPC service automatically gets its own breaker, created on
// demand with the registry's default configuration.
type BreakerRegistry struct {
	breakers      sync.Map
	defaultConfig Config
}

// NewRegistry creates a new BreakerRegistry that uses the provided default
// configuration when creating breakers for previously unseen services.
func NewRegistry(defaultConfig Config) *BreakerRegistry {
	return &BreakerRegistry{
		defaultConfig: defaultConfig,
	}
}

// GetBreaker returns the circuit breaker for the given name, creating one
// with the default configuration if it does not already exist.
func (r *BreakerRegistry) GetBreaker(name string) *Breaker {
	if val, ok := r.breakers.Load(name); ok {
		return val.(*Breaker)
	}

	cfg := r.defaultConfig
	cfg.Name = name

	breaker := New(cfg)
	actual, _ := r.breakers.LoadOrStore(name, breaker)
	return actual.(*Breaker)
}

// UnaryClientInterceptor returns a grpc.UnaryClientInterceptor that wraps
// each outgoing gRPC call with a circuit breaker. The breaker is selected
// based on the service name extracted from the gRPC method string.
//
// If the circuit breaker is open, the call is rejected immediately without
// making a network request.
func (r *BreakerRegistry) UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		serviceName := extractServiceName(method)
		breaker := r.GetBreaker(serviceName)

		_, err := breaker.Execute(func() (interface{}, error) {
			err := invoker(ctx, method, req, reply, cc, opts...)
			return nil, err
		})
		return err
	}
}

// extractServiceName parses the service name from a gRPC full method string.
// A full method has the format "/package.ServiceName/MethodName". This
// function returns "package.ServiceName".
func extractServiceName(fullMethod string) string {
	if len(fullMethod) == 0 {
		return "unknown"
	}
	// Remove leading slash.
	name := fullMethod
	if name[0] == '/' {
		name = name[1:]
	}
	// Split on the last slash to separate service from method.
	idx := strings.LastIndex(name, "/")
	if idx < 0 {
		return name
	}
	return name[:idx]
}
