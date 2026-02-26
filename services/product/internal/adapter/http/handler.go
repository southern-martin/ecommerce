package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/southern-martin/ecommerce/services/product/internal/domain"
	"github.com/southern-martin/ecommerce/services/product/internal/usecase"
)

// Handler holds HTTP handlers for the product service.
type Handler struct {
	productUC   *usecase.ProductUseCase
	categoryUC  *usecase.CategoryUseCase
	attributeUC *usecase.AttributeUseCase
	variantUC   *usecase.VariantUseCase
}

// NewHandler creates a new Handler.
func NewHandler(
	productUC *usecase.ProductUseCase,
	categoryUC *usecase.CategoryUseCase,
	attributeUC *usecase.AttributeUseCase,
	variantUC *usecase.VariantUseCase,
) *Handler {
	return &Handler{
		productUC:   productUC,
		categoryUC:  categoryUC,
		attributeUC: attributeUC,
		variantUC:   variantUC,
	}
}

// --- Health ---

// Health returns service health status.
func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// --- Public Product Endpoints ---

// ListProducts handles GET /api/v1/products
func (h *Handler) ListProducts(c *gin.Context) {
	filter := domain.ProductFilter{
		SellerID:   c.Query("seller_id"),
		CategoryID: c.Query("category_id"),
		Status:     c.Query("status"),
		Query:      c.Query("q"),
		SortBy:     c.Query("sort_by"),
	}

	if v := c.Query("min_price"); v != "" {
		if price, err := strconv.ParseInt(v, 10, 64); err == nil {
			filter.MinPrice = price
		}
	}
	if v := c.Query("max_price"); v != "" {
		if price, err := strconv.ParseInt(v, 10, 64); err == nil {
			filter.MaxPrice = price
		}
	}
	if v := c.Query("page"); v != "" {
		if page, err := strconv.Atoi(v); err == nil {
			filter.Page = page
		}
	}
	if v := c.Query("page_size"); v != "" {
		if ps, err := strconv.Atoi(v); err == nil {
			filter.PageSize = ps
		}
	}

	// Default to active products for public listing
	if filter.Status == "" {
		filter.Status = string(domain.ProductStatusActive)
	}

	products, total, err := h.productUC.ListProducts(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"products": products,
		"total":    total,
		"page":     filter.Page,
		"pageSize": filter.PageSize,
	})
}

// GetProduct handles GET /api/v1/products/:id
func (h *Handler) GetProduct(c *gin.Context) {
	id := c.Param("id")
	product, err := h.productUC.GetProduct(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, product)
}

// GetProductBySlug handles GET /api/v1/products/slug/:slug
func (h *Handler) GetProductBySlug(c *gin.Context) {
	slug := c.Param("slug")
	product, err := h.productUC.GetProductBySlug(c.Request.Context(), slug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, product)
}

// --- Public Category Endpoints ---

// ListCategories handles GET /api/v1/categories
func (h *Handler) ListCategories(c *gin.Context) {
	categories, err := h.categoryUC.GetCategories(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"categories": categories})
}

// --- Seller Product Endpoints ---

type createProductRequest struct {
	CategoryID     string                        `json:"category_id"`
	Name           string                        `json:"name" binding:"required"`
	Description    string                        `json:"description"`
	BasePriceCents int64                         `json:"base_price_cents" binding:"required"`
	Currency       string                        `json:"currency"`
	Tags           []string                      `json:"tags"`
	ImageURLs      []string                      `json:"image_urls"`
	Attributes     []attributeValueInputRequest  `json:"attributes"`
}

type attributeValueInputRequest struct {
	AttributeID string   `json:"attribute_id"`
	Value       string   `json:"value"`
	Values      []string `json:"values"`
}

// CreateProduct handles POST /api/v1/seller/products
func (h *Handler) CreateProduct(c *gin.Context) {
	sellerID := c.GetHeader("X-User-ID")
	if sellerID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing X-User-ID header"})
		return
	}

	var req createProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var attrs []usecase.AttributeValueInput
	for _, a := range req.Attributes {
		attrs = append(attrs, usecase.AttributeValueInput{
			AttributeID: a.AttributeID,
			Value:       a.Value,
			Values:      a.Values,
		})
	}

	input := usecase.CreateProductInput{
		SellerID:       sellerID,
		CategoryID:     req.CategoryID,
		Name:           req.Name,
		Description:    req.Description,
		BasePriceCents: req.BasePriceCents,
		Currency:       req.Currency,
		Tags:           req.Tags,
		ImageURLs:      req.ImageURLs,
		Attributes:     attrs,
	}

	product, err := h.productUC.CreateProduct(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, product)
}

type updateProductRequest struct {
	Name           *string            `json:"name"`
	Description    *string            `json:"description"`
	BasePriceCents *int64             `json:"base_price_cents"`
	Currency       *string            `json:"currency"`
	Status         *domain.ProductStatus `json:"status"`
	Tags           []string           `json:"tags"`
	ImageURLs      []string           `json:"image_urls"`
	CategoryID     *string            `json:"category_id"`
}

