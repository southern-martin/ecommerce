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

// CreatePaymentInput holds the input for creating a payment intent.
type CreatePaymentInput struct {
	OrderID     string                   `json:"order_id" binding:"required"`
	BuyerID     string                   `json:"buyer_id"`
	AmountCents int64                    `json:"amount_cents" binding:"required"`
	Currency    string                   `json:"currency"`
	Method      domain.PaymentMethod     `json:"method"`
	SellerItems []domain.OrderSellerItem `json:"seller_items"`
}

// CreatePaymentOutput holds the output of creating a payment intent.
type CreatePaymentOutput struct {
	PaymentID    string `json:"payment_id"`
	ClientSecret string `json:"client_secret"`
	Status       string `json:"status"`
}

// CreatePaymentUseCase handles creating payment intents.
type CreatePaymentUseCase struct {
	paymentRepo domain.PaymentRepository
	stripe      stripe.StripeClient
	publisher   domain.EventPublisher
}

// NewCreatePaymentUseCase creates a new CreatePaymentUseCase.
func NewCreatePaymentUseCase(
	paymentRepo domain.PaymentRepository,
	stripeClient stripe.StripeClient,
	publisher domain.EventPublisher,
) *CreatePaymentUseCase {
	return &CreatePaymentUseCase{
		paymentRepo: paymentRepo,
		stripe:      stripeClient,
		publisher:   publisher,
	}
}

// Execute creates a new payment intent.
func (uc *CreatePaymentUseCase) Execute(ctx context.Context, input CreatePaymentInput) (*CreatePaymentOutput, error) {
	if input.Currency == "" {
		input.Currency = "usd"
	}
	if input.Method == "" {
		input.Method = domain.PaymentMethodCard
	}

	// Create payment record with pending status.
	payment := &domain.Payment{
		ID:          uuid.New().String(),
		OrderID:     input.OrderID,
		BuyerID:     input.BuyerID,
		AmountCents: input.AmountCents,
		Currency:    input.Currency,
		Status:      domain.PaymentStatusPending,
		Method:      input.Method,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := uc.paymentRepo.Create(ctx, payment); err != nil {
		return nil, fmt.Errorf("failed to create payment record: %w", err)
	}

	// Create Stripe PaymentIntent.
	metadata := map[string]string{
		"payment_id": payment.ID,
		"order_id":   payment.OrderID,
		"buyer_id":   payment.BuyerID,
	}

	stripeID, clientSecret, err := uc.stripe.CreatePaymentIntent(
		payment.AmountCents,
		payment.Currency,
		metadata,
	)
	if err != nil {
		_ = uc.paymentRepo.UpdateStatus(ctx, payment.ID, domain.PaymentStatusFailed, err.Error())
		return nil, fmt.Errorf("failed to create stripe payment intent: %w", err)
	}

	// Update payment with Stripe ID.
	payment.StripePaymentID = stripeID
	if err := uc.paymentRepo.UpdateStatus(ctx, payment.ID, domain.PaymentStatusPending, ""); err != nil {
		log.Error().Err(err).Str("payment_id", payment.ID).Msg("Failed to update payment with stripe ID")
	}

	// Publish payment.initiated event.
	event := domain.PaymentEvent{
		PaymentID:   payment.ID,
		OrderID:     payment.OrderID,
		BuyerID:     payment.BuyerID,
		AmountCents: payment.AmountCents,
		Currency:    payment.Currency,
		Status:      string(domain.PaymentStatusPending),
	}
	if err := uc.publisher.Publish(ctx, domain.EventPaymentInitiated, event); err != nil {
		log.Error().Err(err).Msg("Failed to publish payment.initiated event")
	}

	return &CreatePaymentOutput{
		PaymentID:    payment.ID,
		ClientSecret: clientSecret,
		Status:       string(payment.Status),
	}, nil
}
