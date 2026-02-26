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

// RefundUseCase handles refund processing.
type RefundUseCase struct {
	paymentRepo domain.PaymentRepository
	walletRepo  domain.WalletRepository
	stripe      stripe.StripeClient
	publisher   domain.EventPublisher
}

// NewRefundUseCase creates a new RefundUseCase.
func NewRefundUseCase(
	paymentRepo domain.PaymentRepository,
	walletRepo domain.WalletRepository,
	stripeClient stripe.StripeClient,
	publisher domain.EventPublisher,
) *RefundUseCase {
	return &RefundUseCase{
		paymentRepo: paymentRepo,
		walletRepo:  walletRepo,
		stripe:      stripeClient,
		publisher:   publisher,
	}
}

// RefundInput holds the input for processing a refund.
type RefundInput struct {
	OrderID     string `json:"order_id" binding:"required"`
	AmountCents int64  `json:"amount_cents"` // 0 means full refund
	SellerID    string `json:"seller_id"`
}

// ProcessRefund processes a refund for a payment.
func (uc *RefundUseCase) ProcessRefund(ctx context.Context, input RefundInput) error {
	payment, err := uc.paymentRepo.GetByOrderID(ctx, input.OrderID)
	if err != nil {
		return fmt.Errorf("payment not found for order %s: %w", input.OrderID, err)
	}

	if payment.Status != domain.PaymentStatusCompleted {
		return fmt.Errorf("cannot refund payment with status %s", payment.Status)
	}

	refundAmount := input.AmountCents
	if refundAmount == 0 {
		refundAmount = payment.AmountCents
	}

	if refundAmount > payment.AmountCents {
		return fmt.Errorf("refund amount %d exceeds payment amount %d", refundAmount, payment.AmountCents)
	}

	// Create refund in Stripe.
	refundID, err := uc.stripe.CreateRefund(payment.StripePaymentID, refundAmount)
	if err != nil {
		return fmt.Errorf("failed to create stripe refund: %w", err)
	}

	log.Info().
		Str("refund_id", refundID).
		Str("payment_id", payment.ID).
		Int64("amount_cents", refundAmount).
		Msg("Refund created")

	// Update payment status.
	if err := uc.paymentRepo.UpdateStatus(ctx, payment.ID, domain.PaymentStatusRefunded, ""); err != nil {
		return fmt.Errorf("failed to update payment status: %w", err)
	}

	// Debit seller wallet if seller ID is provided.
	if input.SellerID != "" {
		if err := uc.walletRepo.DebitAvailable(ctx, input.SellerID, refundAmount); err != nil {
			log.Error().Err(err).Str("seller_id", input.SellerID).Msg("Failed to debit seller wallet for refund")
		}

		tx := &domain.WalletTransaction{
			ID:            uuid.New().String(),
			SellerID:      input.SellerID,
			Type:          domain.WalletTxRefundDebit,
			AmountCents:   -refundAmount,
			ReferenceType: "refund",
			ReferenceID:   payment.OrderID,
			Description:   fmt.Sprintf("Refund for order %s", payment.OrderID),
			CreatedAt:     time.Now(),
		}
		if err := uc.walletRepo.CreateTransaction(ctx, tx); err != nil {
			log.Error().Err(err).Msg("Failed to create refund wallet transaction")
		}
	}

	// Publish payment.refunded event.
	evt := domain.PaymentEvent{
		PaymentID:   payment.ID,
		OrderID:     payment.OrderID,
		BuyerID:     payment.BuyerID,
		AmountCents: refundAmount,
		Currency:    payment.Currency,
		Status:      string(domain.PaymentStatusRefunded),
	}
	if err := uc.publisher.Publish(ctx, domain.EventPaymentRefunded, evt); err != nil {
		log.Error().Err(err).Msg("Failed to publish payment.refunded event")
	}

	return nil
}
