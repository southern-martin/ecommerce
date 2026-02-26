package http

import (
	"github.com/gin-gonic/gin"
)

// NewRouter creates and configures the Gin router with all shipping service routes.
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
		// Shipping rates
		v1.POST("/shipping/rates", handler.GetShippingRates)

		// Shipments
		shipments := v1.Group("/shipments")
		{
			shipments.POST("", handler.CreateShipment)
			shipments.GET("/:id", handler.GetShipment)
			shipments.POST("/:id/label", handler.GenerateLabel)
			shipments.POST("/:id/tracking", handler.AddTrackingEvent)
		}

		// Public tracking
		v1.GET("/tracking/:tracking_number", handler.GetTracking)

		// Seller routes
		seller := v1.Group("/seller")
		{
			seller.GET("/shipments", handler.ListSellerShipments)
			seller.POST("/carriers", handler.SetupSellerCarrier)
			seller.GET("/carriers", handler.GetSellerCarriers)
		}

		// Admin routes
		admin := v1.Group("/admin")
		{
			admin.POST("/carriers", handler.CreateCarrier)
			admin.PATCH("/carriers/:code", handler.UpdateCarrier)
			admin.GET("/carriers", handler.ListCarriers)
		}
	}

	return router
}
