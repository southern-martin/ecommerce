package http

import (
	"github.com/gin-gonic/gin"
)

// NewRouter creates and configures the Gin router with all media service routes.
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
		media := v1.Group("/media")
		{
			media.POST("", handler.CreateMedia)
			media.GET("", handler.ListMedia)
			media.GET("/:id", handler.GetMedia)
			media.DELETE("/:id", handler.DeleteMedia)
			media.POST("/upload-url", handler.GetUploadURL)
			media.GET("/:id/download-url", handler.GetDownloadURL)
		}
	}

	return router
}
