package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	pkgerrors "github.com/southern-martin/ecommerce/pkg/errors"
	"github.com/southern-martin/ecommerce/services/user/internal/domain"
)

func TestAddItem_Success(t *testing.T) {
	var created *domain.WishlistItem
	repo := &mockWishlistRepo{
		existsFn: func(_ context.Context, _, _ string) (bool, error) {
			return false, nil
		},
		createFn: func(_ context.Context, item *domain.WishlistItem) error {
			created = item
			return nil
		},
	}

	uc := NewWishlistUseCase(repo, zerolog.Nop())
	err := uc.AddItem(context.Background(), "user-1", "product-1")

	require.NoError(t, err)
	require.NotNil(t, created)
	assert.Equal(t, "user-1", created.UserID)
	assert.Equal(t, "product-1", created.ProductID)
}

func TestAddItem_AlreadyInWishlist(t *testing.T) {
	repo := &mockWishlistRepo{
		existsFn: func(_ context.Context, _, _ string) (bool, error) {
			return true, nil
		},
	}

	uc := NewWishlistUseCase(repo, zerolog.Nop())
	err := uc.AddItem(context.Background(), "user-1", "product-1")

	require.Error(t, err)
	var conflictErr *pkgerrors.ConflictError
	assert.ErrorAs(t, err, &conflictErr)
	assert.Equal(t, "ALREADY_IN_WISHLIST", conflictErr.Code)
}

func TestRemoveItem_Success(t *testing.T) {
	deleted := false
	repo := &mockWishlistRepo{
		deleteFn: func(_ context.Context, userID, productID string) error {
			assert.Equal(t, "user-1", userID)
			assert.Equal(t, "product-1", productID)
			deleted = true
			return nil
		},
	}

	uc := NewWishlistUseCase(repo, zerolog.Nop())
	err := uc.RemoveItem(context.Background(), "user-1", "product-1")

	require.NoError(t, err)
	assert.True(t, deleted)
}

func TestListItems_Success(t *testing.T) {
	now := time.Now()
	expected := []domain.WishlistItem{
		{ID: "w1", UserID: "user-1", ProductID: "p1", CreatedAt: now},
		{ID: "w2", UserID: "user-1", ProductID: "p2", CreatedAt: now},
		{ID: "w3", UserID: "user-1", ProductID: "p3", CreatedAt: now},
	}
	repo := &mockWishlistRepo{
		listByUserIDFn: func(_ context.Context, userID string) ([]domain.WishlistItem, error) {
			assert.Equal(t, "user-1", userID)
			return expected, nil
		},
	}

	uc := NewWishlistUseCase(repo, zerolog.Nop())
	items, err := uc.ListItems(context.Background(), "user-1")

	require.NoError(t, err)
	assert.Len(t, items, 3)
	assert.Equal(t, "p1", items[0].ProductID)
	assert.Equal(t, "p2", items[1].ProductID)
	assert.Equal(t, "p3", items[2].ProductID)
}

func TestAddItem_ExistsCheckError(t *testing.T) {
	repo := &mockWishlistRepo{
		existsFn: func(_ context.Context, _, _ string) (bool, error) {
			return false, assert.AnError
		},
	}

	uc := NewWishlistUseCase(repo, zerolog.Nop())
	err := uc.AddItem(context.Background(), "user-1", "product-1")

	require.Error(t, err)
	assert.ErrorIs(t, err, assert.AnError)
}
