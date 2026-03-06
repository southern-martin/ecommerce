package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/southern-martin/ecommerce/services/notification/internal/domain"
	"github.com/southern-martin/ecommerce/services/notification/internal/usecase"
)

// Handler holds the HTTP handlers for the notification service.
type Handler struct {
	notificationUC *usecase.NotificationUseCase
	preferenceUC   *usecase.PreferenceUseCase
	db             *gorm.DB
}

// NewHandler creates a new Handler.
func NewHandler(notificationUC *usecase.NotificationUseCase, preferenceUC *usecase.PreferenceUseCase, db *gorm.DB) *Handler {
	return &Handler{
		notificationUC: notificationUC,
		preferenceUC:   preferenceUC,
		db:             db,
	}
}

// Health handles health check requests.
func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "notification"})
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

// SendNotification godoc
// @Summary      Send a notification
// @Tags         Notifications
// @Accept       json
// @Produce      json
// @Param        body  body  usecase.SendNotificationRequest  true  "Notification payload"
// @Success      201  {object}  domain.Notification
// @Failure      400  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /notifications [post]
func (h *Handler) SendNotification(c *gin.Context) {
	var req usecase.SendNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	notification, err := h.notificationUC.SendNotification(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, notification)
}

// ListNotifications godoc
// @Summary      List notifications for the authenticated user
// @Tags         Notifications
// @Produce      json
// @Param        X-User-ID  header  string  true  "User ID"
// @Param        type       query   string  false  "Filter by notification type"
// @Param        channel    query   string  false  "Filter by channel"
// @Param        status     query   string  false  "Filter by status"
// @Param        page       query   int     false  "Page number"   default(1)
// @Param        page_size  query   int     false  "Page size"     default(20)
// @Success      200  {object}  object{notifications=[]domain.Notification,total=int64,page=int,page_size=int}
// @Failure      400  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /notifications [get]
// @Security     BearerAuth
func (h *Handler) ListNotifications(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-User-ID header is required"})
		return
	}

	filter := domain.NotificationFilter{
		UserID:  userID,
		Type:    domain.NotificationType(c.Query("type")),
		Channel: domain.NotificationChannel(c.Query("channel")),
		Status:  domain.NotificationStatus(c.Query("status")),
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	filter.Page = page
	filter.PageSize = pageSize

	notifications, total, err := h.notificationUC.ListUserNotifications(c.Request.Context(), userID, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"notifications": notifications,
		"total":         total,
		"page":          page,
		"page_size":     pageSize,
	})
}

// GetNotification godoc
// @Summary      Get a notification by ID
// @Tags         Notifications
// @Produce      json
// @Param        id  path  string  true  "Notification ID"
// @Success      200  {object}  domain.Notification
// @Failure      404  {object}  object{error=string}
// @Router       /notifications/{id} [get]
func (h *Handler) GetNotification(c *gin.Context) {
	id := c.Param("id")

	notification, err := h.notificationUC.GetNotification(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "notification not found"})
		return
	}

	c.JSON(http.StatusOK, notification)
}

// MarkAsRead godoc
// @Summary      Mark a notification as read
// @Tags         Notifications
// @Produce      json
// @Param        id  path  string  true  "Notification ID"
// @Success      200  {object}  object{message=string}
// @Failure      500  {object}  object{error=string}
// @Router       /notifications/{id}/read [patch]
func (h *Handler) MarkAsRead(c *gin.Context) {
	id := c.Param("id")

	if err := h.notificationUC.MarkAsRead(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "notification marked as read"})
}

// MarkAllAsRead godoc
// @Summary      Mark all notifications as read for the authenticated user
// @Tags         Notifications
// @Produce      json
// @Param        X-User-ID  header  string  true  "User ID"
// @Success      200  {object}  object{message=string}
// @Failure      400  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /notifications/read-all [patch]
// @Security     BearerAuth
func (h *Handler) MarkAllAsRead(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-User-ID header is required"})
		return
	}

	if err := h.notificationUC.MarkAllAsRead(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "all notifications marked as read"})
}

// GetUnreadCount godoc
// @Summary      Get unread notification count for the authenticated user
// @Tags         Notifications
// @Produce      json
// @Param        X-User-ID  header  string  true  "User ID"
// @Success      200  {object}  object{unread_count=int64}
// @Failure      400  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /notifications/unread-count [get]
// @Security     BearerAuth
func (h *Handler) GetUnreadCount(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-User-ID header is required"})
		return
	}

	count, err := h.notificationUC.GetUnreadCount(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"unread_count": count})
}

// GetPreferences godoc
// @Summary      Get notification preferences for the authenticated user
// @Tags         Notifications
// @Produce      json
// @Param        X-User-ID  header  string  true  "User ID"
// @Success      200  {object}  object{preferences=[]domain.NotificationPreference}
// @Failure      400  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /notifications/preferences [get]
// @Security     BearerAuth
func (h *Handler) GetPreferences(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-User-ID header is required"})
		return
	}

	preferences, err := h.preferenceUC.GetPreferences(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"preferences": preferences})
}

// UpdatePreference godoc
// @Summary      Update a notification preference
// @Tags         Notifications
// @Accept       json
// @Produce      json
// @Param        X-User-ID  header  string                          true  "User ID"
// @Param        body       body    usecase.UpdatePreferenceRequest  true  "Preference update payload"
// @Success      200  {object}  object{message=string}
// @Failure      400  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /notifications/preferences [patch]
// @Security     BearerAuth
func (h *Handler) UpdatePreference(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-User-ID header is required"})
		return
	}

	var req usecase.UpdatePreferenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.UserID = userID

	if err := h.preferenceUC.UpdatePreference(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "preference updated"})
}
