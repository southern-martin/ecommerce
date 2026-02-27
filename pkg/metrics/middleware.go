package metrics

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// GinMiddleware returns a Gin middleware that records request count and latency.
func GinMiddleware(serviceName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start).Seconds()
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}
		status := fmt.Sprintf("%d", c.Writer.Status())

		HTTPRequestsTotal.WithLabelValues(serviceName, c.Request.Method, path, status).Inc()
		HTTPRequestDuration.WithLabelValues(serviceName, c.Request.Method, path).Observe(duration)
	}
}

// Handler returns a Gin handler that serves Prometheus metrics.
func Handler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
