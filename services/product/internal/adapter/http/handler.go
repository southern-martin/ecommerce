package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/southern-martin/ecommerce/services/product/internal/domain"
	"github.com/southern-martin/ecommerce/services/product/internal/usecase"
)

// Handler holds HTTP handlers for the product service.
type Handler struct {
	productUC        *usecase.ProductUseCase
	categoryUC       *usecase.CategoryUseCase
	attributeUC      *usecase.AttributeUseCase
	variantUC        *usecase.VariantUseCase
	attributeGroupUC *usecase.AttributeGroupUseCase
	db               *gorm.DB
}

// NewHandler creates a new Handler.
func NewHandler(
	productUC *usecase.ProductUseCase,
	categoryUC *usecase.CategoryUseCase,
	attributeUC *usecase.AttributeUseCase,
	variantUC *usecase.VariantUseCase,
	attributeGroupUC *usecase.AttributeGroupUseCase,
	db *gorm.DB,
) *Handler {
	return &Handler{
		productUC:        productUC,
		categoryUC:       categoryUC,
		attributeUC:      attributeUC,
		variantUC:        variantUC,
		attributeGroupUC: attributeGroupUC,
		db:               db,
	}
}

// --- Health ---

// Health returns service health status.
func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// Ready handles GET /ready — deep health check including database connectivity.
func (h *Handler) Ready(c *gin.Context) {
	sqlDB, err := h.db.DB()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not ready", "error": "db connection lost"})
		return
	}
	if err := sqlDB.Ping(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not ready", "error": "db ping failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ready"})
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
	CategoryID       string                        `json:"category_id"`
	AttributeGroupID string                        `json:"attribute_group_id"`
	Name             string                        `json:"name" binding:"required"`
	Description      string                        `json:"description"`
	BasePriceCents   int64                         `json:"base_price_cents" binding:"required"`
	Currency         string                        `json:"currency"`
	ProductType      string                        `json:"product_type"`
	StockQuantity    int                           `json:"stock_quantity"`
	Tags             []string                      `json:"tags"`
	ImageURLs        []string                      `json:"image_urls"`
	Attributes       []attributeValueInputRequest  `json:"attributes"`
}

type attributeValueInputRequest struct {
	AttributeID    string   `json:"attribute_id"`
	Value          string   `json:"value"`
	Values         []string `json:"values"`
	OptionValueID  string   `json:"option_value_id"`
	OptionValueIDs []string `json:"option_value_ids"`
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
		SellerID:         sellerID,
		CategoryID:       req.CategoryID,
		AttributeGroupID: req.AttributeGroupID,
		Name:             req.Name,
		Description:      req.Description,
		BasePriceCents:   req.BasePriceCents,
		Currency:         req.Currency,
		ProductType:      domain.ProductType(req.ProductType),
		StockQuantity:    req.StockQuantity,
		Tags:             req.Tags,
		ImageURLs:        req.ImageURLs,
		Attributes:       attrs,
	}

	product, err := h.productUC.CreateProduct(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, product)
}

