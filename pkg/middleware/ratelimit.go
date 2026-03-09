package middleware

import (
	"fmt"
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

// remaining returns the number of requests remaining and the unix timestamp
// when the current window resets. The caller must hold cw.mu.
func (cw *clientWindow) remaining(limit int, window time.Duration) (int, int64) {
	rem := limit - cw.count
	if rem < 0 {
		rem = 0
	}
	resetAt := cw.start.Add(window).Unix()
	return rem, resetAt
}

// rateLimitResponse is the JSON body returned when a client exceeds the limit.
type rateLimitResponse struct {
	Error      string `json:"error"`
	RetryAfter int    `json:"retry_after"`
}

// startCleanup launches a background goroutine that removes stale entries from
// the sync.Map every 5 minutes.
func startCleanup(clients *sync.Map, window time.Duration) {
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
}

// setRateLimitHeaders writes standard rate-limit headers on every response.
func setRateLimitHeaders(c *gin.Context, limit, rem int, resetAt int64) {
	c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", limit))
	c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", rem))
	c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", resetAt))
}

// RateLimiterByIP returns a Gin middleware that limits requests per client IP
// using a sliding window counter stored in memory. When a client exceeds
// requestsPerMinute within a rolling 60-second window the middleware
// responds with 429 Too Many Requests.
//
// This is the original IP-only implementation preserved for backwards
// compatibility. Prefer RateLimiter for new code.
func RateLimiterByIP(requestsPerMinute int) gin.HandlerFunc {
	var clients sync.Map // map[string]*clientWindow
	window := time.Minute

	startCleanup(&clients, window)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		key := "ip:" + ip
		now := time.Now()

		val, _ := clients.LoadOrStore(key, &clientWindow{start: now})
		cw := val.(*clientWindow)

		cw.mu.Lock()
		if now.Sub(cw.start) >= window {
			cw.count = 0
			cw.start = now
		}
		cw.count++
		allowed := cw.count <= requestsPerMinute
		rem, resetAt := cw.remaining(requestsPerMinute, window)
		cw.mu.Unlock()

		setRateLimitHeaders(c, requestsPerMinute, rem, resetAt)

		if !allowed {
			retryAfter := int(time.Until(time.Unix(resetAt, 0)).Seconds())
			if retryAfter < 1 {
				retryAfter = 1
			}
			c.Header("Retry-After", fmt.Sprintf("%d", retryAfter))
			c.AbortWithStatusJSON(http.StatusTooManyRequests, rateLimitResponse{
				Error:      "rate limit exceeded",
				RetryAfter: retryAfter,
			})
			return
		}

		c.Next()
	}
}

// RateLimiter returns a Gin middleware that applies per-user rate limiting for
// authenticated requests and per-IP rate limiting for anonymous requests.
//
// Authenticated users (identified by ContextKeyUserID in the Gin context) are
// allowed authRequestsPerMinute within a rolling 60-second window. Anonymous
// clients are identified by IP address and allowed anonRequestsPerMinute.
//
// Both buckets use the same sliding window algorithm and share one cleanup
// goroutine. Standard rate-limit response headers (X-RateLimit-Limit,
// X-RateLimit-Remaining, X-RateLimit-Reset) are set on every response.
// A Retry-After header is added when the limit is exceeded.
func RateLimiter(authRequestsPerMinute, anonRequestsPerMinute int) gin.HandlerFunc {
	var clients sync.Map // map[string]*clientWindow
	window := time.Minute

	startCleanup(&clients, window)

	return func(c *gin.Context) {
		now := time.Now()

		// Determine the bucket key and applicable limit.
		var key string
		var limit int

		userID := c.GetString(ContextKeyUserID)
		if userID != "" {
			key = "user:" + userID
			limit = authRequestsPerMinute
		} else {
			key = "ip:" + c.ClientIP()
			limit = anonRequestsPerMinute
		}

		val, _ := clients.LoadOrStore(key, &clientWindow{start: now})
		cw := val.(*clientWindow)

		cw.mu.Lock()
		if now.Sub(cw.start) >= window {
			cw.count = 0
			cw.start = now
		}
		cw.count++
		allowed := cw.count <= limit
		rem, resetAt := cw.remaining(limit, window)
		cw.mu.Unlock()

		setRateLimitHeaders(c, limit, rem, resetAt)

		if !allowed {
			retryAfter := int(time.Until(time.Unix(resetAt, 0)).Seconds())
			if retryAfter < 1 {
				retryAfter = 1
			}
			c.Header("Retry-After", fmt.Sprintf("%d", retryAfter))
			c.AbortWithStatusJSON(http.StatusTooManyRequests, rateLimitResponse{
				Error:      "rate limit exceeded",
				RetryAfter: retryAfter,
			})
			return
		}

		c.Next()
	}
}
