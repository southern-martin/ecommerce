package http

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/southern-martin/ecommerce/pkg/metrics"
	"github.com/southern-martin/ecommerce/pkg/tracing"
	_ "github.com/southern-martin/ecommerce/services/notification/docs"
)

// NewRouter creates and configures the Gin router with all notification service routes.
func NewRouter(handler *Handler) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	router.Use(tracing.GinMiddleware("notification-service"))
	router.Use(metrics.GinMiddleware("notification-service"))
	router.GET("/metrics", metrics.Handler())

	// Health check
	router.GET("/health", handler.Health)
	router.GET("/ready", handler.Ready)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
