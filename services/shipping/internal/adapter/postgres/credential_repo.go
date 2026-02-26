package postgres

import (
	"context"

	"github.com/southern-martin/ecommerce/services/shipping/internal/domain"
	"gorm.io/gorm"
)

// CredentialRepo implements domain.CarrierCredentialRepository.
type CredentialRepo struct {
	db *gorm.DB
}

// NewCredentialRepo creates a new CredentialRepo.
func NewCredentialRepo(db *gorm.DB) *CredentialRepo {
	return &CredentialRepo{db: db}
}

func (r *CredentialRepo) GetBySellerAndCarrier(ctx context.Context, sellerID, carrierCode string) (*domain.CarrierCredential, error) {
	var model CarrierCredentialModel
	if err := r.db.WithContext(ctx).Where("seller_id = ? AND carrier_code = ?", sellerID, carrierCode).First(&model).Error; err != nil {
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *CredentialRepo) ListBySeller(ctx context.Context, sellerID string) ([]domain.CarrierCredential, error) {
	var models []CarrierCredentialModel
	if err := r.db.WithContext(ctx).Where("seller_id = ?", sellerID).Find(&models).Error; err != nil {
		return nil, err
	}
	creds := make([]domain.CarrierCredential, len(models))
	for i, m := range models {
		creds[i] = *m.ToDomain()
	}
	return creds, nil
}

func (r *CredentialRepo) Create(ctx context.Context, cred *domain.CarrierCredential) error {
	model := ToCarrierCredentialModel(cred)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *CredentialRepo) Update(ctx context.Context, cred *domain.CarrierCredential) error {
	model := ToCarrierCredentialModel(cred)
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *CredentialRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&CarrierCredentialModel{}, "id = ?", id).Error
}
