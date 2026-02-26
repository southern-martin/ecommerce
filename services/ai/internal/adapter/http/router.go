package http

import "github.com/gin-gonic/gin"

// NewRouter creates a new Gin router with all AI service routes.
func NewRouter(handler *Handler) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery(), gin.Logger())

	r.GET("/health", handler.Health)

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
