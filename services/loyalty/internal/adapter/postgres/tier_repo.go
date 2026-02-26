package postgres

import (
	"context"

	"github.com/southern-martin/ecommerce/services/loyalty/internal/domain"
	"gorm.io/gorm"
)

// TierRepo implements domain.TierRepository.
type TierRepo struct {
	db *gorm.DB
}

// NewTierRepo creates a new TierRepo.
func NewTierRepo(db *gorm.DB) *TierRepo {
	return &TierRepo{db: db}
}

func (r *TierRepo) GetAll(ctx context.Context) ([]domain.Tier, error) {
	var models []TierModel
	if err := r.db.WithContext(ctx).Order("min_points ASC").Find(&models).Error; err != nil {
		return nil, err
	}

	tiers := make([]domain.Tier, len(models))
	for i, m := range models {
		tiers[i] = *m.ToDomain()
	}
	return tiers, nil
}

func (r *TierRepo) GetByName(ctx context.Context, name string) (*domain.Tier, error) {
	var model TierModel
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&model).Error; err != nil {
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *TierRepo) GetTierForPoints(ctx context.Context, lifetimePoints int64) (*domain.Tier, error) {
	var model TierModel
	if err := r.db.WithContext(ctx).Where("min_points <= ?", lifetimePoints).
		Order("min_points DESC").First(&model).Error; err != nil {
		return nil, err
	}
	return model.ToDomain(), nil
}
