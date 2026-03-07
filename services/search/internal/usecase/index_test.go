package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/southern-martin/ecommerce/services/search/internal/domain"
)

// Mocks are defined in search_test.go (shared within same package).

// ===========================================================================
// IndexProduct tests
// ===========================================================================

func TestIndexProduct_Success(t *testing.T) {
	repo := &mockSearchRepo{}
	pub := &mockSearchEventPub{}
	var indexed *domain.SearchIndex
	repo.indexFn = func(_ context.Context, idx *domain.SearchIndex) error {
		indexed = idx
		return nil
	}

	uc := NewIndexUseCase(repo, pub)
	idx := &domain.SearchIndex{
		ProductID:   "prod-1",
		Name:        "Wireless Mouse",
		Description: "A nice mouse",
		PriceCents:  2999,
		Currency:    "USD",
		CategoryID:  "cat-1",
		SellerID:    "seller-1",
		InStock:     true,
		Tags:        []string{"electronics", "mouse"},
	}

	err := uc.IndexProduct(context.Background(), idx)

	require.NoError(t, err)
	require.NotNil(t, indexed)
	assert.NotEmpty(t, indexed.ID)
	assert.Equal(t, "prod-1", indexed.ProductID)
	assert.Equal(t, "Wireless Mouse", indexed.Name)
	assert.False(t, indexed.CreatedAt.IsZero())
	assert.False(t, indexed.UpdatedAt.IsZero())
}

func TestIndexProduct_GeneratesIDWhenEmpty(t *testing.T) {
	repo := &mockSearchRepo{}
	pub := &mockSearchEventPub{}
	repo.indexFn = func(_ context.Context, _ *domain.SearchIndex) error { return nil }

	uc := NewIndexUseCase(repo, pub)
	idx := &domain.SearchIndex{
		ID:        "", // empty, should be generated
		ProductID: "prod-1",
		Name:      "Widget",
	}

	err := uc.IndexProduct(context.Background(), idx)

	require.NoError(t, err)
	assert.NotEmpty(t, idx.ID)
}

func TestIndexProduct_PreservesExistingID(t *testing.T) {
	repo := &mockSearchRepo{}
	pub := &mockSearchEventPub{}
	repo.indexFn = func(_ context.Context, _ *domain.SearchIndex) error { return nil }

	uc := NewIndexUseCase(repo, pub)
	idx := &domain.SearchIndex{
		ID:        "existing-id",
		ProductID: "prod-1",
		Name:      "Widget",
	}

	err := uc.IndexProduct(context.Background(), idx)

	require.NoError(t, err)
	assert.Equal(t, "existing-id", idx.ID)
}

func TestIndexProduct_RepoError(t *testing.T) {
	repo := &mockSearchRepo{}
	pub := &mockSearchEventPub{}
	repo.indexFn = func(_ context.Context, _ *domain.SearchIndex) error {
		return errors.New("index failed")
	}

	uc := NewIndexUseCase(repo, pub)
	idx := &domain.SearchIndex{
		ProductID: "prod-1",
		Name:      "Widget",
	}

	err := uc.IndexProduct(context.Background(), idx)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "index failed")
}

func TestIndexProduct_PublishEventErrorDoesNotFail(t *testing.T) {
	repo := &mockSearchRepo{}
	pub := &mockSearchEventPub{}
	repo.indexFn = func(_ context.Context, _ *domain.SearchIndex) error { return nil }
	pub.publishFn = func(_ context.Context, _ string, _ interface{}) error {
		return errors.New("event bus down")
	}

	uc := NewIndexUseCase(repo, pub)
	idx := &domain.SearchIndex{
		ProductID: "prod-1",
		Name:      "Widget",
	}

	err := uc.IndexProduct(context.Background(), idx)

	// Event publish failure should not cause IndexProduct to fail
	require.NoError(t, err)
}

func TestIndexProduct_PublishesCorrectSubject(t *testing.T) {
	repo := &mockSearchRepo{}
	pub := &mockSearchEventPub{}
	repo.indexFn = func(_ context.Context, _ *domain.SearchIndex) error { return nil }
	var capturedSubject string
	pub.publishFn = func(_ context.Context, subject string, _ interface{}) error {
		capturedSubject = subject
		return nil
	}

	uc := NewIndexUseCase(repo, pub)
	idx := &domain.SearchIndex{
		ProductID: "prod-1",
		Name:      "Widget",
	}

	err := uc.IndexProduct(context.Background(), idx)

	require.NoError(t, err)
	assert.Equal(t, "search.product.indexed", capturedSubject)
}

// ===========================================================================
// RemoveProduct (DeleteProduct) tests
// ===========================================================================

func TestRemoveProduct_Success(t *testing.T) {
	repo := &mockSearchRepo{}
	pub := &mockSearchEventPub{}
	var deletedProductID string
	repo.deleteFn = func(_ context.Context, productID string) error {
		deletedProductID = productID
		return nil
	}

	uc := NewIndexUseCase(repo, pub)
	err := uc.RemoveProduct(context.Background(), "prod-1")

	require.NoError(t, err)
	assert.Equal(t, "prod-1", deletedProductID)
}

func TestRemoveProduct_RepoError(t *testing.T) {
	repo := &mockSearchRepo{}
	pub := &mockSearchEventPub{}
	repo.deleteFn = func(_ context.Context, _ string) error {
		return errors.New("delete failed")
	}

	uc := NewIndexUseCase(repo, pub)
	err := uc.RemoveProduct(context.Background(), "prod-1")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "delete failed")
}

func TestRemoveProduct_PublishEventErrorDoesNotFail(t *testing.T) {
	repo := &mockSearchRepo{}
	pub := &mockSearchEventPub{}
	repo.deleteFn = func(_ context.Context, _ string) error { return nil }
	pub.publishFn = func(_ context.Context, _ string, _ interface{}) error {
		return errors.New("event bus down")
	}

	uc := NewIndexUseCase(repo, pub)
	err := uc.RemoveProduct(context.Background(), "prod-1")

	// Event publish failure should not cause RemoveProduct to fail
	require.NoError(t, err)
}

func TestRemoveProduct_PublishesCorrectSubject(t *testing.T) {
	repo := &mockSearchRepo{}
	pub := &mockSearchEventPub{}
	repo.deleteFn = func(_ context.Context, _ string) error { return nil }
	var capturedSubject string
	pub.publishFn = func(_ context.Context, subject string, _ interface{}) error {
		capturedSubject = subject
		return nil
	}

	uc := NewIndexUseCase(repo, pub)
	err := uc.RemoveProduct(context.Background(), "prod-1")

	require.NoError(t, err)
	assert.Equal(t, "search.product.deindexed", capturedSubject)
}
