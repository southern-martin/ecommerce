package postgres

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/southern-martin/ecommerce/services/product/internal/domain"
)

// AttributeGroupRepo implements domain.AttributeGroupRepository using GORM.
type AttributeGroupRepo struct {
	db *gorm.DB
}

// NewAttributeGroupRepo creates a new AttributeGroupRepo.
func NewAttributeGroupRepo(db *gorm.DB) *AttributeGroupRepo {
	return &AttributeGroupRepo{db: db}
}

func (r *AttributeGroupRepo) Create(ctx context.Context, group *domain.AttributeGroup) error {
	model := AttributeGroupModelFromDomain(group)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *AttributeGroupRepo) GetByID(ctx context.Context, id string) (*domain.AttributeGroup, error) {
	var model AttributeGroupModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		return nil, fmt.Errorf("attribute group not found: %w", err)
	}
	return model.ToDomain(), nil
}

func (r *AttributeGroupRepo) List(ctx context.Context) ([]*domain.AttributeGroup, error) {
	var models []AttributeGroupModel
	if err := r.db.WithContext(ctx).Order("sort_order ASC, name ASC").Find(&models).Error; err != nil {
		return nil, err
	}

	groups := make([]*domain.AttributeGroup, len(models))
	for i, m := range models {
		groups[i] = m.ToDomain()
	}
	return groups, nil
}

func (r *AttributeGroupRepo) Update(ctx context.Context, group *domain.AttributeGroup) error {
	model := AttributeGroupModelFromDomain(group)
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *AttributeGroupRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&AttributeGroupModel{}, "id = ?", id).Error
}

func (r *AttributeGroupRepo) AddAttribute(ctx context.Context, groupID, attributeID string, sortOrder int) error {
	model := AttributeGroupItemModel{
		GroupID:     groupID,
		AttributeID: attributeID,
		SortOrder:   sortOrder,
	}
	return r.db.WithContext(ctx).Create(&model).Error
}

func (r *AttributeGroupRepo) RemoveAttribute(ctx context.Context, groupID, attributeID string) error {
	return r.db.WithContext(ctx).
		Where("group_id = ? AND attribute_id = ?", groupID, attributeID).
		Delete(&AttributeGroupItemModel{}).Error
}

func (r *AttributeGroupRepo) ListAttributes(ctx context.Context, groupID string) ([]*domain.AttributeDefinition, error) {
	var models []AttributeDefinitionModel
	err := r.db.WithContext(ctx).
		Preload("OptionValues", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC")
		}).
		Joins("JOIN attribute_group_items agi ON agi.attribute_id = attribute_definitions.id").
		Where("agi.group_id = ?", groupID).
		Order("agi.sort_order ASC").
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
