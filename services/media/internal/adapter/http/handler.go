package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/southern-martin/ecommerce/services/media/internal/usecase"
	"gorm.io/gorm"
)

// Handler holds all HTTP handlers for the media service.
type Handler struct {
	mediaUC *usecase.MediaUseCase
	db      *gorm.DB
}

// NewHandler creates a new Handler.
func NewHandler(mediaUC *usecase.MediaUseCase, db *gorm.DB) *Handler {
	return &Handler{
		mediaUC: mediaUC,
		db:      db,
	}
}

// Health returns a health check response.
func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "media"})
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

// --- Media Handlers ---

type createMediaRequest struct {
	OriginalName string `json:"original_name" binding:"required"`
	ContentType  string `json:"content_type" binding:"required"`
	SizeBytes    int64  `json:"size_bytes"`
	OwnerType    string `json:"owner_type"`
}

// CreateMedia godoc
// @Summary      Create media metadata and generate upload URL
// @Tags         Media
// @Accept       json
// @Produce      json
// @Param        X-User-ID  header  string             true  "User ID"
// @Param        body       body    createMediaRequest  true  "Media metadata"
// @Success      201  {object}  object{media=object{id=string,owner_id=string,file_name=string,original_name=string,content_type=string,size_bytes=int,url=string,status=string,created_at=string},upload_url=object{url=string,method=string,expires_at=string}}
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /media [post]
// @Security     BearerAuth
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

// GetMedia godoc
// @Summary      Get media by ID
// @Tags         Media
// @Accept       json
// @Produce      json
// @Param        id  path  string  true  "Media ID"
// @Success      200  {object}  object{media=object{id=string,owner_id=string,file_name=string,original_name=string,content_type=string,size_bytes=int,url=string,status=string,created_at=string}}
// @Failure      404  {object}  object{error=string}
// @Router       /media/{id} [get]
// @Security     BearerAuth
func (h *Handler) GetMedia(c *gin.Context) {
	id := c.Param("id")
	media, err := h.mediaUC.GetMedia(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "media not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"media": media})
}

// ListMedia godoc
// @Summary      List media with optional filters and pagination
// @Tags         Media
// @Accept       json
// @Produce      json
// @Param        owner_id    query  string  false  "Owner ID"
// @Param        owner_type  query  string  false  "Owner type"
// @Param        page        query  int     false  "Page number"  default(1)
// @Param        page_size   query  int     false  "Page size"    default(20)
// @Success      200  {object}  object{media=[]object{id=string,owner_id=string,file_name=string,original_name=string,content_type=string,size_bytes=int,url=string,status=string,created_at=string},total=int,page=int,page_size=int}
// @Failure      500  {object}  object{error=string}
// @Router       /media [get]
// @Security     BearerAuth
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

// DeleteMedia godoc
// @Summary      Delete media by ID
// @Tags         Media
// @Accept       json
// @Produce      json
// @Param        id  path  string  true  "Media ID"
// @Success      200  {object}  object{message=string}
// @Failure      404  {object}  object{error=string}
// @Router       /media/{id} [delete]
// @Security     BearerAuth
func (h *Handler) DeleteMedia(c *gin.Context) {
	id := c.Param("id")
	if err := h.mediaUC.DeleteMedia(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "media not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "media deleted"})
}

// UploadMedia godoc
// @Summary      Upload a file directly
// @Tags         Media
// @Accept       multipart/form-data
// @Produce      json
// @Param        X-User-ID   header    string  true   "User ID"
// @Param        file        formData  file    true   "File to upload"
// @Param        owner_type  formData  string  false  "Owner type (default: product)"
// @Success      201  {object}  object{media=object{id=string,owner_id=string,file_name=string,original_name=string,content_type=string,size_bytes=int,url=string,status=string,created_at=string}}
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /media/upload [post]
// @Security     BearerAuth
func (h *Handler) UploadMedia(c *gin.Context) {
	ownerID := c.GetHeader("X-User-ID")
	if ownerID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user ID"})
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}
	defer file.Close()

	if header.Size > 10*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file size exceeds 10MB limit"})
		return
	}

	ownerType := c.PostForm("owner_type")
	if ownerType == "" {
		ownerType = "product"
	}

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	media, err := h.mediaUC.UploadMedia(c.Request.Context(), usecase.UploadMediaRequest{
		OwnerID:      ownerID,
		OwnerType:    ownerType,
		OriginalName: header.Filename,
		ContentType:  contentType,
		SizeBytes:    header.Size,
		Reader:       file,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"media": media})
}

// --- Upload/Download URL Handlers ---

type getUploadURLRequest struct {
	Key         string `json:"key" binding:"required"`
	ContentType string `json:"content_type" binding:"required"`
}

// GetUploadURL godoc
// @Summary      Generate a presigned upload URL
// @Tags         Media
// @Accept       json
// @Produce      json
// @Param        body  body  getUploadURLRequest  true  "Upload URL parameters"
// @Success      200  {object}  object{upload_url=object{url=string,method=string,expires_at=string}}
// @Failure      400  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /media/upload-url [post]
// @Security     BearerAuth
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

// GetDownloadURL godoc
// @Summary      Generate a presigned download URL for a media file
// @Tags         Media
// @Accept       json
// @Produce      json
// @Param        id  path  string  true  "Media ID"
// @Success      200  {object}  object{download_url=object{url=string,method=string,expires_at=string}}
// @Failure      404  {object}  object{error=string}
// @Router       /media/{id}/download-url [get]
// @Security     BearerAuth
func (h *Handler) GetDownloadURL(c *gin.Context) {
	id := c.Param("id")
	url, err := h.mediaUC.GenerateDownloadURL(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "media not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"download_url": url})
}
