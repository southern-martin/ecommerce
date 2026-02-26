package postgres

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/southern-martin/ecommerce/services/payment/internal/domain"
)

// WalletRepo implements domain.WalletRepository using PostgreSQL via GORM.
type WalletRepo struct {
	db *gorm.DB
}

// NewWalletRepo creates a new WalletRepo.
func NewWalletRepo(db *gorm.DB) *WalletRepo {
	return &WalletRepo{db: db}
}

// GetOrCreate retrieves a seller's wallet, creating one if it doesn't exist.
func (r *WalletRepo) GetOrCreate(ctx context.Context, sellerID string) (*domain.SellerWallet, error) {
	var model SellerWalletModel

	// Try to find existing wallet.
	err := r.db.WithContext(ctx).Where("seller_id = ?", sellerID).First(&model).Error
	if err == nil {
		return model.ToDomain(), nil
	}

	if err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}

	// Create new wallet.
	model = SellerWalletModel{
		SellerID:         sellerID,
		AvailableBalance: 0,
		PendingBalance:   0,
		Currency:         "usd",
		UpdatedAt:        time.Now(),
	}

	if err := r.db.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(&model).Error; err != nil {
		return nil, fmt.Errorf("failed to create wallet: %w", err)
	}

	// Re-fetch to handle race condition.
	if err := r.db.WithContext(ctx).Where("seller_id = ?", sellerID).First(&model).Error; err != nil {
		return nil, fmt.Errorf("failed to get wallet after create: %w", err)
	}

	return model.ToDomain(), nil
}

// CreditPending adds to a seller's pending balance.
func (r *WalletRepo) CreditPending(ctx context.Context, sellerID string, amountCents int64) error {
	// Ensure wallet exists.
	if _, err := r.GetOrCreate(ctx, sellerID); err != nil {
		return err
	}

	result := r.db.WithContext(ctx).Model(&SellerWalletModel{}).
		Where("seller_id = ?", sellerID).
		Updates(map[string]interface{}{
			"pending_balance": gorm.Expr("pending_balance + ?", amountCents),
			"updated_at":      time.Now(),
		})
	if result.Error != nil {
		return fmt.Errorf("failed to credit pending balance: %w", result.Error)
	}
	return nil
}

// MovePendingToAvailable moves funds from pending to available balance.
func (r *WalletRepo) MovePendingToAvailable(ctx context.Context, sellerID string, amountCents int64) error {
	result := r.db.WithContext(ctx).Model(&SellerWalletModel{}).
		Where("seller_id = ? AND pending_balance >= ?", sellerID, amountCents).
		Updates(map[string]interface{}{
			"pending_balance":   gorm.Expr("pending_balance - ?", amountCents),
			"available_balance": gorm.Expr("available_balance + ?", amountCents),
			"updated_at":        time.Now(),
		})
	if result.Error != nil {
		return fmt.Errorf("failed to move pending to available: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("insufficient pending balance for seller %s", sellerID)
	}
	return nil
}

// DebitAvailable deducts from a seller's available balance.
func (r *WalletRepo) DebitAvailable(ctx context.Context, sellerID string, amountCents int64) error {
	result := r.db.WithContext(ctx).Model(&SellerWalletModel{}).
		Where("seller_id = ? AND available_balance >= ?", sellerID, amountCents).
		Updates(map[string]interface{}{
			"available_balance": gorm.Expr("available_balance - ?", amountCents),
			"updated_at":        time.Now(),
		})
	if result.Error != nil {
		return fmt.Errorf("failed to debit available balance: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("insufficient available balance for seller %s", sellerID)
	}
	return nil
}

// CreateTransaction records a wallet transaction.
func (r *WalletRepo) CreateTransaction(ctx context.Context, tx *domain.WalletTransaction) error {
	model := WalletTransactionModelFromDomain(tx)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return fmt.Errorf("failed to create wallet transaction: %w", err)
	}
	return nil
}

// ListTransactions lists wallet transactions for a seller with pagination.
func (r *WalletRepo) ListTransactions(ctx context.Context, sellerID string, page, pageSize int) ([]*domain.WalletTransaction, int64, error) {
	var total int64
	query := r.db.WithContext(ctx).Model(&WalletTransactionModel{}).Where("seller_id = ?", sellerID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count transactions: %w", err)
	}

	var models []WalletTransactionModel
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&models).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list transactions: %w", err)
	}

	transactions := make([]*domain.WalletTransaction, len(models))
	for i, m := range models {
		transactions[i] = m.ToDomain()
	}
	return transactions, total, nil
}

// Ensure WalletRepo implements domain.WalletRepository.
var _ domain.WalletRepository = (*WalletRepo)(nil)
