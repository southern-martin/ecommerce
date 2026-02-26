// Package circuitbreaker provides a production-ready circuit breaker
// implementation built on top of sony/gobreaker. It includes a registry
// for managing per-service breakers and a gRPC client interceptor.
//
// See breaker.go for the core Breaker type and grpc_interceptor.go for
// the BreakerRegistry and gRPC interceptor.
package circuitbreaker
