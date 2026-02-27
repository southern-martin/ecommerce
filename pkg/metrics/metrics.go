package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTPRequestsTotal counts HTTP requests by service, method, path, and status code.
	HTTPRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests",
	}, []string{"service", "method", "path", "status"})

	// HTTPRequestDuration records HTTP request latency in seconds.
	HTTPRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "HTTP request duration in seconds",
		Buckets: prometheus.DefBuckets,
	}, []string{"service", "method", "path"})

	// GRPCRequestsTotal counts gRPC calls by service, method, and status code.
	GRPCRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "grpc_requests_total",
		Help: "Total number of gRPC requests",
	}, []string{"service", "method", "status"})

	// GRPCRequestDuration records gRPC call latency in seconds.
	GRPCRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "grpc_request_duration_seconds",
		Help:    "gRPC request duration in seconds",
		Buckets: prometheus.DefBuckets,
	}, []string{"service", "method"})

	// CircuitBreakerState tracks circuit breaker state per service and name.
	// 0=closed, 1=half-open, 2=open.
	CircuitBreakerState = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "circuit_breaker_state",
		Help: "Circuit breaker state (0=closed, 1=half-open, 2=open)",
	}, []string{"service", "name"})
)
