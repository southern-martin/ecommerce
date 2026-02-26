package postgres

import (
	"context"
	"errors"

	apperrors "github.com/southern-martin/ecommerce/pkg/errors"
	"github.com/southern-martin/ecommerce/services/user/internal/domain"
	"gorm.io/gorm"
)

// FollowRepository implements domain.FollowRepository using GORM.
type FollowRepository struct {
	db *gorm.DB
}

// NewFollowRepository creates a new FollowRepository.
func NewFollowRepository(db *gorm.DB) *FollowRepository {
	return &FollowRepository{db: db}
}

// Create inserts a new follow relationship.
func (r *FollowRepository) Create(ctx context.Context, follow *domain.UserFollow) error {
	if err := r.db.WithContext(ctx).Create(follow).Error; err != nil {
		return apperrors.NewInternalError("DB_ERROR", "failed to create follow")
	}
	return nil
}

// Delete removes a follow relationship by follower and seller IDs.
func (r *FollowRepository) Delete(ctx context.Context, followerID, sellerID string) error {
	result := r.db.WithContext(ctx).Where("follower_id = ? AND seller_id = ?", followerID, sellerID).Delete(&domain.UserFollow{})
	if result.Error != nil {
		return apperrors.NewInternalError("DB_ERROR", "failed to delete follow")
	}
	if result.RowsAffected == 0 {
		return apperrors.NewNotFoundError("FOLLOW_NOT_FOUND", "follow relationship not found")
	}
	return nil
}

// ListByFollowerID retrieves a paginated list of sellers followed by a user.
func (r *FollowRepository) ListByFollowerID(ctx context.Context, followerID string, page, size int) ([]domain.SellerProfile, int64, error) {
	var sellers []domain.SellerProfile
	var total int64

	subQuery := r.db.WithContext(ctx).Model(&domain.UserFollow{}).Select("seller_id").Where("follower_id = ?", followerID)

	if err := r.db.WithContext(ctx).Model(&domain.SellerProfile{}).Where("id IN (?)", subQuery).Count(&total).Error; err != nil {
		return nil, 0, apperrors.NewInternalError("DB_ERROR", "failed to count followed sellers")
	}

	offset := (page - 1) * size
	if err := r.db.WithContext(ctx).Where("id IN (?)", subQuery).Offset(offset).Limit(size).Find(&sellers).Error; err != nil {
		return nil, 0, apperrors.NewInternalError("DB_ERROR", "failed to list followed sellers")
	}

	return sellers, total, nil
}

// CountBySellerID counts the number of followers for a seller.
func (r *FollowRepository) CountBySellerID(ctx context.Context, sellerID string) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&domain.UserFollow{}).Where("seller_id = ?", sellerID).Count(&count).Error; err != nil {
		return 0, apperrors.NewInternalError("DB_ERROR", "failed to count followers")
	}
	return count, nil
}

// Exists checks whether a follow relationship exists.
func (r *FollowRepository) Exists(ctx context.Context, followerID, sellerID string) (bool, error) {
	var follow domain.UserFollow
	err := r.db.WithContext(ctx).Where("follower_id = ? AND seller_id = ?", followerID, sellerID).First(&follow).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, apperrors.NewInternalError("DB_ERROR", "failed to check follow existence")
	}
	return true, nil
}
