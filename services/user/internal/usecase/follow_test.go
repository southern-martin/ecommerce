package usecase

import (
	"context"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	pkgerrors "github.com/southern-martin/ecommerce/pkg/errors"
	"github.com/southern-martin/ecommerce/services/user/internal/domain"
)

func TestFollow_Success(t *testing.T) {
	var created *domain.UserFollow
	repo := &mockFollowRepo{
		existsFn: func(_ context.Context, _, _ string) (bool, error) {
			return false, nil
		},
		createFn: func(_ context.Context, f *domain.UserFollow) error {
			created = f
			return nil
		},
	}

	uc := NewFollowUseCase(repo, zerolog.Nop())
	err := uc.Follow(context.Background(), "follower-1", "seller-1")

	require.NoError(t, err)
	require.NotNil(t, created)
	assert.Equal(t, "follower-1", created.FollowerID)
	assert.Equal(t, "seller-1", created.SellerID)
}

func TestFollow_AlreadyFollowing(t *testing.T) {
	repo := &mockFollowRepo{
		existsFn: func(_ context.Context, _, _ string) (bool, error) {
			return true, nil
		},
	}

	uc := NewFollowUseCase(repo, zerolog.Nop())
	err := uc.Follow(context.Background(), "follower-1", "seller-1")

	require.Error(t, err)
	var conflictErr *pkgerrors.ConflictError
	assert.ErrorAs(t, err, &conflictErr)
	assert.Equal(t, "ALREADY_FOLLOWING", conflictErr.Code)
}

func TestUnfollow_Success(t *testing.T) {
	deleted := false
	repo := &mockFollowRepo{
		deleteFn: func(_ context.Context, followerID, sellerID string) error {
			assert.Equal(t, "follower-1", followerID)
			assert.Equal(t, "seller-1", sellerID)
			deleted = true
			return nil
		},
	}

	uc := NewFollowUseCase(repo, zerolog.Nop())
	err := uc.Unfollow(context.Background(), "follower-1", "seller-1")

	require.NoError(t, err)
	assert.True(t, deleted)
}

func TestListFollowed_Success(t *testing.T) {
	expected := []domain.SellerProfile{
		{ID: "s1", StoreName: "Store 1"},
		{ID: "s2", StoreName: "Store 2"},
	}
	repo := &mockFollowRepo{
		listByFollowerFn: func(_ context.Context, followerID string, page, size int) ([]domain.SellerProfile, int64, error) {
			assert.Equal(t, "follower-1", followerID)
			assert.Equal(t, 1, page)
			assert.Equal(t, 20, size)
			return expected, 2, nil
		},
	}

	uc := NewFollowUseCase(repo, zerolog.Nop())
	sellers, total, err := uc.ListFollowed(context.Background(), "follower-1", 1, 20)

	require.NoError(t, err)
	assert.Len(t, sellers, 2)
	assert.Equal(t, int64(2), total)
	assert.Equal(t, "Store 1", sellers[0].StoreName)
}

func TestGetFollowerCount_Success(t *testing.T) {
	repo := &mockFollowRepo{
		countBySellerFn: func(_ context.Context, sellerID string) (int64, error) {
			assert.Equal(t, "seller-1", sellerID)
			return 42, nil
		},
	}

	uc := NewFollowUseCase(repo, zerolog.Nop())
	count, err := uc.GetFollowerCount(context.Background(), "seller-1")

	require.NoError(t, err)
	assert.Equal(t, int64(42), count)
}
