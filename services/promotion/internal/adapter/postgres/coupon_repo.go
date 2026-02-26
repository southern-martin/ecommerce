package postgres

import (
	"context"
	"errors"

	"github.com/southern-martin/ecommerce/services/promotion/internal/domain"
	"gorm.io/gorm"
)

// CouponRepo implements domain.CouponRepository using GORM/Postgres.
type CouponRepo struct {
	db *gorm.DB
}

// NewCouponRepo creates a new CouponRepo.
func NewCouponRepo(db *gorm.DB) *CouponRepo {
	return &CouponRepo{db: db}
}

// GetByID retrieves a coupon by its UUID.
func (r *CouponRepo) GetByID(ctx context.Context, id string) (*domain.Coupon, error) {
	var model CouponModel
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("coupon not found")
		}
		return nil, err
	}
	return model.ToDomain(), nil
}

// GetByCode retrieves a coupon by its unique code.
func (r *CouponRepo) GetByCode(ctx context.Context, code string) (*domain.Coupon, error) {
	var model CouponModel
	err := r.db.WithContext(ctx).
		Where("code = ?", code).
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("coupon not found")
		}
		return nil, err
	}
	return model.ToDomain(), nil
}

// ListAll retrieves a paginated list of all coupons.
func (r *CouponRepo) ListAll(ctx context.Context, page, pageSize int) ([]*domain.Coupon, int64, error) {
	var models []CouponModel
	var total int64

	query := r.db.WithContext(ctx).Model(&CouponModel{})

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&models).Error
	if err != nil {
		return nil, 0, err
	}

	var coupons []*domain.Coupon
	for i := range models {
		coupons = append(coupons, models[i].ToDomain())
	}
	return coupons, total, nil
}

// ListBySeller retrieves a paginated list of coupons created by a seller.
func (r *CouponRepo) ListBySeller(ctx context.Context, sellerID string, page, pageSize int) ([]*domain.Coupon, int64, error) {
	var models []CouponModel
	var total int64

	query := r.db.WithContext(ctx).Model(&CouponModel{}).Where("created_by = ?", sellerID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&models).Error
	if err != nil {
		return nil, 0, err
	}

	var coupons []*domain.Coupon
	for i := range models {
		coupons = append(coupons, models[i].ToDomain())
	}
	return coupons, total, nil
}

// Create persists a new coupon.
func (r *CouponRepo) Create(ctx context.Context, coupon *domain.Coupon) error {
	model := ToCouponModel(coupon)
	return r.db.WithContext(ctx).Create(model).Error
}

// Update persists all changes to an existing coupon.
func (r *CouponRepo) Update(ctx context.Context, coupon *domain.Coupon) error {
	model := ToCouponModel(coupon)
	return r.db.WithContext(ctx).Save(model).Error
}

// IncrementUsageCount atomically increments the usage count of a coupon.
func (r *CouponRepo) IncrementUsageCount(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).
		Model(&CouponModel{}).
		Where("id = ?", id).
		UpdateColumn("usage_count", gorm.Expr("usage_count + 1"))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("coupon not found")
	}
	return nil
}
