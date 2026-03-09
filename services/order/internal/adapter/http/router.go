package http

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/southern-martin/ecommerce/pkg/currency"
	"github.com/southern-martin/ecommerce/pkg/i18n"
	"github.com/southern-martin/ecommerce/pkg/metrics"
	"github.com/southern-martin/ecommerce/pkg/middleware"
	"github.com/southern-martin/ecommerce/pkg/tracing"
	_ "github.com/southern-martin/ecommerce/services/order/docs"
)

// NewRouter creates and configures the Gin router with all order service routes.
func NewRouter(handler *Handler, logger zerolog.Logger) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(middleware.RecoveryWithLogger(logger))
	router.Use(middleware.RequestLogging(logger))
	router.Use(middleware.CorrelationID())
	router.Use(middleware.ExtractUserID())
	router.Use(tracing.GinMiddleware("order-service"))
	router.Use(metrics.GinMiddleware("order-service"))

	// i18n: detect Accept-Language header and store resolved locale in context
	bundle := i18n.NewBundle()
	bundle.SetupDefaults()
	router.Use(i18n.GinMiddleware(bundle))

	// Multi-currency: read X-Currency header and store resolved currency in context
	router.Use(currency.GinMiddleware())
	router.GET("/metrics", metrics.Handler())

	// Health check
	router.GET("/health", handler.Health)
	router.GET("/ready", handler.Ready)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
