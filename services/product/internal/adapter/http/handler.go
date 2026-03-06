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

// ListProducts godoc
// @Summary      List products
// @Tags         Public Products
// @Produce      json
// @Param        seller_id   query  string  false  "Filter by seller ID"
// @Param        category_id query  string  false  "Filter by category ID"
// @Param        status      query  string  false  "Filter by status (default: active)"
// @Param        q           query  string  false  "Search query"
// @Param        sort_by     query  string  false  "Sort field"
// @Param        min_price   query  int     false  "Minimum price in cents"
// @Param        max_price   query  int     false  "Maximum price in cents"
// @Param        page        query  int     false  "Page number"
// @Param        page_size   query  int     false  "Page size"
// @Success      200  {object}  object{products=[]domain.Product,total=int,page=int,pageSize=int}
// @Failure      500  {object}  object{error=string}
// @Router       /products [get]
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

// GetProduct godoc
// @Summary      Get product by ID
// @Tags         Public Products
// @Produce      json
// @Param        id  path  string  true  "Product ID"
// @Success      200  {object}  domain.Product
// @Failure      404  {object}  object{error=string}
// @Router       /products/{id} [get]
func (h *Handler) GetProduct(c *gin.Context) {
	id := c.Param("id")
	product, err := h.productUC.GetProduct(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, product)
}

// GetProductBySlug godoc
// @Summary      Get product by slug
// @Tags         Public Products
// @Produce      json
// @Param        slug  path  string  true  "Product slug"
// @Success      200  {object}  domain.Product
// @Failure      404  {object}  object{error=string}
// @Router       /products/slug/{slug} [get]
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

// ListCategories godoc
// @Summary      List all categories
// @Tags         Categories
// @Produce      json
// @Success      200  {object}  object{categories=[]domain.Category}
// @Failure      500  {object}  object{error=string}
// @Router       /categories [get]
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

// CreateProduct godoc
// @Summary      Create a new product
// @Tags         Seller Products
// @Accept       json
// @Produce      json
// @Param        X-User-ID  header  string               true  "Seller ID"
// @Param        body       body    createProductRequest  true  "Product data"
// @Success      201  {object}  domain.Product
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Router       /seller/products [post]
// @Security     BearerAuth
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

// UpdateProduct godoc
// @Summary      Update a product (seller)
// @Tags         Seller Products
// @Accept       json
// @Produce      json
// @Param        X-User-ID  header  string               true  "Seller ID"
// @Param        id         path    string               true  "Product ID"
// @Param        body       body    updateProductRequest  true  "Fields to update"
// @Success      200  {object}  domain.Product
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Router       /seller/products/{id} [patch]
// @Security     BearerAuth
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

// DeleteProduct godoc
// @Summary      Delete a product (seller)
// @Tags         Seller Products
// @Produce      json
// @Param        X-User-ID  header  string  true  "Seller ID"
// @Param        id         path    string  true  "Product ID"
// @Success      200  {object}  object{message=string}
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Router       /seller/products/{id} [delete]
// @Security     BearerAuth
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

// AddOption godoc
// @Summary      Add an option to a product (seller)
// @Tags         Seller Products
// @Accept       json
// @Produce      json
// @Param        X-User-ID  header  string           true  "Seller ID"
// @Param        id         path    string           true  "Product ID"
// @Param        body       body    addOptionRequest true  "Option data"
// @Success      201  {object}  domain.ProductOption
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Router       /seller/products/{id}/options [post]
// @Security     BearerAuth
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

// RemoveOption godoc
// @Summary      Remove an option from a product (seller)
// @Tags         Seller Products
// @Produce      json
// @Param        X-User-ID  header  string  true  "Seller ID"
// @Param        id         path    string  true  "Product ID"
// @Param        optionId   path    string  true  "Option ID"
// @Success      200  {object}  object{message=string}
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Router       /seller/products/{id}/options/{optionId} [delete]
// @Security     BearerAuth
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

// GenerateVariants godoc
// @Summary      Generate variants from options (seller)
// @Tags         Seller Products
// @Produce      json
// @Param        X-User-ID  header  string  true  "Seller ID"
// @Param        id         path    string  true  "Product ID"
// @Success      201  {object}  object{variants=[]domain.Variant}
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Router       /seller/products/{id}/variants/generate [post]
// @Security     BearerAuth
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

// UpdateVariant godoc
// @Summary      Update a variant (seller)
// @Tags         Seller Products
// @Accept       json
// @Produce      json
// @Param        X-User-ID  header  string                true  "Seller ID"
// @Param        id         path    string                true  "Product ID"
// @Param        variantId  path    string                true  "Variant ID"
// @Param        body       body    updateVariantRequest  true  "Variant fields to update"
// @Success      200  {object}  domain.Variant
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Router       /seller/products/{id}/variants/{variantId} [patch]
// @Security     BearerAuth
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

