package http

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/southern-martin/ecommerce/pkg/metrics"
	"github.com/southern-martin/ecommerce/pkg/middleware"
	"github.com/southern-martin/ecommerce/pkg/tracing"

	_ "github.com/southern-martin/ecommerce/services/user/docs"
)

// NewRouter creates a new Gin engine with all user service routes.
func NewRouter(h *Handler, logger zerolog.Logger) *gin.Engine {
	r := gin.New()
	r.Use(middleware.RecoveryWithLogger(logger))
	r.Use(middleware.RequestLogging(logger))
	r.Use(middleware.CorrelationID())
	r.Use(middleware.ExtractUserID())
	r.Use(tracing.GinMiddleware("user-service"))
	r.Use(metrics.GinMiddleware("user-service"))
	r.GET("/metrics", metrics.Handler())

	// Health check
	r.GET("/health", h.Health)
	r.GET("/ready", h.Ready)

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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

	// Wishlist routes
	wishlist := v1.Group("/wishlist")
	wishlist.Use(middleware.RequireAuth())
	{
		wishlist.GET("", h.GetWishlist)
		wishlist.POST("", h.AddToWishlist)
		wishlist.DELETE("/:productId", h.RemoveFromWishlist)
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
