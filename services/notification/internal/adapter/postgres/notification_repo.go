package postgres

import (
	"context"
	"time"

	"github.com/southern-martin/ecommerce/services/notification/internal/domain"
	"gorm.io/gorm"
)

// NotificationRepo implements domain.NotificationRepository using GORM.
type NotificationRepo struct {
	db *gorm.DB
}

// NewNotificationRepo creates a new NotificationRepo.
func NewNotificationRepo(db *gorm.DB) *NotificationRepo {
	return &NotificationRepo{db: db}
}

// GetByID retrieves a notification by its ID.
func (r *NotificationRepo) GetByID(ctx context.Context, id string) (*domain.Notification, error) {
	var model NotificationModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		return nil, err
	}
	return model.ToDomain(), nil
}

// ListByUser retrieves notifications for a user with optional filters and pagination.
func (r *NotificationRepo) ListByUser(ctx context.Context, userID string, filter domain.NotificationFilter) ([]domain.Notification, int64, error) {
	var models []NotificationModel
	var total int64

	query := r.db.WithContext(ctx).Where("user_id = ?", userID)

	if filter.Type != "" {
		query = query.Where("type = ?", string(filter.Type))
	}
	if filter.Channel != "" {
		query = query.Where("channel = ?", string(filter.Channel))
	}
	if filter.Status != "" {
		query = query.Where("status = ?", string(filter.Status))
	}

	// Count total before pagination
	if err := query.Model(&NotificationModel{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	page := filter.Page
	if page < 1 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize < 1 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	notifications := make([]domain.Notification, len(models))
	for i, m := range models {
		notifications[i] = *m.ToDomain()
	}

	return notifications, total, nil
}

// Create inserts a new notification.
func (r *NotificationRepo) Create(ctx context.Context, notification *domain.Notification) error {
	model := ToNotificationModel(notification)
	return r.db.WithContext(ctx).Create(model).Error
}

// Update updates an existing notification.
func (r *NotificationRepo) Update(ctx context.Context, notification *domain.Notification) error {
	model := ToNotificationModel(notification)
	return r.db.WithContext(ctx).Save(model).Error
}

// MarkAsRead marks a single notification as read.
func (r *NotificationRepo) MarkAsRead(ctx context.Context, id string) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&NotificationModel{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":  string(domain.StatusRead),
			"read_at": now,
		}).Error
}

// MarkAllAsRead marks all notifications for a user as read.
func (r *NotificationRepo) MarkAllAsRead(ctx context.Context, userID string) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&NotificationModel{}).
		Where("user_id = ? AND read_at IS NULL", userID).
		Updates(map[string]interface{}{
			"status":  string(domain.StatusRead),
			"read_at": now,
		}).Error
}

// CountUnread returns the count of unread notifications for a user.
func (r *NotificationRepo) CountUnread(ctx context.Context, userID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&NotificationModel{}).
		Where("user_id = ? AND read_at IS NULL", userID).
		Count(&count).Error
	return count, err
}
