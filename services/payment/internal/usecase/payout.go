package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/southern-martin/ecommerce/services/payment/internal/domain"
	"github.com/southern-martin/ecommerce/services/payment/internal/infrastructure/stripe"
)

// RequestPayoutInput holds the input for requesting a payout.
type RequestPayoutInput struct {
	SellerID    string `json:"seller_id"`
	AmountCents int64  `json:"amount_cents" binding:"required"`
	Currency    string `json:"currency"`
	Method      string `json:"method"`
}

// PayoutUseCase handles payout-related operations.
type PayoutUseCase struct {
	payoutRepo domain.PayoutRepository
	walletRepo domain.WalletRepository
	stripe     stripe.StripeClient
}

// NewPayoutUseCase creates a new PayoutUseCase.
func NewPayoutUseCase(
	payoutRepo domain.PayoutRepository,
	walletRepo domain.WalletRepository,
	stripeClient stripe.StripeClient,
) *PayoutUseCase {
	return &PayoutUseCase{
		payoutRepo: payoutRepo,
		walletRepo: walletRepo,
		stripe:     stripeClient,
	}
}

// RequestPayout creates a new payout request.
func (uc *PayoutUseCase) RequestPayout(ctx context.Context, input RequestPayoutInput) (*domain.Payout, error) {
	if input.Currency == "" {
		input.Currency = "usd"
	}
	if input.Method == "" {
		input.Method = "stripe_connect"
	}

	// Verify seller has sufficient available balance.
	wallet, err := uc.walletRepo.GetOrCreate(ctx, input.SellerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}

	if wallet.AvailableBalance < input.AmountCents {
		return nil, fmt.Errorf("insufficient available balance: have %d, need %d", wallet.AvailableBalance, input.AmountCents)
	}

	// Debit available balance.
	if err := uc.walletRepo.DebitAvailable(ctx, input.SellerID, input.AmountCents); err != nil {
		return nil, fmt.Errorf("failed to debit wallet: %w", err)
	}

	// Record payout transaction.
	payoutID := uuid.New().String()
	tx := &domain.WalletTransaction{
		ID:            uuid.New().String(),
		SellerID:      input.SellerID,
		Type:          domain.WalletTxPayout,
		AmountCents:   -input.AmountCents,
		ReferenceType: "payout",
		ReferenceID:   payoutID,
		Description:   fmt.Sprintf("Payout request %s", payoutID),
		CreatedAt:     time.Now(),
	}
	if err := uc.walletRepo.CreateTransaction(ctx, tx); err != nil {
		log.Error().Err(err).Msg("Failed to create payout wallet transaction")
	}

	payout := &domain.Payout{
		ID:          payoutID,
		SellerID:    input.SellerID,
		AmountCents: input.AmountCents,
		Currency:    input.Currency,
		Method:      input.Method,
		Status:      domain.PayoutStatusRequested,
		RequestedAt: time.Now(),
	}

	if err := uc.payoutRepo.Create(ctx, payout); err != nil {
		return nil, fmt.Errorf("failed to create payout: %w", err)
	}

	return payout, nil
}

// ProcessPayout processes a pending payout (called by admin or cron).
func (uc *PayoutUseCase) ProcessPayout(ctx context.Context, payoutID string) error {
	payout, err := uc.payoutRepo.GetByID(ctx, payoutID)
	if err != nil {
		return fmt.Errorf("payout not found: %w", err)
	}

	if payout.Status != domain.PayoutStatusRequested {
		return fmt.Errorf("payout is not in requested status: %s", payout.Status)
	}

	// Update to processing.
	if err := uc.payoutRepo.UpdateStatus(ctx, payoutID, domain.PayoutStatusProcessing); err != nil {
		return fmt.Errorf("failed to update payout status: %w", err)
	}

	// Create Stripe transfer.
	metadata := map[string]string{
		"payout_id": payoutID,
		"seller_id": payout.SellerID,
	}
	transferID, err := uc.stripe.CreateTransfer(payout.AmountCents, payout.SellerID, metadata)
	if err != nil {
		_ = uc.payoutRepo.UpdateStatus(ctx, payoutID, domain.PayoutStatusFailed)
		return fmt.Errorf("failed to create stripe transfer: %w", err)
	}

	log.Info().
		Str("payout_id", payoutID).
		Str("transfer_id", transferID).
		Msg("Payout processed successfully")

	return uc.payoutRepo.UpdateStatus(ctx, payoutID, domain.PayoutStatusCompleted)
}

// ListPayouts lists payouts for a seller.
func (uc *PayoutUseCase) ListPayouts(ctx context.Context, sellerID string, page, pageSize int) ([]*domain.Payout, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return uc.payoutRepo.ListBySeller(ctx, sellerID, page, pageSize)
}
