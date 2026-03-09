package http

import (
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/southern-martin/ecommerce/pkg/cache"
	"github.com/southern-martin/ecommerce/pkg/metrics"
	"github.com/southern-martin/ecommerce/pkg/middleware"
	"github.com/southern-martin/ecommerce/pkg/tracing"

	_ "github.com/southern-martin/ecommerce/services/search/docs"
)

// NewRouter creates and configures the Gin router with all search service routes.
func NewRouter(handler *Handler, cacheClient *cache.Client) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	router.Use(middleware.CorrelationID())
	router.Use(middleware.ExtractUserID())
	router.Use(tracing.GinMiddleware("search-service"))
	router.Use(metrics.GinMiddleware("search-service"))
	router.GET("/metrics", metrics.Handler())

	// Health check
	router.GET("/health", handler.Health)
	router.GET("/ready", handler.Ready)

	// Swagger UI
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Search routes
		search := v1.Group("/search")
		if cacheClient != nil {
			search.Use(cache.CacheResponse(cacheClient, 2*time.Minute, func(c *gin.Context) string {
				return "search:" + c.Request.URL.RequestURI()
			}))
		}
		{
			search.GET("", handler.Search)
			search.GET("/suggest", handler.Suggest)
		}

		// Admin index management routes
		admin := v1.Group("/admin/search")
		if cacheClient != nil {
			admin.Use(cache.InvalidateCache(cacheClient, "search:*"))
		}
		{
			admin.POST("/index", handler.IndexProduct)
			admin.DELETE("/index/:product_id", handler.DeleteProduct)
		}
	}

	return router
}
