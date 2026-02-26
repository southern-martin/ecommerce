package domain

import "time"

// MemberTier represents the tier level of a loyalty member.
type MemberTier string

const (
	TierBronze   MemberTier = "bronze"
	TierSilver   MemberTier = "silver"
	TierGold     MemberTier = "gold"
	TierPlatinum MemberTier = "platinum"
)

// TransactionType represents the type of a points transaction.
type TransactionType string

const (
	TransactionEarn   TransactionType = "earn"
	TransactionRedeem TransactionType = "redeem"
	TransactionExpire TransactionType = "expire"
	TransactionAdjust TransactionType = "adjust"
)

// PointsSource represents the source of a points transaction.
type PointsSource string

const (
	SourceOrder     PointsSource = "order"
	SourceReview    PointsSource = "review"
	SourceReferral  PointsSource = "referral"
	SourcePromotion PointsSource = "promotion"
	SourceSignup    PointsSource = "signup"
)

// Membership represents a user's loyalty membership.
type Membership struct {
	UserID         string
	Tier           MemberTier
	PointsBalance  int64
	LifetimePoints int64
	TierExpiresAt  *time.Time
	JoinedAt       time.Time
}

// PointsTransaction represents a points transaction.
type PointsTransaction struct {
	ID          string
	UserID      string
	Type        TransactionType
	Points      int64
	Source      PointsSource
	ReferenceID string
	Description string
	CreatedAt   time.Time
}

// Tier represents a loyalty tier definition.
type Tier struct {
	Name                 string
	MinPoints            int64
	CashbackRate         float64
	PointsMultiplier     float64
	FreeShipping         bool
	PrioritySupportHours int
}
