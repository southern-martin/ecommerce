package http

import (
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/southern-martin/ecommerce/pkg/cache"
	"github.com/southern-martin/ecommerce/pkg/currency"
	"github.com/southern-martin/ecommerce/pkg/i18n"
	"github.com/southern-martin/ecommerce/pkg/metrics"
	"github.com/southern-martin/ecommerce/pkg/middleware"
	"github.com/southern-martin/ecommerce/pkg/tracing"

	_ "github.com/southern-martin/ecommerce/services/promotion/docs"
)

// NewRouter creates and configures the Gin router with all promotion service routes.
func NewRouter(handler *Handler, cacheClient *cache.Client) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	router.Use(middleware.CorrelationID())
	router.Use(middleware.ExtractUserID())

	router.Use(tracing.GinMiddleware("promotion-service"))
	router.Use(metrics.GinMiddleware("promotion-service"))

	// i18n: detect Accept-Language header and store resolved locale in context
	bundle := i18n.NewBundle()
	bundle.SetupDefaults()
	router.Use(i18n.GinMiddleware(bundle))

	// Multi-currency: read X-Currency header and store resolved currency in context
	router.Use(currency.GinMiddleware())
	router.GET("/metrics", metrics.Handler())

	// Health check
	router.GET("/health", handler.Health)
	router.GET("/ready", handler.Ready)

	// Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public coupon routes
		coupons := v1.Group("/coupons")
		if cacheClient != nil {
			coupons.Use(cache.CacheResponse(cacheClient, 5*time.Minute, func(c *gin.Context) string {
				return "promo:" + c.Request.URL.RequestURI()
			}))
		}
		{
			coupons.POST("/validate", handler.ValidateCoupon)
			coupons.GET("", handler.ListActiveCoupons)
		}

		// Seller coupon routes
		seller := v1.Group("/seller")
		if cacheClient != nil {
			seller.Use(cache.InvalidateCache(cacheClient, "promo:*"))
		}
		{
			seller.POST("/coupons", handler.CreateSellerCoupon)
			seller.GET("/coupons", handler.ListSellerCoupons)
			seller.GET("/coupons/:id", handler.GetSellerCoupon)
			seller.PATCH("/coupons/:id", handler.UpdateSellerCoupon)
			seller.DELETE("/coupons/:id", handler.DeleteSellerCoupon)
		}

		// Public flash sale routes
		flashSales := v1.Group("")
		if cacheClient != nil {
			flashSales.Use(cache.CacheResponse(cacheClient, 3*time.Minute, func(c *gin.Context) string {
				return "promo:" + c.Request.URL.RequestURI()
			}))
		}
		flashSales.GET("/flash-sales", handler.ListActiveFlashSales)

		// Public bundle routes
		bundles := v1.Group("")
		if cacheClient != nil {
			bundles.Use(cache.CacheResponse(cacheClient, 5*time.Minute, func(c *gin.Context) string {
				return "promo:" + c.Request.URL.RequestURI()
			}))
		}
		bundles.GET("/bundles", handler.ListActiveBundles)

		// Admin routes
		admin := v1.Group("/admin/promotions")
		if cacheClient != nil {
			admin.Use(cache.InvalidateCache(cacheClient, "promo:*"))
		}
		{
			// Admin coupon routes
			adminCoupons := admin.Group("/coupons")
			{
				adminCoupons.POST("", handler.AdminCreateCoupon)
				adminCoupons.GET("", handler.AdminListCoupons)
				adminCoupons.GET("/:id", handler.AdminGetCoupon)
				adminCoupons.PATCH("/:id", handler.AdminUpdateCoupon)
				adminCoupons.DELETE("/:id", handler.AdminDeleteCoupon)
			}

			// Admin flash sale routes
			adminFlashSales := admin.Group("/flash-sales")
			{
				adminFlashSales.POST("", handler.AdminCreateFlashSale)
				adminFlashSales.GET("", handler.AdminListFlashSales)
				adminFlashSales.GET("/:id", handler.AdminGetFlashSale)
				adminFlashSales.PATCH("/:id", handler.AdminUpdateFlashSale)
			}

			// Admin bundle routes
			adminBundles := admin.Group("/bundles")
			{
				adminBundles.POST("", handler.AdminCreateBundle)
				adminBundles.GET("", handler.AdminListBundles)
				adminBundles.GET("/:id", handler.AdminGetBundle)
				adminBundles.PATCH("/:id", handler.AdminUpdateBundle)
				adminBundles.DELETE("/:id", handler.AdminDeleteBundle)
			}
		}
	}

	return router
}
