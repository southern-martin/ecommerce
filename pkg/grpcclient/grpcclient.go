// Package grpcclient provides a shared connection factory for gRPC clients
// with circuit breaker and tracing interceptors. It also registers a JSON
// codec required for services that use manually registered ServiceDescs
// with plain Go structs instead of proto-generated types.
package grpcclient

import (
	"encoding/json"
	"fmt"

	"github.com/southern-martin/ecommerce/pkg/circuitbreaker"
	"github.com/southern-martin/ecommerce/pkg/tracing"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding"
)

func init() {
	encoding.RegisterCodec(JSONCodec{})
}

// JSONCodec implements grpc encoding.Codec using JSON marshaling.
// This is required because the gRPC services use manually registered
// ServiceDescs with plain Go structs instead of proto-generated types.
type JSONCodec struct{}

// Marshal serializes v into a JSON-encoded byte slice.
func (JSONCodec) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// Unmarshal deserializes data into v using JSON decoding.
func (JSONCodec) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// Name returns the codec name registered with gRPC.
func (JSONCodec) Name() string {
	return "json"
}

// Dial creates a gRPC client connection with circuit breaker and tracing
// interceptors. The connection uses insecure credentials (suitable for
// internal service-to-service calls) and forces the JSON codec for
// compatibility with manually registered gRPC services.
//
// If registry is nil, no circuit breaker or tracing interceptors are added.
func Dial(addr string, registry *circuitbreaker.BreakerRegistry) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.CallContentSubtype("json")),
	}

	if registry != nil {
		opts = append(opts, grpc.WithChainUnaryInterceptor(
			registry.UnaryClientInterceptor(),
			tracing.GRPCUnaryClientInterceptor(),
		))
	}

	conn, err := grpc.NewClient(addr, opts...)
	if err != nil {
		return nil, fmt.Errorf("grpcclient: failed to connect to %s: %w", addr, err)
	}

	return conn, nil
}
