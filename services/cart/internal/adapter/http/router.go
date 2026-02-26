package http

import (
	"github.com/gin-gonic/gin"
)

// NewRouter creates a new Gin router with all cart routes registered.
func NewRouter(handler *CartHandler) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Recovery())

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
