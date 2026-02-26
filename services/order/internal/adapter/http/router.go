package http

import (
	"github.com/gin-gonic/gin"
)

// NewRouter creates and configures the Gin router with all order service routes.
func NewRouter(handler *Handler) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// Health check
	router.GET("/health", handler.Health)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Buyer order routes
		orders := v1.Group("/orders")
		{
			orders.POST("", handler.CreateOrder)
			orders.GET("", handler.ListOrders)
			orders.GET("/:id", handler.GetOrder)
			orders.POST("/:id/cancel", handler.CancelOrder)
		}

		// Seller order routes
		seller := v1.Group("/seller")
		{
			seller.GET("/orders", handler.ListSellerOrders)
			seller.GET("/orders/:id", handler.GetSellerOrder)
			seller.PATCH("/orders/:id/status", handler.UpdateSellerOrderStatus)
		}
	}

	return router
}
