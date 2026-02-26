package tracing

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	grpccodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// GinMiddleware returns a gin.HandlerFunc that creates spans for HTTP requests.
// It extracts trace context from incoming headers using W3C TraceContext propagation,
// and sets standard span attributes: http.method, http.url, http.status_code,
// and http.route.
func GinMiddleware(serviceName string) gin.HandlerFunc {
	tracer := otel.Tracer(serviceName)

	return func(c *gin.Context) {
		// Extract trace context from incoming request headers.
		propagator := otel.GetTextMapPropagator()
		ctx := propagator.Extract(c.Request.Context(), propagation.HeaderCarrier(c.Request.Header))

		spanName := c.Request.Method + " " + c.FullPath()
		if c.FullPath() == "" {
			spanName = c.Request.Method + " " + c.Request.URL.Path
		}

		ctx, span := tracer.Start(ctx, spanName,
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(
				semconv.HTTPMethod(c.Request.Method),
				semconv.HTTPTarget(c.Request.URL.String()),
				attribute.String("http.route", c.FullPath()),
				attribute.String("net.host.name", c.Request.Host),
				attribute.String("http.user_agent", c.Request.UserAgent()),
				attribute.String("net.peer.ip", c.ClientIP()),
			),
		)
		defer span.End()

		// Attach the span context to the gin context so downstream handlers
		// can create child spans.
		c.Request = c.Request.WithContext(ctx)

		c.Next()

		statusCode := c.Writer.Status()
		span.SetAttributes(semconv.HTTPStatusCode(statusCode))

		if statusCode >= 500 {
			span.SetStatus(codes.Error, fmt.Sprintf("HTTP %d", statusCode))
		} else {
			span.SetStatus(codes.Ok, "")
		}

		if len(c.Errors) > 0 {
			span.SetAttributes(attribute.String("gin.errors", c.Errors.String()))
			span.SetStatus(codes.Error, c.Errors.String())
		}
	}
}

// GRPCUnaryInterceptor returns a grpc.UnaryServerInterceptor that creates spans
// for each incoming gRPC call. It sets rpc.system, rpc.service, and rpc.method
// as span attributes.
func GRPCUnaryInterceptor() grpc.UnaryServerInterceptor {
	tracer := otel.Tracer("grpc-server")

	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Extract trace context from incoming gRPC metadata.
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			propagator := otel.GetTextMapPropagator()
			ctx = propagator.Extract(ctx, metadataCarrier(md))
		}

		serviceName, methodName := parseFullMethod(info.FullMethod)

		ctx, span := tracer.Start(ctx, info.FullMethod,
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(
				semconv.RPCSystemGRPC,
				semconv.RPCService(serviceName),
				semconv.RPCMethod(methodName),
			),
		)
		defer span.End()

		resp, err := handler(ctx, req)
		if err != nil {
			st, _ := status.FromError(err)
			span.SetAttributes(attribute.Int("rpc.grpc.status_code", int(st.Code())))
			if st.Code() != grpccodes.OK {
				span.SetStatus(codes.Error, st.Message())
			}
		} else {
			span.SetAttributes(attribute.Int("rpc.grpc.status_code", int(grpccodes.OK)))
			span.SetStatus(codes.Ok, "")
		}

		return resp, err
	}
}

// parseFullMethod splits a gRPC full method string like "/package.Service/Method"
// into the service and method components.
func parseFullMethod(fullMethod string) (string, string) {
	if len(fullMethod) == 0 {
		return "", ""
	}
	// Remove leading slash.
	name := fullMethod
	if name[0] == '/' {
		name = name[1:]
	}
	for i := len(name) - 1; i >= 0; i-- {
		if name[i] == '/' {
			return name[:i], name[i+1:]
		}
	}
	return name, ""
}

// metadataCarrier adapts gRPC metadata to the TextMapCarrier interface
// used by OpenTelemetry propagators.
type metadataCarrier metadata.MD

func (mc metadataCarrier) Get(key string) string {
	vals := metadata.MD(mc).Get(key)
	if len(vals) == 0 {
		return ""
	}
	return vals[0]
}

func (mc metadataCarrier) Set(key, value string) {
	metadata.MD(mc).Set(key, value)
}

func (mc metadataCarrier) Keys() []string {
	keys := make([]string, 0, len(mc))
	for k := range mc {
		keys = append(keys, k)
	}
	return keys
}
