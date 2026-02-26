package http

import (
	"github.com/gin-gonic/gin"
)

// NewRouter creates and configures the Gin router with all chat service routes.
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
		conversations := v1.Group("/conversations")
		{
			conversations.POST("", handler.CreateConversation)
			conversations.GET("", handler.ListConversations)
			conversations.GET("/:id", handler.GetConversation)
			conversations.PATCH("/:id/archive", handler.ArchiveConversation)

			// Messages
			conversations.POST("/:id/messages", handler.SendMessage)
			conversations.GET("/:id/messages", handler.ListMessages)
			conversations.PATCH("/:id/messages/read", handler.MarkAsRead)
			conversations.GET("/:id/unread", handler.GetUnreadCount)
		}
	}

	return router
}
