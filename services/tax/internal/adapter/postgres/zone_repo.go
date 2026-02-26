package postgres

import (
	"context"

	"github.com/southern-martin/ecommerce/services/tax/internal/domain"
	"gorm.io/gorm"
)

type taxZoneRepository struct {
	db *gorm.DB
}

// NewTaxZoneRepository creates a new TaxZoneRepository backed by PostgreSQL.
func NewTaxZoneRepository(db *gorm.DB) domain.TaxZoneRepository {
	return &taxZoneRepository{db: db}
}

func (r *taxZoneRepository) Create(ctx context.Context, zone *domain.TaxZone) error {
	model := TaxZoneModelFromDomain(zone)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *taxZoneRepository) GetByID(ctx context.Context, id string) (*domain.TaxZone, error) {
	var model TaxZoneModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *taxZoneRepository) GetByLocation(ctx context.Context, countryCode, stateCode string) (*domain.TaxZone, error) {
	var model TaxZoneModel
	query := r.db.WithContext(ctx).Where("country_code = ?", countryCode)

	if stateCode != "" {
		// Try exact match with state first
		err := query.Where("state_code = ?", stateCode).First(&model).Error
		if err == nil {
			return model.ToDomain(), nil
		}
	}

	// Fall back to country-level zone (state_code is empty)
	err := r.db.WithContext(ctx).
		Where("country_code = ? AND (state_code = '' OR state_code IS NULL)", countryCode).
		First(&model).Error
	if err != nil {
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *taxZoneRepository) List(ctx context.Context) ([]*domain.TaxZone, error) {
	var models []TaxZoneModel
	if err := r.db.WithContext(ctx).Order("country_code, state_code").Find(&models).Error; err != nil {
		return nil, err
	}

	zones := make([]*domain.TaxZone, len(models))
	for i := range models {
		zones[i] = models[i].ToDomain()
	}
	return zones, nil
}
