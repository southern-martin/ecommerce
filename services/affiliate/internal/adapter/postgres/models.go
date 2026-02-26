package postgres

import (
	"time"

	"github.com/southern-martin/ecommerce/services/affiliate/internal/domain"
)

// AffiliateProgramModel is the GORM model for the affiliate_programs table.
type AffiliateProgramModel struct {
	ID                 string    `gorm:"type:uuid;primaryKey"`
	CommissionRate     float64   `gorm:"type:decimal(5,4);not null;default:0.05"`
	MinPayoutCents     int64     `gorm:"not null;default:5000"`
	CookieDays         int       `gorm:"not null;default:30"`
	ReferrerBonusCents int64     `gorm:"not null;default:0"`
	ReferredBonusCents int64     `gorm:"not null;default:0"`
	IsActive           bool      `gorm:"default:true"`
	CreatedAt          time.Time `gorm:"autoCreateTime"`
	UpdatedAt          time.Time `gorm:"autoUpdateTime"`
}

func (AffiliateProgramModel) TableName() string { return "affiliate_programs" }

func (m *AffiliateProgramModel) ToDomain() *domain.AffiliateProgram {
	return &domain.AffiliateProgram{
		ID:                 m.ID,
		CommissionRate:     m.CommissionRate,
		MinPayoutCents:     m.MinPayoutCents,
		CookieDays:         m.CookieDays,
		ReferrerBonusCents: m.ReferrerBonusCents,
		ReferredBonusCents: m.ReferredBonusCents,
		IsActive:           m.IsActive,
		CreatedAt:          m.CreatedAt,
		UpdatedAt:          m.UpdatedAt,
	}
}

func ToAffiliateProgramModel(p *domain.AffiliateProgram) *AffiliateProgramModel {
	return &AffiliateProgramModel{
		ID:                 p.ID,
		CommissionRate:     p.CommissionRate,
		MinPayoutCents:     p.MinPayoutCents,
		CookieDays:         p.CookieDays,
		ReferrerBonusCents: p.ReferrerBonusCents,
		ReferredBonusCents: p.ReferredBonusCents,
		IsActive:           p.IsActive,
		CreatedAt:          p.CreatedAt,
		UpdatedAt:          p.UpdatedAt,
	}
}

// AffiliateLinkModel is the GORM model for the affiliate_links table.
type AffiliateLinkModel struct {
	ID                 string    `gorm:"type:uuid;primaryKey"`
	UserID             string    `gorm:"type:uuid;index;not null"`
	Code               string    `gorm:"type:varchar(8);uniqueIndex;not null"`
	TargetURL          string    `gorm:"type:text;not null"`
	ClickCount         int64     `gorm:"default:0"`
	ConversionCount    int64     `gorm:"default:0"`
	TotalEarningsCents int64     `gorm:"default:0"`
	CreatedAt          time.Time `gorm:"autoCreateTime"`
}

func (AffiliateLinkModel) TableName() string { return "affiliate_links" }

func (m *AffiliateLinkModel) ToDomain() *domain.AffiliateLink {
	return &domain.AffiliateLink{
		ID:                 m.ID,
		UserID:             m.UserID,
		Code:               m.Code,
		TargetURL:          m.TargetURL,
		ClickCount:         m.ClickCount,
		ConversionCount:    m.ConversionCount,
		TotalEarningsCents: m.TotalEarningsCents,
		CreatedAt:          m.CreatedAt,
	}
}

func ToAffiliateLinkModel(l *domain.AffiliateLink) *AffiliateLinkModel {
	return &AffiliateLinkModel{
		ID:                 l.ID,
		UserID:             l.UserID,
		Code:               l.Code,
		TargetURL:          l.TargetURL,
		ClickCount:         l.ClickCount,
		ConversionCount:    l.ConversionCount,
		TotalEarningsCents: l.TotalEarningsCents,
		CreatedAt:          l.CreatedAt,
	}
}

// ReferralModel is the GORM model for the referrals table.
type ReferralModel struct {
	ID              string    `gorm:"type:uuid;primaryKey"`
	ReferrerID      string    `gorm:"type:uuid;index;not null"`
	ReferredID      string    `gorm:"type:uuid;index;not null"`
	OrderID         string    `gorm:"type:uuid;index;not null"`
	OrderTotalCents int64     `gorm:"not null;default:0"`
	CommissionCents int64     `gorm:"not null;default:0"`
	Status          string    `gorm:"type:varchar(20);default:'pending'"`
	CreatedAt       time.Time `gorm:"autoCreateTime"`
}

func (ReferralModel) TableName() string { return "referrals" }

func (m *ReferralModel) ToDomain() *domain.Referral {
	return &domain.Referral{
		ID:              m.ID,
		ReferrerID:      m.ReferrerID,
		ReferredID:      m.ReferredID,
		OrderID:         m.OrderID,
		OrderTotalCents: m.OrderTotalCents,
		CommissionCents: m.CommissionCents,
		Status:          domain.ReferralStatus(m.Status),
		CreatedAt:       m.CreatedAt,
	}
}

func ToReferralModel(r *domain.Referral) *ReferralModel {
	return &ReferralModel{
		ID:              r.ID,
		ReferrerID:      r.ReferrerID,
		ReferredID:      r.ReferredID,
		OrderID:         r.OrderID,
		OrderTotalCents: r.OrderTotalCents,
		CommissionCents: r.CommissionCents,
		Status:          string(r.Status),
		CreatedAt:       r.CreatedAt,
	}
}

// AffiliatePayoutModel is the GORM model for the affiliate_payouts table.
type AffiliatePayoutModel struct {
	ID           string     `gorm:"type:uuid;primaryKey"`
	UserID       string     `gorm:"type:uuid;index;not null"`
	AmountCents  int64      `gorm:"not null;default:0"`
	Status       string     `gorm:"type:varchar(20);default:'requested'"`
	PayoutMethod string     `gorm:"type:varchar(20);not null"`
	CreatedAt    time.Time  `gorm:"autoCreateTime"`
	CompletedAt  *time.Time `gorm:""`
}

func (AffiliatePayoutModel) TableName() string { return "affiliate_payouts" }

func (m *AffiliatePayoutModel) ToDomain() *domain.AffiliatePayout {
	return &domain.AffiliatePayout{
		ID:           m.ID,
		UserID:       m.UserID,
		AmountCents:  m.AmountCents,
		Status:       domain.PayoutStatus(m.Status),
		PayoutMethod: domain.PayoutMethod(m.PayoutMethod),
		CreatedAt:    m.CreatedAt,
		CompletedAt:  m.CompletedAt,
	}
}

func ToAffiliatePayoutModel(p *domain.AffiliatePayout) *AffiliatePayoutModel {
	return &AffiliatePayoutModel{
		ID:           p.ID,
		UserID:       p.UserID,
		AmountCents:  p.AmountCents,
		Status:       string(p.Status),
		PayoutMethod: string(p.PayoutMethod),
		CreatedAt:    p.CreatedAt,
		CompletedAt:  p.CompletedAt,
	}
}
