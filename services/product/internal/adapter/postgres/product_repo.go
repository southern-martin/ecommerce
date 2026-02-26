package postgres

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/southern-martin/ecommerce/services/product/internal/domain"
)

// ProductRepo implements domain.ProductRepository using GORM.
type ProductRepo struct {
	db *gorm.DB
}

// NewProductRepo creates a new ProductRepo.
func NewProductRepo(db *gorm.DB) *ProductRepo {
	return &ProductRepo{db: db}
}

func (r *ProductRepo) Create(ctx context.Context, p *domain.Product) error {
	model := ProductModelFromDomain(p)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *ProductRepo) GetByID(ctx context.Context, id string) (*domain.Product, error) {
	var model ProductModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}
	return model.ToDomain(), nil
}

func (r *ProductRepo) GetBySlug(ctx context.Context, slug string) (*domain.Product, error) {
	var model ProductModel
	if err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&model).Error; err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}
	return model.ToDomain(), nil
}

func (r *ProductRepo) List(ctx context.Context, filter domain.ProductFilter) ([]*domain.Product, int64, error) {
	query := r.db.WithContext(ctx).Model(&ProductModel{})

	if filter.SellerID != "" {
		query = query.Where("seller_id = ?", filter.SellerID)
	}
	if filter.CategoryID != "" {
		query = query.Where("category_id = ?", filter.CategoryID)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.Query != "" {
		like := "%" + filter.Query + "%"
		query = query.Where("name ILIKE ? OR description ILIKE ?", like, like)
	}
	if filter.MinPrice > 0 {
		query = query.Where("base_price_cents >= ?", filter.MinPrice)
	}
	if filter.MaxPrice > 0 {
		query = query.Where("base_price_cents <= ?", filter.MaxPrice)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Sorting
	switch filter.SortBy {
	case "price_asc":
		query = query.Order("base_price_cents ASC")
	case "price_desc":
		query = query.Order("base_price_cents DESC")
	case "name_asc":
		query = query.Order("name ASC")
	case "name_desc":
		query = query.Order("name DESC")
	case "newest":
		query = query.Order("created_at DESC")
	case "rating":
		query = query.Order("rating_avg DESC")
	default:
		query = query.Order("created_at DESC")
	}

	offset := (filter.Page - 1) * filter.PageSize
	var models []ProductModel
	if err := query.Offset(offset).Limit(filter.PageSize).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	products := make([]*domain.Product, len(models))
	for i, m := range models {
		products[i] = m.ToDomain()
	}

	return products, total, nil
}

func (r *ProductRepo) Update(ctx context.Context, p *domain.Product) error {
	model := ProductModelFromDomain(p)
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *ProductRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&ProductModel{}, "id = ?", id).Error
}
