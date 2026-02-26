package postgres

import (
	"context"

	"github.com/southern-martin/ecommerce/services/return/internal/domain"
	"gorm.io/gorm"
)

// ReturnRepo implements domain.ReturnRepository.
type ReturnRepo struct {
	db *gorm.DB
}

// NewReturnRepo creates a new ReturnRepo.
func NewReturnRepo(db *gorm.DB) *ReturnRepo {
	return &ReturnRepo{db: db}
}

func (r *ReturnRepo) GetByID(ctx context.Context, id string) (*domain.Return, error) {
	var model ReturnModel
	if err := r.db.WithContext(ctx).Preload("Items").Where("id = ?", id).First(&model).Error; err != nil {
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *ReturnRepo) GetByOrderID(ctx context.Context, orderID string) ([]domain.Return, error) {
	var models []ReturnModel
	if err := r.db.WithContext(ctx).Preload("Items").Where("order_id = ?", orderID).Find(&models).Error; err != nil {
		return nil, err
	}
	returns := make([]domain.Return, len(models))
	for i, m := range models {
		returns[i] = *m.ToDomain()
	}
	return returns, nil
}

func (r *ReturnRepo) ListByBuyer(ctx context.Context, buyerID string, page, pageSize int) ([]domain.Return, int64, error) {
	var total int64
	r.db.WithContext(ctx).Model(&ReturnModel{}).Where("buyer_id = ?", buyerID).Count(&total)

	var models []ReturnModel
	offset := (page - 1) * pageSize
	if err := r.db.WithContext(ctx).Preload("Items").Where("buyer_id = ?", buyerID).
		Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	returns := make([]domain.Return, len(models))
	for i, m := range models {
		returns[i] = *m.ToDomain()
	}
	return returns, total, nil
}

func (r *ReturnRepo) ListBySeller(ctx context.Context, sellerID string, page, pageSize int) ([]domain.Return, int64, error) {
	var total int64
	r.db.WithContext(ctx).Model(&ReturnModel{}).Where("seller_id = ?", sellerID).Count(&total)

	var models []ReturnModel
	offset := (page - 1) * pageSize
	if err := r.db.WithContext(ctx).Preload("Items").Where("seller_id = ?", sellerID).
		Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	returns := make([]domain.Return, len(models))
	for i, m := range models {
		returns[i] = *m.ToDomain()
	}
	return returns, total, nil
}

func (r *ReturnRepo) Create(ctx context.Context, ret *domain.Return) error {
	model := ToReturnModel(ret)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *ReturnRepo) Update(ctx context.Context, ret *domain.Return) error {
	return r.db.WithContext(ctx).Model(&ReturnModel{}).Where("id = ?", ret.ID).Updates(map[string]interface{}{
		"status":              string(ret.Status),
		"refund_amount_cents": ret.RefundAmountCents,
		"refund_method":       ret.RefundMethod,
		"return_tracking":     ret.ReturnTracking,
	}).Error
}
