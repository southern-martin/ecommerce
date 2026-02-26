package usecase

import (
	"context"
	"fmt"

	"github.com/southern-martin/ecommerce/services/payment/internal/domain"
)

// WalletUseCase handles wallet-related operations.
type WalletUseCase struct {
	walletRepo domain.WalletRepository
}

// NewWalletUseCase creates a new WalletUseCase.
func NewWalletUseCase(walletRepo domain.WalletRepository) *WalletUseCase {
	return &WalletUseCase{walletRepo: walletRepo}
}

// GetBalance returns the wallet balance for a seller.
func (uc *WalletUseCase) GetBalance(ctx context.Context, sellerID string) (*domain.SellerWallet, error) {
	wallet, err := uc.walletRepo.GetOrCreate(ctx, sellerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}
	return wallet, nil
}

// ListTransactions returns wallet transactions for a seller.
func (uc *WalletUseCase) ListTransactions(ctx context.Context, sellerID string, page, pageSize int) ([]*domain.WalletTransaction, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	transactions, total, err := uc.walletRepo.ListTransactions(ctx, sellerID, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list transactions: %w", err)
	}
	return transactions, total, nil
}

// CreditPending adds to a seller's pending balance.
func (uc *WalletUseCase) CreditPending(ctx context.Context, sellerID string, amountCents int64) error {
	return uc.walletRepo.CreditPending(ctx, sellerID, amountCents)
}

// MovePendingToAvailable moves funds from pending to available balance.
func (uc *WalletUseCase) MovePendingToAvailable(ctx context.Context, sellerID string, amountCents int64) error {
	return uc.walletRepo.MovePendingToAvailable(ctx, sellerID, amountCents)
}

// DebitAvailable deducts from a seller's available balance.
func (uc *WalletUseCase) DebitAvailable(ctx context.Context, sellerID string, amountCents int64) error {
	return uc.walletRepo.DebitAvailable(ctx, sellerID, amountCents)
}
