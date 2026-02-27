package postgres

import (
	"context"
	"errors"

	apperrors "github.com/southern-martin/ecommerce/pkg/errors"
	"github.com/southern-martin/ecommerce/services/user/internal/domain"
	"gorm.io/gorm"
)

type WishlistRepository struct {
	db *gorm.DB
}

func NewWishlistRepository(db *gorm.DB) *WishlistRepository {
	return &WishlistRepository{db: db}
}

func (r *WishlistRepository) Create(ctx context.Context, item *domain.WishlistItem) error {
	if err := r.db.WithContext(ctx).Create(item).Error; err != nil {
		return apperrors.NewInternalError("DB_ERROR", "failed to add item to wishlist")
	}
	return nil
}

func (r *WishlistRepository) Delete(ctx context.Context, userID, productID string) error {
	result := r.db.WithContext(ctx).Where("user_id = ? AND product_id = ?", userID, productID).Delete(&domain.WishlistItem{})
	if result.Error != nil {
		return apperrors.NewInternalError("DB_ERROR", "failed to remove item from wishlist")
	}
	if result.RowsAffected == 0 {
		return apperrors.NewNotFoundError("WISHLIST_ITEM_NOT_FOUND", "item not found in wishlist")
	}
	return nil
}

func (r *WishlistRepository) ListByUserID(ctx context.Context, userID string) ([]domain.WishlistItem, error) {
	var items []domain.WishlistItem
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Find(&items).Error; err != nil {
		return nil, apperrors.NewInternalError("DB_ERROR", "failed to list wishlist items")
	}
	return items, nil
}

func (r *WishlistRepository) Exists(ctx context.Context, userID, productID string) (bool, error) {
	var item domain.WishlistItem
	err := r.db.WithContext(ctx).Where("user_id = ? AND product_id = ?", userID, productID).First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, apperrors.NewInternalError("DB_ERROR", "failed to check wishlist item existence")
	}
	return true, nil
}
