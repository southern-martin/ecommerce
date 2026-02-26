package usecase

import (
	"context"

	"github.com/southern-martin/ecommerce/services/search/internal/domain"
)

// SearchUseCase handles search query operations.
type SearchUseCase struct {
	repo domain.SearchRepository
}

// NewSearchUseCase creates a new SearchUseCase.
func NewSearchUseCase(repo domain.SearchRepository) *SearchUseCase {
	return &SearchUseCase{repo: repo}
}

// Search performs a search with the given filter, normalizing pagination params.
func (uc *SearchUseCase) Search(ctx context.Context, filter domain.SearchFilter) ([]domain.SearchResult, int64, error) {
	// Normalize pagination
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PageSize < 1 {
		filter.PageSize = 20
	}
	if filter.PageSize > 100 {
		filter.PageSize = 100
	}

	return uc.repo.Search(ctx, filter)
}

// Suggest returns search autocomplete suggestions.
func (uc *SearchUseCase) Suggest(ctx context.Context, query string, limit int) ([]domain.SearchSuggestion, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}

	return uc.repo.Suggest(ctx, query, limit)
}
