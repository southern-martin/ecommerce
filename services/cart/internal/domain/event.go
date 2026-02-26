package domain

import "context"

// EventPublisher defines the interface for publishing domain events.
type EventPublisher interface {
	// Publish publishes an event with the given subject and payload.
	Publish(ctx context.Context, subject string, payload interface{}) error
}

// Event subjects for cart domain events.
const (
	EventCartItemAdded   = "cart.item.added"
	EventCartItemRemoved = "cart.item.removed"
	EventCartItemUpdated = "cart.item.updated"
	EventCartCleared     = "cart.cleared"
)

// CartEvent represents a cart domain event payload.
type CartEvent struct {
	UserID    string    `json:"user_id"`
	ProductID string    `json:"product_id,omitempty"`
	VariantID string    `json:"variant_id,omitempty"`
	Quantity  int       `json:"quantity,omitempty"`
	Cart      *Cart     `json:"cart,omitempty"`
}
