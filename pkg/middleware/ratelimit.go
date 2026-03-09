package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// clientWindow tracks the request count and window start for a single client.
type clientWindow struct {
	mu    sync.Mutex
	count int
	start time.Time
}

// RateLimiter returns a Gin middleware that limits requests per client IP
// using a sliding window counter stored in memory. When a client exceeds
// requestsPerMinute within a rolling 60-second window the middleware
// responds with 429 Too Many Requests using the standardized ErrorResponse.
//
// Old entries are cleaned up every 5 minutes to prevent unbounded growth.
func RateLimiter(requestsPerMinute int) gin.HandlerFunc {
	var clients sync.Map // map[string]*clientWindow
	window := time.Minute

	// Background cleanup of stale entries.
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			now := time.Now()
			clients.Range(func(key, value any) bool {
				cw := value.(*clientWindow)
				cw.mu.Lock()
				if now.Sub(cw.start) > 2*window {
					clients.Delete(key)
				}
				cw.mu.Unlock()
				return true
			})
		}
	}()

	return func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now()

		val, _ := clients.LoadOrStore(ip, &clientWindow{start: now})
		cw := val.(*clientWindow)

		cw.mu.Lock()

		// Reset window if it has elapsed.
		if now.Sub(cw.start) >= window {
			cw.count = 0
			cw.start = now
		}

		cw.count++
		allowed := cw.count <= requestsPerMinute
		cw.mu.Unlock()

		if !allowed {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, ErrorResponse{
				Error: ErrorDetail{
					Code:    "RATE_LIMIT_EXCEEDED",
					Message: "too many requests, please try again later",
					Status:  http.StatusTooManyRequests,
				},
			})
			return
		}

		c.Next()
	}
}
