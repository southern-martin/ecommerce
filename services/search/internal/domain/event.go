package domain

import "context"

// EventPublisher defines the interface for publishing domain events.
type EventPublisher interface {
	Publish(ctx context.Context, subject string, data interface{}) error
}
