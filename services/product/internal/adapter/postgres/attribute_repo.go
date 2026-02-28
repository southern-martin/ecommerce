package postgres

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/southern-martin/ecommerce/services/product/internal/domain"
)

// AttributeRepo implements domain.AttributeRepository using GORM.
type AttributeRepo struct {
	db *gorm.DB
}

// NewAttributeRepo creates a new AttributeRepo.
func NewAttributeRepo(db *gorm.DB) *AttributeRepo {
	return &AttributeRepo{db: db}
}

func (r *AttributeRepo) CreateDefinition(ctx context.Context, attr *domain.AttributeDefinition) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		model := AttributeDefinitionModelFromDomain(attr)
		if err := tx.Create(model).Error; err != nil {
			return err
		}
		// Create option values
		for _, ov := range attr.OptionValues {
			ovModel := AttributeOptionValueModelFromDomain(ov)
			if err := tx.Create(ovModel).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *AttributeRepo) GetDefinitionByID(ctx context.Context, id string) (*domain.AttributeDefinition, error) {
	var model AttributeDefinitionModel
	if err := r.db.WithContext(ctx).
		Preload("OptionValues", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC")
		}).
		Where("id = ?", id).
		First(&model).Error; err != nil {
		return nil, fmt.Errorf("attribute definition not found: %w", err)
	}
	return model.ToDomain(), nil
}

func (r *AttributeRepo) ListDefinitions(ctx context.Context) ([]*domain.AttributeDefinition, error) {
	var models []AttributeDefinitionModel
	if err := r.db.WithContext(ctx).
		Preload("OptionValues", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC")
		}).
		Order("sort_order ASC, name ASC").
		Find(&models).Error; err != nil {
		return nil, err
	}

	attrs := make([]*domain.AttributeDefinition, len(models))
	for i, m := range models {
		attrs[i] = m.ToDomain()
	}
	return attrs, nil
}

func (r *AttributeRepo) UpdateDefinition(ctx context.Context, attr *domain.AttributeDefinition) error {
	model := AttributeDefinitionModelFromDomain(attr)
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *AttributeRepo) DeleteDefinition(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&AttributeDefinitionModel{}, "id = ?", id).Error
}

func (r *AttributeRepo) AssignToCategory(ctx context.Context, categoryID, attributeID string, sortOrder int) error {
	model := CategoryAttributeModel{
		CategoryID:  categoryID,
		AttributeID: attributeID,
		SortOrder:   sortOrder,
	}
	return r.db.WithContext(ctx).Create(&model).Error
}

func (r *AttributeRepo) RemoveFromCategory(ctx context.Context, categoryID, attributeID string) error {
	return r.db.WithContext(ctx).
		Where("category_id = ? AND attribute_id = ?", categoryID, attributeID).
		Delete(&CategoryAttributeModel{}).Error
}

func (r *AttributeRepo) ListByCategory(ctx context.Context, categoryID string) ([]*domain.AttributeDefinition, error) {
	var models []AttributeDefinitionModel
	err := r.db.WithContext(ctx).
		Preload("OptionValues", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC")
		}).
		Joins("JOIN category_attributes ca ON ca.attribute_id = attribute_definitions.id").
		Where("ca.category_id = ?", categoryID).
		Order("ca.sort_order ASC").
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	attrs := make([]*domain.AttributeDefinition, len(models))
	for i, m := range models {
		attrs[i] = m.ToDomain()
	}
	return attrs, nil
}

func (r *AttributeRepo) SetProductValues(ctx context.Context, productID string, values []domain.ProductAttributeValue) error {
	// Delete existing values
	if err := r.db.WithContext(ctx).Where("product_id = ?", productID).Delete(&ProductAttributeValueModel{}).Error; err != nil {
		return err
	}

	if len(values) == 0 {
		return nil
	}

	// Insert new values
	var models []ProductAttributeValueModel
	for _, v := range values {
		models = append(models, *ProductAttributeValueModelFromDomain(v))
	}
	return r.db.WithContext(ctx).Create(&models).Error
}

func (r *AttributeRepo) GetProductValues(ctx context.Context, productID string) ([]domain.ProductAttributeValue, error) {
	var models []ProductAttributeValueModel
	if err := r.db.WithContext(ctx).Where("product_id = ?", productID).Find(&models).Error; err != nil {
		return nil, err
	}

	values := make([]domain.ProductAttributeValue, len(models))
	for i, m := range models {
		values[i] = m.ToDomain()
	}
	return values, nil
}

// SetOptionValues replaces all option values for an attribute definition.
func (r *AttributeRepo) SetOptionValues(ctx context.Context, attributeID string, opts []domain.AttributeOptionValue) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Delete existing
		if err := tx.Where("attribute_id = ?", attributeID).Delete(&AttributeOptionValueModel{}).Error; err != nil {
			return err
		}
		// Insert new
		for _, ov := range opts {
			model := AttributeOptionValueModelFromDomain(ov)
			model.AttributeID = attributeID
			if err := tx.Create(model).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// ListOptionValues returns all option values for an attribute definition.
func (r *AttributeRepo) ListOptionValues(ctx context.Context, attributeID string) ([]domain.AttributeOptionValue, error) {
	var models []AttributeOptionValueModel
	if err := r.db.WithContext(ctx).
		Where("attribute_id = ?", attributeID).
		Order("sort_order ASC").
		Find(&models).Error; err != nil {
		return nil, err
	}

	values := make([]domain.AttributeOptionValue, len(models))
	for i, m := range models {
		values[i] = m.ToDomain()
	}
	return values, nil
}
