package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// defaultAllowedOrigins used when CORS_ALLOWED_ORIGINS is not set.
var defaultAllowedOrigins = []string{
	"http://localhost:28080",
	"http://localhost:3000",
	"http://localhost:5173",
}

// CORS returns a Gin middleware that handles Cross-Origin Resource Sharing.
// It reads allowed origins from the CORS_ALLOWED_ORIGINS env var (comma-separated).
// When credentials are enabled, the spec forbids using "*" as the origin,
// so we reflect the request Origin if it matches the allowed list.
func CORS() gin.HandlerFunc {
	allowed := parseAllowedOrigins()
	allowedSet := make(map[string]struct{}, len(allowed))
	for _, o := range allowed {
		allowedSet[o] = struct{}{}
	}

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		if _, ok := allowedSet[origin]; ok {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-User-ID, X-User-Role, X-Request-ID, X-Currency, Accept-Language")
		c.Header("Access-Control-Expose-Headers", "Content-Length, X-Request-ID")
		c.Header("Access-Control-Max-Age", "86400")
		c.Header("Vary", "Origin")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func parseAllowedOrigins() []string {
	env := os.Getenv("CORS_ALLOWED_ORIGINS")
	if env == "" {
		return defaultAllowedOrigins
	}
	parts := strings.Split(env, ",")
	origins := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			origins = append(origins, p)
		}
	}
	return origins
}
