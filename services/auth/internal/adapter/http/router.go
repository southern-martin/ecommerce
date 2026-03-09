package http

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/southern-martin/ecommerce/pkg/metrics"
	"github.com/southern-martin/ecommerce/pkg/middleware"
	"github.com/southern-martin/ecommerce/pkg/tracing"

	_ "github.com/southern-martin/ecommerce/services/auth/docs"
)

// SetupRouter configures and returns the Gin router with all auth routes.
func SetupRouter(handler *Handler, logger zerolog.Logger) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(middleware.RecoveryWithLogger(logger))
	r.Use(middleware.RequestLogging(logger))
	r.Use(middleware.CorrelationID())
	r.Use(middleware.ExtractUserID())
	r.Use(middleware.APIVersion())
	r.Use(tracing.GinMiddleware("auth-service"))
	r.Use(metrics.GinMiddleware("auth-service"))
	r.GET("/metrics", metrics.Handler())

	r.GET("/health", handler.Health)
	r.GET("/ready", handler.Ready)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// --- v1 routes (current) ---------------------------------------------------
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

	// --- Versioning guide -------------------------------------------------------
	// When v2 routes are introduced:
	//
	// 1. Add deprecation middleware to the v1 group:
	//    v1.Use(middleware.DeprecateVersion(
	//        "Sat, 01 Jun 2027 00:00:00 GMT", // sunset date (RFC 7231)
	//        "/api/v2/auth",                   // successor path
	//    ))
	//
	// 2. Create the v2 route group:
	//    v2 := r.Group("/api/v2/auth")
	//    {
	//        v2.POST("/register", handler.RegisterV2)
	//        // ... v2 handlers
	//    }
	//
	// 3. Optionally, catch unsupported versions:
	//    r.Any("/api/v3/*path", middleware.VersionNotFound([]string{"v1", "v2"}))
	// ---------------------------------------------------------------------------

	return r
}
