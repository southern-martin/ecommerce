package domain

import "context"

// PaymentRepository defines the interface for payment persistence.
type PaymentRepository interface {
	Create(ctx context.Context, payment *Payment) error
	GetByID(ctx context.Context, id string) (*Payment, error)
	GetByOrderID(ctx context.Context, orderID string) (*Payment, error)
	GetByStripeID(ctx context.Context, stripePaymentID string) (*Payment, error)
	UpdateStatus(ctx context.Context, id string, status PaymentStatus, failureReason string) error
	List(ctx context.Context, buyerID string, page, pageSize int) ([]*Payment, int64, error)
}

// WalletRepository defines the interface for seller wallet persistence.
type WalletRepository interface {
	GetOrCreate(ctx context.Context, sellerID string) (*SellerWallet, error)
	CreditPending(ctx context.Context, sellerID string, amountCents int64) error
	MovePendingToAvailable(ctx context.Context, sellerID string, amountCents int64) error
	DebitAvailable(ctx context.Context, sellerID string, amountCents int64) error
	CreateTransaction(ctx context.Context, tx *WalletTransaction) error
	ListTransactions(ctx context.Context, sellerID string, page, pageSize int) ([]*WalletTransaction, int64, error)
}

// PayoutRepository defines the interface for payout persistence.
type PayoutRepository interface {
	Create(ctx context.Context, payout *Payout) error
	GetByID(ctx context.Context, id string) (*Payout, error)
	ListBySeller(ctx context.Context, sellerID string, page, pageSize int) ([]*Payout, int64, error)
	UpdateStatus(ctx context.Context, id string, status PayoutStatus) error
}
