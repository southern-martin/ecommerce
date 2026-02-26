package postgres

import (
	"context"

	"github.com/southern-martin/ecommerce/services/affiliate/internal/domain"
	"gorm.io/gorm"
)

// ProgramRepo implements domain.AffiliateProgramRepository.
type ProgramRepo struct {
	db *gorm.DB
}

// NewProgramRepo creates a new ProgramRepo.
func NewProgramRepo(db *gorm.DB) *ProgramRepo {
	return &ProgramRepo{db: db}
}

func (r *ProgramRepo) Get(ctx context.Context) (*domain.AffiliateProgram, error) {
	var model AffiliateProgramModel
	if err := r.db.WithContext(ctx).Where("is_active = ?", true).First(&model).Error; err != nil {
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *ProgramRepo) Create(ctx context.Context, program *domain.AffiliateProgram) error {
	model := ToAffiliateProgramModel(program)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *ProgramRepo) Update(ctx context.Context, program *domain.AffiliateProgram) error {
	return r.db.WithContext(ctx).Model(&AffiliateProgramModel{}).Where("id = ?", program.ID).Updates(map[string]interface{}{
		"commission_rate":      program.CommissionRate,
		"min_payout_cents":     program.MinPayoutCents,
		"cookie_days":          program.CookieDays,
		"referrer_bonus_cents": program.ReferrerBonusCents,
		"referred_bonus_cents": program.ReferredBonusCents,
		"is_active":            program.IsActive,
	}).Error
}
