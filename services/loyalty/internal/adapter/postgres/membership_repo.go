package postgres

import (
	"context"

	"github.com/southern-martin/ecommerce/services/loyalty/internal/domain"
	"gorm.io/gorm"
)

// MembershipRepo implements domain.MembershipRepository.
type MembershipRepo struct {
	db *gorm.DB
}

// NewMembershipRepo creates a new MembershipRepo.
func NewMembershipRepo(db *gorm.DB) *MembershipRepo {
	return &MembershipRepo{db: db}
}

func (r *MembershipRepo) GetByUserID(ctx context.Context, userID string) (*domain.Membership, error) {
	var model MembershipModel
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&model).Error; err != nil {
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *MembershipRepo) Create(ctx context.Context, membership *domain.Membership) error {
	model := ToMembershipModel(membership)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *MembershipRepo) Update(ctx context.Context, membership *domain.Membership) error {
	model := ToMembershipModel(membership)
	return r.db.WithContext(ctx).Model(&MembershipModel{}).Where("user_id = ?", membership.UserID).Updates(model).Error
}

func (r *MembershipRepo) UpdateTier(ctx context.Context, userID string, tier domain.MemberTier) error {
	return r.db.WithContext(ctx).Model(&MembershipModel{}).Where("user_id = ?", userID).Update("tier", string(tier)).Error
}

func (r *MembershipRepo) UpdatePoints(ctx context.Context, userID string, pointsBalance, lifetimePoints int64) error {
	return r.db.WithContext(ctx).Model(&MembershipModel{}).Where("user_id = ?", userID).Updates(map[string]interface{}{
		"points_balance":  pointsBalance,
		"lifetime_points": lifetimePoints,
	}).Error
}
