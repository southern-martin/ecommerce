package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/southern-martin/ecommerce/services/affiliate/internal/domain"
)

// PayoutUseCase handles payout operations.
type PayoutUseCase struct {
	payoutRepo  domain.PayoutRepository
	programRepo domain.AffiliateProgramRepository
	linkRepo    domain.AffiliateLinkRepository
	publisher   domain.EventPublisher
}

// NewPayoutUseCase creates a new PayoutUseCase.
func NewPayoutUseCase(
	payoutRepo domain.PayoutRepository,
	programRepo domain.AffiliateProgramRepository,
	linkRepo domain.AffiliateLinkRepository,
	publisher domain.EventPublisher,
) *PayoutUseCase {
	return &PayoutUseCase{
		payoutRepo:  payoutRepo,
		programRepo: programRepo,
		linkRepo:    linkRepo,
		publisher:   publisher,
	}
}

// RequestPayoutRequest is the input for requesting a payout.
type RequestPayoutRequest struct {
	UserID       string
	AmountCents  int64
	PayoutMethod domain.PayoutMethod
}

// RequestPayout creates a new payout request after validating minimum payout threshold.
func (uc *PayoutUseCase) RequestPayout(ctx context.Context, req RequestPayoutRequest) (*domain.AffiliatePayout, error) {
	program, err := uc.programRepo.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get affiliate program: %w", err)
	}

	// Validate minimum payout
	if req.AmountCents < program.MinPayoutCents {
		return nil, fmt.Errorf("payout amount %d cents is below minimum of %d cents", req.AmountCents, program.MinPayoutCents)
	}

	// Calculate total earnings from user's links
	links, _, err := uc.linkRepo.ListByUser(ctx, req.UserID, 1, 1000)
	if err != nil {
		return nil, fmt.Errorf("failed to get user links: %w", err)
	}

	var totalEarnings int64
	for _, link := range links {
		totalEarnings += link.TotalEarningsCents
	}

	if req.AmountCents > totalEarnings {
		return nil, fmt.Errorf("requested amount %d cents exceeds total earnings of %d cents", req.AmountCents, totalEarnings)
	}

	payout := &domain.AffiliatePayout{
		ID:           uuid.New().String(),
		UserID:       req.UserID,
		AmountCents:  req.AmountCents,
		Status:       domain.PayoutStatusRequested,
		PayoutMethod: req.PayoutMethod,
	}

	if err := uc.payoutRepo.Create(ctx, payout); err != nil {
		return nil, fmt.Errorf("failed to create payout: %w", err)
	}

	_ = uc.publisher.Publish(ctx, "affiliate.payout.requested", map[string]interface{}{
		"payout_id":     payout.ID,
		"user_id":       payout.UserID,
		"amount_cents":  payout.AmountCents,
		"payout_method": string(payout.PayoutMethod),
	})

	return payout, nil
}

// ListPayouts lists payouts for a user with pagination.
func (uc *PayoutUseCase) ListPayouts(ctx context.Context, userID string, page, pageSize int) ([]domain.AffiliatePayout, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return uc.payoutRepo.ListByUser(ctx, userID, page, pageSize)
}

// ListAllPayouts lists all payouts with pagination (admin).
func (uc *PayoutUseCase) ListAllPayouts(ctx context.Context, page, pageSize int) ([]domain.AffiliatePayout, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return uc.payoutRepo.ListAll(ctx, page, pageSize)
}

// GetPayout retrieves a payout by ID.
func (uc *PayoutUseCase) GetPayout(ctx context.Context, id string) (*domain.AffiliatePayout, error) {
	return uc.payoutRepo.GetByID(ctx, id)
}

// UpdatePayoutStatus updates the status of a payout.
func (uc *PayoutUseCase) UpdatePayoutStatus(ctx context.Context, id string, status domain.PayoutStatus) (*domain.AffiliatePayout, error) {
	var completedAt *time.Time
	if status == domain.PayoutStatusCompleted {
		now := time.Now()
		completedAt = &now
	}

	if err := uc.payoutRepo.UpdateStatus(ctx, id, status, completedAt); err != nil {
		return nil, fmt.Errorf("failed to update payout status: %w", err)
	}

	return uc.payoutRepo.GetByID(ctx, id)
}
