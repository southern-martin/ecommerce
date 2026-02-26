package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/southern-martin/ecommerce/services/cms/internal/domain"
	"github.com/southern-martin/ecommerce/services/cms/internal/usecase"
)

// Handler holds all HTTP handlers for the CMS service.
type Handler struct {
	bannerUC   *usecase.BannerUseCase
	pageUC     *usecase.PageUseCase
	scheduleUC *usecase.ScheduleUseCase
}

// NewHandler creates a new Handler.
func NewHandler(
	bannerUC *usecase.BannerUseCase,
	pageUC *usecase.PageUseCase,
	scheduleUC *usecase.ScheduleUseCase,
) *Handler {
	return &Handler{
		bannerUC:   bannerUC,
		pageUC:     pageUC,
		scheduleUC: scheduleUC,
	}
}

// Health returns a health check response.
func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "cms"})
}

// --- Public Banner Handlers ---

// ListActiveBanners returns active banners, optionally filtered by position.
func (h *Handler) ListActiveBanners(c *gin.Context) {
	position := c.Query("position")

	banners, err := h.bannerUC.ListActiveBanners(c.Request.Context(), position)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"banners": banners})
}

// --- Public Page Handlers ---

// GetPageBySlug returns a page by its slug.
func (h *Handler) GetPageBySlug(c *gin.Context) {
	slug := c.Param("slug")

	page, err := h.pageUC.GetPageBySlug(c.Request.Context(), slug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "page not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"page": page})
}

// --- Admin Banner Handlers ---

type createBannerRequest struct {
	Title          string     `json:"title" binding:"required"`
	ImageURL       string     `json:"image_url" binding:"required"`
	LinkURL        string     `json:"link_url"`
	Position       string     `json:"position"`
	SortOrder      int        `json:"sort_order"`
	TargetAudience string     `json:"target_audience"`
	StartsAt       time.Time  `json:"starts_at" binding:"required"`
	EndsAt         *time.Time `json:"ends_at"`
	IsActive       *bool      `json:"is_active"`
}

func (h *Handler) CreateBanner(c *gin.Context) {
	var req createBannerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	banner := &domain.Banner{
		Title:          req.Title,
		ImageURL:       req.ImageURL,
		LinkURL:        req.LinkURL,
		Position:       req.Position,
		SortOrder:      req.SortOrder,
		TargetAudience: req.TargetAudience,
		StartsAt:       req.StartsAt,
		EndsAt:         req.EndsAt,
		IsActive:       isActive,
	}

	if err := h.bannerUC.CreateBanner(c.Request.Context(), banner); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"banner": banner})
}

type updateBannerRequest struct {
	Title          string     `json:"title"`
	ImageURL       string     `json:"image_url"`
	LinkURL        string     `json:"link_url"`
	Position       string     `json:"position"`
	SortOrder      *int       `json:"sort_order"`
	TargetAudience string     `json:"target_audience"`
	StartsAt       *time.Time `json:"starts_at"`
	EndsAt         *time.Time `json:"ends_at"`
	IsActive       *bool      `json:"is_active"`
}

func (h *Handler) UpdateBanner(c *gin.Context) {
	id := c.Param("id")
	var req updateBannerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	banner := &domain.Banner{
		ID:             id,
		Title:          req.Title,
		ImageURL:       req.ImageURL,
		LinkURL:        req.LinkURL,
		Position:       req.Position,
		TargetAudience: req.TargetAudience,
	}

	if req.SortOrder != nil {
		banner.SortOrder = *req.SortOrder
	}
	if req.StartsAt != nil {
		banner.StartsAt = *req.StartsAt
	}
	if req.EndsAt != nil {
		banner.EndsAt = req.EndsAt
	}
	if req.IsActive != nil {
		banner.IsActive = *req.IsActive
	}

	if err := h.bannerUC.UpdateBanner(c.Request.Context(), banner); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"banner": banner})
}

func (h *Handler) DeleteBanner(c *gin.Context) {
	id := c.Param("id")

	if err := h.bannerUC.DeleteBanner(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "banner deleted"})
}

func (h *Handler) ListAllBanners(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	banners, total, err := h.bannerUC.ListAllBanners(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"banners": banners, "total": total, "page": page, "page_size": pageSize})
}

// --- Admin Page Handlers ---

type createPageRequest struct {
	Title           string `json:"title" binding:"required"`
	ContentHTML     string `json:"content_html"`
	MetaTitle       string `json:"meta_title"`
	MetaDescription string `json:"meta_description"`
}

func (h *Handler) CreatePage(c *gin.Context) {
	var req createPageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	page := &domain.Page{
		Title:           req.Title,
		ContentHTML:     req.ContentHTML,
		MetaTitle:       req.MetaTitle,
		MetaDescription: req.MetaDescription,
	}

	if err := h.pageUC.CreatePage(c.Request.Context(), page); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"page": page})
}

type updatePageRequest struct {
	Title           string `json:"title"`
	ContentHTML     string `json:"content_html"`
	MetaTitle       string `json:"meta_title"`
	MetaDescription string `json:"meta_description"`
}

func (h *Handler) UpdatePage(c *gin.Context) {
	id := c.Param("id")
	var req updatePageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	page := &domain.Page{
		ID:              id,
		Title:           req.Title,
		ContentHTML:     req.ContentHTML,
		MetaTitle:       req.MetaTitle,
		MetaDescription: req.MetaDescription,
	}

	if err := h.pageUC.UpdatePage(c.Request.Context(), page); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"page": page})
}

func (h *Handler) DeletePage(c *gin.Context) {
	id := c.Param("id")

	if err := h.pageUC.DeletePage(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "page deleted"})
}

func (h *Handler) ListAllPages(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	pages, total, err := h.pageUC.ListAllPages(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pages": pages, "total": total, "page": page, "page_size": pageSize})
}

func (h *Handler) PublishPage(c *gin.Context) {
	id := c.Param("id")

	page, err := h.pageUC.PublishPage(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"page": page})
}

// --- Admin Schedule Handlers ---

type scheduleContentRequest struct {
	ContentType string    `json:"content_type" binding:"required"`
	ContentID   string    `json:"content_id" binding:"required"`
	Action      string    `json:"action" binding:"required"`
	ScheduledAt time.Time `json:"scheduled_at" binding:"required"`
}

func (h *Handler) ScheduleContent(c *gin.Context) {
	var req scheduleContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	schedule := &domain.ContentSchedule{
		ContentType: req.ContentType,
		ContentID:   req.ContentID,
		Action:      req.Action,
		ScheduledAt: req.ScheduledAt,
	}

	if err := h.scheduleUC.ScheduleContent(c.Request.Context(), schedule); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"schedule": schedule})
}
