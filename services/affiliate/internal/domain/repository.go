package domain

import (
	"context"
	"time"
)

// AffiliateProgramRepository defines the interface for affiliate program persistence.
type AffiliateProgramRepository interface {
	Get(ctx context.Context) (*AffiliateProgram, error)
	Create(ctx context.Context, program *AffiliateProgram) error
	Update(ctx context.Context, program *AffiliateProgram) error
}

// AffiliateLinkRepository defines the interface for affiliate link persistence.
type AffiliateLinkRepository interface {
	GetByID(ctx context.Context, id string) (*AffiliateLink, error)
	GetByCode(ctx context.Context, code string) (*AffiliateLink, error)
	ListByUser(ctx context.Context, userID string, page, pageSize int) ([]AffiliateLink, int64, error)
	Create(ctx context.Context, link *AffiliateLink) error
	IncrementClicks(ctx context.Context, id string) error
	IncrementConversions(ctx context.Context, id string) error
	AddEarnings(ctx context.Context, id string, amountCents int64) error
}

// ReferralRepository defines the interface for referral persistence.
type ReferralRepository interface {
	GetByID(ctx context.Context, id string) (*Referral, error)
	ListByReferrer(ctx context.Context, referrerID string, page, pageSize int) ([]Referral, int64, error)
	ListByReferred(ctx context.Context, referredID string) ([]Referral, error)
	Create(ctx context.Context, referral *Referral) error
	UpdateStatus(ctx context.Context, id string, status ReferralStatus) error
}

// PayoutRepository defines the interface for payout persistence.
type PayoutRepository interface {
	GetByID(ctx context.Context, id string) (*AffiliatePayout, error)
	ListByUser(ctx context.Context, userID string, page, pageSize int) ([]AffiliatePayout, int64, error)
	ListAll(ctx context.Context, page, pageSize int) ([]AffiliatePayout, int64, error)
	Create(ctx context.Context, payout *AffiliatePayout) error
	UpdateStatus(ctx context.Context, id string, status PayoutStatus, completedAt *time.Time) error
}
