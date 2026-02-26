package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/southern-martin/ecommerce/services/review/internal/domain"
	"github.com/southern-martin/ecommerce/services/review/internal/usecase"
)

// Handler holds all HTTP handlers for the review service.
type Handler struct {
	reviewUC *usecase.ReviewUseCase
}

// NewHandler creates a new Handler.
func NewHandler(reviewUC *usecase.ReviewUseCase) *Handler {
	return &Handler{reviewUC: reviewUC}
}

// Health returns a health check response.
func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "review"})
}

// --- Review Handlers ---

type createReviewRequest struct {
	ProductID          string   `json:"product_id" binding:"required"`
	UserName           string   `json:"user_name"`
	Rating             int      `json:"rating" binding:"required"`
	Title              string   `json:"title"`
	Content            string   `json:"content"`
	Pros               []string `json:"pros"`
	Cons               []string `json:"cons"`
	Images             []string `json:"images"`
	IsVerifiedPurchase bool     `json:"is_verified_purchase"`
}

func (h *Handler) CreateReview(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user ID"})
		return
	}

	var req createReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	review, err := h.reviewUC.CreateReview(c.Request.Context(), usecase.CreateReviewRequest{
		ProductID:          req.ProductID,
		UserID:             userID,
		UserName:           req.UserName,
		Rating:             req.Rating,
		Title:              req.Title,
		Content:            req.Content,
		Pros:               req.Pros,
		Cons:               req.Cons,
		Images:             req.Images,
		IsVerifiedPurchase: req.IsVerifiedPurchase,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"review": review})
}

func (h *Handler) ListProductReviews(c *gin.Context) {
	productID := c.Query("product_id")
	if productID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product_id query parameter is required"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	minRating, _ := strconv.Atoi(c.DefaultQuery("min_rating", "0"))

	filter := domain.ReviewFilter{
		ProductID: productID,
		MinRating: minRating,
		Status:    domain.ReviewStatus(c.DefaultQuery("status", "")),
		SortBy:    c.DefaultQuery("sort_by", "created_at"),
		SortOrder: c.DefaultQuery("sort_order", "desc"),
		Page:      page,
		PageSize:  pageSize,
	}

	reviews, total, err := h.reviewUC.ListProductReviews(c.Request.Context(), productID, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"reviews":   reviews,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (h *Handler) GetReview(c *gin.Context) {
	id := c.Param("id")
	review, err := h.reviewUC.GetReview(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "review not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"review": review})
}

type updateReviewRequest struct {
	Rating  *int     `json:"rating"`
	Title   *string  `json:"title"`
	Content *string  `json:"content"`
	Pros    []string `json:"pros"`
	Cons    []string `json:"cons"`
	Images  []string `json:"images"`
}

func (h *Handler) UpdateReview(c *gin.Context) {
	id := c.Param("id")

	var req updateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	review, err := h.reviewUC.UpdateReview(c.Request.Context(), id, usecase.UpdateReviewRequest{
		Rating:  req.Rating,
		Title:   req.Title,
		Content: req.Content,
		Pros:    req.Pros,
		Cons:    req.Cons,
		Images:  req.Images,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"review": review})
}

func (h *Handler) DeleteReview(c *gin.Context) {
	id := c.Param("id")
	if err := h.reviewUC.DeleteReview(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "review deleted"})
}

func (h *Handler) GetProductSummary(c *gin.Context) {
	productID := c.Param("product_id")
	summary, err := h.reviewUC.GetProductSummary(c.Request.Context(), productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"summary": summary})
}

func (h *Handler) ApproveReview(c *gin.Context) {
	id := c.Param("id")
	review, err := h.reviewUC.ApproveReview(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"review": review})
}
