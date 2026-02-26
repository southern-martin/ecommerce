package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/southern-martin/ecommerce/services/search/internal/domain"
	"github.com/southern-martin/ecommerce/services/search/internal/usecase"
)

// Handler holds the HTTP handlers for the search service.
type Handler struct {
	searchUC *usecase.SearchUseCase
	indexUC  *usecase.IndexUseCase
}

// NewHandler creates a new Handler.
func NewHandler(searchUC *usecase.SearchUseCase, indexUC *usecase.IndexUseCase) *Handler {
	return &Handler{
		searchUC: searchUC,
		indexUC:  indexUC,
	}
}

// Health handles health check requests.
func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "search"})
}

// Search handles GET /api/v1/search
func (h *Handler) Search(c *gin.Context) {
	filter := domain.SearchFilter{
		Query:      c.Query("q"),
		CategoryID: c.Query("category_id"),
		SellerID:   c.Query("seller_id"),
		SortBy:     c.Query("sort_by"),
		SortOrder:  c.Query("sort_order"),
	}

	if v := c.Query("min_price"); v != "" {
		if parsed, err := strconv.ParseInt(v, 10, 64); err == nil {
			filter.MinPrice = parsed
		}
	}

	if v := c.Query("max_price"); v != "" {
		if parsed, err := strconv.ParseInt(v, 10, 64); err == nil {
			filter.MaxPrice = parsed
		}
	}

	if v := c.Query("in_stock"); v != "" {
		if parsed, err := strconv.ParseBool(v); err == nil {
			filter.InStock = &parsed
		}
	}

	if v := c.Query("page"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			filter.Page = parsed
		}
	}

	if v := c.Query("page_size"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			filter.PageSize = parsed
		}
	}

	results, total, err := h.searchUC.Search(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "search failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"results":   results,
		"total":     total,
		"page":      filter.Page,
		"page_size": filter.PageSize,
	})
}

// Suggest handles GET /api/v1/search/suggest
func (h *Handler) Suggest(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'q' is required"})
		return
	}

	limit := 10
	if v := c.Query("limit"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			limit = parsed
		}
	}

	suggestions, err := h.searchUC.Suggest(c.Request.Context(), query, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "suggest failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"suggestions": suggestions,
	})
}

// IndexProductRequest is the request body for indexing a product.
type IndexProductRequest struct {
	ProductID   string            `json:"product_id" binding:"required"`
	Name        string            `json:"name" binding:"required"`
	Slug        string            `json:"slug"`
	Description string            `json:"description"`
	PriceCents  int64             `json:"price_cents"`
	Currency    string            `json:"currency"`
	CategoryID  string            `json:"category_id"`
	SellerID    string            `json:"seller_id"`
	ImageURL    string            `json:"image_url"`
	Rating      float64           `json:"rating"`
	ReviewCount int               `json:"review_count"`
	InStock     bool              `json:"in_stock"`
	Tags        []string          `json:"tags"`
	Attributes  map[string]string `json:"attributes"`
}

// IndexProduct handles POST /api/v1/admin/search/index
func (h *Handler) IndexProduct(c *gin.Context) {
	var req IndexProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	idx := &domain.SearchIndex{
		ProductID:   req.ProductID,
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		PriceCents:  req.PriceCents,
		Currency:    req.Currency,
		CategoryID:  req.CategoryID,
		SellerID:    req.SellerID,
		ImageURL:    req.ImageURL,
		Rating:      req.Rating,
		ReviewCount: req.ReviewCount,
		InStock:     req.InStock,
		Tags:        req.Tags,
		Attributes:  req.Attributes,
	}

	if err := h.indexUC.IndexProduct(c.Request.Context(), idx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to index product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "product indexed successfully",
		"product_id": req.ProductID,
	})
}

// DeleteProduct handles DELETE /api/v1/admin/search/index/:product_id
func (h *Handler) DeleteProduct(c *gin.Context) {
	productID := c.Param("product_id")
	if productID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product_id is required"})
		return
	}

	if err := h.indexUC.RemoveProduct(c.Request.Context(), productID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to remove product from index"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "product removed from index",
		"product_id": productID,
	})
}