type updateProductRequest struct {
	Name             *string               `json:"name"`
	Description      *string               `json:"description"`
	BasePriceCents   *int64                `json:"base_price_cents"`
	Currency         *string               `json:"currency"`
	Status           *domain.ProductStatus `json:"status"`
	StockQuantity    *int                  `json:"stock_quantity"`
	Tags             []string              `json:"tags"`
	ImageURLs        []string              `json:"image_urls"`
	CategoryID       *string               `json:"category_id"`
	AttributeGroupID *string               `json:"attribute_group_id"`
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
		Name:             req.Name,
		Description:      req.Description,
		BasePriceCents:   req.BasePriceCents,
		Currency:         req.Currency,
		Status:           req.Status,
		StockQuantity:    req.StockQuantity,
		Tags:             req.Tags,
		ImageURLs:        req.ImageURLs,
		CategoryID:       req.CategoryID,
		AttributeGroupID: req.AttributeGroupID,
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
	Stock          *int     `json:"stock"`
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
		Stock:          req.Stock,
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

// --- Admin Product Endpoints ---

// AdminListProducts handles GET /api/v1/admin/products — lists ALL products (all statuses, all sellers).
func (h *Handler) AdminListProducts(c *gin.Context) {
	filter := domain.ProductFilter{
		SellerID:   c.Query("seller_id"),
		CategoryID: c.Query("category_id"),
		Status:     c.Query("status"), // admin can filter by any status; empty = all statuses
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

	// Admin sees ALL statuses by default (no forced "active" filter)

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

// AdminUpdateProduct handles PATCH /api/v1/admin/products/:id — updates any product regardless of seller.
func (h *Handler) AdminUpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var req updateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input := usecase.UpdateProductInput{
		Name:             req.Name,
		Description:      req.Description,
		BasePriceCents:   req.BasePriceCents,
		Currency:         req.Currency,
		Status:           req.Status,
		StockQuantity:    req.StockQuantity,
		Tags:             req.Tags,
		ImageURLs:        req.ImageURLs,
		CategoryID:       req.CategoryID,
		AttributeGroupID: req.AttributeGroupID,
	}

	product, err := h.productUC.AdminUpdateProduct(c.Request.Context(), id, input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

// AdminDeleteProduct handles DELETE /api/v1/admin/products/:id — deletes any product regardless of seller.
func (h *Handler) AdminDeleteProduct(c *gin.Context) {
	id := c.Param("id")
	if err := h.productUC.AdminDeleteProduct(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "product deleted"})
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

type updateCategoryRequest struct {
	Name      *string `json:"name"`
	ParentID  *string `json:"parent_id"`
	SortOrder *int    `json:"sort_order"`
	ImageURL  *string `json:"image_url"`
	IsActive  *bool   `json:"is_active"`
}

// UpdateCategory handles PATCH /api/v1/admin/categories/:id
func (h *Handler) UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var req updateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cat, err := h.categoryUC.UpdateCategory(c.Request.Context(), id, usecase.UpdateCategoryInput{
		Name:      req.Name,
		ParentID:  req.ParentID,
		SortOrder: req.SortOrder,
		ImageURL:  req.ImageURL,
		IsActive:  req.IsActive,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cat)
}

// DeleteCategory handles DELETE /api/v1/admin/categories/:id
func (h *Handler) DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	if err := h.categoryUC.DeleteCategory(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "category deleted"})
}

// --- Admin Attribute Endpoints ---

type optionValueInput struct {
	Value     string `json:"value" binding:"required"`
	ColorHex  string `json:"color_hex"`
	SortOrder int    `json:"sort_order"`
}

type createAttributeRequest struct {
	Name         string               `json:"name" binding:"required"`
	Type         domain.AttributeType `json:"type" binding:"required"`
	Required     bool                 `json:"required"`
	Filterable   bool                 `json:"filterable"`
	OptionValues []optionValueInput   `json:"option_values"`
	Unit         string               `json:"unit"`
	SortOrder    int                  `json:"sort_order"`
}

// CreateAttributeDefinition handles POST /api/v1/admin/attributes
func (h *Handler) CreateAttributeDefinition(c *gin.Context) {
	var req createAttributeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var ovInputs []usecase.AttributeOptionValueInput
	for _, ov := range req.OptionValues {
		ovInputs = append(ovInputs, usecase.AttributeOptionValueInput{
			Value:     ov.Value,
			ColorHex:  ov.ColorHex,
			SortOrder: ov.SortOrder,
		})
	}

	attr, err := h.attributeUC.CreateAttributeDefinition(c.Request.Context(), usecase.CreateAttributeInput{
		Name:         req.Name,
		Type:         req.Type,
		Required:     req.Required,
		Filterable:   req.Filterable,
		OptionValues: ovInputs,
		Unit:         req.Unit,
		SortOrder:    req.SortOrder,
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
	Name         *string            `json:"name"`
	Type         *string            `json:"type"`
	Required     *bool              `json:"required"`
	Filterable   *bool              `json:"filterable"`
	OptionValues []optionValueInput `json:"option_values"`
	Unit         *string            `json:"unit"`
	SortOrder    *int               `json:"sort_order"`
}

// UpdateAttributeDefinition handles PATCH /api/v1/admin/attributes/:id
func (h *Handler) UpdateAttributeDefinition(c *gin.Context) {
	id := c.Param("id")
	var req updateAttributeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var ovInputs []usecase.AttributeOptionValueInput
	for _, ov := range req.OptionValues {
		ovInputs = append(ovInputs, usecase.AttributeOptionValueInput{
			Value:     ov.Value,
			ColorHex:  ov.ColorHex,
			SortOrder: ov.SortOrder,
		})
	}

	attr, err := h.attributeUC.UpdateAttributeDefinition(c.Request.Context(), id, usecase.UpdateAttributeInput{
		Name:         req.Name,
		Type:         req.Type,
		Required:     req.Required,
		Filterable:   req.Filterable,
		OptionValues: ovInputs,
		Unit:         req.Unit,
		SortOrder:    req.SortOrder,
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

// --- Public Option & Variant Read Endpoints ---

// ListProductOptions handles GET /api/v1/products/:id/options
func (h *Handler) ListProductOptions(c *gin.Context) {
	productID := c.Param("id")
	options, err := h.variantUC.ListOptions(c.Request.Context(), productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"options": options})
}

// ListProductVariants handles GET /api/v1/products/:id/variants
func (h *Handler) ListProductVariants(c *gin.Context) {
	productID := c.Param("id")
	variants, err := h.variantUC.ListVariantsByProduct(c.Request.Context(), productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"variants": variants})
}

// GetVariant handles GET /api/v1/products/:id/variants/:variantId
func (h *Handler) GetVariant(c *gin.Context) {
	variantID := c.Param("variantId")
	variant, err := h.variantUC.GetVariant(c.Request.Context(), variantID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, variant)
}

// --- Product Attribute Value Endpoints ---

type setProductAttributesRequest struct {
	Attributes []attributeValueInputRequest `json:"attributes" binding:"required"`
}

// SetProductAttributes handles PUT /api/v1/seller/products/:id/attributes
func (h *Handler) SetProductAttributes(c *gin.Context) {
	sellerID := c.GetHeader("X-User-ID")
	if sellerID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing X-User-ID header"})
		return
	}

	productID := c.Param("id")

	// Verify seller owns this product
	product, err := h.productUC.GetProduct(c.Request.Context(), productID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}
	if product.SellerID != sellerID {
		c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized: product belongs to another seller"})
		return
	}

	var req setProductAttributesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var values []domain.ProductAttributeValue
	for _, a := range req.Attributes {
		values = append(values, domain.ProductAttributeValue{
			ProductID:      productID,
			AttributeID:    a.AttributeID,
			Value:          a.Value,
			Values:         a.Values,
			OptionValueID:  a.OptionValueID,
			OptionValueIDs: a.OptionValueIDs,
		})
	}

	if err := h.attributeUC.SetProductAttributeValues(c.Request.Context(), productID, values); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "attributes updated"})
}

// GetProductAttributes handles GET /api/v1/products/:id/attributes
func (h *Handler) GetProductAttributes(c *gin.Context) {
	productID := c.Param("id")
	attrs, err := h.attributeUC.GetProductAttributeValues(c.Request.Context(), productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"attributes": attrs})
}

// --- Attribute Group Endpoints ---

type createAttributeGroupRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	SortOrder   int    `json:"sort_order"`
}