// UpdateProduct handles PATCH /api/v1/seller/products/:id
func (h *Handler) UpdateProduct(c *gin.Context) {
	sellerID := c.GetHeader("X-User-ID")
	if sellerID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing X-User-ID header"})
		return
	}

	id := c.Param("id")
	var req updateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input := usecase.UpdateProductInput{
		Name:           req.Name,
		Description:    req.Description,
		BasePriceCents: req.BasePriceCents,
		Currency:       req.Currency,
		Status:         req.Status,
		Tags:           req.Tags,
		ImageURLs:      req.ImageURLs,
		CategoryID:     req.CategoryID,
	}

	product, err := h.productUC.UpdateProduct(c.Request.Context(), id, sellerID, input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

// DeleteProduct handles DELETE /api/v1/seller/products/:id
func (h *Handler) DeleteProduct(c *gin.Context) {
	sellerID := c.GetHeader("X-User-ID")
	if sellerID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing X-User-ID header"})
		return
	}

	id := c.Param("id")
	if err := h.productUC.DeleteProduct(c.Request.Context(), id, sellerID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "product deleted"})
}

// --- Seller Option Endpoints ---

type addOptionRequest struct {
	Name      string                    `json:"name" binding:"required"`
	SortOrder int                       `json:"sort_order"`
	Values    []optionValueInputRequest `json:"values" binding:"required"`
}

type optionValueInputRequest struct {
	Value     string `json:"value" binding:"required"`
	ColorHex  string `json:"color_hex"`
	SortOrder int    `json:"sort_order"`
}

// AddOption handles POST /api/v1/seller/products/:id/options
func (h *Handler) AddOption(c *gin.Context) {
	sellerID := c.GetHeader("X-User-ID")
	if sellerID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing X-User-ID header"})
		return
	}

	productID := c.Param("id")
	var req addOptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var values []usecase.OptionValueInput
	for _, v := range req.Values {
		values = append(values, usecase.OptionValueInput{
			Value:     v.Value,
			ColorHex:  v.ColorHex,
			SortOrder: v.SortOrder,
		})
	}

	option, err := h.variantUC.AddOption(c.Request.Context(), productID, sellerID, usecase.AddOptionInput{
		Name:      req.Name,
		SortOrder: req.SortOrder,
		Values:    values,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, option)
}

// RemoveOption handles DELETE /api/v1/seller/products/:id/options/:optionId
func (h *Handler) RemoveOption(c *gin.Context) {
	sellerID := c.GetHeader("X-User-ID")
	if sellerID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing X-User-ID header"})
		return
	}

	productID := c.Param("id")
	optionID := c.Param("optionId")

	if err := h.variantUC.RemoveOption(c.Request.Context(), productID, optionID, sellerID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "option removed"})
}

// --- Seller Variant Endpoints ---

// GenerateVariants handles POST /api/v1/seller/products/:id/variants/generate
func (h *Handler) GenerateVariants(c *gin.Context) {
	sellerID := c.GetHeader("X-User-ID")
	if sellerID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing X-User-ID header"})
		return
	}

	productID := c.Param("id")
	variants, err := h.variantUC.GenerateVariants(c.Request.Context(), productID, sellerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"variants": variants})
}

type updateVariantRequest struct {
	Name           *string  `json:"name"`
	PriceCents     *int64   `json:"price_cents"`
	CompareAtCents *int64   `json:"compare_at_cents"`
	CostCents      *int64   `json:"cost_cents"`
	WeightGrams    *int     `json:"weight_grams"`
	IsActive       *bool    `json:"is_active"`
	ImageURLs      []string `json:"image_urls"`
	Barcode        *string  `json:"barcode"`
	LowStockAlert  *int     `json:"low_stock_alert"`
}

// UpdateVariant handles PATCH /api/v1/seller/products/:id/variants/:variantId
func (h *Handler) UpdateVariant(c *gin.Context) {
	sellerID := c.GetHeader("X-User-ID")
	if sellerID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing X-User-ID header"})
		return
	}

	productID := c.Param("id")
	variantID := c.Param("variantId")

	var req updateVariantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input := usecase.UpdateVariantInput{
		Name:           req.Name,
		PriceCents:     req.PriceCents,
		CompareAtCents: req.CompareAtCents,
		CostCents:      req.CostCents,
		WeightGrams:    req.WeightGrams,
		IsActive:       req.IsActive,
		ImageURLs:      req.ImageURLs,
		Barcode:        req.Barcode,
		LowStockAlert:  req.LowStockAlert,
	}

	variant, err := h.variantUC.UpdateVariant(c.Request.Context(), productID, variantID, sellerID, input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, variant)
}

