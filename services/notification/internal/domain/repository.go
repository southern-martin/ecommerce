package domain

import "context"

// NotificationRepository defines the interface for notification persistence.
type NotificationRepository interface {
	GetByID(ctx context.Context, id string) (*Notification, error)
	ListByUser(ctx context.Context, userID string, filter NotificationFilter) ([]Notification, int64, error)
	Create(ctx context.Context, notification *Notification) error
	Update(ctx context.Context, notification *Notification) error
	MarkAsRead(ctx context.Context, id string) error
	MarkAllAsRead(ctx context.Context, userID string) error
	CountUnread(ctx context.Context, userID string) (int64, error)
}

// PreferenceRepository defines the interface for notification preference persistence.
type PreferenceRepository interface {
	GetByUser(ctx context.Context, userID string) ([]NotificationPreference, error)
	Upsert(ctx context.Context, preference *NotificationPreference) error
}
