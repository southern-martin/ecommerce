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
// RefundUseCase tests
// ---------------------------------------------------------------------------

func refundablePayment() *domain.Payment {
	return &domain.Payment{
		ID:              "pay-r1",
		OrderID:         "order-r1",
		BuyerID:         "buyer-r1",
		AmountCents:     8000,
		Currency:        "usd",
		Status:          domain.PaymentStatusCompleted,
		StripePaymentID: "pi_refund_1",
	}
}

func TestRefund_FullRefund(t *testing.T) {
	var refundedAmount int64
	var updatedStatus domain.PaymentStatus
	var publishedSubject string

	repo := &mockPaymentRepo{
		getByOrderIDFn: func(_ context.Context, _ string) (*domain.Payment, error) {
			return refundablePayment(), nil
		},
		updateStatusFn: func(_ context.Context, _ string, status domain.PaymentStatus, _ string) error {
			updatedStatus = status
			return nil
		},
	}
	sc := &mockStripeClient{
		createRefundFn: func(_ string, amount int64) (string, error) {
			refundedAmount = amount
			return "re_full", nil
		},
	}
	pub := &mockEventPublisher{
		publishFn: func(_ context.Context, subj string, _ interface{}) error {
			publishedSubject = subj
			return nil
		},
	}

	uc := NewRefundUseCase(repo, &mockWalletRepo{}, sc, pub)
	err := uc.ProcessRefund(context.Background(), RefundInput{
		OrderID:     "order-r1",
		AmountCents: 0, // 0 means full refund
	})

	require.NoError(t, err)
	assert.Equal(t, int64(8000), refundedAmount) // full payment amount
	assert.Equal(t, domain.PaymentStatusRefunded, updatedStatus)
	assert.Equal(t, domain.EventPaymentRefunded, publishedSubject)
}

func TestRefund_PartialRefund(t *testing.T) {
	var refundedAmount int64

	repo := &mockPaymentRepo{
		getByOrderIDFn: func(_ context.Context, _ string) (*domain.Payment, error) {
			return refundablePayment(), nil
		},
	}
	sc := &mockStripeClient{
		createRefundFn: func(_ string, amount int64) (string, error) {
			refundedAmount = amount
			return "re_partial", nil
		},
	}

	uc := NewRefundUseCase(repo, &mockWalletRepo{}, sc, &mockEventPublisher{})
	err := uc.ProcessRefund(context.Background(), RefundInput{
		OrderID:     "order-r1",
		AmountCents: 3000, // partial
	})

	require.NoError(t, err)
	assert.Equal(t, int64(3000), refundedAmount)
}

func TestRefund_ExceedsPaymentAmount(t *testing.T) {
	repo := &mockPaymentRepo{
		getByOrderIDFn: func(_ context.Context, _ string) (*domain.Payment, error) {
			return refundablePayment(), nil // AmountCents = 8000
		},
	}

	uc := NewRefundUseCase(repo, &mockWalletRepo{}, &mockStripeClient{}, &mockEventPublisher{})
	err := uc.ProcessRefund(context.Background(), RefundInput{
		OrderID:     "order-r1",
		AmountCents: 9000, // exceeds 8000
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "exceeds payment amount")
}

func TestRefund_NonCompletedStatus(t *testing.T) {
	pending := refundablePayment()
	pending.Status = domain.PaymentStatusPending

	repo := &mockPaymentRepo{
		getByOrderIDFn: func(_ context.Context, _ string) (*domain.Payment, error) {
			return pending, nil
		},
	}

	uc := NewRefundUseCase(repo, &mockWalletRepo{}, &mockStripeClient{}, &mockEventPublisher{})
	err := uc.ProcessRefund(context.Background(), RefundInput{
		OrderID:     "order-r1",
		AmountCents: 1000,
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "cannot refund payment with status")
}

func TestRefund_StripeError(t *testing.T) {
	repo := &mockPaymentRepo{
		getByOrderIDFn: func(_ context.Context, _ string) (*domain.Payment, error) {
			return refundablePayment(), nil
		},
	}
	sc := &mockStripeClient{
		createRefundFn: func(_ string, _ int64) (string, error) {
			return "", errors.New("stripe refund failed")
		},
	}

	uc := NewRefundUseCase(repo, &mockWalletRepo{}, sc, &mockEventPublisher{})
	err := uc.ProcessRefund(context.Background(), RefundInput{
		OrderID:     "order-r1",
		AmountCents: 2000,
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create stripe refund")
}

func TestRefund_SellerWalletDebit(t *testing.T) {
	var debitedSeller string
	var debitedAmount int64
	var createdTx *domain.WalletTransaction

	repo := &mockPaymentRepo{
		getByOrderIDFn: func(_ context.Context, _ string) (*domain.Payment, error) {
			return refundablePayment(), nil
		},
	}
	walletRepo := &mockWalletRepo{
		debitAvailableFn: func(_ context.Context, sellerID string, amount int64) error {
			debitedSeller = sellerID
			debitedAmount = amount
			return nil
		},
		createTransactionFn: func(_ context.Context, tx *domain.WalletTransaction) error {
			createdTx = tx
			return nil
		},
	}

	uc := NewRefundUseCase(repo, walletRepo, &mockStripeClient{}, &mockEventPublisher{})
	err := uc.ProcessRefund(context.Background(), RefundInput{
		OrderID:     "order-r1",
		AmountCents: 4000,
		SellerID:    "seller-refund-1",
	})

	require.NoError(t, err)
	assert.Equal(t, "seller-refund-1", debitedSeller)
	assert.Equal(t, int64(4000), debitedAmount)
	require.NotNil(t, createdTx)
	assert.Equal(t, domain.WalletTxRefundDebit, createdTx.Type)
	assert.Equal(t, int64(-4000), createdTx.AmountCents)
	assert.Equal(t, "seller-refund-1", createdTx.SellerID)
}

func TestRefund_NoSeller_SkipsWalletDebit(t *testing.T) {
	debitCalled := false

	repo := &mockPaymentRepo{
		getByOrderIDFn: func(_ context.Context, _ string) (*domain.Payment, error) {
			return refundablePayment(), nil
		},
	}
	walletRepo := &mockWalletRepo{
		debitAvailableFn: func(_ context.Context, _ string, _ int64) error {
			debitCalled = true
			return nil
		},
	}

	uc := NewRefundUseCase(repo, walletRepo, &mockStripeClient{}, &mockEventPublisher{})
	err := uc.ProcessRefund(context.Background(), RefundInput{
		OrderID:     "order-r1",
		AmountCents: 2000,
		SellerID:    "", // no seller
	})

	require.NoError(t, err)
	assert.False(t, debitCalled, "wallet debit should not be called when no seller ID")
}

func TestRefund_PaymentNotFound(t *testing.T) {
	repo := &mockPaymentRepo{
		getByOrderIDFn: func(_ context.Context, _ string) (*domain.Payment, error) {
			return nil, errors.New("not found")
		},
	}

	uc := NewRefundUseCase(repo, &mockWalletRepo{}, &mockStripeClient{}, &mockEventPublisher{})
	err := uc.ProcessRefund(context.Background(), RefundInput{
		OrderID: "order-missing",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "payment not found")
}
