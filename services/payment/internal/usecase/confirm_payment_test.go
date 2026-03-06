package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/southern-martin/ecommerce/services/payment/internal/domain"
)

// ---------------------------------------------------------------------------
// ConfirmPaymentUseCase tests
// ---------------------------------------------------------------------------

func completedPayment() *domain.Payment {
	return &domain.Payment{
		ID:              "pay-1",
		OrderID:         "order-1",
		BuyerID:         "buyer-1",
		AmountCents:     10000,
		Currency:        "usd",
		Status:          domain.PaymentStatusPending,
		StripePaymentID: "pi_stripe_1",
	}
}

func TestConfirmPayment_Success(t *testing.T) {
	var updatedStatus domain.PaymentStatus
	var creditedSeller string
	var creditedAmount int64
	var publishedSubject string
	var createdTxCount int

	repo := &mockPaymentRepo{
		getByStripeIDFn: func(_ context.Context, _ string) (*domain.Payment, error) {
			return completedPayment(), nil
		},
		updateStatusFn: func(_ context.Context, _ string, status domain.PaymentStatus, _ string) error {
			updatedStatus = status
			return nil
		},
	}
	walletRepo := &mockWalletRepo{
		creditPendingFn: func(_ context.Context, sellerID string, amount int64) error {
			creditedSeller = sellerID
			creditedAmount = amount
			return nil
		},
		createTransactionFn: func(_ context.Context, _ *domain.WalletTransaction) error {
			createdTxCount++
			return nil
		},
	}
	pub := &mockEventPublisher{
		publishFn: func(_ context.Context, subj string, _ interface{}) error {
			publishedSubject = subj
			return nil
		},
	}

	uc := NewConfirmPaymentUseCase(repo, walletRepo, pub, 0.10) // 10% commission
	err := uc.Execute(context.Background(), WebhookEvent{
		Type:            "payment_intent.succeeded",
		StripePaymentID: "pi_stripe_1",
		SellerItems: []domain.OrderSellerItem{
			{SellerID: "seller-1", AmountCents: 10000},
		},
	})

	require.NoError(t, err)
	assert.Equal(t, domain.PaymentStatusCompleted, updatedStatus)
	assert.Equal(t, "seller-1", creditedSeller)
	assert.Equal(t, int64(9000), creditedAmount) // 10000 - 10% commission
	assert.Equal(t, 2, createdTxCount)           // sale tx + commission tx
	assert.Equal(t, domain.EventPaymentCompleted, publishedSubject)
}

func TestConfirmPayment_CommissionCalculation(t *testing.T) {
	var creditedAmounts []int64

	repo := &mockPaymentRepo{
		getByStripeIDFn: func(_ context.Context, _ string) (*domain.Payment, error) {
			return completedPayment(), nil
		},
	}
	walletRepo := &mockWalletRepo{
		creditPendingFn: func(_ context.Context, _ string, amount int64) error {
			creditedAmounts = append(creditedAmounts, amount)
			return nil
		},
	}

	uc := NewConfirmPaymentUseCase(repo, walletRepo, &mockEventPublisher{}, 0.15) // 15%
	err := uc.Execute(context.Background(), WebhookEvent{
		Type:            "payment_intent.succeeded",
		StripePaymentID: "pi_stripe_1",
		SellerItems: []domain.OrderSellerItem{
			{SellerID: "seller-a", AmountCents: 5000},
			{SellerID: "seller-b", AmountCents: 3000},
		},
	})

	require.NoError(t, err)
	require.Len(t, creditedAmounts, 2)
	// seller-a: 5000 - round(5000*0.15) = 5000 - 750 = 4250
	assert.Equal(t, int64(4250), creditedAmounts[0])
	// seller-b: 3000 - round(3000*0.15) = 3000 - 450 = 2550
	assert.Equal(t, int64(2550), creditedAmounts[1])
}

func TestConfirmPayment_ZeroCommission(t *testing.T) {
	var creditedAmount int64

	repo := &mockPaymentRepo{
		getByStripeIDFn: func(_ context.Context, _ string) (*domain.Payment, error) {
			return completedPayment(), nil
		},
	}
	walletRepo := &mockWalletRepo{
		creditPendingFn: func(_ context.Context, _ string, amount int64) error {
			creditedAmount = amount
			return nil
		},
	}

	uc := NewConfirmPaymentUseCase(repo, walletRepo, &mockEventPublisher{}, 0.0) // 0% commission
	err := uc.Execute(context.Background(), WebhookEvent{
		Type:            "payment_intent.succeeded",
		StripePaymentID: "pi_stripe_1",
		SellerItems: []domain.OrderSellerItem{
			{SellerID: "seller-1", AmountCents: 7000},
		},
	})

	require.NoError(t, err)
	assert.Equal(t, int64(7000), creditedAmount) // full amount, no commission
}

