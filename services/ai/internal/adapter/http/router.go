package http

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/southern-martin/ecommerce/pkg/metrics"
	"github.com/southern-martin/ecommerce/pkg/middleware"
	"github.com/southern-martin/ecommerce/pkg/tracing"

	_ "github.com/southern-martin/ecommerce/services/ai/docs"
)

// NewRouter creates a new Gin router with all AI service routes.
func NewRouter(handler *Handler) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery(), gin.Logger())
	r.Use(middleware.CorrelationID())
	r.Use(middleware.ExtractUserID())

	r.Use(tracing.GinMiddleware("ai-service"))
	r.Use(metrics.GinMiddleware("ai-service"))
	r.GET("/metrics", metrics.Handler())

	r.GET("/health", handler.Health)
	r.GET("/ready", handler.Ready)

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")
	{
		// AI Chat
		ai := v1.Group("/ai")
		{
			ai.POST("/chat", handler.Chat)
			ai.GET("/chat", handler.ListConversations)
			ai.GET("/chat/:id", handler.GetConversation)

			// Recommendations
			ai.GET("/recommendations", handler.GetRecommendations)

			// Content generation
			ai.POST("/generate-description", handler.GenerateDescription)

			// Embeddings (admin/internal)
			ai.POST("/embeddings", handler.GenerateEmbedding)
		}

		// Image search
		v1.POST("/search/image", handler.ImageSearch)
	}

	return r
}
