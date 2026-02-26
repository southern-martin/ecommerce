package postgres

import (
	"context"

	"github.com/southern-martin/ecommerce/services/return/internal/domain"
	"gorm.io/gorm"
)

// DisputeRepo implements domain.DisputeRepository.
type DisputeRepo struct {
	db *gorm.DB
}

// NewDisputeRepo creates a new DisputeRepo.
func NewDisputeRepo(db *gorm.DB) *DisputeRepo {
	return &DisputeRepo{db: db}
}

func (r *DisputeRepo) GetByID(ctx context.Context, id string) (*domain.Dispute, error) {
	var model DisputeModel
	if err := r.db.WithContext(ctx).Preload("Messages").Where("id = ?", id).First(&model).Error; err != nil {
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *DisputeRepo) GetByOrderID(ctx context.Context, orderID string) ([]domain.Dispute, error) {
	var models []DisputeModel
	if err := r.db.WithContext(ctx).Preload("Messages").Where("order_id = ?", orderID).Find(&models).Error; err != nil {
		return nil, err
	}
	disputes := make([]domain.Dispute, len(models))
	for i, m := range models {
		disputes[i] = *m.ToDomain()
	}
	return disputes, nil
}

func (r *DisputeRepo) ListAll(ctx context.Context, page, pageSize int) ([]domain.Dispute, int64, error) {
	var total int64
	r.db.WithContext(ctx).Model(&DisputeModel{}).Count(&total)

	var models []DisputeModel
	offset := (page - 1) * pageSize
	if err := r.db.WithContext(ctx).Preload("Messages").
		Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	disputes := make([]domain.Dispute, len(models))
	for i, m := range models {
		disputes[i] = *m.ToDomain()
	}
	return disputes, total, nil
}

func (r *DisputeRepo) ListByBuyer(ctx context.Context, buyerID string, page, pageSize int) ([]domain.Dispute, int64, error) {
	var total int64
	r.db.WithContext(ctx).Model(&DisputeModel{}).Where("buyer_id = ?", buyerID).Count(&total)

	var models []DisputeModel
	offset := (page - 1) * pageSize
	if err := r.db.WithContext(ctx).Preload("Messages").Where("buyer_id = ?", buyerID).
		Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	disputes := make([]domain.Dispute, len(models))
	for i, m := range models {
		disputes[i] = *m.ToDomain()
	}
	return disputes, total, nil
}

func (r *DisputeRepo) Create(ctx context.Context, dispute *domain.Dispute) error {
	model := ToDisputeModel(dispute)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *DisputeRepo) Update(ctx context.Context, dispute *domain.Dispute) error {
	return r.db.WithContext(ctx).Model(&DisputeModel{}).Where("id = ?", dispute.ID).Updates(map[string]interface{}{
		"status":      string(dispute.Status),
		"resolution":  dispute.Resolution,
		"resolved_by": dispute.ResolvedBy,
		"resolved_at": dispute.ResolvedAt,
	}).Error
}