// UpdateVariantStock godoc
// @Summary      Update variant stock (seller)
// @Tags         Seller Products
// @Accept       json
// @Produce      json
// @Param        X-User-ID  header  string              true  "Seller ID"
// @Param        id         path    string              true  "Product ID"
// @Param        variantId  path    string              true  "Variant ID"
// @Param        body       body    updateStockRequest  true  "Stock delta"
// @Success      200  {object}  object{message=string}
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Router       /seller/products/{id}/variants/{variantId}/stock [patch]
// @Security     BearerAuth
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

// AdminListProducts godoc
// @Summary      List all products (admin)
// @Description  Lists ALL products across all statuses and sellers
// @Tags         Admin Products
// @Produce      json
// @Param        seller_id   query  string  false  "Filter by seller ID"
// @Param        category_id query  string  false  "Filter by category ID"
// @Param        status      query  string  false  "Filter by status"
// @Param        q           query  string  false  "Search query"
// @Param        sort_by     query  string  false  "Sort field"
// @Param        min_price   query  int     false  "Minimum price in cents"
// @Param        max_price   query  int     false  "Maximum price in cents"
// @Param        page        query  int     false  "Page number"
// @Param        page_size   query  int     false  "Page size"
// @Success      200  {object}  object{products=[]domain.Product,total=int,page=int,pageSize=int}
// @Failure      500  {object}  object{error=string}
// @Router       /admin/products [get]
// @Security     BearerAuth
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

// AdminUpdateProduct godoc
// @Summary      Update any product (admin)
// @Description  Updates a product regardless of seller ownership
// @Tags         Admin Products
// @Accept       json
// @Produce      json
// @Param        id    path  string               true  "Product ID"
// @Param        body  body  updateProductRequest  true  "Fields to update"
// @Success      200  {object}  domain.Product
// @Failure      400  {object}  object{error=string}
// @Router       /admin/products/{id} [patch]
// @Security     BearerAuth
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

// AdminDeleteProduct godoc
// @Summary      Delete any product (admin)
// @Description  Deletes a product regardless of seller ownership
// @Tags         Admin Products
// @Produce      json
// @Param        id  path  string  true  "Product ID"
// @Success      200  {object}  object{message=string}
// @Failure      400  {object}  object{error=string}
// @Router       /admin/products/{id} [delete]
// @Security     BearerAuth
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

// CreateCategory godoc
// @Summary      Create a category
// @Tags         Admin Categories
// @Accept       json
// @Produce      json
// @Param        body  body  createCategoryRequest  true  "Category data"
// @Success      201  {object}  domain.Category
// @Failure      400  {object}  object{error=string}
// @Router       /admin/categories [post]
// @Security     BearerAuth
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

// UpdateCategory godoc
// @Summary      Update a category
// @Tags         Admin Categories
// @Accept       json
// @Produce      json
// @Param        id    path  string                true  "Category ID"
// @Param        body  body  updateCategoryRequest  true  "Fields to update"
// @Success      200  {object}  domain.Category
// @Failure      400  {object}  object{error=string}
// @Router       /admin/categories/{id} [patch]
// @Security     BearerAuth
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

// DeleteCategory godoc
// @Summary      Delete a category
// @Tags         Admin Categories
// @Produce      json
// @Param        id  path  string  true  "Category ID"
// @Success      200  {object}  object{message=string}
// @Failure      400  {object}  object{error=string}
// @Router       /admin/categories/{id} [delete]
// @Security     BearerAuth
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

// CreateAttributeDefinition godoc
// @Summary      Create an attribute definition
// @Tags         Admin Attributes
// @Accept       json
// @Produce      json
// @Param        body  body  createAttributeRequest  true  "Attribute definition data"
// @Success      201  {object}  domain.AttributeDefinition
// @Failure      400  {object}  object{error=string}
// @Router       /admin/attributes [post]
// @Security     BearerAuth
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

// ListAttributeDefinitions godoc
// @Summary      List all attribute definitions
// @Tags         Admin Attributes
// @Produce      json
// @Success      200  {object}  object{attributes=[]domain.AttributeDefinition}
// @Failure      500  {object}  object{error=string}
// @Router       /admin/attributes [get]
// @Security     BearerAuth
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