type updateStockRequest struct {
	Delta int `json:"delta" binding:"required"`
}

// UpdateVariantStock handles PATCH /api/v1/seller/products/:id/variants/:variantId/stock
func (h *Handler) UpdateVariantStock(c *gin.Context) {
	sellerID := c.GetHeader("X-User-ID")
	if sellerID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing X-User-ID header"})
		return
	}

	productID := c.Param("id")
	variantID := c.Param("variantId")

	var req updateStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.variantUC.UpdateStock(c.Request.Context(), productID, variantID, sellerID, req.Delta); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "stock updated"})
}

// --- Admin Category Endpoints ---

type createCategoryRequest struct {
	Name      string `json:"name" binding:"required"`
	ParentID  string `json:"parent_id"`
	SortOrder int    `json:"sort_order"`
	ImageURL  string `json:"image_url"`
}

// CreateCategory handles POST /api/v1/admin/categories
func (h *Handler) CreateCategory(c *gin.Context) {
	var req createCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category, err := h.categoryUC.CreateCategory(c.Request.Context(), usecase.CreateCategoryInput{
		Name:      req.Name,
		ParentID:  req.ParentID,
		SortOrder: req.SortOrder,
		ImageURL:  req.ImageURL,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, category)
}

// --- Admin Attribute Endpoints ---

type createAttributeRequest struct {
	Name       string              `json:"name" binding:"required"`
	Type       domain.AttributeType `json:"type" binding:"required"`
	Required   bool                `json:"required"`
	Filterable bool                `json:"filterable"`
	Options    []string            `json:"options"`
	Unit       string              `json:"unit"`
	SortOrder  int                 `json:"sort_order"`
}

// CreateAttributeDefinition handles POST /api/v1/admin/attributes
func (h *Handler) CreateAttributeDefinition(c *gin.Context) {
	var req createAttributeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	attr, err := h.attributeUC.CreateAttributeDefinition(c.Request.Context(), usecase.CreateAttributeInput{
		Name:       req.Name,
		Type:       req.Type,
		Required:   req.Required,
		Filterable: req.Filterable,
		Options:    req.Options,
		Unit:       req.Unit,
		SortOrder:  req.SortOrder,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, attr)
}

// ListAttributeDefinitions handles GET /api/v1/admin/attributes
func (h *Handler) ListAttributeDefinitions(c *gin.Context) {
	attrs, err := h.attributeUC.ListAttributeDefinitions(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"attributes": attrs})
}

type updateAttributeRequest struct {
	Name       *string  `json:"name"`
	Required   *bool    `json:"required"`
	Filterable *bool    `json:"filterable"`
	Options    []string `json:"options"`
	Unit       *string  `json:"unit"`
	SortOrder  *int     `json:"sort_order"`
}

// UpdateAttributeDefinition handles PATCH /api/v1/admin/attributes/:id
func (h *Handler) UpdateAttributeDefinition(c *gin.Context) {
	id := c.Param("id")
	var req updateAttributeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	attr, err := h.attributeUC.UpdateAttributeDefinition(c.Request.Context(), id, usecase.UpdateAttributeInput{
		Name:       req.Name,
		Required:   req.Required,
		Filterable: req.Filterable,
		Options:    req.Options,
		Unit:       req.Unit,
		SortOrder:  req.SortOrder,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, attr)
}

// DeleteAttributeDefinition handles DELETE /api/v1/admin/attributes/:id
func (h *Handler) DeleteAttributeDefinition(c *gin.Context) {
	id := c.Param("id")
	if err := h.attributeUC.DeleteAttributeDefinition(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "attribute definition deleted"})
}

// --- Category Attribute Assignment Endpoints ---

type assignAttributeRequest struct {
	AttributeID string `json:"attribute_id" binding:"required"`
	SortOrder   int    `json:"sort_order"`
}

// AssignAttributeToCategory handles POST /api/v1/categories/:id/attributes
func (h *Handler) AssignAttributeToCategory(c *gin.Context) {
	categoryID := c.Param("id")
	var req assignAttributeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.attributeUC.AssignAttributeToCategory(c.Request.Context(), categoryID, req.AttributeID, req.SortOrder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "attribute assigned to category"})
}

// RemoveAttributeFromCategory handles DELETE /api/v1/categories/:id/attributes/:attrId
func (h *Handler) RemoveAttributeFromCategory(c *gin.Context) {
	categoryID := c.Param("id")
	attrID := c.Param("attrId")

	if err := h.attributeUC.RemoveAttributeFromCategory(c.Request.Context(), categoryID, attrID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "attribute removed from category"})
}

// ListCategoryAttributes handles GET /api/v1/categories/:id/attributes
func (h *Handler) ListCategoryAttributes(c *gin.Context) {
	categoryID := c.Param("id")
	attrs, err := h.attributeUC.ListCategoryAttributes(c.Request.Context(), categoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"attributes": attrs})
}
