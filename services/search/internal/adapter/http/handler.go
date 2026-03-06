package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/southern-martin/ecommerce/services/search/internal/domain"
	"github.com/southern-martin/ecommerce/services/search/internal/usecase"
)

// Handler holds the HTTP handlers for the search service.
type Handler struct {
	searchUC *usecase.SearchUseCase
	indexUC  *usecase.IndexUseCase
	db       *gorm.DB
}

// NewHandler creates a new Handler.
func NewHandler(searchUC *usecase.SearchUseCase, indexUC *usecase.IndexUseCase, db *gorm.DB) *Handler {
	return &Handler{
		searchUC: searchUC,
		indexUC:  indexUC,
		db:       db,
	}
}

// Health handles health check requests.
func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "search"})
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

// Search godoc
// @Summary      Search products
// @Description  Full-text search across products with filtering and pagination.
// @Tags         Search
// @Produce      json
// @Param        q            query  string   false  "Search query"
// @Param        category_id  query  string   false  "Filter by category"
// @Param        seller_id    query  string   false  "Filter by seller"
// @Param        min_price    query  integer  false  "Minimum price in cents"
// @Param        max_price    query  integer  false  "Maximum price in cents"
// @Param        in_stock     query  boolean  false  "Filter by stock availability"
// @Param        sort_by      query  string   false  "Sort field"
// @Param        sort_order   query  string   false  "Sort direction (asc/desc)"
// @Param        page         query  integer  false  "Page number"
// @Param        page_size    query  integer  false  "Items per page"
// @Success      200  {object}  object{results=[]domain.SearchIndex,total=int,page=int,page_size=int}
// @Failure      500  {object}  object{error=string}
// @Router       /search [get]
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

// Suggest godoc
// @Summary      Autocomplete suggestions
// @Description  Returns search suggestions based on a partial query string.
// @Tags         Search
// @Produce      json
// @Param        q      query  string   true   "Search query prefix"
// @Param        limit  query  integer  false  "Max suggestions (default 10)"
// @Success      200  {object}  object{suggestions=[]string}
// @Failure      400  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /search/suggest [get]
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

// IndexProduct godoc
// @Summary      Index a product
// @Description  Add or update a product in the search index.
// @Tags         Admin Search Index
// @Accept       json
// @Produce      json
// @Param        X-User-ID  header  string               true  "Admin user ID"
// @Param        body       body    IndexProductRequest   true  "Product data to index"
// @Success      200  {object}  object{message=string,product_id=string}
// @Failure      400  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /admin/search/index [post]
// @Security     BearerAuth
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

// DeleteProduct godoc
// @Summary      Remove product from index
// @Description  Delete a product from the search index.
// @Tags         Admin Search Index
// @Produce      json
// @Param        X-User-ID   header  string  true  "Admin user ID"
// @Param        product_id   path    string  true  "Product ID"
// @Success      200  {object}  object{message=string,product_id=string}
// @Failure      400  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /admin/search/index/{product_id} [delete]
// @Security     BearerAuth
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
