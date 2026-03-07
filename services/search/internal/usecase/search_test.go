package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/southern-martin/ecommerce/services/search/internal/domain"
)

// ---------------------------------------------------------------------------
// Hand-written function-field mocks (shared with index_test.go)
// ---------------------------------------------------------------------------

// --- SearchRepository mock ---

type mockSearchRepo struct {
	indexFn   func(ctx context.Context, index *domain.SearchIndex) error
	deleteFn  func(ctx context.Context, productID string) error
	searchFn  func(ctx context.Context, filter domain.SearchFilter) ([]domain.SearchResult, int64, error)
	suggestFn func(ctx context.Context, query string, limit int) ([]domain.SearchSuggestion, error)
}

func (m *mockSearchRepo) Index(ctx context.Context, index *domain.SearchIndex) error {
	if m.indexFn != nil {
		return m.indexFn(ctx, index)
	}
	return nil
}
func (m *mockSearchRepo) Delete(ctx context.Context, productID string) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, productID)
	}
	return nil
}
func (m *mockSearchRepo) Search(ctx context.Context, filter domain.SearchFilter) ([]domain.SearchResult, int64, error) {
	if m.searchFn != nil {
		return m.searchFn(ctx, filter)
	}
	return nil, 0, nil
}
func (m *mockSearchRepo) Suggest(ctx context.Context, query string, limit int) ([]domain.SearchSuggestion, error) {
	if m.suggestFn != nil {
		return m.suggestFn(ctx, query, limit)
	}
	return nil, nil
}

// --- EventPublisher mock ---

type mockSearchEventPub struct {
	publishFn func(ctx context.Context, subject string, data interface{}) error
}

func (m *mockSearchEventPub) Publish(ctx context.Context, subject string, data interface{}) error {
	if m.publishFn != nil {
		return m.publishFn(ctx, subject, data)
	}
	return nil
}

// ===========================================================================
// SearchUseCase tests
// ===========================================================================

func TestSearch_Success(t *testing.T) {
	repo := &mockSearchRepo{}
	results := []domain.SearchResult{
		{ID: "idx-1", ProductID: "prod-1", Name: "Widget", PriceCents: 999},
		{ID: "idx-2", ProductID: "prod-2", Name: "Gadget", PriceCents: 1999},
	}
	repo.searchFn = func(_ context.Context, filter domain.SearchFilter) ([]domain.SearchResult, int64, error) {
		assert.Equal(t, "widget", filter.Query)
		return results, 2, nil
	}

	uc := NewSearchUseCase(repo)
	res, total, err := uc.Search(context.Background(), domain.SearchFilter{
		Query:    "widget",
		Page:     1,
		PageSize: 20,
	})

	require.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, int64(2), total)
	assert.Equal(t, "Widget", res[0].Name)
}

func TestSearch_EmptyQuery(t *testing.T) {
	repo := &mockSearchRepo{}
	repo.searchFn = func(_ context.Context, filter domain.SearchFilter) ([]domain.SearchResult, int64, error) {
		assert.Empty(t, filter.Query)
		return nil, 0, nil
	}

	uc := NewSearchUseCase(repo)
	res, total, err := uc.Search(context.Background(), domain.SearchFilter{
		Query:    "",
		Page:     1,
		PageSize: 20,
	})

	require.NoError(t, err)
	assert.Empty(t, res)
	assert.Equal(t, int64(0), total)
}

func TestSearch_NormalizesPageDefaults(t *testing.T) {
	repo := &mockSearchRepo{}
	var capturedFilter domain.SearchFilter
	repo.searchFn = func(_ context.Context, filter domain.SearchFilter) ([]domain.SearchResult, int64, error) {
		capturedFilter = filter
		return nil, 0, nil
	}

	uc := NewSearchUseCase(repo)
	_, _, err := uc.Search(context.Background(), domain.SearchFilter{
		Query:    "test",
		Page:     0,
		PageSize: 0,
	})

	require.NoError(t, err)
	assert.Equal(t, 1, capturedFilter.Page)
	assert.Equal(t, 20, capturedFilter.PageSize)
}

func TestSearch_CapsPageSizeMax(t *testing.T) {
	repo := &mockSearchRepo{}
	var capturedFilter domain.SearchFilter
	repo.searchFn = func(_ context.Context, filter domain.SearchFilter) ([]domain.SearchResult, int64, error) {
		capturedFilter = filter
		return nil, 0, nil
	}

	uc := NewSearchUseCase(repo)
	_, _, err := uc.Search(context.Background(), domain.SearchFilter{
		Query:    "test",
		Page:     1,
		PageSize: 200,
	})

	require.NoError(t, err)
	assert.Equal(t, 100, capturedFilter.PageSize)
}

