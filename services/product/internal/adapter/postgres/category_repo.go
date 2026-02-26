package postgres

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/southern-martin/ecommerce/services/product/internal/domain"
)

// CategoryRepo implements domain.CategoryRepository using GORM.
type CategoryRepo struct {
	db *gorm.DB
}

// NewCategoryRepo creates a new CategoryRepo.
func NewCategoryRepo(db *gorm.DB) *CategoryRepo {
	return &CategoryRepo{db: db}
}

func (r *CategoryRepo) Create(ctx context.Context, c *domain.Category) error {
	model := CategoryModelFromDomain(c)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *CategoryRepo) GetByID(ctx context.Context, id string) (*domain.Category, error) {
	var model CategoryModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		return nil, fmt.Errorf("category not found: %w", err)
	}
	return model.ToDomain(), nil
}

func (r *CategoryRepo) List(ctx context.Context) ([]*domain.Category, error) {
	var models []CategoryModel
	if err := r.db.WithContext(ctx).Where("is_active = ?", true).Order("sort_order ASC, name ASC").Find(&models).Error; err != nil {
		return nil, err
	}

	categories := make([]*domain.Category, len(models))
	for i, m := range models {
		categories[i] = m.ToDomain()
	}
	return categories, nil
}

func (r *CategoryRepo) Update(ctx context.Context, c *domain.Category) error {
	model := CategoryModelFromDomain(c)
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *CategoryRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&CategoryModel{}, "id = ?", id).Error
}