// UpdateAttributeDefinition godoc
// @Summary      Update an attribute definition
// @Tags         Admin Attributes
// @Accept       json
// @Produce      json
// @Param        id    path  string                  true  "Attribute definition ID"
// @Param        body  body  updateAttributeRequest  true  "Fields to update"
// @Success      200  {object}  domain.AttributeDefinition
// @Failure      400  {object}  object{error=string}
// @Router       /admin/attributes/{id} [patch]
// @Security     BearerAuth
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

// DeleteAttributeDefinition godoc
// @Summary      Delete an attribute definition
// @Tags         Admin Attributes
// @Produce      json
// @Param        id  path  string  true  "Attribute definition ID"
// @Success      200  {object}  object{message=string}
// @Failure      400  {object}  object{error=string}
// @Router       /admin/attributes/{id} [delete]
// @Security     BearerAuth
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

// AssignAttributeToCategory godoc
// @Summary      Assign an attribute to a category
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Param        id    path  string                  true  "Category ID"
// @Param        body  body  assignAttributeRequest  true  "Attribute assignment"
// @Success      200  {object}  object{message=string}
// @Failure      400  {object}  object{error=string}
// @Router       /categories/{id}/attributes [post]
// @Security     BearerAuth
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

// RemoveAttributeFromCategory godoc
// @Summary      Remove an attribute from a category
// @Tags         Categories
// @Produce      json
// @Param        id      path  string  true  "Category ID"
// @Param        attrId  path  string  true  "Attribute definition ID"
// @Success      200  {object}  object{message=string}
// @Failure      400  {object}  object{error=string}
// @Router       /categories/{id}/attributes/{attrId} [delete]
// @Security     BearerAuth
func (h *Handler) RemoveAttributeFromCategory(c *gin.Context) {
	categoryID := c.Param("id")
	attrID := c.Param("attrId")

	if err := h.attributeUC.RemoveAttributeFromCategory(c.Request.Context(), categoryID, attrID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "attribute removed from category"})
}

// ListCategoryAttributes godoc
// @Summary      List attributes assigned to a category
// @Tags         Categories
// @Produce      json
// @Param        id  path  string  true  "Category ID"
// @Success      200  {object}  object{attributes=[]domain.AttributeDefinition}
// @Failure      500  {object}  object{error=string}
// @Router       /categories/{id}/attributes [get]
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

// ListProductOptions godoc
// @Summary      List product options
// @Tags         Public Products
// @Produce      json
// @Param        id  path  string  true  "Product ID"
// @Success      200  {object}  object{options=[]domain.ProductOption}
// @Failure      500  {object}  object{error=string}
// @Router       /products/{id}/options [get]
func (h *Handler) ListProductOptions(c *gin.Context) {
	productID := c.Param("id")
	options, err := h.variantUC.ListOptions(c.Request.Context(), productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"options": options})
}

// ListProductVariants godoc
// @Summary      List product variants
// @Tags         Public Products
// @Produce      json
// @Param        id  path  string  true  "Product ID"
// @Success      200  {object}  object{variants=[]domain.Variant}
// @Failure      500  {object}  object{error=string}
// @Router       /products/{id}/variants [get]
func (h *Handler) ListProductVariants(c *gin.Context) {
	productID := c.Param("id")
	variants, err := h.variantUC.ListVariantsByProduct(c.Request.Context(), productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"variants": variants})
}

// GetVariant godoc
// @Summary      Get a variant by ID
// @Tags         Public Products
// @Produce      json
// @Param        id         path  string  true  "Product ID"
// @Param        variantId  path  string  true  "Variant ID"
// @Success      200  {object}  domain.Variant
// @Failure      404  {object}  object{error=string}
// @Router       /products/{id}/variants/{variantId} [get]
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

// SetProductAttributes godoc
// @Summary      Set product attribute values (seller)
// @Tags         Seller Products
// @Accept       json
// @Produce      json
// @Param        X-User-ID  header  string                       true  "Seller ID"
// @Param        id         path    string                       true  "Product ID"
// @Param        body       body    setProductAttributesRequest   true  "Attribute values"
// @Success      200  {object}  object{message=string}
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Failure      404  {object}  object{error=string}
// @Router       /seller/products/{id}/attributes [put]
// @Security     BearerAuth
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

// GetProductAttributes godoc
// @Summary      Get product attribute values
// @Tags         Public Products
// @Produce      json
// @Param        id  path  string  true  "Product ID"
// @Success      200  {object}  object{attributes=[]domain.ProductAttributeValue}
// @Failure      500  {object}  object{error=string}
// @Router       /products/{id}/attributes [get]
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

