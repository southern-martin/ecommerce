package postgres

import (
	"context"
	"errors"

	"github.com/southern-martin/ecommerce/services/promotion/internal/domain"
	"gorm.io/gorm"
)

// BundleRepo implements domain.BundleRepository using GORM/Postgres.
type BundleRepo struct {
	db *gorm.DB
}

// NewBundleRepo creates a new BundleRepo.
func NewBundleRepo(db *gorm.DB) *BundleRepo {
	return &BundleRepo{db: db}
}

// GetByID retrieves a bundle by its UUID.
func (r *BundleRepo) GetByID(ctx context.Context, id string) (*domain.Bundle, error) {
	var model BundleModel
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("bundle not found")
		}
		return nil, err
	}
	return model.ToDomain(), nil
}

// ListBySeller retrieves a paginated list of bundles by seller.
func (r *BundleRepo) ListBySeller(ctx context.Context, sellerID string, page, pageSize int) ([]*domain.Bundle, int64, error) {
	var models []BundleModel
	var total int64

	query := r.db.WithContext(ctx).Model(&BundleModel{}).Where("seller_id = ?", sellerID)

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

	var bundles []*domain.Bundle
	for i := range models {
		bundles = append(bundles, models[i].ToDomain())
	}
	return bundles, total, nil
}

// ListActive retrieves a paginated list of active bundles.
func (r *BundleRepo) ListActive(ctx context.Context, page, pageSize int) ([]*domain.Bundle, int64, error) {
	var models []BundleModel
	var total int64

	query := r.db.WithContext(ctx).Model(&BundleModel{}).Where("is_active = ?", true)

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

	var bundles []*domain.Bundle
	for i := range models {
		bundles = append(bundles, models[i].ToDomain())
	}
	return bundles, total, nil
}

// Create persists a new bundle.
func (r *BundleRepo) Create(ctx context.Context, bundle *domain.Bundle) error {
	model := ToBundleModel(bundle)
	return r.db.WithContext(ctx).Create(model).Error
}

// Update persists all changes to an existing bundle.
func (r *BundleRepo) Update(ctx context.Context, bundle *domain.Bundle) error {
	model := ToBundleModel(bundle)
	return r.db.WithContext(ctx).Save(model).Error
}