// CreateAttributeGroup handles POST /api/v1/admin/attribute-groups
func (h *Handler) CreateAttributeGroup(c *gin.Context) {
	var req createAttributeGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	group, err := h.attributeGroupUC.CreateAttributeGroup(c.Request.Context(), usecase.CreateAttributeGroupInput{
		Name:        req.Name,
		Description: req.Description,
		SortOrder:   req.SortOrder,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, group)
}

// ListAttributeGroups handles GET /api/v1/admin/attribute-groups and GET /api/v1/attribute-groups
func (h *Handler) ListAttributeGroups(c *gin.Context) {
	groups, err := h.attributeGroupUC.ListAttributeGroups(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"attribute_groups": groups})
}

// GetAttributeGroup handles GET /api/v1/attribute-groups/:id
func (h *Handler) GetAttributeGroup(c *gin.Context) {
	id := c.Param("id")
	group, err := h.attributeGroupUC.GetAttributeGroup(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, group)
}

type updateAttributeGroupRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	SortOrder   *int    `json:"sort_order"`
}

// UpdateAttributeGroup handles PATCH /api/v1/admin/attribute-groups/:id
func (h *Handler) UpdateAttributeGroup(c *gin.Context) {
	id := c.Param("id")
	var req updateAttributeGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	group, err := h.attributeGroupUC.UpdateAttributeGroup(c.Request.Context(), id, usecase.UpdateAttributeGroupInput{
		Name:        req.Name,
		Description: req.Description,
		SortOrder:   req.SortOrder,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, group)
}

// DeleteAttributeGroup handles DELETE /api/v1/admin/attribute-groups/:id
func (h *Handler) DeleteAttributeGroup(c *gin.Context) {
	id := c.Param("id")
	if err := h.attributeGroupUC.DeleteAttributeGroup(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "attribute group deleted"})
}

type addAttributeToGroupRequest struct {
	AttributeID string `json:"attribute_id" binding:"required"`
	SortOrder   int    `json:"sort_order"`
}

// AddAttributeToGroup handles POST /api/v1/admin/attribute-groups/:id/attributes
func (h *Handler) AddAttributeToGroup(c *gin.Context) {
	groupID := c.Param("id")
	var req addAttributeToGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.attributeGroupUC.AddAttributeToGroup(c.Request.Context(), groupID, req.AttributeID, req.SortOrder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "attribute added to group"})
}

