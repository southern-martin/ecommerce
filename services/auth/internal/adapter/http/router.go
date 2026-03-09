package http

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/southern-martin/ecommerce/pkg/metrics"
	"github.com/southern-martin/ecommerce/pkg/middleware"
	"github.com/southern-martin/ecommerce/pkg/tracing"

	_ "github.com/southern-martin/ecommerce/services/auth/docs"
)

// SetupRouter configures and returns the Gin router with all auth routes.
func SetupRouter(handler *Handler) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.Use(middleware.CorrelationID())
	r.Use(middleware.ExtractUserID())
	r.Use(tracing.GinMiddleware("auth-service"))
	r.Use(metrics.GinMiddleware("auth-service"))
	r.GET("/metrics", metrics.Handler())

	r.GET("/health", handler.Health)
	r.GET("/ready", handler.Ready)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
