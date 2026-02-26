package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/southern-martin/ecommerce/services/media/internal/usecase"
)

// Handler holds all HTTP handlers for the media service.
type Handler struct {
	mediaUC *usecase.MediaUseCase
}

// NewHandler creates a new Handler.
func NewHandler(mediaUC *usecase.MediaUseCase) *Handler {
	return &Handler{
		mediaUC: mediaUC,
	}
}

// Health returns a health check response.
func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "media"})
}

// --- Media Handlers ---

type createMediaRequest struct {
	OriginalName string `json:"original_name" binding:"required"`
	ContentType  string `json:"content_type" binding:"required"`
	SizeBytes    int64  `json:"size_bytes"`
	OwnerType    string `json:"owner_type"`
}

// CreateMedia handles POST /api/v1/media.
func (h *Handler) CreateMedia(c *gin.Context) {
	ownerID := c.GetHeader("X-User-ID")
	if ownerID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user ID"})
		return
	}

	var req createMediaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ownerType := req.OwnerType
	if ownerType == "" {
		ownerType = "user"
	}

	result, err := h.mediaUC.CreateMedia(c.Request.Context(), usecase.CreateMediaRequest{
		OwnerID:      ownerID,
		OwnerType:    ownerType,
		OriginalName: req.OriginalName,
		ContentType:  req.ContentType,
		SizeBytes:    req.SizeBytes,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"media": result.Media, "upload_url": result.UploadURL})
}

// GetMedia handles GET /api/v1/media/:id.
func (h *Handler) GetMedia(c *gin.Context) {
	id := c.Param("id")
	media, err := h.mediaUC.GetMedia(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "media not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"media": media})
}

// ListMedia handles GET /api/v1/media.
func (h *Handler) ListMedia(c *gin.Context) {
	ownerID := c.Query("owner_id")
	ownerType := c.Query("owner_type")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	media, total, err := h.mediaUC.ListMedia(c.Request.Context(), ownerID, ownerType, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"media": media, "total": total, "page": page, "page_size": pageSize})
}

// DeleteMedia handles DELETE /api/v1/media/:id.
func (h *Handler) DeleteMedia(c *gin.Context) {
	id := c.Param("id")
	if err := h.mediaUC.DeleteMedia(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "media not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "media deleted"})
}

// --- Upload/Download URL Handlers ---

type getUploadURLRequest struct {
	Key         string `json:"key" binding:"required"`
	ContentType string `json:"content_type" binding:"required"`
}

// GetUploadURL handles POST /api/v1/media/upload-url.
func (h *Handler) GetUploadURL(c *gin.Context) {
	var req getUploadURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	url, err := h.mediaUC.GenerateUploadURL(c.Request.Context(), req.Key, req.ContentType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"upload_url": url})
}

// GetDownloadURL handles GET /api/v1/media/:id/download-url.
func (h *Handler) GetDownloadURL(c *gin.Context) {
	id := c.Param("id")
	url, err := h.mediaUC.GenerateDownloadURL(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "media not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"download_url": url})
}
