package postgres

import (
	"time"

	"github.com/southern-martin/ecommerce/services/loyalty/internal/domain"
)

// MembershipModel is the GORM model for the memberships table.
type MembershipModel struct {
	UserID         string     `gorm:"type:varchar(36);primaryKey"`
	Tier           string     `gorm:"type:varchar(20);not null;default:'bronze'"`
	PointsBalance  int64      `gorm:"not null;default:0"`
	LifetimePoints int64      `gorm:"not null;default:0"`
	TierExpiresAt  *time.Time `gorm:"index"`
	JoinedAt       time.Time  `gorm:"autoCreateTime"`
}

func (MembershipModel) TableName() string { return "memberships" }

func (m *MembershipModel) ToDomain() *domain.Membership {
	return &domain.Membership{
		UserID:         m.UserID,
		Tier:           domain.MemberTier(m.Tier),
		PointsBalance:  m.PointsBalance,
		LifetimePoints: m.LifetimePoints,
		TierExpiresAt:  m.TierExpiresAt,
		JoinedAt:       m.JoinedAt,
	}
}

func ToMembershipModel(m *domain.Membership) *MembershipModel {
	return &MembershipModel{
		UserID:         m.UserID,
		Tier:           string(m.Tier),
		PointsBalance:  m.PointsBalance,
		LifetimePoints: m.LifetimePoints,
		TierExpiresAt:  m.TierExpiresAt,
		JoinedAt:       m.JoinedAt,
	}
}

// PointsTransactionModel is the GORM model for the points_transactions table.
type PointsTransactionModel struct {
	ID          string    `gorm:"type:uuid;primaryKey"`
	UserID      string    `gorm:"type:varchar(36);index;not null"`
	Type        string    `gorm:"type:varchar(20);not null"`
	Points      int64     `gorm:"not null"`
	Source      string    `gorm:"type:varchar(20);not null"`
	ReferenceID string    `gorm:"type:varchar(255)"`
	Description string    `gorm:"type:text"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}

func (PointsTransactionModel) TableName() string { return "points_transactions" }

func (m *PointsTransactionModel) ToDomain() *domain.PointsTransaction {
	return &domain.PointsTransaction{
		ID:          m.ID,
		UserID:      m.UserID,
		Type:        domain.TransactionType(m.Type),
		Points:      m.Points,
		Source:      domain.PointsSource(m.Source),
		ReferenceID: m.ReferenceID,
		Description: m.Description,
		CreatedAt:   m.CreatedAt,
	}
}

func ToPointsTransactionModel(t *domain.PointsTransaction) *PointsTransactionModel {
	return &PointsTransactionModel{
		ID:          t.ID,
		UserID:      t.UserID,
		Type:        string(t.Type),
		Points:      t.Points,
		Source:      string(t.Source),
		ReferenceID: t.ReferenceID,
		Description: t.Description,
		CreatedAt:   t.CreatedAt,
	}
}

// TierModel is the GORM model for the tiers table.
type TierModel struct {
	Name                 string  `gorm:"type:varchar(20);primaryKey"`
	MinPoints            int64   `gorm:"not null;default:0"`
	CashbackRate         float64 `gorm:"not null;default:0"`
	PointsMultiplier     float64 `gorm:"not null;default:1"`
	FreeShipping         bool    `gorm:"not null;default:false"`
	PrioritySupportHours int     `gorm:"not null;default:48"`
}

func (TierModel) TableName() string { return "tiers" }

func (m *TierModel) ToDomain() *domain.Tier {
	return &domain.Tier{
		Name:                 m.Name,
		MinPoints:            m.MinPoints,
		CashbackRate:         m.CashbackRate,
		PointsMultiplier:     m.PointsMultiplier,
		FreeShipping:         m.FreeShipping,
		PrioritySupportHours: m.PrioritySupportHours,
	}
}

func ToTierModel(t *domain.Tier) *TierModel {
	return &TierModel{
		Name:                 t.Name,
		MinPoints:            t.MinPoints,
		CashbackRate:         t.CashbackRate,
		PointsMultiplier:     t.PointsMultiplier,
		FreeShipping:         t.FreeShipping,
		PrioritySupportHours: t.PrioritySupportHours,
	}
}
