package postgres

import (
	"context"
	"errors"

	"github.com/southern-martin/ecommerce/services/promotion/internal/domain"
	"gorm.io/gorm"
)

// FlashSaleItemRepo implements domain.FlashSaleItemRepository using GORM/Postgres.
type FlashSaleItemRepo struct {
	db *gorm.DB
}

// NewFlashSaleItemRepo creates a new FlashSaleItemRepo.
func NewFlashSaleItemRepo(db *gorm.DB) *FlashSaleItemRepo {
	return &FlashSaleItemRepo{db: db}
}

// GetByFlashSaleID retrieves all items for a given flash sale.
func (r *FlashSaleItemRepo) GetByFlashSaleID(ctx context.Context, flashSaleID string) ([]*domain.FlashSaleItem, error) {
	var models []FlashSaleItemModel
	err := r.db.WithContext(ctx).
		Where("flash_sale_id = ?", flashSaleID).
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	var items []*domain.FlashSaleItem
	for i := range models {
		items = append(items, models[i].ToDomain())
	}
	return items, nil
}

// Create persists a new flash sale item.
func (r *FlashSaleItemRepo) Create(ctx context.Context, item *domain.FlashSaleItem) error {
	model := ToFlashSaleItemModel(item)
	return r.db.WithContext(ctx).Create(model).Error
}

// IncrementSoldCount atomically increments the sold count of a flash sale item.
func (r *FlashSaleItemRepo) IncrementSoldCount(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).
		Model(&FlashSaleItemModel{}).
		Where("id = ?", id).
		UpdateColumn("sold_count", gorm.Expr("sold_count + 1"))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("flash sale item not found")
	}
	return nil
}
