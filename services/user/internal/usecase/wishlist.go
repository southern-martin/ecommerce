package usecase

import (
	"context"

	"github.com/rs/zerolog"

	apperrors "github.com/southern-martin/ecommerce/pkg/errors"
	"github.com/southern-martin/ecommerce/services/user/internal/domain"
)

type WishlistUseCase struct {
	repo   domain.WishlistRepository
	logger zerolog.Logger
}

func NewWishlistUseCase(repo domain.WishlistRepository, logger zerolog.Logger) *WishlistUseCase {
	return &WishlistUseCase{
		repo:   repo,
		logger: logger,
	}
}

func (uc *WishlistUseCase) AddItem(ctx context.Context, userID, productID string) error {
	exists, err := uc.repo.Exists(ctx, userID, productID)
	if err != nil {
		return err
	}
	if exists {
		return apperrors.NewConflictError("ALREADY_IN_WISHLIST", "product already in wishlist")
	}

	item := &domain.WishlistItem{
		UserID:    userID,
		ProductID: productID,
	}
	return uc.repo.Create(ctx, item)
}

func (uc *WishlistUseCase) RemoveItem(ctx context.Context, userID, productID string) error {
	return uc.repo.Delete(ctx, userID, productID)
}

func (uc *WishlistUseCase) ListItems(ctx context.Context, userID string) ([]domain.WishlistItem, error) {
	return uc.repo.ListByUserID(ctx, userID)
}
