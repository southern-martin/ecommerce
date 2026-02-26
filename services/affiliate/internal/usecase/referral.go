package usecase

import (
	"context"
	"fmt"
	"math"

	"github.com/google/uuid"
	"github.com/southern-martin/ecommerce/services/affiliate/internal/domain"
)

// ReferralUseCase handles referral operations.
type ReferralUseCase struct {
	referralRepo domain.ReferralRepository
	linkRepo     domain.AffiliateLinkRepository
	programRepo  domain.AffiliateProgramRepository
	publisher    domain.EventPublisher
}

// NewReferralUseCase creates a new ReferralUseCase.
func NewReferralUseCase(
	referralRepo domain.ReferralRepository,
	linkRepo domain.AffiliateLinkRepository,
	programRepo domain.AffiliateProgramRepository,
	publisher domain.EventPublisher,
) *ReferralUseCase {
	return &ReferralUseCase{
		referralRepo: referralRepo,
		linkRepo:     linkRepo,
		programRepo:  programRepo,
		publisher:    publisher,
	}
}

// TrackConversionRequest is the input for tracking a conversion.
type TrackConversionRequest struct {
	LinkID          string
	ReferredID      string
	OrderID         string
	OrderTotalCents int64
}

// TrackConversion records a referral conversion, calculates commission, and updates link stats.
func (uc *ReferralUseCase) TrackConversion(ctx context.Context, req TrackConversionRequest) (*domain.Referral, error) {
	link, err := uc.linkRepo.GetByID(ctx, req.LinkID)
	if err != nil {
		return nil, fmt.Errorf("affiliate link not found: %w", err)
	}

	program, err := uc.programRepo.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get affiliate program: %w", err)
	}

	// Calculate commission from program rate
	commissionCents := int64(math.Round(float64(req.OrderTotalCents) * program.CommissionRate))

	referral := &domain.Referral{
		ID:              uuid.New().String(),
		ReferrerID:      link.UserID,
		ReferredID:      req.ReferredID,
		OrderID:         req.OrderID,
		OrderTotalCents: req.OrderTotalCents,
		CommissionCents: commissionCents,
		Status:          domain.ReferralStatusPending,
	}

	if err := uc.referralRepo.Create(ctx, referral); err != nil {
		return nil, fmt.Errorf("failed to create referral: %w", err)
	}

	// Increment link conversions and add earnings
	_ = uc.linkRepo.IncrementConversions(ctx, link.ID)
	_ = uc.linkRepo.AddEarnings(ctx, link.ID, commissionCents)

	_ = uc.publisher.Publish(ctx, "affiliate.conversion.tracked", map[string]interface{}{
		"referral_id":      referral.ID,
		"referrer_id":      referral.ReferrerID,
		"referred_id":      referral.ReferredID,
		"order_id":         referral.OrderID,
		"order_total_cents": referral.OrderTotalCents,
		"commission_cents": referral.CommissionCents,
	})

	return referral, nil
}

// ListReferrals lists referrals for a referrer with pagination.
func (uc *ReferralUseCase) ListReferrals(ctx context.Context, referrerID string, page, pageSize int) ([]domain.Referral, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return uc.referralRepo.ListByReferrer(ctx, referrerID, page, pageSize)
}
