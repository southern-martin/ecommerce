package postgres

import (
	"context"
	"errors"

	apperrors "github.com/southern-martin/ecommerce/pkg/errors"
	"github.com/southern-martin/ecommerce/services/user/internal/domain"
	"gorm.io/gorm"
)

// AddressRepository implements domain.AddressRepository using GORM.
type AddressRepository struct {
	db *gorm.DB
}

// NewAddressRepository creates a new AddressRepository.
func NewAddressRepository(db *gorm.DB) *AddressRepository {
	return &AddressRepository{db: db}
}

// Create inserts a new address.
func (r *AddressRepository) Create(ctx context.Context, addr *domain.Address) error {
	if err := r.db.WithContext(ctx).Create(addr).Error; err != nil {
		return apperrors.NewInternalError("DB_ERROR", "failed to create address")
	}
	return nil
}

// GetByID retrieves an address by its ID.
func (r *AddressRepository) GetByID(ctx context.Context, id string) (*domain.Address, error) {
	var addr domain.Address
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&addr).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError("ADDRESS_NOT_FOUND", "address not found")
		}
		return nil, apperrors.NewInternalError("DB_ERROR", "failed to get address")
	}
	return &addr, nil
}

// ListByUserID retrieves all addresses for a user.
func (r *AddressRepository) ListByUserID(ctx context.Context, userID string) ([]domain.Address, error) {
	var addresses []domain.Address
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("is_default DESC, created_at ASC").Find(&addresses).Error; err != nil {
		return nil, apperrors.NewInternalError("DB_ERROR", "failed to list addresses")
	}
	return addresses, nil
}

// Update persists changes to an existing address.
func (r *AddressRepository) Update(ctx context.Context, addr *domain.Address) error {
	if err := r.db.WithContext(ctx).Save(addr).Error; err != nil {
		return apperrors.NewInternalError("DB_ERROR", "failed to update address")
	}
	return nil
}

// Delete removes an address by its ID.
func (r *AddressRepository) Delete(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.Address{}).Error; err != nil {
		return apperrors.NewInternalError("DB_ERROR", "failed to delete address")
	}
	return nil
}

// CountByUserID counts the number of addresses for a user.
func (r *AddressRepository) CountByUserID(ctx context.Context, userID string) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&domain.Address{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return 0, apperrors.NewInternalError("DB_ERROR", "failed to count addresses")
	}
	return count, nil
}

// ClearDefaultByUserID clears the is_default flag for all addresses of a user.
func (r *AddressRepository) ClearDefaultByUserID(ctx context.Context, userID string) error {
	if err := r.db.WithContext(ctx).Model(&domain.Address{}).Where("user_id = ?", userID).Update("is_default", false).Error; err != nil {
		return apperrors.NewInternalError("DB_ERROR", "failed to clear default addresses")
	}
	return nil
}
