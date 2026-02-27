package metrics

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// GRPCUnaryInterceptor returns a gRPC unary server interceptor that records
// call count and latency.
func GRPCUnaryInterceptor(serviceName string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()

		resp, err := handler(ctx, req)

		duration := time.Since(start).Seconds()
		st, _ := status.FromError(err)
		code := st.Code().String()

		GRPCRequestsTotal.WithLabelValues(serviceName, info.FullMethod, code).Inc()
		GRPCRequestDuration.WithLabelValues(serviceName, info.FullMethod).Observe(duration)

		return resp, err
	}
}
