package http

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/southern-martin/ecommerce/pkg/metrics"
	"github.com/southern-martin/ecommerce/pkg/middleware"
	"github.com/southern-martin/ecommerce/pkg/tracing"

	_ "github.com/southern-martin/ecommerce/services/return/docs"
)

// NewRouter creates and configures the Gin router with all return service routes.
func NewRouter(handler *Handler, logger zerolog.Logger) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(middleware.RecoveryWithLogger(logger))
	router.Use(middleware.RequestLogging(logger))
	router.Use(middleware.CorrelationID())
	router.Use(middleware.ExtractUserID())
	router.Use(tracing.GinMiddleware("return-service"))
	router.Use(metrics.GinMiddleware("return-service"))
	router.GET("/metrics", metrics.Handler())

	// Health check
	router.GET("/health", handler.Health)
	router.GET("/ready", handler.Ready)

	// Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
