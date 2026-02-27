package http

import (
	"github.com/gin-gonic/gin"

	"github.com/southern-martin/ecommerce/pkg/metrics"
	"github.com/southern-martin/ecommerce/pkg/tracing"
)

// NewRouter creates a new Gin router with all cart routes registered.
func NewRouter(handler *CartHandler) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(tracing.GinMiddleware("cart-service"))
	router.Use(metrics.GinMiddleware("cart-service"))
	router.GET("/metrics", metrics.Handler())

	// Health check
	router.GET("/health", handler.Health)

	// Cart API v1 routes
	v1 := router.Group("/api/v1")
	{
		cart := v1.Group("/cart")
		{
			cart.GET("", handler.GetCart)
			cart.DELETE("", handler.ClearCart)

			items := cart.Group("/items")
			{
				items.POST("", handler.AddItem)
				items.PATCH("", handler.UpdateQuantity)
				items.DELETE("", handler.RemoveItem)
			}

			cart.POST("/merge", handler.MergeCart)
		}
	}

	return router
}
