package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/southern-martin/ecommerce/services/payment/internal/domain"
)

// ---------------------------------------------------------------------------
// PayoutUseCase tests
// ---------------------------------------------------------------------------

func TestRequestPayout_Success(t *testing.T) {
	var debitedAmount int64
	var createdPayout *domain.Payout
	var createdTx *domain.WalletTransaction

	walletRepo := &mockWalletRepo{
		getOrCreateFn: func(_ context.Context, _ string) (*domain.SellerWallet, error) {
			return &domain.SellerWallet{
				SellerID:         "seller-p1",
				AvailableBalance: 50000,
				Currency:         "usd",
			}, nil
		},
		debitAvailableFn: func(_ context.Context, _ string, amount int64) error {
			debitedAmount = amount
			return nil
		},
		createTransactionFn: func(_ context.Context, tx *domain.WalletTransaction) error {
			createdTx = tx
			return nil
		},
	}
	payoutRepo := &mockPayoutRepo{
		createFn: func(_ context.Context, p *domain.Payout) error {
			createdPayout = p
			return nil
		},
	}

	uc := NewPayoutUseCase(payoutRepo, walletRepo, &mockStripeClient{})
	payout, err := uc.RequestPayout(context.Background(), RequestPayoutInput{
		SellerID:    "seller-p1",
		AmountCents: 20000,
		Currency:    "usd",
		Method:      "stripe_connect",
	})

	require.NoError(t, err)
	require.NotNil(t, payout)
	assert.NotEmpty(t, payout.ID)
	assert.Equal(t, "seller-p1", payout.SellerID)
	assert.Equal(t, int64(20000), payout.AmountCents)
	assert.Equal(t, "usd", payout.Currency)
	assert.Equal(t, "stripe_connect", payout.Method)
	assert.Equal(t, domain.PayoutStatusRequested, payout.Status)

	assert.Equal(t, int64(20000), debitedAmount)
	require.NotNil(t, createdPayout)
	require.NotNil(t, createdTx)
	assert.Equal(t, domain.WalletTxPayout, createdTx.Type)
	assert.Equal(t, int64(-20000), createdTx.AmountCents)
}

func TestRequestPayout_InsufficientBalance(t *testing.T) {
	walletRepo := &mockWalletRepo{
		getOrCreateFn: func(_ context.Context, _ string) (*domain.SellerWallet, error) {
			return &domain.SellerWallet{
				SellerID:         "seller-p2",
				AvailableBalance: 5000, // only 5000 available
			}, nil
		},
	}

	uc := NewPayoutUseCase(&mockPayoutRepo{}, walletRepo, &mockStripeClient{})
	payout, err := uc.RequestPayout(context.Background(), RequestPayoutInput{
		SellerID:    "seller-p2",
		AmountCents: 10000, // needs 10000
	})

	require.Error(t, err)
	assert.Nil(t, payout)
	assert.Contains(t, err.Error(), "insufficient available balance")
}

func TestRequestPayout_DefaultCurrencyAndMethod(t *testing.T) {
	walletRepo := &mockWalletRepo{
		getOrCreateFn: func(_ context.Context, _ string) (*domain.SellerWallet, error) {
			return &domain.SellerWallet{
				SellerID:         "seller-p3",
				AvailableBalance: 100000,
			}, nil
		},
	}
	var createdPayout *domain.Payout
	payoutRepo := &mockPayoutRepo{
		createFn: func(_ context.Context, p *domain.Payout) error {
			createdPayout = p
			return nil
		},
	}

	uc := NewPayoutUseCase(payoutRepo, walletRepo, &mockStripeClient{})
	payout, err := uc.RequestPayout(context.Background(), RequestPayoutInput{
		SellerID:    "seller-p3",
		AmountCents: 5000,
		Currency:    "", // should default to "usd"
		Method:      "", // should default to "stripe_connect"
	})

	require.NoError(t, err)
	require.NotNil(t, payout)
	assert.Equal(t, "usd", payout.Currency)
	assert.Equal(t, "stripe_connect", payout.Method)
	require.NotNil(t, createdPayout)
	assert.Equal(t, "usd", createdPayout.Currency)
	assert.Equal(t, "stripe_connect", createdPayout.Method)
}

func TestRequestPayout_WalletGetError(t *testing.T) {
	walletRepo := &mockWalletRepo{
		getOrCreateFn: func(_ context.Context, _ string) (*domain.SellerWallet, error) {
			return nil, errors.New("wallet service down")
		},
	}

	uc := NewPayoutUseCase(&mockPayoutRepo{}, walletRepo, &mockStripeClient{})
	payout, err := uc.RequestPayout(context.Background(), RequestPayoutInput{
		SellerID:    "seller-p4",
		AmountCents: 1000,
	})

	require.Error(t, err)
	assert.Nil(t, payout)
	assert.Contains(t, err.Error(), "failed to get wallet")
}

func TestRequestPayout_PayoutRepoError(t *testing.T) {
	walletRepo := &mockWalletRepo{
		getOrCreateFn: func(_ context.Context, _ string) (*domain.SellerWallet, error) {
			return &domain.SellerWallet{AvailableBalance: 99999}, nil
		},
	}
	payoutRepo := &mockPayoutRepo{
		createFn: func(_ context.Context, _ *domain.Payout) error {
			return errors.New("payout db error")
		},
	}

	uc := NewPayoutUseCase(payoutRepo, walletRepo, &mockStripeClient{})
	payout, err := uc.RequestPayout(context.Background(), RequestPayoutInput{
		SellerID:    "seller-p5",
		AmountCents: 1000,
	})

	require.Error(t, err)
	assert.Nil(t, payout)
	assert.Contains(t, err.Error(), "failed to create payout")
}

