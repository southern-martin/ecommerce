package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	// ContextKeyUserID is the context key for the user ID.
	ContextKeyUserID = "user_id"
	// ContextKeyUserRole is the context key for the user role.
	ContextKeyUserRole = "user_role"
)

// ExtractUserID reads the X-User-ID and X-User-Role headers and sets them in the Gin context.
func ExtractUserID() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetHeader("X-User-ID")
		if userID != "" {
			c.Set(ContextKeyUserID, userID)
		}

		userRole := c.GetHeader("X-User-Role")
		if userRole != "" {
			c.Set(ContextKeyUserRole, userRole)
		}

		c.Next()
	}
}

// RequireAuth returns 401 if X-User-ID header is missing from the request.
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetHeader("X-User-ID")
		if userID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized: missing X-User-ID header",
			})
			return
		}

		c.Set(ContextKeyUserID, userID)

		userRole := c.GetHeader("X-User-Role")
		if userRole != "" {
			c.Set(ContextKeyUserRole, userRole)
		}

		c.Next()
	}
}

// RequireRole checks that the X-User-Role header matches one of the allowed roles.
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := c.GetHeader("X-User-Role")
		if userRole == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "forbidden: missing X-User-Role header",
			})
			return
		}

		for _, role := range roles {
			if userRole == role {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "forbidden: insufficient role",
		})
	}
}
