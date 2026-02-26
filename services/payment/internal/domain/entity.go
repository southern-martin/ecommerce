package domain

import "time"

// PaymentStatus represents the status of a payment.
type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusCompleted PaymentStatus = "completed"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusRefunded  PaymentStatus = "refunded"
)

// PaymentMethod represents the method of payment.
type PaymentMethod string

const (
	PaymentMethodCard   PaymentMethod = "card"
	PaymentMethodWallet PaymentMethod = "wallet"
)

// Payment represents a payment transaction.
type Payment struct {
	ID              string
	OrderID         string
	BuyerID         string
	AmountCents     int64
	Currency        string
	Status          PaymentStatus
	Method          PaymentMethod
	StripePaymentID string
	FailureReason   string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// SellerWallet represents a seller's wallet balance.
type SellerWallet struct {
	SellerID         string
	AvailableBalance int64
	PendingBalance   int64
	Currency         string
	UpdatedAt        time.Time
}

// WalletTransactionType represents the type of wallet transaction.
type WalletTransactionType string

const (
	WalletTxSale               WalletTransactionType = "sale"
	WalletTxCommissionDeducted WalletTransactionType = "commission_deducted"
	WalletTxPayout             WalletTransactionType = "payout"
	WalletTxRefundDebit        WalletTransactionType = "refund_debit"
	WalletTxAdjustment         WalletTransactionType = "adjustment"
)

// WalletTransaction represents a transaction in a seller's wallet.
type WalletTransaction struct {
	ID            string
	SellerID      string
	Type          WalletTransactionType
	AmountCents   int64 // positive = credit, negative = debit
	ReferenceType string
	ReferenceID   string
	Description   string
	CreatedAt     time.Time
}

// PayoutStatus represents the status of a payout.
type PayoutStatus string

const (
	PayoutStatusRequested  PayoutStatus = "requested"
	PayoutStatusProcessing PayoutStatus = "processing"
	PayoutStatusCompleted  PayoutStatus = "completed"
	PayoutStatusFailed     PayoutStatus = "failed"
)

// Payout represents a payout request from a seller.
type Payout struct {
	ID               string
	SellerID         string
	AmountCents      int64
	Currency         string
	Method           string // "stripe_connect", "bank_transfer"
	StripeTransferID string
	Status           PayoutStatus
	RequestedAt      time.Time
	CompletedAt      *time.Time
}
