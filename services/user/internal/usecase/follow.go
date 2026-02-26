package usecase

import (
	"context"

	"github.com/rs/zerolog"

	apperrors "github.com/southern-martin/ecommerce/pkg/errors"
	"github.com/southern-martin/ecommerce/services/user/internal/domain"
)

// FollowUseCase handles user follow business logic.
type FollowUseCase struct {
	repo   domain.FollowRepository
	logger zerolog.Logger
}

// NewFollowUseCase creates a new FollowUseCase.
func NewFollowUseCase(repo domain.FollowRepository, logger zerolog.Logger) *FollowUseCase {
	return &FollowUseCase{
		repo:   repo,
		logger: logger,
	}
}

// Follow creates a follow relationship between a user and a seller.
func (uc *FollowUseCase) Follow(ctx context.Context, followerID, sellerID string) error {
	exists, err := uc.repo.Exists(ctx, followerID, sellerID)
	if err != nil {
		return err
	}
	if exists {
		return apperrors.NewConflictError("ALREADY_FOLLOWING", "already following this seller")
	}

	follow := &domain.UserFollow{
		FollowerID: followerID,
		SellerID:   sellerID,
	}

	return uc.repo.Create(ctx, follow)
}

// Unfollow removes a follow relationship between a user and a seller.
func (uc *FollowUseCase) Unfollow(ctx context.Context, followerID, sellerID string) error {
	return uc.repo.Delete(ctx, followerID, sellerID)
}

// ListFollowed lists sellers that a user follows with pagination.
func (uc *FollowUseCase) ListFollowed(ctx context.Context, followerID string, page, size int) ([]domain.SellerProfile, int64, error) {
	return uc.repo.ListByFollowerID(ctx, followerID, page, size)
}

// GetFollowerCount returns the number of followers for a seller.
func (uc *FollowUseCase) GetFollowerCount(ctx context.Context, sellerID string) (int64, error) {
	return uc.repo.CountBySellerID(ctx, sellerID)
}