// CreateAttributeGroup godoc
// @Summary      Create an attribute group
// @Tags         Admin Attribute Groups
// @Accept       json
// @Produce      json
// @Param        body  body  createAttributeGroupRequest  true  "Attribute group data"
// @Success      201  {object}  domain.AttributeGroup
// @Failure      400  {object}  object{error=string}
// @Router       /admin/attribute-groups [post]
// @Security     BearerAuth
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

// ListAttributeGroups godoc
// @Summary      List all attribute groups
// @Tags         Attribute Groups
// @Produce      json
// @Success      200  {object}  object{attribute_groups=[]domain.AttributeGroup}
// @Failure      500  {object}  object{error=string}
// @Router       /attribute-groups [get]
func (h *Handler) ListAttributeGroups(c *gin.Context) {
	groups, err := h.attributeGroupUC.ListAttributeGroups(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"attribute_groups": groups})
}

// GetAttributeGroup godoc
// @Summary      Get an attribute group by ID
// @Tags         Attribute Groups
// @Produce      json
// @Param        id  path  string  true  "Attribute group ID"
// @Success      200  {object}  domain.AttributeGroup
// @Failure      404  {object}  object{error=string}
// @Router       /attribute-groups/{id} [get]
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

// UpdateAttributeGroup godoc
// @Summary      Update an attribute group
// @Tags         Admin Attribute Groups
// @Accept       json
// @Produce      json
// @Param        id    path  string                       true  "Attribute group ID"
// @Param        body  body  updateAttributeGroupRequest  true  "Fields to update"
// @Success      200  {object}  domain.AttributeGroup
// @Failure      400  {object}  object{error=string}
// @Router       /admin/attribute-groups/{id} [patch]
// @Security     BearerAuth
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

// DeleteAttributeGroup godoc
// @Summary      Delete an attribute group
// @Tags         Admin Attribute Groups
// @Produce      json
// @Param        id  path  string  true  "Attribute group ID"
// @Success      200  {object}  object{message=string}
// @Failure      400  {object}  object{error=string}
// @Router       /admin/attribute-groups/{id} [delete]
// @Security     BearerAuth
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

// AddAttributeToGroup godoc
// @Summary      Add an attribute to a group
// @Tags         Admin Attribute Groups
// @Accept       json
// @Produce      json
// @Param        id    path  string                      true  "Attribute group ID"
// @Param        body  body  addAttributeToGroupRequest  true  "Attribute to add"
// @Success      200  {object}  object{message=string}
// @Failure      400  {object}  object{error=string}
// @Router       /admin/attribute-groups/{id}/attributes [post]
// @Security     BearerAuth
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

// RemoveAttributeFromGroup godoc
// @Summary      Remove an attribute from a group
// @Tags         Admin Attribute Groups
// @Produce      json
// @Param        id      path  string  true  "Attribute group ID"
// @Param        attrId  path  string  true  "Attribute definition ID"
// @Success      200  {object}  object{message=string}
// @Failure      400  {object}  object{error=string}
// @Router       /admin/attribute-groups/{id}/attributes/{attrId} [delete]
// @Security     BearerAuth
func (h *Handler) RemoveAttributeFromGroup(c *gin.Context) {
	groupID := c.Param("id")
	attrID := c.Param("attrId")

	if err := h.attributeGroupUC.RemoveAttributeFromGroup(c.Request.Context(), groupID, attrID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "attribute removed from group"})
}

// ListGroupAttributes godoc
// @Summary      List attributes in a group
// @Tags         Attribute Groups
// @Produce      json
// @Param        id  path  string  true  "Attribute group ID"
// @Success      200  {object}  object{attributes=[]domain.AttributeDefinition}
// @Failure      500  {object}  object{error=string}
// @Router       /attribute-groups/{id}/attributes [get]
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

// AdminGetProduct godoc
// @Summary      Get any product by ID (admin)
// @Description  Returns full product with preloaded relations
// @Tags         Admin Products
// @Produce      json
// @Param        id  path  string  true  "Product ID"
// @Success      200  {object}  domain.Product
// @Failure      404  {object}  object{error=string}
// @Router       /admin/products/{id} [get]
// @Security     BearerAuth
func (h *Handler) AdminGetProduct(c *gin.Context) {
	id := c.Param("id")
	product, err := h.productUC.GetProduct(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, product)
}

// AdminListOptions godoc
// @Summary      List product options (admin)
// @Tags         Admin Products
// @Produce      json
// @Param        id  path  string  true  "Product ID"
// @Success      200  {object}  object{options=[]domain.ProductOption}
// @Failure      500  {object}  object{error=string}
// @Router       /admin/products/{id}/options [get]
// @Security     BearerAuth
func (h *Handler) AdminListOptions(c *gin.Context) {
	productID := c.Param("id")
	options, err := h.variantUC.ListOptions(c.Request.Context(), productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"options": options})
}

