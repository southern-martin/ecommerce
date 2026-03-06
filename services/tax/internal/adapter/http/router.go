package http

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/southern-martin/ecommerce/pkg/metrics"
	"github.com/southern-martin/ecommerce/pkg/tracing"

	_ "github.com/southern-martin/ecommerce/services/tax/docs"
)

// NewRouter creates and configures the Gin router with all routes.
func NewRouter(handler *Handler) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(tracing.GinMiddleware("tax-service"))
	r.Use(metrics.GinMiddleware("tax-service"))
	r.GET("/metrics", metrics.Handler())

	// Health check
	r.GET("/health", handler.Health)
	r.GET("/ready", handler.Ready)

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Public routes
	v1 := r.Group("/api/v1")
	{
		tax := v1.Group("/tax")
		{
			tax.GET("/zones", handler.ListZones)
			tax.POST("/calculate", handler.CalculateTax)
		}

		// Admin routes
		admin := v1.Group("/admin/tax")
		{
			admin.GET("/rules", handler.ListRules)
			admin.POST("/rules", handler.CreateRule)
			admin.PATCH("/rules/:id", handler.UpdateRule)
			admin.DELETE("/rules/:id", handler.DeleteRule)
		}
	}

	return r
}
