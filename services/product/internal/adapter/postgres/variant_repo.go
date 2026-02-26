package postgres

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/southern-martin/ecommerce/services/product/internal/domain"
)

// VariantRepo implements domain.VariantRepository using GORM.
type VariantRepo struct {
	db *gorm.DB
}

// NewVariantRepo creates a new VariantRepo.
func NewVariantRepo(db *gorm.DB) *VariantRepo {
	return &VariantRepo{db: db}
}

func (r *VariantRepo) Create(ctx context.Context, v *domain.Variant) error {
	model := VariantModelFromDomain(v)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}

	// Save option values
	for _, ov := range v.OptionValues {
		ovModel := &VariantOptionValueModel{
			VariantID:     ov.VariantID,
			OptionID:      ov.OptionID,
			OptionValueID: ov.OptionValueID,
			OptionName:    ov.OptionName,
			Value:         ov.Value,
		}
		if err := r.db.WithContext(ctx).Create(ovModel).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *VariantRepo) GetByID(ctx context.Context, id string) (*domain.Variant, error) {
	var model VariantModel
	if err := r.db.WithContext(ctx).
		Preload("OptionValues").
		Where("id = ?", id).
		First(&model).Error; err != nil {
		return nil, fmt.Errorf("variant not found: %w", err)
	}
	return model.ToDomain(), nil
}

func (r *VariantRepo) GetBySKU(ctx context.Context, sku string) (*domain.Variant, error) {
	var model VariantModel
	if err := r.db.WithContext(ctx).
		Preload("OptionValues").
		Where("sku = ?", sku).
		First(&model).Error; err != nil {
		return nil, fmt.Errorf("variant not found: %w", err)
	}
	return model.ToDomain(), nil
}

func (r *VariantRepo) ListByProduct(ctx context.Context, productID string) ([]domain.Variant, error) {
	var models []VariantModel
	if err := r.db.WithContext(ctx).
		Preload("OptionValues").
		Where("product_id = ?", productID).
		Order("is_default DESC, created_at ASC").
		Find(&models).Error; err != nil {
		return nil, err
	}

	variants := make([]domain.Variant, len(models))
	for i, m := range models {
		variants[i] = *m.ToDomain()
	}
	return variants, nil
}

func (r *VariantRepo) Update(ctx context.Context, v *domain.Variant) error {
	model := VariantModelFromDomain(v)
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *VariantRepo) Delete(ctx context.Context, id string) error {
	// Delete option values first
	if err := r.db.WithContext(ctx).Where("variant_id = ?", id).Delete(&VariantOptionValueModel{}).Error; err != nil {
		return err
	}
	return r.db.WithContext(ctx).Delete(&VariantModel{}, "id = ?", id).Error
}

func (r *VariantRepo) BulkCreate(ctx context.Context, variants []domain.Variant) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, v := range variants {
			model := VariantModelFromDomain(&v)
			if err := tx.Create(model).Error; err != nil {
				return err
			}

			for _, ov := range v.OptionValues {
				ovModel := &VariantOptionValueModel{
					VariantID:     ov.VariantID,
					OptionID:      ov.OptionID,
					OptionValueID: ov.OptionValueID,
					OptionName:    ov.OptionName,
					Value:         ov.Value,
				}
				if err := tx.Create(ovModel).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (r *VariantRepo) UpdateStock(ctx context.Context, variantID string, delta int) error {
	result := r.db.WithContext(ctx).
		Model(&VariantModel{}).
		Where("id = ? AND stock + ? >= 0", variantID, delta).
		Update("stock", gorm.Expr("stock + ?", delta))

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("insufficient stock or variant not found")
	}
	return nil
}

func (r *VariantRepo) SetOptionValues(ctx context.Context, variantID string, values []domain.VariantOptionValue) error {
	// Delete existing
	if err := r.db.WithContext(ctx).Where("variant_id = ?", variantID).Delete(&VariantOptionValueModel{}).Error; err != nil {
		return err
	}

	for _, ov := range values {
		model := &VariantOptionValueModel{
			VariantID:     variantID,
			OptionID:      ov.OptionID,
			OptionValueID: ov.OptionValueID,
			OptionName:    ov.OptionName,
			Value:         ov.Value,
		}
		if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
			return err
		}
	}
	return nil
}
