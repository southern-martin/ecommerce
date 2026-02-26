package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/southern-martin/ecommerce/services/notification/internal/domain"
	"github.com/southern-martin/ecommerce/services/notification/internal/usecase"
)

// Handler holds the HTTP handlers for the notification service.
type Handler struct {
	notificationUC *usecase.NotificationUseCase
	preferenceUC   *usecase.PreferenceUseCase
}

// NewHandler creates a new Handler.
func NewHandler(notificationUC *usecase.NotificationUseCase, preferenceUC *usecase.PreferenceUseCase) *Handler {
	return &Handler{
		notificationUC: notificationUC,
		preferenceUC:   preferenceUC,
	}
}

// Health handles health check requests.
func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "notification"})
}

// SendNotification handles POST /api/v1/notifications.
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

// ListNotifications handles GET /api/v1/notifications.
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

// GetNotification handles GET /api/v1/notifications/:id.
func (h *Handler) GetNotification(c *gin.Context) {
	id := c.Param("id")

	notification, err := h.notificationUC.GetNotification(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "notification not found"})
		return
	}

	c.JSON(http.StatusOK, notification)
}

// MarkAsRead handles PATCH /api/v1/notifications/:id/read.
func (h *Handler) MarkAsRead(c *gin.Context) {
	id := c.Param("id")

	if err := h.notificationUC.MarkAsRead(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "notification marked as read"})
}

// MarkAllAsRead handles PATCH /api/v1/notifications/read-all.
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

// GetUnreadCount handles GET /api/v1/notifications/unread-count.
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

// GetPreferences handles GET /api/v1/notifications/preferences.
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

// UpdatePreference handles PATCH /api/v1/notifications/preferences.
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