func TestConfirmPayment_FailureEvent(t *testing.T) {
	var updatedStatus domain.PaymentStatus
	var failureReason string
	var publishedSubject string

	repo := &mockPaymentRepo{
		getByStripeIDFn: func(_ context.Context, _ string) (*domain.Payment, error) {
			return completedPayment(), nil
		},
		updateStatusFn: func(_ context.Context, _ string, status domain.PaymentStatus, reason string) error {
			updatedStatus = status
			failureReason = reason
			return nil
		},
	}
	pub := &mockEventPublisher{
		publishFn: func(_ context.Context, subj string, _ interface{}) error {
			publishedSubject = subj
			return nil
		},
	}

	uc := NewConfirmPaymentUseCase(repo, &mockWalletRepo{}, pub, 0.10)
	err := uc.Execute(context.Background(), WebhookEvent{
		Type:            "payment_intent.payment_failed",
		StripePaymentID: "pi_stripe_1",
		FailureReason:   "card_declined",
	})

	require.NoError(t, err)
	assert.Equal(t, domain.PaymentStatusFailed, updatedStatus)
	assert.Equal(t, "card_declined", failureReason)
	assert.Equal(t, domain.EventPaymentFailed, publishedSubject)
}

func TestConfirmPayment_FailureDefaultReason(t *testing.T) {
	var failureReason string

	repo := &mockPaymentRepo{
		getByStripeIDFn: func(_ context.Context, _ string) (*domain.Payment, error) {
			return completedPayment(), nil
		},
		updateStatusFn: func(_ context.Context, _ string, _ domain.PaymentStatus, reason string) error {
			failureReason = reason
			return nil
		},
	}

	uc := NewConfirmPaymentUseCase(repo, &mockWalletRepo{}, &mockEventPublisher{}, 0.10)
	err := uc.Execute(context.Background(), WebhookEvent{
		Type:            "payment_intent.payment_failed",
		StripePaymentID: "pi_stripe_1",
		FailureReason:   "", // empty reason should default
	})

	require.NoError(t, err)
	assert.Equal(t, "payment failed", failureReason)
}

func TestConfirmPayment_UnknownEventType(t *testing.T) {
	repo := &mockPaymentRepo{
		getByStripeIDFn: func(_ context.Context, _ string) (*domain.Payment, error) {
			return completedPayment(), nil
		},
	}

	uc := NewConfirmPaymentUseCase(repo, &mockWalletRepo{}, &mockEventPublisher{}, 0.10)
	err := uc.Execute(context.Background(), WebhookEvent{
		Type:            "charge.dispute.created",
		StripePaymentID: "pi_stripe_1",
	})

	// Unknown event types should be silently ignored (no error).
	require.NoError(t, err)
}

func TestConfirmPayment_PaymentNotFound(t *testing.T) {
	repo := &mockPaymentRepo{
		getByStripeIDFn: func(_ context.Context, _ string) (*domain.Payment, error) {
			return nil, errors.New("record not found")
		},
	}

	uc := NewConfirmPaymentUseCase(repo, &mockWalletRepo{}, &mockEventPublisher{}, 0.10)
	err := uc.Execute(context.Background(), WebhookEvent{
		Type:            "payment_intent.succeeded",
		StripePaymentID: "pi_unknown",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "payment not found")
}

func TestConfirmPayment_UpdateStatusError(t *testing.T) {
	repo := &mockPaymentRepo{
		getByStripeIDFn: func(_ context.Context, _ string) (*domain.Payment, error) {
			return completedPayment(), nil
		},
		updateStatusFn: func(_ context.Context, _ string, _ domain.PaymentStatus, _ string) error {
			return errors.New("db write error")
		},
	}

	uc := NewConfirmPaymentUseCase(repo, &mockWalletRepo{}, &mockEventPublisher{}, 0.10)
	err := uc.Execute(context.Background(), WebhookEvent{
		Type:            "payment_intent.succeeded",
		StripePaymentID: "pi_stripe_1",
		SellerItems:     []domain.OrderSellerItem{{SellerID: "s1", AmountCents: 1000}},
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to update payment status")
}
