package domain

import "context"

// CartRepository defines the interface for cart persistence.
type CartRepository interface {
	// GetCart retrieves the cart for a given user. Returns an empty cart if none exists.
	GetCart(ctx context.Context, userID string) (*Cart, error)
	// SaveCart persists the cart state.
	SaveCart(ctx context.Context, cart *Cart) error
	// DeleteCart removes the cart for a given user.
	DeleteCart(ctx context.Context, userID string) error
}
