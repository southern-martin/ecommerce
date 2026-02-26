package postgres

import (
	"time"

	"github.com/southern-martin/ecommerce/services/notification/internal/domain"
)

// NotificationModel is the GORM model for notifications.
type NotificationModel struct {
	ID        string     `gorm:"type:uuid;primaryKey"`
	UserID    string     `gorm:"type:uuid;index;not null"`
	Type      string     `gorm:"type:varchar(50);not null"`
	Channel   string     `gorm:"type:varchar(20);not null"`
	Subject   string     `gorm:"type:varchar(500);not null"`
	Body      string     `gorm:"type:text;not null"`
	Data      string     `gorm:"type:text"`
	Status    string     `gorm:"type:varchar(20);not null;default:'queued'"`
	SentAt    *time.Time `gorm:"type:timestamptz"`
	ReadAt    *time.Time `gorm:"type:timestamptz"`
	CreatedAt time.Time  `gorm:"type:timestamptz;not null;default:now()"`
}

// TableName returns the table name for the notification model.
func (NotificationModel) TableName() string {
	return "notifications"
}

// ToDomain converts a NotificationModel to a domain Notification.
func (m *NotificationModel) ToDomain() *domain.Notification {
	return &domain.Notification{
		ID:        m.ID,
		UserID:    m.UserID,
		Type:      domain.NotificationType(m.Type),
		Channel:   domain.NotificationChannel(m.Channel),
		Subject:   m.Subject,
		Body:      m.Body,
		Data:      m.Data,
		Status:    domain.NotificationStatus(m.Status),
		SentAt:    m.SentAt,
		ReadAt:    m.ReadAt,
		CreatedAt: m.CreatedAt,
	}
}

// ToModel converts a domain Notification to a NotificationModel.
func ToNotificationModel(n *domain.Notification) *NotificationModel {
	return &NotificationModel{
		ID:        n.ID,
		UserID:    n.UserID,
		Type:      string(n.Type),
		Channel:   string(n.Channel),
		Subject:   n.Subject,
		Body:      n.Body,
		Data:      n.Data,
		Status:    string(n.Status),
		SentAt:    n.SentAt,
		ReadAt:    n.ReadAt,
		CreatedAt: n.CreatedAt,
	}
}

// PreferenceModel is the GORM model for notification preferences.
type PreferenceModel struct {
	ID        string    `gorm:"type:uuid;primaryKey"`
	UserID    string    `gorm:"type:uuid;index;not null"`
	Channel   string    `gorm:"type:varchar(20);not null"`
	Enabled   bool      `gorm:"type:bool;not null;default:true"`
	CreatedAt time.Time `gorm:"type:timestamptz;not null;default:now()"`
}

// TableName returns the table name for the preference model.
func (PreferenceModel) TableName() string {
	return "notification_preferences"
}

// ToDomain converts a PreferenceModel to a domain NotificationPreference.
func (m *PreferenceModel) ToDomain() *domain.NotificationPreference {
	return &domain.NotificationPreference{
		ID:        m.ID,
		UserID:    m.UserID,
		Channel:   domain.NotificationChannel(m.Channel),
		Enabled:   m.Enabled,
		CreatedAt: m.CreatedAt,
	}
}

// ToPreferenceModel converts a domain NotificationPreference to a PreferenceModel.
func ToPreferenceModel(p *domain.NotificationPreference) *PreferenceModel {
	return &PreferenceModel{
		ID:        p.ID,
		UserID:    p.UserID,
		Channel:   string(p.Channel),
		Enabled:   p.Enabled,
		CreatedAt: p.CreatedAt,
	}
}
