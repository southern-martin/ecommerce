package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/southern-martin/ecommerce/services/notification/internal/domain"
)

// NotificationUseCase handles notification business logic.
type NotificationUseCase struct {
	repo      domain.NotificationRepository
	publisher domain.EventPublisher
}

// NewNotificationUseCase creates a new NotificationUseCase.
func NewNotificationUseCase(repo domain.NotificationRepository, publisher domain.EventPublisher) *NotificationUseCase {
	return &NotificationUseCase{
		repo:      repo,
		publisher: publisher,
	}
}

// SendNotificationRequest holds the data needed to send a notification.
type SendNotificationRequest struct {
	UserID  string `json:"user_id" binding:"required"`
	Type    string `json:"type" binding:"required"`
	Channel string `json:"channel" binding:"required"`
	Subject string `json:"subject" binding:"required"`
	Body    string `json:"body" binding:"required"`
	Data    string `json:"data"`
}

// SendNotification creates and sends a notification.
func (uc *NotificationUseCase) SendNotification(ctx context.Context, req SendNotificationRequest) (*domain.Notification, error) {
	notification := &domain.Notification{
		ID:        uuid.New().String(),
		UserID:    req.UserID,
		Type:      domain.NotificationType(req.Type),
		Channel:   domain.NotificationChannel(req.Channel),
		Subject:   req.Subject,
		Body:      req.Body,
		Data:      req.Data,
		Status:    domain.StatusQueued,
		CreatedAt: time.Now(),
	}

	if err := uc.repo.Create(ctx, notification); err != nil {
		return nil, err
	}

	// Mock send: simulate sending the notification
	now := time.Now()
	notification.SentAt = &now
	notification.Status = domain.StatusSent

	if err := uc.repo.Update(ctx, notification); err != nil {
		log.Error().Err(err).Str("id", notification.ID).Msg("failed to update notification status after send")
		return nil, err
	}

	// Publish event
	_ = uc.publisher.Publish(ctx, "notification.sent", map[string]string{
		"notification_id": notification.ID,
		"user_id":         notification.UserID,
		"type":            string(notification.Type),
		"channel":         string(notification.Channel),
	})

	log.Info().Str("id", notification.ID).Str("channel", string(notification.Channel)).Msg("notification sent")

	return notification, nil
}

// GetNotification retrieves a notification by ID.
func (uc *NotificationUseCase) GetNotification(ctx context.Context, id string) (*domain.Notification, error) {
	return uc.repo.GetByID(ctx, id)
}

// ListUserNotifications retrieves notifications for a user with filters.
func (uc *NotificationUseCase) ListUserNotifications(ctx context.Context, userID string, filter domain.NotificationFilter) ([]domain.Notification, int64, error) {
	return uc.repo.ListByUser(ctx, userID, filter)
}

// MarkAsRead marks a notification as read.
func (uc *NotificationUseCase) MarkAsRead(ctx context.Context, id string) error {
	return uc.repo.MarkAsRead(ctx, id)
}

// MarkAllAsRead marks all notifications for a user as read.
func (uc *NotificationUseCase) MarkAllAsRead(ctx context.Context, userID string) error {
	return uc.repo.MarkAllAsRead(ctx, userID)
}

// GetUnreadCount returns the count of unread notifications for a user.
func (uc *NotificationUseCase) GetUnreadCount(ctx context.Context, userID string) (int64, error) {
	return uc.repo.CountUnread(ctx, userID)
}