func TestProcessPayout_Success(t *testing.T) {
	var statusUpdates []domain.PayoutStatus

	payoutRepo := &mockPayoutRepo{
		getByIDFn: func(_ context.Context, _ string) (*domain.Payout, error) {
			return &domain.Payout{
				ID:          "payout-proc-1",
				SellerID:    "seller-proc-1",
				AmountCents: 15000,
				Status:      domain.PayoutStatusRequested,
				RequestedAt: time.Now(),
			}, nil
		},
		updateStatusFn: func(_ context.Context, _ string, status domain.PayoutStatus) error {
			statusUpdates = append(statusUpdates, status)
			return nil
		},
	}
	sc := &mockStripeClient{
		createTransferFn: func(amount int64, dest string, meta map[string]string) (string, error) {
			assert.Equal(t, int64(15000), amount)
			assert.Equal(t, "seller-proc-1", dest)
			return "tr_success_1", nil
		},
	}

	uc := NewPayoutUseCase(payoutRepo, &mockWalletRepo{}, sc)
	err := uc.ProcessPayout(context.Background(), "payout-proc-1")

	require.NoError(t, err)
	require.Len(t, statusUpdates, 2)
	assert.Equal(t, domain.PayoutStatusProcessing, statusUpdates[0])
	assert.Equal(t, domain.PayoutStatusCompleted, statusUpdates[1])
}

func TestProcessPayout_StripeTransferFails(t *testing.T) {
	var statusUpdates []domain.PayoutStatus

	payoutRepo := &mockPayoutRepo{
		getByIDFn: func(_ context.Context, _ string) (*domain.Payout, error) {
			return &domain.Payout{
				ID:          "payout-fail-1",
				SellerID:    "seller-fail-1",
				AmountCents: 5000,
				Status:      domain.PayoutStatusRequested,
			}, nil
		},
		updateStatusFn: func(_ context.Context, _ string, status domain.PayoutStatus) error {
			statusUpdates = append(statusUpdates, status)
			return nil
		},
	}
	sc := &mockStripeClient{
		createTransferFn: func(_ int64, _ string, _ map[string]string) (string, error) {
			return "", errors.New("stripe transfer error")
		},
	}

	uc := NewPayoutUseCase(payoutRepo, &mockWalletRepo{}, sc)
	err := uc.ProcessPayout(context.Background(), "payout-fail-1")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create stripe transfer")
	require.Len(t, statusUpdates, 2)
	assert.Equal(t, domain.PayoutStatusProcessing, statusUpdates[0])
	assert.Equal(t, domain.PayoutStatusFailed, statusUpdates[1])
}

func TestProcessPayout_NotRequested(t *testing.T) {
	payoutRepo := &mockPayoutRepo{
		getByIDFn: func(_ context.Context, _ string) (*domain.Payout, error) {
			return &domain.Payout{
				ID:     "payout-wrong-1",
				Status: domain.PayoutStatusCompleted, // already completed
			}, nil
		},
	}

	uc := NewPayoutUseCase(payoutRepo, &mockWalletRepo{}, &mockStripeClient{})
	err := uc.ProcessPayout(context.Background(), "payout-wrong-1")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "not in requested status")
}

func TestProcessPayout_NotFound(t *testing.T) {
	payoutRepo := &mockPayoutRepo{
		getByIDFn: func(_ context.Context, _ string) (*domain.Payout, error) {
			return nil, errors.New("record not found")
		},
	}

	uc := NewPayoutUseCase(payoutRepo, &mockWalletRepo{}, &mockStripeClient{})
	err := uc.ProcessPayout(context.Background(), "payout-missing")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "payout not found")
}

func TestListPayouts_Success(t *testing.T) {
	now := time.Now()
	payoutRepo := &mockPayoutRepo{
		listBySellerFn: func(_ context.Context, sellerID string, page, pageSize int) ([]*domain.Payout, int64, error) {
			return []*domain.Payout{
				{ID: "po-1", SellerID: sellerID, AmountCents: 10000, Status: domain.PayoutStatusCompleted, RequestedAt: now},
				{ID: "po-2", SellerID: sellerID, AmountCents: 5000, Status: domain.PayoutStatusRequested, RequestedAt: now},
			}, 2, nil
		},
	}

	uc := NewPayoutUseCase(payoutRepo, &mockWalletRepo{}, &mockStripeClient{})
	payouts, total, err := uc.ListPayouts(context.Background(), "seller-list-1", 1, 10)

	require.NoError(t, err)
	assert.Equal(t, int64(2), total)
	require.Len(t, payouts, 2)
	assert.Equal(t, "po-1", payouts[0].ID)
	assert.Equal(t, "po-2", payouts[1].ID)
}

func TestListPayouts_PaginationDefaults(t *testing.T) {
	var capturedPage, capturedPageSize int

	payoutRepo := &mockPayoutRepo{
		listBySellerFn: func(_ context.Context, _ string, page, pageSize int) ([]*domain.Payout, int64, error) {
			capturedPage = page
			capturedPageSize = pageSize
			return nil, 0, nil
		},
	}

	uc := NewPayoutUseCase(payoutRepo, &mockWalletRepo{}, &mockStripeClient{})

	// page < 1 defaults to 1
	_, _, err := uc.ListPayouts(context.Background(), "s1", 0, 10)
	require.NoError(t, err)
	assert.Equal(t, 1, capturedPage)

	// pageSize > 100 defaults to 20
	_, _, err = uc.ListPayouts(context.Background(), "s1", 1, 150)
	require.NoError(t, err)
	assert.Equal(t, 20, capturedPageSize)

	// pageSize < 1 defaults to 20
	_, _, err = uc.ListPayouts(context.Background(), "s1", 1, -5)
	require.NoError(t, err)
	assert.Equal(t, 20, capturedPageSize)
}
