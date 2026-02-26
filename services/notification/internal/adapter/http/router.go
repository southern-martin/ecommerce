package http

import (
	"github.com/gin-gonic/gin"
)

// NewRouter creates and configures the Gin router with all notification service routes.
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
		notifications := v1.Group("/notifications")
		{
			notifications.POST("", handler.SendNotification)
			notifications.GET("", handler.ListNotifications)
			notifications.GET("/unread-count", handler.GetUnreadCount)
			notifications.GET("/:id", handler.GetNotification)
			notifications.PATCH("/:id/read", handler.MarkAsRead)
			notifications.PATCH("/read-all", handler.MarkAllAsRead)
			notifications.GET("/preferences", handler.GetPreferences)
			notifications.PATCH("/preferences", handler.UpdatePreference)
		}
	}

	return router
}
