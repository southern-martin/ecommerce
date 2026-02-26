package http

import (
	"github.com/gin-gonic/gin"
)

// NewRouter creates and configures the Gin router with all return service routes.
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
		// Buyer return routes
		returns := v1.Group("/returns")
		{
			returns.POST("", handler.CreateReturn)
			returns.GET("", handler.ListBuyerReturns)
			returns.GET("/:id", handler.GetReturn)
		}

		// Seller return routes
		seller := v1.Group("/seller")
		{
			seller.GET("/returns", handler.ListSellerReturns)
			seller.PATCH("/returns/:id/approve", handler.ApproveReturn)
			seller.PATCH("/returns/:id/reject", handler.RejectReturn)
			seller.PATCH("/returns/:id/status", handler.UpdateReturnStatus)
		}

		// Dispute routes
		disputes := v1.Group("/disputes")
		{
			disputes.POST("", handler.CreateDispute)
			disputes.GET("", handler.ListBuyerDisputes)
			disputes.GET("/:id", handler.GetDispute)
			disputes.POST("/:id/messages", handler.AddMessage)
		}

		// Admin dispute routes
		admin := v1.Group("/admin")
		{
			admin.GET("/disputes", handler.ListAllDisputes)
			admin.PATCH("/disputes/:id/resolve", handler.ResolveDispute)
		}
	}

	return router
}
