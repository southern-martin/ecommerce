package http

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/southern-martin/ecommerce/pkg/metrics"
	"github.com/southern-martin/ecommerce/pkg/tracing"

	_ "github.com/southern-martin/ecommerce/services/affiliate/docs"
)

// NewRouter creates and configures the Gin router with all affiliate service routes.
func NewRouter(handler *Handler) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	router.Use(tracing.GinMiddleware("affiliate-service"))
	router.Use(metrics.GinMiddleware("affiliate-service"))
	router.GET("/metrics", metrics.Handler())

	// Health check
	router.GET("/health", handler.Health)
	router.GET("/ready", handler.Ready)

	// Swagger UI
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Affiliate links
		affiliate := v1.Group("/affiliate")
		{
			affiliate.POST("/links", handler.CreateLink)
			affiliate.GET("/links", handler.ListLinks)
			affiliate.GET("/stats", handler.GetStats)
			affiliate.GET("/referrals", handler.ListReferrals)
			affiliate.POST("/payout", handler.RequestPayout)
		}

		// Public referral redirect
		v1.GET("/r/:code", handler.TrackClick)

		// Admin routes
		admin := v1.Group("/admin/affiliates")
		{
			admin.GET("/program", handler.GetProgram)
			admin.PATCH("/program", handler.UpdateProgram)
			admin.GET("/payouts", handler.ListAllPayouts)
			admin.PATCH("/payouts/:id", handler.UpdatePayoutStatus)
		}
	}

	return router
}