func TestSearch_RepoError(t *testing.T) {
	repo := &mockSearchRepo{}
	repo.searchFn = func(_ context.Context, _ domain.SearchFilter) ([]domain.SearchResult, int64, error) {
		return nil, 0, errors.New("elasticsearch unavailable")
	}

	uc := NewSearchUseCase(repo)
	_, _, err := uc.Search(context.Background(), domain.SearchFilter{
		Query:    "test",
		Page:     1,
		PageSize: 10,
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "elasticsearch unavailable")
}

func TestSearch_WithFilters(t *testing.T) {
	repo := &mockSearchRepo{}
	inStock := true
	var capturedFilter domain.SearchFilter
	repo.searchFn = func(_ context.Context, filter domain.SearchFilter) ([]domain.SearchResult, int64, error) {
		capturedFilter = filter
		return nil, 0, nil
	}

	uc := NewSearchUseCase(repo)
	_, _, err := uc.Search(context.Background(), domain.SearchFilter{
		Query:      "laptop",
		CategoryID: "cat-1",
		MinPrice:   5000,
		MaxPrice:   100000,
		InStock:    &inStock,
		SellerID:   "seller-1",
		SortBy:     "price",
		SortOrder:  "asc",
		Page:       2,
		PageSize:   10,
	})

	require.NoError(t, err)
	assert.Equal(t, "cat-1", capturedFilter.CategoryID)
	assert.Equal(t, int64(5000), capturedFilter.MinPrice)
	assert.Equal(t, int64(100000), capturedFilter.MaxPrice)
	require.NotNil(t, capturedFilter.InStock)
	assert.True(t, *capturedFilter.InStock)
	assert.Equal(t, "seller-1", capturedFilter.SellerID)
	assert.Equal(t, "price", capturedFilter.SortBy)
	assert.Equal(t, "asc", capturedFilter.SortOrder)
}

// ===========================================================================
// Suggest tests
// ===========================================================================

func TestSuggest_Success(t *testing.T) {
	repo := &mockSearchRepo{}
	suggestions := []domain.SearchSuggestion{
		{Text: "wireless mouse", Type: "product", ProductID: "prod-1"},
		{Text: "wireless keyboard", Type: "product", ProductID: "prod-2"},
	}
	repo.suggestFn = func(_ context.Context, query string, limit int) ([]domain.SearchSuggestion, error) {
		assert.Equal(t, "wire", query)
		assert.Equal(t, 5, limit)
		return suggestions, nil
	}

	uc := NewSearchUseCase(repo)
	result, err := uc.Suggest(context.Background(), "wire", 5)

	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "wireless mouse", result[0].Text)
}

func TestSuggest_DefaultLimit(t *testing.T) {
	repo := &mockSearchRepo{}
	var capturedLimit int
	repo.suggestFn = func(_ context.Context, _ string, limit int) ([]domain.SearchSuggestion, error) {
		capturedLimit = limit
		return nil, nil
	}

	uc := NewSearchUseCase(repo)
	_, err := uc.Suggest(context.Background(), "test", 0)

	require.NoError(t, err)
	assert.Equal(t, 10, capturedLimit)
}

func TestSuggest_CapsLimitMax(t *testing.T) {
	repo := &mockSearchRepo{}
	var capturedLimit int
	repo.suggestFn = func(_ context.Context, _ string, limit int) ([]domain.SearchSuggestion, error) {
		capturedLimit = limit
		return nil, nil
	}

	uc := NewSearchUseCase(repo)
	_, err := uc.Suggest(context.Background(), "test", 100)

	require.NoError(t, err)
	assert.Equal(t, 50, capturedLimit)
}

func TestSuggest_NegativeLimit(t *testing.T) {
	repo := &mockSearchRepo{}
	var capturedLimit int
	repo.suggestFn = func(_ context.Context, _ string, limit int) ([]domain.SearchSuggestion, error) {
		capturedLimit = limit
		return nil, nil
	}

	uc := NewSearchUseCase(repo)
	_, err := uc.Suggest(context.Background(), "test", -5)

	require.NoError(t, err)
	assert.Equal(t, 10, capturedLimit)
}

func TestSuggest_RepoError(t *testing.T) {
	repo := &mockSearchRepo{}
	repo.suggestFn = func(_ context.Context, _ string, _ int) ([]domain.SearchSuggestion, error) {
		return nil, errors.New("suggest failed")
	}

	uc := NewSearchUseCase(repo)
	_, err := uc.Suggest(context.Background(), "test", 5)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "suggest failed")
}
