package http

import (
	"github.com/gin-gonic/gin"

	"github.com/southern-martin/ecommerce/pkg/metrics"
	"github.com/southern-martin/ecommerce/pkg/middleware"
	"github.com/southern-martin/ecommerce/pkg/tracing"
)

// NewRouter creates a new Gin engine with all user service routes.
func NewRouter(h *Handler) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.ExtractUserID())
	r.Use(tracing.GinMiddleware("user-service"))
	r.Use(metrics.GinMiddleware("user-service"))
	r.GET("/metrics", metrics.Handler())

	// Health check
	r.GET("/health", h.Health)

	v1 := r.Group("/api/v1")

	// User profile routes (require auth)
	users := v1.Group("/users")
	users.Use(middleware.RequireAuth())
	{
		users.GET("/me", h.GetProfile)
		users.PATCH("/me", h.UpdateProfile)

		// Address routes
		users.POST("/me/addresses", h.CreateAddress)
		users.GET("/me/addresses", h.ListAddresses)
		users.PATCH("/me/addresses/:id", h.UpdateAddress)
		users.DELETE("/me/addresses/:id", h.DeleteAddress)
		users.PATCH("/me/addresses/:id/default", h.SetDefaultAddress)

		// Following routes
		users.GET("/me/following", h.ListFollowed)

		// Follow/unfollow a seller
		users.POST("/:id/follow", h.FollowSeller)
		users.DELETE("/:id/follow", h.UnfollowSeller)
	}

	// Seller routes
	sellers := v1.Group("/sellers")
	sellers.Use(middleware.RequireAuth())
	{
		sellers.POST("", h.CreateSeller)
		sellers.GET("/:id", h.GetSeller)
		sellers.PATCH("/me", h.UpdateSeller)
		sellers.GET("/:id/followers/count", h.GetFollowerCount)
	}

	// Admin routes
	admin := v1.Group("/admin")
	admin.Use(middleware.RequireAuth())
	admin.Use(middleware.RequireRole("admin"))
	{
		admin.POST("/sellers/:id/approve", h.ApproveSeller)
	}

	return r
}
