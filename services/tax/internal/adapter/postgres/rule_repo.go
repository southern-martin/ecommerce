package postgres

import (
	"context"

	"github.com/southern-martin/ecommerce/services/tax/internal/domain"
	"gorm.io/gorm"
)

type taxRuleRepository struct {
	db *gorm.DB
}

// NewTaxRuleRepository creates a new TaxRuleRepository backed by PostgreSQL.
func NewTaxRuleRepository(db *gorm.DB) domain.TaxRuleRepository {
	return &taxRuleRepository{db: db}
}

func (r *taxRuleRepository) Create(ctx context.Context, rule *domain.TaxRule) error {
	model := TaxRuleModelFromDomain(rule)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *taxRuleRepository) GetByID(ctx context.Context, id string) (*domain.TaxRule, error) {
	var model TaxRuleModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *taxRuleRepository) ListByZone(ctx context.Context, zoneID string) ([]*domain.TaxRule, error) {
	var models []TaxRuleModel
	if err := r.db.WithContext(ctx).Where("zone_id = ?", zoneID).Order("tax_name").Find(&models).Error; err != nil {
		return nil, err
	}

	rules := make([]*domain.TaxRule, len(models))
	for i := range models {
		rules[i] = models[i].ToDomain()
	}
	return rules, nil
}

func (r *taxRuleRepository) ListActive(ctx context.Context) ([]*domain.TaxRule, error) {
	var models []TaxRuleModel
	if err := r.db.WithContext(ctx).
		Where("is_active = ? AND starts_at <= NOW() AND (expires_at IS NULL OR expires_at > NOW())", true).
		Order("zone_id, tax_name").
		Find(&models).Error; err != nil {
		return nil, err
	}

	rules := make([]*domain.TaxRule, len(models))
	for i := range models {
		rules[i] = models[i].ToDomain()
	}
	return rules, nil
}

func (r *taxRuleRepository) Update(ctx context.Context, rule *domain.TaxRule) error {
	model := TaxRuleModelFromDomain(rule)
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *taxRuleRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&TaxRuleModel{}).Error
}

func (r *taxRuleRepository) GetByZoneAndCategory(ctx context.Context, zoneID, category string) ([]*domain.TaxRule, error) {
	var models []TaxRuleModel

	query := r.db.WithContext(ctx).
		Where("zone_id = ? AND is_active = ? AND starts_at <= NOW() AND (expires_at IS NULL OR expires_at > NOW())", zoneID, true)

	if category != "" {
		// Get rules matching the specific category OR rules with empty category (applies to all)
		query = query.Where("category = ? OR category = ''", category)
	} else {
		// Get rules with empty category (applies to all)
		query = query.Where("category = ''")
	}

	if err := query.Order("tax_name").Find(&models).Error; err != nil {
		return nil, err
	}

	rules := make([]*domain.TaxRule, len(models))
	for i := range models {
		rules[i] = models[i].ToDomain()
	}
	return rules, nil
}
