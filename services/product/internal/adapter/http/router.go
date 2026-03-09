package http

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/southern-martin/ecommerce/pkg/cache"
	"github.com/southern-martin/ecommerce/pkg/currency"
	"github.com/southern-martin/ecommerce/pkg/i18n"
	"github.com/southern-martin/ecommerce/pkg/metrics"
	"github.com/southern-martin/ecommerce/pkg/middleware"
	"github.com/southern-martin/ecommerce/pkg/tracing"

	_ "github.com/southern-martin/ecommerce/services/product/docs"
)

// NewRouter creates and configures the Gin router with all product service routes.
func NewRouter(h *Handler, cacheClient *cache.Client, logger zerolog.Logger) *gin.Engine {
	r := gin.New()
	r.Use(middleware.RequestLogging(logger))
	r.Use(middleware.RecoveryWithLogger(logger))
	r.Use(middleware.CorrelationID())
	r.Use(middleware.ExtractUserID())
	r.Use(tracing.GinMiddleware("product-service"))
	r.Use(metrics.GinMiddleware("product-service"))

	// i18n: detect Accept-Language header and store resolved locale in context
	bundle := i18n.NewBundle()
	bundle.SetupDefaults()
	r.Use(i18n.GinMiddleware(bundle))

	// Multi-currency: read X-Currency header and store resolved currency in context
	r.Use(currency.GinMiddleware())
	r.GET("/metrics", metrics.Handler())

	// Health check
	r.GET("/health", h.Health)
	r.GET("/ready", h.Ready)

	// Swagger docs
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")
	{
		// Public product endpoints
		products := v1.Group("/products")
		if cacheClient != nil {
			products.Use(cache.CacheResponse(cacheClient, 5*time.Minute, func(c *gin.Context) string {
				return "product:" + c.Request.URL.RequestURI()
			}))
		}
		{
			products.GET("", h.ListProducts)
			products.GET("/:id", h.GetProduct)
			products.GET("/slug/:slug", h.GetProductBySlug)
			products.GET("/:id/options", h.ListProductOptions)
			products.GET("/:id/variants", h.ListProductVariants)
			products.GET("/:id/variants/:variantId", h.GetVariant)
			products.GET("/:id/attributes", h.GetProductAttributes)
		}

		// Public category endpoints
		categories_public := v1.Group("")
		if cacheClient != nil {
			categories_public.Use(cache.CacheResponse(cacheClient, 10*time.Minute, func(c *gin.Context) string {
				return "category:" + c.Request.URL.RequestURI()
			}))
		}
		categories_public.GET("/categories", h.ListCategories)

		// Seller endpoints (X-User-ID header required via Kong)
		seller := v1.Group("/seller")
		if cacheClient != nil {
			seller.Use(cache.InvalidateCache(cacheClient, "product:*", "category:*"))
		}
		{
			sellerProducts := seller.Group("/products")
			{
				sellerProducts.POST("", h.CreateProduct)
				sellerProducts.GET("/:id", h.GetProduct)
				sellerProducts.PATCH("/:id", h.UpdateProduct)
				sellerProducts.DELETE("/:id", h.DeleteProduct)
				sellerProducts.POST("/:id/options", h.AddOption)
				sellerProducts.DELETE("/:id/options/:optionId", h.RemoveOption)
				sellerProducts.POST("/:id/variants/generate", h.GenerateVariants)
				sellerProducts.PATCH("/:id/variants/:variantId", h.UpdateVariant)
				sellerProducts.PATCH("/:id/variants/:variantId/stock", h.UpdateVariantStock)
				sellerProducts.GET("/:id/options", h.ListProductOptions)
				sellerProducts.GET("/:id/variants", h.ListProductVariants)
				sellerProducts.PUT("/:id/attributes", h.SetProductAttributes)
				sellerProducts.GET("/:id/attributes", h.GetProductAttributes)
			}
		}

		// Public attribute group endpoints
		attrGroups := v1.Group("/attribute-groups")
		if cacheClient != nil {
			attrGroups.Use(cache.CacheResponse(cacheClient, 10*time.Minute, func(c *gin.Context) string {
				return "attr-group:" + c.Request.URL.RequestURI()
			}))
		}
		{
			attrGroups.GET("", h.ListAttributeGroups)
			attrGroups.GET("/:id", h.GetAttributeGroup)
			attrGroups.GET("/:id/attributes", h.ListGroupAttributes)
		}

		// Admin endpoints
		admin := v1.Group("/admin")
		if cacheClient != nil {
			admin.Use(cache.InvalidateCache(cacheClient, "product:*", "category:*", "attr-group:*"))
		}
		{
			// Admin product management (no seller ownership check)
			adminProducts := admin.Group("/products")
			{
				adminProducts.GET("", h.AdminListProducts)
				adminProducts.GET("/:id", h.AdminGetProduct)
				adminProducts.PATCH("/:id", h.AdminUpdateProduct)
				adminProducts.DELETE("/:id", h.AdminDeleteProduct)
				adminProducts.GET("/:id/options", h.AdminListOptions)
				adminProducts.POST("/:id/options", h.AdminAddOption)
				adminProducts.DELETE("/:id/options/:optionId", h.AdminRemoveOption)
				adminProducts.GET("/:id/variants", h.AdminListVariants)
				adminProducts.POST("/:id/variants/generate", h.AdminGenerateVariants)
				adminProducts.PATCH("/:id/variants/:variantId", h.AdminUpdateVariant)
				adminProducts.PATCH("/:id/variants/:variantId/stock", h.AdminUpdateVariantStock)
				adminProducts.PUT("/:id/attributes", h.AdminSetProductAttributes)
				adminProducts.GET("/:id/attributes", h.AdminGetProductAttributes)
			}

			admin.POST("/categories", h.CreateCategory)
			admin.PATCH("/categories/:id", h.UpdateCategory)
			admin.DELETE("/categories/:id", h.DeleteCategory)
			admin.POST("/attributes", h.CreateAttributeDefinition)
			admin.GET("/attributes", h.ListAttributeDefinitions)
			admin.PATCH("/attributes/:id", h.UpdateAttributeDefinition)
			admin.DELETE("/attributes/:id", h.DeleteAttributeDefinition)

			// Admin attribute group management
			adminAttrGroups := admin.Group("/attribute-groups")
			{
				adminAttrGroups.POST("", h.CreateAttributeGroup)
				adminAttrGroups.GET("", h.ListAttributeGroups)
				adminAttrGroups.GET("/:id", h.GetAttributeGroup)
				adminAttrGroups.PATCH("/:id", h.UpdateAttributeGroup)
				adminAttrGroups.DELETE("/:id", h.DeleteAttributeGroup)
				adminAttrGroups.POST("/:id/attributes", h.AddAttributeToGroup)
				adminAttrGroups.DELETE("/:id/attributes/:attrId", h.RemoveAttributeFromGroup)
				adminAttrGroups.GET("/:id/attributes", h.ListGroupAttributes)
			}
		}

		// Category attribute assignment endpoints (legacy — kept for backward compatibility)
		categories := v1.Group("/categories")
		{
			categories.POST("/:id/attributes", h.AssignAttributeToCategory)
			categories.DELETE("/:id/attributes/:attrId", h.RemoveAttributeFromCategory)
			categories.GET("/:id/attributes", h.ListCategoryAttributes)
		}
	}

	return r
}
