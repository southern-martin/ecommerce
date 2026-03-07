package http

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/southern-martin/ecommerce/pkg/currency"
	"github.com/southern-martin/ecommerce/pkg/i18n"
	"github.com/southern-martin/ecommerce/pkg/metrics"
	"github.com/southern-martin/ecommerce/pkg/tracing"

	_ "github.com/southern-martin/ecommerce/services/cart/docs"
)

// NewRouter creates a new Gin router with all cart routes registered.
func NewRouter(handler *CartHandler) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(tracing.GinMiddleware("cart-service"))
	router.Use(metrics.GinMiddleware("cart-service"))

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

	// Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
