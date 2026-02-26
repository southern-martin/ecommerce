package domain

import "time"

// Notification represents a notification sent to a user.
type Notification struct {
	ID        string
	UserID    string
	Type      NotificationType
	Channel   NotificationChannel
	Subject   string
	Body      string
	Data      string
	Status    NotificationStatus
	SentAt    *time.Time
	ReadAt    *time.Time
	CreatedAt time.Time
}

// NotificationType represents the type of notification.
type NotificationType string

const (
	TypeOrderUpdate    NotificationType = "order_update"
	TypePaymentUpdate  NotificationType = "payment_update"
	TypeShipmentUpdate NotificationType = "shipment_update"
	TypeReturnUpdate   NotificationType = "return_update"
	TypePromotion      NotificationType = "promotion"
	TypeSystem         NotificationType = "system"
)

// NotificationChannel represents the delivery channel for a notification.
type NotificationChannel string

const (
	ChannelEmail NotificationChannel = "email"
	ChannelPush  NotificationChannel = "push"
	ChannelInApp NotificationChannel = "in_app"
)

// NotificationStatus represents the current status of a notification.
type NotificationStatus string

const (
	StatusQueued NotificationStatus = "queued"
	StatusSent   NotificationStatus = "sent"
	StatusFailed NotificationStatus = "failed"
	StatusRead   NotificationStatus = "read"
)

// NotificationPreference represents a user's notification preference for a channel.
type NotificationPreference struct {
	ID        string
	UserID    string
	Channel   NotificationChannel
	Enabled   bool
	CreatedAt time.Time
}

// NotificationFilter holds filter criteria for listing notifications.
type NotificationFilter struct {
	UserID   string
	Type     NotificationType
	Channel  NotificationChannel
	Status   NotificationStatus
	Page     int
	PageSize int
}
