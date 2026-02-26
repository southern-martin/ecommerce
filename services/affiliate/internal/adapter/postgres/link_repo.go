package postgres

import (
	"context"

	"github.com/southern-martin/ecommerce/services/affiliate/internal/domain"
	"gorm.io/gorm"
)

// LinkRepo implements domain.AffiliateLinkRepository.
type LinkRepo struct {
	db *gorm.DB
}

// NewLinkRepo creates a new LinkRepo.
func NewLinkRepo(db *gorm.DB) *LinkRepo {
	return &LinkRepo{db: db}
}

func (r *LinkRepo) GetByID(ctx context.Context, id string) (*domain.AffiliateLink, error) {
	var model AffiliateLinkModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *LinkRepo) GetByCode(ctx context.Context, code string) (*domain.AffiliateLink, error) {
	var model AffiliateLinkModel
	if err := r.db.WithContext(ctx).Where("code = ?", code).First(&model).Error; err != nil {
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *LinkRepo) ListByUser(ctx context.Context, userID string, page, pageSize int) ([]domain.AffiliateLink, int64, error) {
	var total int64
	r.db.WithContext(ctx).Model(&AffiliateLinkModel{}).Where("user_id = ?", userID).Count(&total)

	var models []AffiliateLinkModel
	offset := (page - 1) * pageSize
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).
		Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	links := make([]domain.AffiliateLink, len(models))
	for i, m := range models {
		links[i] = *m.ToDomain()
	}
	return links, total, nil
}

func (r *LinkRepo) Create(ctx context.Context, link *domain.AffiliateLink) error {
	model := ToAffiliateLinkModel(link)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *LinkRepo) IncrementClicks(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&AffiliateLinkModel{}).Where("id = ?", id).
		UpdateColumn("click_count", gorm.Expr("click_count + 1")).Error
}

func (r *LinkRepo) IncrementConversions(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&AffiliateLinkModel{}).Where("id = ?", id).
		UpdateColumn("conversion_count", gorm.Expr("conversion_count + 1")).Error
}

func (r *LinkRepo) AddEarnings(ctx context.Context, id string, amountCents int64) error {
	return r.db.WithContext(ctx).Model(&AffiliateLinkModel{}).Where("id = ?", id).
		UpdateColumn("total_earnings_cents", gorm.Expr("total_earnings_cents + ?", amountCents)).Error
}
