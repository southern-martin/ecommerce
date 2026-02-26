package http

import (
	"github.com/gin-gonic/gin"
)

// SetupRouter configures and returns the Gin router with all auth routes.
func SetupRouter(handler *Handler) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	r.GET("/health", handler.Health)

	v1 := r.Group("/api/v1/auth")
	{
		v1.POST("/register", handler.Register)
		v1.POST("/login", handler.Login)
		v1.POST("/refresh", handler.RefreshToken)
		v1.POST("/logout", handler.Logout)
		v1.POST("/forgot-password", handler.ForgotPassword)
		v1.POST("/reset-password", handler.ResetPassword)
		v1.POST("/oauth/:provider", handler.OAuthLogin)
	}

	return r
}
