package domain

import "time"

// ReferralStatus represents the status of a referral.
type ReferralStatus string

const (
	ReferralStatusPending   ReferralStatus = "pending"
	ReferralStatusConfirmed ReferralStatus = "confirmed"
	ReferralStatusPaid      ReferralStatus = "paid"
)

// PayoutStatus represents the status of a payout.
type PayoutStatus string

const (
	PayoutStatusRequested  PayoutStatus = "requested"
	PayoutStatusProcessing PayoutStatus = "processing"
	PayoutStatusCompleted  PayoutStatus = "completed"
	PayoutStatusFailed     PayoutStatus = "failed"
)

// PayoutMethod represents the payout method.
type PayoutMethod string

const (
	PayoutMethodBankTransfer PayoutMethod = "bank_transfer"
	PayoutMethodStripe       PayoutMethod = "stripe"
)

// AffiliateProgram represents the affiliate program configuration.
type AffiliateProgram struct {
	ID                string
	CommissionRate    float64
	MinPayoutCents    int64
	CookieDays        int
	ReferrerBonusCents int64
	ReferredBonusCents int64
	IsActive          bool
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// AffiliateLink represents an affiliate referral link.
type AffiliateLink struct {
	ID                string
	UserID            string
	Code              string
	TargetURL         string
	ClickCount        int64
	ConversionCount   int64
	TotalEarningsCents int64
	CreatedAt         time.Time
}

// Referral represents a referral conversion.
type Referral struct {
	ID              string
	ReferrerID      string
	ReferredID      string
	OrderID         string
	OrderTotalCents int64
	CommissionCents int64
	Status          ReferralStatus
	CreatedAt       time.Time
}

// AffiliatePayout represents a payout request.
type AffiliatePayout struct {
	ID           string
	UserID       string
	AmountCents  int64
	Status       PayoutStatus
	PayoutMethod PayoutMethod
	CreatedAt    time.Time
	CompletedAt  *time.Time
}
