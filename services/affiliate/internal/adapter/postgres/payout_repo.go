package postgres

import (
	"context"
	"time"

	"github.com/southern-martin/ecommerce/services/affiliate/internal/domain"
	"gorm.io/gorm"
)

// PayoutRepo implements domain.PayoutRepository.
type PayoutRepo struct {
	db *gorm.DB
}

// NewPayoutRepo creates a new PayoutRepo.
func NewPayoutRepo(db *gorm.DB) *PayoutRepo {
	return &PayoutRepo{db: db}
}

func (r *PayoutRepo) GetByID(ctx context.Context, id string) (*domain.AffiliatePayout, error) {
	var model AffiliatePayoutModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *PayoutRepo) ListByUser(ctx context.Context, userID string, page, pageSize int) ([]domain.AffiliatePayout, int64, error) {
	var total int64
	r.db.WithContext(ctx).Model(&AffiliatePayoutModel{}).Where("user_id = ?", userID).Count(&total)

	var models []AffiliatePayoutModel
	offset := (page - 1) * pageSize
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).
		Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	payouts := make([]domain.AffiliatePayout, len(models))
	for i, m := range models {
		payouts[i] = *m.ToDomain()
	}
	return payouts, total, nil
}

func (r *PayoutRepo) ListAll(ctx context.Context, page, pageSize int) ([]domain.AffiliatePayout, int64, error) {
	var total int64
	r.db.WithContext(ctx).Model(&AffiliatePayoutModel{}).Count(&total)

	var models []AffiliatePayoutModel
	offset := (page - 1) * pageSize
	if err := r.db.WithContext(ctx).Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	payouts := make([]domain.AffiliatePayout, len(models))
	for i, m := range models {
		payouts[i] = *m.ToDomain()
	}
	return payouts, total, nil
}

func (r *PayoutRepo) Create(ctx context.Context, payout *domain.AffiliatePayout) error {
	model := ToAffiliatePayoutModel(payout)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *PayoutRepo) UpdateStatus(ctx context.Context, id string, status domain.PayoutStatus, completedAt *time.Time) error {
	updates := map[string]interface{}{
		"status": string(status),
	}
	if completedAt != nil {
		updates["completed_at"] = completedAt
	}
	return r.db.WithContext(ctx).Model(&AffiliatePayoutModel{}).Where("id = ?", id).Updates(updates).Error
}
