package http

import (
	"github.com/gin-gonic/gin"

	"github.com/southern-martin/ecommerce/pkg/metrics"
	"github.com/southern-martin/ecommerce/pkg/tracing"
)

// NewRouter creates and configures the Gin router with all product service routes.
func NewRouter(h *Handler) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(tracing.GinMiddleware("product-service"))
	r.Use(metrics.GinMiddleware("product-service"))
	r.GET("/metrics", metrics.Handler())

	// Health check
	r.GET("/health", h.Health)

	v1 := r.Group("/api/v1")
	{
		// Public product endpoints
		products := v1.Group("/products")
		{
			products.GET("", h.ListProducts)
			products.GET("/:id", h.GetProduct)
			products.GET("/slug/:slug", h.GetProductBySlug)
		}

		// Public category endpoints
		v1.GET("/categories", h.ListCategories)

		// Seller endpoints (X-User-ID header required via Kong)
		seller := v1.Group("/seller")
		{
			sellerProducts := seller.Group("/products")
			{
				sellerProducts.POST("", h.CreateProduct)
				sellerProducts.PATCH("/:id", h.UpdateProduct)
				sellerProducts.DELETE("/:id", h.DeleteProduct)
				sellerProducts.POST("/:id/options", h.AddOption)
				sellerProducts.DELETE("/:id/options/:optionId", h.RemoveOption)
				sellerProducts.POST("/:id/variants/generate", h.GenerateVariants)
				sellerProducts.PATCH("/:id/variants/:variantId", h.UpdateVariant)
				sellerProducts.PATCH("/:id/variants/:variantId/stock", h.UpdateVariantStock)
			}
		}

		// Admin endpoints
		admin := v1.Group("/admin")
		{
			admin.POST("/categories", h.CreateCategory)
			admin.POST("/attributes", h.CreateAttributeDefinition)
			admin.GET("/attributes", h.ListAttributeDefinitions)
			admin.PATCH("/attributes/:id", h.UpdateAttributeDefinition)
			admin.DELETE("/attributes/:id", h.DeleteAttributeDefinition)
		}

		// Category attribute assignment endpoints
		categories := v1.Group("/categories")
		{
			categories.POST("/:id/attributes", h.AssignAttributeToCategory)
			categories.DELETE("/:id/attributes/:attrId", h.RemoveAttributeFromCategory)
			categories.GET("/:id/attributes", h.ListCategoryAttributes)
		}
	}

	return r
}
