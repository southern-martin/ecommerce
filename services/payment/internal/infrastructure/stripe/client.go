package stripe

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// StripeClient abstracts Stripe operations for testability.
type StripeClient interface {
	// CreatePaymentIntent creates a new payment intent and returns (paymentIntentID, clientSecret, error).
	CreatePaymentIntent(amountCents int64, currency string, metadata map[string]string) (string, string, error)
	// ConfirmPaymentIntent confirms a payment intent.
	ConfirmPaymentIntent(paymentIntentID string) error
	// CreateRefund creates a refund for a payment intent and returns (refundID, error).
	CreateRefund(paymentIntentID string, amountCents int64) (string, error)
	// CreateTransfer creates a transfer to a connected account and returns (transferID, error).
	CreateTransfer(amountCents int64, destinationAccountID string, metadata map[string]string) (string, error)
}

// MockStripeClient is used for development without actual Stripe integration.
type MockStripeClient struct{}

// NewMockStripeClient creates a new MockStripeClient.
func NewMockStripeClient() *MockStripeClient {
	log.Info().Msg("Using mock Stripe client")
	return &MockStripeClient{}
}

// CreatePaymentIntent creates a mock payment intent.
func (m *MockStripeClient) CreatePaymentIntent(amountCents int64, currency string, metadata map[string]string) (string, string, error) {
	id := "pi_mock_" + uuid.New().String()[:8]
	secret := "pi_secret_mock_" + uuid.New().String()[:8]
	log.Debug().
		Str("payment_intent_id", id).
		Int64("amount_cents", amountCents).
		Str("currency", currency).
		Msg("Mock: Created payment intent")
	return id, secret, nil
}

// ConfirmPaymentIntent confirms a mock payment intent.
func (m *MockStripeClient) ConfirmPaymentIntent(paymentIntentID string) error {
	log.Debug().
		Str("payment_intent_id", paymentIntentID).
		Msg("Mock: Confirmed payment intent")
	return nil
}

// CreateRefund creates a mock refund.
func (m *MockStripeClient) CreateRefund(paymentIntentID string, amountCents int64) (string, error) {
	refundID := fmt.Sprintf("re_mock_%s", uuid.New().String()[:8])
	log.Debug().
		Str("payment_intent_id", paymentIntentID).
		Str("refund_id", refundID).
		Int64("amount_cents", amountCents).
		Msg("Mock: Created refund")
	return refundID, nil
}

// CreateTransfer creates a mock transfer.
func (m *MockStripeClient) CreateTransfer(amountCents int64, destinationAccountID string, metadata map[string]string) (string, error) {
	transferID := fmt.Sprintf("tr_mock_%s", uuid.New().String()[:8])
	log.Debug().
		Str("transfer_id", transferID).
		Str("destination", destinationAccountID).
		Int64("amount_cents", amountCents).
		Msg("Mock: Created transfer")
	return transferID, nil
}

// Ensure MockStripeClient implements StripeClient.
var _ StripeClient = (*MockStripeClient)(nil)
