package http

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/southern-martin/ecommerce/pkg/metrics"
	"github.com/southern-martin/ecommerce/pkg/tracing"
	_ "github.com/southern-martin/ecommerce/services/review/docs"
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
	router.GET("/ready", handler.Ready)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
