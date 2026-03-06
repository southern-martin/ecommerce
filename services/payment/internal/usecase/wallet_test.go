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
// WalletUseCase tests
// ---------------------------------------------------------------------------

func TestWallet_GetBalance(t *testing.T) {
	walletRepo := &mockWalletRepo{
		getOrCreateFn: func(_ context.Context, sellerID string) (*domain.SellerWallet, error) {
			return &domain.SellerWallet{
				SellerID:         sellerID,
				AvailableBalance: 50000,
				PendingBalance:   12000,
				Currency:         "usd",
			}, nil
		},
	}

	uc := NewWalletUseCase(walletRepo)
	wallet, err := uc.GetBalance(context.Background(), "seller-w1")

	require.NoError(t, err)
	require.NotNil(t, wallet)
	assert.Equal(t, "seller-w1", wallet.SellerID)
	assert.Equal(t, int64(50000), wallet.AvailableBalance)
	assert.Equal(t, int64(12000), wallet.PendingBalance)
}

func TestWallet_GetBalance_Error(t *testing.T) {
	walletRepo := &mockWalletRepo{
		getOrCreateFn: func(_ context.Context, _ string) (*domain.SellerWallet, error) {
			return nil, errors.New("db unavailable")
		},
	}

	uc := NewWalletUseCase(walletRepo)
	wallet, err := uc.GetBalance(context.Background(), "seller-w2")

	require.Error(t, err)
	assert.Nil(t, wallet)
	assert.Contains(t, err.Error(), "failed to get wallet")
}

func TestWallet_ListTransactions(t *testing.T) {
	now := time.Now()
	walletRepo := &mockWalletRepo{
		listTransactionsFn: func(_ context.Context, sellerID string, page, pageSize int) ([]*domain.WalletTransaction, int64, error) {
			assert.Equal(t, "seller-w3", sellerID)
			assert.Equal(t, 1, page)
			assert.Equal(t, 20, pageSize)
			return []*domain.WalletTransaction{
				{ID: "tx-1", SellerID: sellerID, Type: domain.WalletTxSale, AmountCents: 5000, CreatedAt: now},
				{ID: "tx-2", SellerID: sellerID, Type: domain.WalletTxPayout, AmountCents: -3000, CreatedAt: now},
			}, 2, nil
		},
	}

	uc := NewWalletUseCase(walletRepo)
	txns, total, err := uc.ListTransactions(context.Background(), "seller-w3", 1, 20)

	require.NoError(t, err)
	assert.Equal(t, int64(2), total)
	require.Len(t, txns, 2)
	assert.Equal(t, "tx-1", txns[0].ID)
	assert.Equal(t, "tx-2", txns[1].ID)
}

func TestWallet_CreditPending(t *testing.T) {
	var creditedSeller string
	var creditedAmount int64

	walletRepo := &mockWalletRepo{
		creditPendingFn: func(_ context.Context, sellerID string, amount int64) error {
			creditedSeller = sellerID
			creditedAmount = amount
			return nil
		},
	}

	uc := NewWalletUseCase(walletRepo)
	err := uc.CreditPending(context.Background(), "seller-w4", 7500)

	require.NoError(t, err)
	assert.Equal(t, "seller-w4", creditedSeller)
	assert.Equal(t, int64(7500), creditedAmount)
}

func TestWallet_DebitAvailable(t *testing.T) {
	var debitedSeller string
	var debitedAmount int64

	walletRepo := &mockWalletRepo{
		debitAvailableFn: func(_ context.Context, sellerID string, amount int64) error {
			debitedSeller = sellerID
			debitedAmount = amount
			return nil
		},
	}

	uc := NewWalletUseCase(walletRepo)
	err := uc.DebitAvailable(context.Background(), "seller-w5", 3000)

	require.NoError(t, err)
	assert.Equal(t, "seller-w5", debitedSeller)
	assert.Equal(t, int64(3000), debitedAmount)
}

func TestWallet_MovePendingToAvailable(t *testing.T) {
	var movedSeller string
	var movedAmount int64

	walletRepo := &mockWalletRepo{
		movePendingToAvailFn: func(_ context.Context, sellerID string, amount int64) error {
			movedSeller = sellerID
			movedAmount = amount
			return nil
		},
	}

	uc := NewWalletUseCase(walletRepo)
	err := uc.MovePendingToAvailable(context.Background(), "seller-w6", 10000)

	require.NoError(t, err)
	assert.Equal(t, "seller-w6", movedSeller)
	assert.Equal(t, int64(10000), movedAmount)
}

func TestWallet_ListTransactions_PaginationDefaults(t *testing.T) {
	var capturedPage, capturedPageSize int

	walletRepo := &mockWalletRepo{
		listTransactionsFn: func(_ context.Context, _ string, page, pageSize int) ([]*domain.WalletTransaction, int64, error) {
			capturedPage = page
			capturedPageSize = pageSize
			return nil, 0, nil
		},
	}

	uc := NewWalletUseCase(walletRepo)

	// page < 1 should default to 1
	_, _, err := uc.ListTransactions(context.Background(), "seller-x", 0, 10)
	require.NoError(t, err)
	assert.Equal(t, 1, capturedPage)
	assert.Equal(t, 10, capturedPageSize)

	// pageSize > 100 should default to 20
	_, _, err = uc.ListTransactions(context.Background(), "seller-x", 2, 200)
	require.NoError(t, err)
	assert.Equal(t, 2, capturedPage)
	assert.Equal(t, 20, capturedPageSize)

	// pageSize < 1 should default to 20
	_, _, err = uc.ListTransactions(context.Background(), "seller-x", 1, 0)
	require.NoError(t, err)
	assert.Equal(t, 1, capturedPage)
	assert.Equal(t, 20, capturedPageSize)
}
