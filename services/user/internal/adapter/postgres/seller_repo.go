package postgres

import (
	"context"
	"errors"

	apperrors "github.com/southern-martin/ecommerce/pkg/errors"
	"github.com/southern-martin/ecommerce/services/user/internal/domain"
	"gorm.io/gorm"
)

// SellerRepository implements domain.SellerProfileRepository using GORM.
type SellerRepository struct {
	db *gorm.DB
}

// NewSellerRepository creates a new SellerRepository.
func NewSellerRepository(db *gorm.DB) *SellerRepository {
	return &SellerRepository{db: db}
}

// Create inserts a new seller profile.
func (r *SellerRepository) Create(ctx context.Context, seller *domain.SellerProfile) error {
	if err := r.db.WithContext(ctx).Create(seller).Error; err != nil {
		return apperrors.NewInternalError("DB_ERROR", "failed to create seller profile")
	}
	return nil
}

// GetByID retrieves a seller profile by its ID.
func (r *SellerRepository) GetByID(ctx context.Context, id string) (*domain.SellerProfile, error) {
	var seller domain.SellerProfile
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&seller).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError("SELLER_NOT_FOUND", "seller profile not found")
		}
		return nil, apperrors.NewInternalError("DB_ERROR", "failed to get seller profile")
	}
	return &seller, nil
}

// GetByUserID retrieves a seller profile by user ID.
func (r *SellerRepository) GetByUserID(ctx context.Context, userID string) (*domain.SellerProfile, error) {
	var seller domain.SellerProfile
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&seller).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError("SELLER_NOT_FOUND", "seller profile not found")
		}
		return nil, apperrors.NewInternalError("DB_ERROR", "failed to get seller profile")
	}
	return &seller, nil
}

// Update persists changes to an existing seller profile.
func (r *SellerRepository) Update(ctx context.Context, seller *domain.SellerProfile) error {
	if err := r.db.WithContext(ctx).Save(seller).Error; err != nil {
		return apperrors.NewInternalError("DB_ERROR", "failed to update seller profile")
	}
	return nil
}

// List retrieves a paginated list of seller profiles.
func (r *SellerRepository) List(ctx context.Context, page, size int) ([]domain.SellerProfile, int64, error) {
	var sellers []domain.SellerProfile
	var total int64

	if err := r.db.WithContext(ctx).Model(&domain.SellerProfile{}).Count(&total).Error; err != nil {
		return nil, 0, apperrors.NewInternalError("DB_ERROR", "failed to count seller profiles")
	}

	offset := (page - 1) * size
	if err := r.db.WithContext(ctx).Offset(offset).Limit(size).Order("created_at DESC").Find(&sellers).Error; err != nil {
		return nil, 0, apperrors.NewInternalError("DB_ERROR", "failed to list seller profiles")
	}

	return sellers, total, nil
}