// RemoveAttributeFromGroup handles DELETE /api/v1/admin/attribute-groups/:id/attributes/:attrId
func (h *Handler) RemoveAttributeFromGroup(c *gin.Context) {
	groupID := c.Param("id")
	attrID := c.Param("attrId")

	if err := h.attributeGroupUC.RemoveAttributeFromGroup(c.Request.Context(), groupID, attrID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "attribute removed from group"})
}

// ListGroupAttributes handles GET /api/v1/attribute-groups/:id/attributes
func (h *Handler) ListGroupAttributes(c *gin.Context) {
	groupID := c.Param("id")
	attrs, err := h.attributeGroupUC.ListGroupAttributes(c.Request.Context(), groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"attributes": attrs})
}

// --- Admin Product Management Endpoints (options, variants, attributes — no seller check) ---

// AdminGetProduct handles GET /api/v1/admin/products/:id — returns full product with preloaded relations.
func (h *Handler) AdminGetProduct(c *gin.Context) {
	id := c.Param("id")
	product, err := h.productUC.GetProduct(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, product)
}

// AdminListOptions handles GET /api/v1/admin/products/:id/options
func (h *Handler) AdminListOptions(c *gin.Context) {
	productID := c.Param("id")
	options, err := h.variantUC.ListOptions(c.Request.Context(), productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"options": options})
}

// AdminAddOption handles POST /api/v1/admin/products/:id/options
func (h *Handler) AdminAddOption(c *gin.Context) {
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

	option, err := h.variantUC.AdminAddOption(c.Request.Context(), productID, usecase.AddOptionInput{
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

// AdminRemoveOption handles DELETE /api/v1/admin/products/:id/options/:optionId
func (h *Handler) AdminRemoveOption(c *gin.Context) {
	optionID := c.Param("optionId")

	if err := h.variantUC.AdminRemoveOption(c.Request.Context(), optionID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "option removed"})
}

// AdminListVariants handles GET /api/v1/admin/products/:id/variants
func (h *Handler) AdminListVariants(c *gin.Context) {
	productID := c.Param("id")
	variants, err := h.variantUC.ListVariantsByProduct(c.Request.Context(), productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"variants": variants})
}

// AdminGenerateVariants handles POST /api/v1/admin/products/:id/variants/generate
func (h *Handler) AdminGenerateVariants(c *gin.Context) {
	productID := c.Param("id")
	variants, err := h.variantUC.AdminGenerateVariants(c.Request.Context(), productID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"variants": variants})
}

// AdminUpdateVariant handles PATCH /api/v1/admin/products/:id/variants/:variantId
func (h *Handler) AdminUpdateVariant(c *gin.Context) {
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
		Stock:          req.Stock,
		WeightGrams:    req.WeightGrams,
		IsActive:       req.IsActive,
		ImageURLs:      req.ImageURLs,
		Barcode:        req.Barcode,
		LowStockAlert:  req.LowStockAlert,
	}

	variant, err := h.variantUC.AdminUpdateVariant(c.Request.Context(), productID, variantID, input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, variant)
}

// AdminUpdateVariantStock handles PATCH /api/v1/admin/products/:id/variants/:variantId/stock
func (h *Handler) AdminUpdateVariantStock(c *gin.Context) {
	productID := c.Param("id")
	variantID := c.Param("variantId")

	var req updateStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.variantUC.UpdateStock(c.Request.Context(), productID, variantID, "", req.Delta); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "stock updated"})
}

// AdminSetProductAttributes handles PUT /api/v1/admin/products/:id/attributes
func (h *Handler) AdminSetProductAttributes(c *gin.Context) {
	productID := c.Param("id")

	var req setProductAttributesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var values []domain.ProductAttributeValue
	for _, a := range req.Attributes {
		values = append(values, domain.ProductAttributeValue{
			ProductID:      productID,
			AttributeID:    a.AttributeID,
			Value:          a.Value,
			Values:         a.Values,
			OptionValueID:  a.OptionValueID,
			OptionValueIDs: a.OptionValueIDs,
		})
	}

	if err := h.attributeUC.SetProductAttributeValues(c.Request.Context(), productID, values); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "attributes updated"})
}

// AdminGetProductAttributes handles GET /api/v1/admin/products/:id/attributes
func (h *Handler) AdminGetProductAttributes(c *gin.Context) {
	productID := c.Param("id")
	attrs, err := h.attributeUC.GetProductAttributeValues(c.Request.Context(), productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"attributes": attrs})
}
