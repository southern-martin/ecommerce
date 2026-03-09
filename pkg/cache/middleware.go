package cache

import (
	"bytes"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// responseRecorder captures the response body written by downstream handlers.
type responseRecorder struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

// CacheResponse returns a Gin middleware that caches JSON responses.
// keyFn determines the cache key from the request context.
// On a cache hit the cached body is returned immediately with Content-Type application/json.
// On a cache miss the response is captured; if the status is 2xx the body is stored.
func CacheResponse(client *Client, ttl time.Duration, keyFn func(*gin.Context) string) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := keyFn(c)

		// Try cache hit.
		cached, err := client.Get(c.Request.Context(), key)
		if err == nil {
			c.Data(http.StatusOK, "application/json", []byte(cached))
			c.Abort()
			return
		}

		// Cache miss — capture the response.
		rec := &responseRecorder{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = rec

		c.Next()

		// Only cache successful responses.
		if c.Writer.Status() >= 200 && c.Writer.Status() < 300 {
			_ = client.Set(c.Request.Context(), key, rec.body.String(), ttl)
		}
	}
}

// InvalidateCache returns a Gin middleware that deletes cache entries matching
// the given patterns after a successful mutation request (POST/PUT/PATCH/DELETE).
func InvalidateCache(client *Client, patterns ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Only invalidate on successful mutations.
		method := c.Request.Method
		isMutation := method == http.MethodPost || method == http.MethodPut ||
			method == http.MethodPatch || method == http.MethodDelete
		if !isMutation {
			return
		}
		if c.Writer.Status() >= 200 && c.Writer.Status() < 300 {
			for _, p := range patterns {
				_ = client.DeletePattern(c.Request.Context(), p)
			}
		}
	}
}
