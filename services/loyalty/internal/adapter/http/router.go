package http

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/southern-martin/ecommerce/pkg/metrics"
	"github.com/southern-martin/ecommerce/pkg/tracing"
	_ "github.com/southern-martin/ecommerce/services/loyalty/docs"
)

// NewRouter creates and configures the Gin router with all loyalty service routes.
func NewRouter(handler *Handler) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	router.Use(tracing.GinMiddleware("loyalty-service"))
	router.Use(metrics.GinMiddleware("loyalty-service"))
	router.GET("/metrics", metrics.Handler())

	// Health check
	router.GET("/health", handler.Health)
	router.GET("/ready", handler.Ready)

	// Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		loyalty := v1.Group("/loyalty")
		{
			loyalty.GET("/membership", handler.GetMembership)
			loyalty.GET("/points", handler.GetPointsBalance)
			loyalty.GET("/transactions", handler.ListTransactions)
			loyalty.POST("/redeem", handler.RedeemPoints)
			loyalty.GET("/tiers", handler.ListTiers)
		}
	}

	return router
}