// AdminAddOption godoc
// @Summary      Add an option to a product (admin)
// @Tags         Admin Products
// @Accept       json
// @Produce      json
// @Param        id    path  string           true  "Product ID"
// @Param        body  body  addOptionRequest true  "Option data"
// @Success      201  {object}  domain.ProductOption
// @Failure      400  {object}  object{error=string}
// @Router       /admin/products/{id}/options [post]
// @Security     BearerAuth
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

// AdminRemoveOption godoc
// @Summary      Remove an option from a product (admin)
// @Tags         Admin Products
// @Produce      json
// @Param        id        path  string  true  "Product ID"
// @Param        optionId  path  string  true  "Option ID"
// @Success      200  {object}  object{message=string}
// @Failure      400  {object}  object{error=string}
// @Router       /admin/products/{id}/options/{optionId} [delete]
// @Security     BearerAuth
func (h *Handler) AdminRemoveOption(c *gin.Context) {
	optionID := c.Param("optionId")

	if err := h.variantUC.AdminRemoveOption(c.Request.Context(), optionID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "option removed"})
}

// AdminListVariants godoc
// @Summary      List product variants (admin)
// @Tags         Admin Products
// @Produce      json
// @Param        id  path  string  true  "Product ID"
// @Success      200  {object}  object{variants=[]domain.Variant}
// @Failure      500  {object}  object{error=string}
// @Router       /admin/products/{id}/variants [get]
// @Security     BearerAuth
func (h *Handler) AdminListVariants(c *gin.Context) {
	productID := c.Param("id")
	variants, err := h.variantUC.ListVariantsByProduct(c.Request.Context(), productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"variants": variants})
}

// AdminGenerateVariants godoc
// @Summary      Generate variants from options (admin)
// @Tags         Admin Products
// @Produce      json
// @Param        id  path  string  true  "Product ID"
// @Success      201  {object}  object{variants=[]domain.Variant}
// @Failure      400  {object}  object{error=string}
// @Router       /admin/products/{id}/variants/generate [post]
// @Security     BearerAuth
func (h *Handler) AdminGenerateVariants(c *gin.Context) {
	productID := c.Param("id")
	variants, err := h.variantUC.AdminGenerateVariants(c.Request.Context(), productID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"variants": variants})
}

// AdminUpdateVariant godoc
// @Summary      Update a variant (admin)
// @Tags         Admin Products
// @Accept       json
// @Produce      json
// @Param        id         path  string                true  "Product ID"
// @Param        variantId  path  string                true  "Variant ID"
// @Param        body       body  updateVariantRequest  true  "Variant fields to update"
// @Success      200  {object}  domain.Variant
// @Failure      400  {object}  object{error=string}
// @Router       /admin/products/{id}/variants/{variantId} [patch]
// @Security     BearerAuth
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

// AdminUpdateVariantStock godoc
// @Summary      Update variant stock (admin)
// @Tags         Admin Products
// @Accept       json
// @Produce      json
// @Param        id         path  string              true  "Product ID"
// @Param        variantId  path  string              true  "Variant ID"
// @Param        body       body  updateStockRequest  true  "Stock delta"
// @Success      200  {object}  object{message=string}
// @Failure      400  {object}  object{error=string}
// @Router       /admin/products/{id}/variants/{variantId}/stock [patch]
// @Security     BearerAuth
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

// AdminSetProductAttributes godoc
// @Summary      Set product attribute values (admin)
// @Tags         Admin Products
// @Accept       json
// @Produce      json
// @Param        id    path  string                       true  "Product ID"
// @Param        body  body  setProductAttributesRequest   true  "Attribute values"
// @Success      200  {object}  object{message=string}
// @Failure      400  {object}  object{error=string}
// @Router       /admin/products/{id}/attributes [put]
// @Security     BearerAuth
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

// AdminGetProductAttributes godoc
// @Summary      Get product attribute values (admin)
// @Tags         Admin Products
// @Produce      json
// @Param        id  path  string  true  "Product ID"
// @Success      200  {object}  object{attributes=[]domain.ProductAttributeValue}
// @Failure      500  {object}  object{error=string}
// @Router       /admin/products/{id}/attributes [get]
// @Security     BearerAuth
func (h *Handler) AdminGetProductAttributes(c *gin.Context) {
	productID := c.Param("id")
	attrs, err := h.attributeUC.GetProductAttributeValues(c.Request.Context(), productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"attributes": attrs})
}
