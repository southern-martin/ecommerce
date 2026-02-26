package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	// HeaderRequestID is the header key for the request/correlation ID.
	HeaderRequestID = "X-Request-ID"
	// ContextKeyRequestID is the context key for the request ID.
	ContextKeyRequestID = "request_id"
)

// CorrelationID extracts or generates a request ID, sets it in the context and response header.
func CorrelationID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader(HeaderRequestID)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		c.Set(ContextKeyRequestID, requestID)
		c.Header(HeaderRequestID, requestID)

		c.Next()
	}
}
