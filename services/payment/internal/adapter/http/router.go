package http

import (
	"github.com/gin-gonic/gin"
)

// NewRouter creates and configures a new Gin router with all payment routes.
func NewRouter(handler *Handler) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	// Health check.
	router.GET("/health", handler.Health)

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
