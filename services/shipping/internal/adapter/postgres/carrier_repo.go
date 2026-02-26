package postgres

import (
	"context"

	"github.com/southern-martin/ecommerce/services/shipping/internal/domain"
	"gorm.io/gorm"
)

// CarrierRepo implements domain.CarrierRepository.
type CarrierRepo struct {
	db *gorm.DB
}

// NewCarrierRepo creates a new CarrierRepo.
func NewCarrierRepo(db *gorm.DB) *CarrierRepo {
	return &CarrierRepo{db: db}
}

func (r *CarrierRepo) GetAll(ctx context.Context) ([]domain.Carrier, error) {
	var models []CarrierModel
	if err := r.db.WithContext(ctx).Where("is_active = ?", true).Find(&models).Error; err != nil {
		return nil, err
	}
	carriers := make([]domain.Carrier, len(models))
	for i, m := range models {
		carriers[i] = *m.ToDomain()
	}
	return carriers, nil
}

func (r *CarrierRepo) GetByCode(ctx context.Context, code string) (*domain.Carrier, error) {
	var model CarrierModel
	if err := r.db.WithContext(ctx).Where("code = ?", code).First(&model).Error; err != nil {
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *CarrierRepo) Create(ctx context.Context, carrier *domain.Carrier) error {
	model := ToCarrierModel(carrier)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *CarrierRepo) Update(ctx context.Context, carrier *domain.Carrier) error {
	return r.db.WithContext(ctx).Model(&CarrierModel{}).Where("code = ?", carrier.Code).Updates(map[string]interface{}{
		"name":                carrier.Name,
		"is_active":           carrier.IsActive,
		"supported_countries": carrier.SupportedCountries,
		"api_base_url":        carrier.APIBaseURL,
	}).Error
}
