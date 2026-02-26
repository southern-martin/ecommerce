package domain

import "context"

// EventPublisher defines the interface for publishing domain events.
type EventPublisher interface {
	PublishProductCreated(ctx context.Context, product *Product) error
	PublishProductUpdated(ctx context.Context, product *Product) error
	PublishProductDeleted(ctx context.Context, productID string) error
	PublishStockUpdated(ctx context.Context, variantID string, newStock int, delta int) error
}
