package domain

import "context"

// ReviewRepository defines the interface for review persistence.
type ReviewRepository interface {
	GetByID(ctx context.Context, id string) (*Review, error)
	ListByProduct(ctx context.Context, productID string, filter ReviewFilter) ([]Review, int64, error)
	ListByUser(ctx context.Context, userID string, page, pageSize int) ([]Review, int64, error)
	Create(ctx context.Context, review *Review) error
	Update(ctx context.Context, review *Review) error
	Delete(ctx context.Context, id string) error
	GetSummary(ctx context.Context, productID string) (*ReviewSummary, error)
}
