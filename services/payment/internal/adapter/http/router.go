package http

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/southern-martin/ecommerce/pkg/i18n"
	"github.com/southern-martin/ecommerce/pkg/metrics"
	"github.com/southern-martin/ecommerce/pkg/tracing"
	_ "github.com/southern-martin/ecommerce/services/payment/docs"
)

// NewRouter creates and configures a new Gin router with all payment routes.
func NewRouter(handler *Handler) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(tracing.GinMiddleware("payment-service"))
	router.Use(metrics.GinMiddleware("payment-service"))

	// i18n: detect Accept-Language header and store resolved locale in context
	bundle := i18n.NewBundle()
	bundle.SetupDefaults()
	router.Use(i18n.GinMiddleware(bundle))
	router.GET("/metrics", metrics.Handler())

	// Health check.
	router.GET("/health", handler.Health)
	router.GET("/ready", handler.Ready)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := router.Group("/api/v1/payments")
	{
		// Buyer routes.
		v1.POST("/create-intent", handler.CreatePaymentIntent)
		v1.GET("", handler.ListPayments)
		v1.GET("/:id", handler.GetPayment)

		// Webhook routes (no auth).
		v1.POST("/webhooks/stripe", handler.HandleStripeWebhook)

		// Seller wallet routes.
		v1.GET("/wallet", handler.GetWalletBalance)
		v1.GET("/wallet/transactions", handler.ListWalletTransactions)

		// Payout routes.
		v1.POST("/payouts", handler.RequestPayout)
		v1.GET("/payouts", handler.ListPayouts)
	}

	return router
}
