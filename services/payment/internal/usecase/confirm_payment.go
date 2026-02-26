package usecase

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/southern-martin/ecommerce/services/payment/internal/domain"
)

// WebhookEvent represents a simulated Stripe webhook event.
type WebhookEvent struct {
	Type            string `json:"type"`
	StripePaymentID string `json:"stripe_payment_id"`
	// SellerItems is included for crediting seller wallets on success.
	SellerItems []domain.OrderSellerItem `json:"seller_items,omitempty"`
	// FailureReason is populated when the payment fails.
	FailureReason string `json:"failure_reason,omitempty"`
}

// ConfirmPaymentUseCase handles webhook-driven payment confirmation.
type ConfirmPaymentUseCase struct {
	paymentRepo    domain.PaymentRepository
	walletRepo     domain.WalletRepository
	publisher      domain.EventPublisher
	commissionRate float64
}

// NewConfirmPaymentUseCase creates a new ConfirmPaymentUseCase.
func NewConfirmPaymentUseCase(
	paymentRepo domain.PaymentRepository,
	walletRepo domain.WalletRepository,
	publisher domain.EventPublisher,
	commissionRate float64,
) *ConfirmPaymentUseCase {
	return &ConfirmPaymentUseCase{
		paymentRepo:    paymentRepo,
		walletRepo:     walletRepo,
		publisher:      publisher,
		commissionRate: commissionRate,
	}
}

// Execute processes a webhook event for payment confirmation.
func (uc *ConfirmPaymentUseCase) Execute(ctx context.Context, event WebhookEvent) error {
	payment, err := uc.paymentRepo.GetByStripeID(ctx, event.StripePaymentID)
	if err != nil {
		return fmt.Errorf("payment not found for stripe ID %s: %w", event.StripePaymentID, err)
	}

	switch event.Type {
	case "payment_intent.succeeded":
		return uc.handleSuccess(ctx, payment, event.SellerItems)
	case "payment_intent.payment_failed":
		return uc.handleFailure(ctx, payment, event.FailureReason)
	default:
		log.Warn().Str("type", event.Type).Msg("Unhandled webhook event type")
		return nil
	}
}

func (uc *ConfirmPaymentUseCase) handleSuccess(ctx context.Context, payment *domain.Payment, sellerItems []domain.OrderSellerItem) error {
	if err := uc.paymentRepo.UpdateStatus(ctx, payment.ID, domain.PaymentStatusCompleted, ""); err != nil {
		return fmt.Errorf("failed to update payment status: %w", err)
	}

	// Credit seller wallets with platform commission deducted.
	for _, item := range sellerItems {
		commission := int64(math.Round(float64(item.AmountCents) * uc.commissionRate))
		sellerAmount := item.AmountCents - commission

		// Credit seller's pending balance.
		if err := uc.walletRepo.CreditPending(ctx, item.SellerID, sellerAmount); err != nil {
			log.Error().Err(err).Str("seller_id", item.SellerID).Msg("Failed to credit seller wallet")
			continue
		}

		// Record sale transaction.
		saleTx := &domain.WalletTransaction{
			ID:            uuid.New().String(),
			SellerID:      item.SellerID,
			Type:          domain.WalletTxSale,
			AmountCents:   sellerAmount,
			ReferenceType: "order",
			ReferenceID:   payment.OrderID,
			Description:   fmt.Sprintf("Sale from order %s", payment.OrderID),
			CreatedAt:     time.Now(),
		}
		if err := uc.walletRepo.CreateTransaction(ctx, saleTx); err != nil {
			log.Error().Err(err).Msg("Failed to create sale transaction")
		}

		// Record commission deduction transaction.
		commissionTx := &domain.WalletTransaction{
			ID:            uuid.New().String(),
			SellerID:      item.SellerID,
			Type:          domain.WalletTxCommissionDeducted,
			AmountCents:   -commission,
			ReferenceType: "order",
			ReferenceID:   payment.OrderID,
			Description:   fmt.Sprintf("Platform commission (%.0f%%) for order %s", uc.commissionRate*100, payment.OrderID),
			CreatedAt:     time.Now(),
		}
		if err := uc.walletRepo.CreateTransaction(ctx, commissionTx); err != nil {
			log.Error().Err(err).Msg("Failed to create commission transaction")
		}
	}

	// Publish payment.completed event.
	evt := domain.PaymentEvent{
		PaymentID:   payment.ID,
		OrderID:     payment.OrderID,
		BuyerID:     payment.BuyerID,
		AmountCents: payment.AmountCents,
		Currency:    payment.Currency,
		Status:      string(domain.PaymentStatusCompleted),
	}
	if err := uc.publisher.Publish(ctx, domain.EventPaymentCompleted, evt); err != nil {
		log.Error().Err(err).Msg("Failed to publish payment.completed event")
	}

	return nil
}

func (uc *ConfirmPaymentUseCase) handleFailure(ctx context.Context, payment *domain.Payment, reason string) error {
	if reason == "" {
		reason = "payment failed"
	}

	if err := uc.paymentRepo.UpdateStatus(ctx, payment.ID, domain.PaymentStatusFailed, reason); err != nil {
		return fmt.Errorf("failed to update payment status: %w", err)
	}

	// Publish payment.failed event.
	evt := domain.PaymentEvent{
		PaymentID:   payment.ID,
		OrderID:     payment.OrderID,
		BuyerID:     payment.BuyerID,
		AmountCents: payment.AmountCents,
		Currency:    payment.Currency,
		Status:      string(domain.PaymentStatusFailed),
	}
	if err := uc.publisher.Publish(ctx, domain.EventPaymentFailed, evt); err != nil {
		log.Error().Err(err).Msg("Failed to publish payment.failed event")
	}

	return nil
}
