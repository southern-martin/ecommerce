package http

import (
	"github.com/gin-gonic/gin"
	"github.com/southern-martin/ecommerce/pkg/metrics"
	"github.com/southern-martin/ecommerce/pkg/tracing"
)

// NewRouter creates and configures the Gin router with all CMS service routes.
func NewRouter(handler *Handler) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	router.Use(tracing.GinMiddleware("cms-service"))
	router.Use(metrics.GinMiddleware("cms-service"))
	router.GET("/metrics", metrics.Handler())

	// Health check
	router.GET("/health", handler.Health)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public routes
		v1.GET("/banners", handler.ListActiveBanners)
		v1.GET("/pages/:slug", handler.GetPageBySlug)

		// Admin routes
		admin := v1.Group("/admin")
		{
			// Banner management
			admin.POST("/banners", handler.CreateBanner)
			admin.PATCH("/banners/:id", handler.UpdateBanner)
			admin.DELETE("/banners/:id", handler.DeleteBanner)
			admin.GET("/banners", handler.ListAllBanners)

			// Page management
			admin.POST("/pages", handler.CreatePage)
			admin.PATCH("/pages/:id", handler.UpdatePage)
			admin.DELETE("/pages/:id", handler.DeletePage)
			admin.GET("/pages", handler.ListAllPages)
			admin.PATCH("/pages/:id/publish", handler.PublishPage)

			// Content scheduling
			admin.POST("/content/schedule", handler.ScheduleContent)
		}
	}

	return router
}
