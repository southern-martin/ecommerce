package postgres

import (
	"context"

	"gorm.io/gorm"

	"github.com/southern-martin/ecommerce/services/product/internal/domain"
)

// OptionRepo implements domain.OptionRepository using GORM.
type OptionRepo struct {
	db *gorm.DB
}

// NewOptionRepo creates a new OptionRepo.
func NewOptionRepo(db *gorm.DB) *OptionRepo {
	return &OptionRepo{db: db}
}

func (r *OptionRepo) CreateOption(ctx context.Context, option *domain.ProductOption) error {
	model := &ProductOptionModel{
		ID:        option.ID,
		ProductID: option.ProductID,
		Name:      option.Name,
		SortOrder: option.SortOrder,
	}

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}

	// Create option values
	for _, v := range option.Values {
		valModel := &ProductOptionValueModel{
			ID:        v.ID,
			OptionID:  option.ID,
			Value:     v.Value,
			ColorHex:  v.ColorHex,
			SortOrder: v.SortOrder,
		}
		if err := r.db.WithContext(ctx).Create(valModel).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *OptionRepo) UpdateOption(ctx context.Context, option *domain.ProductOption) error {
	model := &ProductOptionModel{
		ID:        option.ID,
		ProductID: option.ProductID,
		Name:      option.Name,
		SortOrder: option.SortOrder,
	}
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *OptionRepo) DeleteOption(ctx context.Context, optionID string) error {
	// Delete option values first
	if err := r.db.WithContext(ctx).Where("option_id = ?", optionID).Delete(&ProductOptionValueModel{}).Error; err != nil {
		return err
	}
	return r.db.WithContext(ctx).Delete(&ProductOptionModel{}, "id = ?", optionID).Error
}

func (r *OptionRepo) ListByProduct(ctx context.Context, productID string) ([]domain.ProductOption, error) {
	var models []ProductOptionModel
	if err := r.db.WithContext(ctx).
		Preload("Values", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC")
		}).
		Where("product_id = ?", productID).
		Order("sort_order ASC").
		Find(&models).Error; err != nil {
		return nil, err
	}

	options := make([]domain.ProductOption, len(models))
	for i, m := range models {
		options[i] = m.ToDomain()
	}
	return options, nil
}

func (r *OptionRepo) CreateOptionValue(ctx context.Context, value *domain.ProductOptionValue) error {
	model := &ProductOptionValueModel{
		ID:        value.ID,
		OptionID:  value.OptionID,
		Value:     value.Value,
		ColorHex:  value.ColorHex,
		SortOrder: value.SortOrder,
	}
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *OptionRepo) UpdateOptionValue(ctx context.Context, value *domain.ProductOptionValue) error {
	model := &ProductOptionValueModel{
		ID:        value.ID,
		OptionID:  value.OptionID,
		Value:     value.Value,
		ColorHex:  value.ColorHex,
		SortOrder: value.SortOrder,
	}
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *OptionRepo) DeleteOptionValue(ctx context.Context, valueID string) error {
	return r.db.WithContext(ctx).Delete(&ProductOptionValueModel{}, "id = ?", valueID).Error
}
