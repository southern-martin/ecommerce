package http

import (
	"github.com/gin-gonic/gin"

	"github.com/southern-martin/ecommerce/pkg/metrics"
	"github.com/southern-martin/ecommerce/pkg/tracing"
)

// NewRouter creates and configures the Gin router with all review service routes.
func NewRouter(handler *Handler) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	router.Use(tracing.GinMiddleware("review-service"))
	router.Use(metrics.GinMiddleware("review-service"))
	router.GET("/metrics", metrics.Handler())

	// Health check
	router.GET("/health", handler.Health)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Reviews
		reviews := v1.Group("/reviews")
		{
			reviews.POST("", handler.CreateReview)
			reviews.GET("", handler.ListProductReviews)
			reviews.GET("/:id", handler.GetReview)
			reviews.PATCH("/:id", handler.UpdateReview)
			reviews.DELETE("/:id", handler.DeleteReview)
		}

		// Product review summary
		v1.GET("/products/:product_id/reviews/summary", handler.GetProductSummary)

		// Admin routes
		admin := v1.Group("/admin")
		{
			admin.PATCH("/reviews/:id/approve", handler.ApproveReview)
		}
	}

	return router
}
