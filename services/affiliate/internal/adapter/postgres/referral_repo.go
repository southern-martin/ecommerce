package postgres

import (
	"context"

	"github.com/southern-martin/ecommerce/services/affiliate/internal/domain"
	"gorm.io/gorm"
)

// ReferralRepo implements domain.ReferralRepository.
type ReferralRepo struct {
	db *gorm.DB
}

// NewReferralRepo creates a new ReferralRepo.
func NewReferralRepo(db *gorm.DB) *ReferralRepo {
	return &ReferralRepo{db: db}
}

func (r *ReferralRepo) GetByID(ctx context.Context, id string) (*domain.Referral, error) {
	var model ReferralModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *ReferralRepo) ListByReferrer(ctx context.Context, referrerID string, page, pageSize int) ([]domain.Referral, int64, error) {
	var total int64
	r.db.WithContext(ctx).Model(&ReferralModel{}).Where("referrer_id = ?", referrerID).Count(&total)

	var models []ReferralModel
	offset := (page - 1) * pageSize
	if err := r.db.WithContext(ctx).Where("referrer_id = ?", referrerID).
		Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	referrals := make([]domain.Referral, len(models))
	for i, m := range models {
		referrals[i] = *m.ToDomain()
	}
	return referrals, total, nil
}

func (r *ReferralRepo) ListByReferred(ctx context.Context, referredID string) ([]domain.Referral, error) {
	var models []ReferralModel
	if err := r.db.WithContext(ctx).Where("referred_id = ?", referredID).
		Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, err
	}

	referrals := make([]domain.Referral, len(models))
	for i, m := range models {
		referrals[i] = *m.ToDomain()
	}
	return referrals, nil
}

func (r *ReferralRepo) Create(ctx context.Context, referral *domain.Referral) error {
	model := ToReferralModel(referral)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *ReferralRepo) UpdateStatus(ctx context.Context, id string, status domain.ReferralStatus) error {
	return r.db.WithContext(ctx).Model(&ReferralModel{}).Where("id = ?", id).
		Update("status", string(status)).Error
}
