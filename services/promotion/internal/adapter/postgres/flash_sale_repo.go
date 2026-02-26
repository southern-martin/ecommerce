package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/southern-martin/ecommerce/services/promotion/internal/domain"
	"gorm.io/gorm"
)

// FlashSaleRepo implements domain.FlashSaleRepository using GORM/Postgres.
type FlashSaleRepo struct {
	db *gorm.DB
}

// NewFlashSaleRepo creates a new FlashSaleRepo.
func NewFlashSaleRepo(db *gorm.DB) *FlashSaleRepo {
	return &FlashSaleRepo{db: db}
}

// GetByID retrieves a flash sale by its UUID, including items.
func (r *FlashSaleRepo) GetByID(ctx context.Context, id string) (*domain.FlashSale, error) {
	var model FlashSaleModel
	err := r.db.WithContext(ctx).
		Preload("Items").
		Where("id = ?", id).
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("flash sale not found")
		}
		return nil, err
	}
	return model.ToDomain(), nil
}

// ListActive retrieves all currently active flash sales.
func (r *FlashSaleRepo) ListActive(ctx context.Context) ([]*domain.FlashSale, error) {
	var models []FlashSaleModel
	now := time.Now()
	err := r.db.WithContext(ctx).
		Preload("Items").
		Where("is_active = ? AND starts_at <= ? AND ends_at >= ?", true, now, now).
		Order("starts_at ASC").
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	var flashSales []*domain.FlashSale
	for i := range models {
		flashSales = append(flashSales, models[i].ToDomain())
	}
	return flashSales, nil
}

// ListAll retrieves a paginated list of all flash sales.
func (r *FlashSaleRepo) ListAll(ctx context.Context, page, pageSize int) ([]*domain.FlashSale, int64, error) {
	var models []FlashSaleModel
	var total int64

	query := r.db.WithContext(ctx).Model(&FlashSaleModel{})

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.
		Preload("Items").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&models).Error
	if err != nil {
		return nil, 0, err
	}

	var flashSales []*domain.FlashSale
	for i := range models {
		flashSales = append(flashSales, models[i].ToDomain())
	}
	return flashSales, total, nil
}

// Create persists a new flash sale.
func (r *FlashSaleRepo) Create(ctx context.Context, flashSale *domain.FlashSale) error {
	model := ToFlashSaleModel(flashSale)
	return r.db.WithContext(ctx).Create(model).Error
}

// Update persists all changes to an existing flash sale.
func (r *FlashSaleRepo) Update(ctx context.Context, flashSale *domain.FlashSale) error {
	model := ToFlashSaleModel(flashSale)
	return r.db.WithContext(ctx).Save(model).Error
}
