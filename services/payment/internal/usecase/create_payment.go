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
	PaymentID       string `json:"payment_id"`
	StripePaymentID string `json:"stripe_payment_id"`
	ClientSecret    string `json:"client_secret"`
	Status          string `json:"status"`
}

// OrderInfo holds minimal order data returned by the order service client.
type OrderInfo struct {
	ID         string
	TotalCents int64
	Currency   string
	Status     string
}

// OrderServiceClient defines the interface for cross-service order calls.
type OrderServiceClient interface {
	GetOrder(ctx context.Context, orderID string) (*OrderInfo, error)
	UpdateOrderStatus(ctx context.Context, orderID string, status string) error
}

// CreatePaymentUseCase handles creating payment intents.
type CreatePaymentUseCase struct {
	paymentRepo domain.PaymentRepository
	stripe      stripe.StripeClient
	publisher   domain.EventPublisher
	orderClient OrderServiceClient
}

// NewCreatePaymentUseCase creates a new CreatePaymentUseCase.
func NewCreatePaymentUseCase(
	paymentRepo domain.PaymentRepository,
	stripeClient stripe.StripeClient,
	publisher domain.EventPublisher,
	orderClient OrderServiceClient,
) *CreatePaymentUseCase {
	return &CreatePaymentUseCase{
		paymentRepo: paymentRepo,
		stripe:      stripeClient,
		publisher:   publisher,
		orderClient: orderClient,
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

	// Validate order exists and amount matches via order service.
	if uc.orderClient != nil {
		orderInfo, err := uc.orderClient.GetOrder(ctx, input.OrderID)
		if err != nil {
			return nil, fmt.Errorf("failed to validate order %s: %w", input.OrderID, err)
		}
		if orderInfo.TotalCents != input.AmountCents {
			return nil, fmt.Errorf("payment amount %d does not match order total %d", input.AmountCents, orderInfo.TotalCents)
		}
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

	// Persist stripe payment ID so webhook lookup works.
	payment.StripePaymentID = stripeID
	if err := uc.paymentRepo.UpdateStripeID(ctx, payment.ID, stripeID); err != nil {
		log.Error().Err(err).Str("payment_id", payment.ID).Msg("Failed to save stripe payment ID")
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
		PaymentID:       payment.ID,
		StripePaymentID: stripeID,
		ClientSecret:    clientSecret,
		Status:          string(payment.Status),
	}, nil
}
