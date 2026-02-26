package http

import (
	"github.com/gin-gonic/gin"
)

// NewRouter creates and configures the Gin router with all routes.
func NewRouter(handler *Handler) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	// Health check
	r.GET("/health", handler.Health)

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
