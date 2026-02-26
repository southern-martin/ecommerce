package postgres

import (
	"context"

	"github.com/southern-martin/ecommerce/services/promotion/internal/domain"
	"gorm.io/gorm"
)

// CouponUsageRepo implements domain.CouponUsageRepository using GORM/Postgres.
type CouponUsageRepo struct {
	db *gorm.DB
}

// NewCouponUsageRepo creates a new CouponUsageRepo.
func NewCouponUsageRepo(db *gorm.DB) *CouponUsageRepo {
	return &CouponUsageRepo{db: db}
}

// GetByUserAndCoupon retrieves all usage records for a user and coupon combination.
func (r *CouponUsageRepo) GetByUserAndCoupon(ctx context.Context, userID, couponID string) ([]*domain.CouponUsage, error) {
	var models []CouponUsageModel
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND coupon_id = ?", userID, couponID).
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	var usages []*domain.CouponUsage
	for i := range models {
		usages = append(usages, models[i].ToDomain())
	}
	return usages, nil
}

// CountByUser counts the number of times a user has used a specific coupon.
func (r *CouponUsageRepo) CountByUser(ctx context.Context, userID, couponID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&CouponUsageModel{}).
		Where("user_id = ? AND coupon_id = ?", userID, couponID).
		Count(&count).Error
	return count, err
}

// Create persists a new coupon usage record.
func (r *CouponUsageRepo) Create(ctx context.Context, usage *domain.CouponUsage) error {
	model := ToCouponUsageModel(usage)
	return r.db.WithContext(ctx).Create(model).Error
}
