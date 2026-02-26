package domain

import "context"

// SearchRepository defines the interface for search index operations.
type SearchRepository interface {
	Index(ctx context.Context, index *SearchIndex) error
	Delete(ctx context.Context, productID string) error
	Search(ctx context.Context, filter SearchFilter) ([]SearchResult, int64, error)
	Suggest(ctx context.Context, query string, limit int) ([]SearchSuggestion, error)
}
