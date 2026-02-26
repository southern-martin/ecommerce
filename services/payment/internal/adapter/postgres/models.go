package postgres

import (
	"time"

	"github.com/southern-martin/ecommerce/services/payment/internal/domain"
)

// PaymentModel is the GORM model for the payments table.
type PaymentModel struct {
	ID              string `gorm:"type:varchar(36);primaryKey"`
	OrderID         string `gorm:"type:varchar(36);index;not null"`
	BuyerID         string `gorm:"type:varchar(36);index;not null"`
	AmountCents     int64  `gorm:"not null"`
	Currency        string `gorm:"type:varchar(3);not null;default:'usd'"`
	Status          string `gorm:"type:varchar(20);not null;default:'pending';index"`
	Method          string `gorm:"type:varchar(20);not null;default:'card'"`
	StripePaymentID string `gorm:"type:varchar(255);index"`
	FailureReason   string `gorm:"type:text"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// TableName returns the table name for PaymentModel.
func (PaymentModel) TableName() string {
	return "payments"
}

// ToDomain converts the GORM model to a domain entity.
func (m *PaymentModel) ToDomain() *domain.Payment {
	return &domain.Payment{
		ID:              m.ID,
		OrderID:         m.OrderID,
		BuyerID:         m.BuyerID,
		AmountCents:     m.AmountCents,
		Currency:        m.Currency,
		Status:          domain.PaymentStatus(m.Status),
		Method:          domain.PaymentMethod(m.Method),
		StripePaymentID: m.StripePaymentID,
		FailureReason:   m.FailureReason,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}

// PaymentModelFromDomain creates a GORM model from a domain entity.
func PaymentModelFromDomain(p *domain.Payment) *PaymentModel {
	return &PaymentModel{
		ID:              p.ID,
		OrderID:         p.OrderID,
		BuyerID:         p.BuyerID,
		AmountCents:     p.AmountCents,
		Currency:        p.Currency,
		Status:          string(p.Status),
		Method:          string(p.Method),
		StripePaymentID: p.StripePaymentID,
		FailureReason:   p.FailureReason,
		CreatedAt:       p.CreatedAt,
		UpdatedAt:       p.UpdatedAt,
	}
}

// SellerWalletModel is the GORM model for the seller_wallets table.
type SellerWalletModel struct {
	SellerID         string `gorm:"type:varchar(36);primaryKey"`
	AvailableBalance int64  `gorm:"not null;default:0"`
	PendingBalance   int64  `gorm:"not null;default:0"`
	Currency         string `gorm:"type:varchar(3);not null;default:'usd'"`
	UpdatedAt        time.Time
}

// TableName returns the table name for SellerWalletModel.
func (SellerWalletModel) TableName() string {
	return "seller_wallets"
}

// ToDomain converts the GORM model to a domain entity.
func (m *SellerWalletModel) ToDomain() *domain.SellerWallet {
	return &domain.SellerWallet{
		SellerID:         m.SellerID,
		AvailableBalance: m.AvailableBalance,
		PendingBalance:   m.PendingBalance,
		Currency:         m.Currency,
		UpdatedAt:        m.UpdatedAt,
	}
}

// WalletTransactionModel is the GORM model for the wallet_transactions table.
type WalletTransactionModel struct {
	ID            string `gorm:"type:varchar(36);primaryKey"`
	SellerID      string `gorm:"type:varchar(36);index;not null"`
	Type          string `gorm:"type:varchar(30);not null"`
	AmountCents   int64  `gorm:"not null"`
	ReferenceType string `gorm:"type:varchar(20)"`
	ReferenceID   string `gorm:"type:varchar(36);index"`
	Description   string `gorm:"type:text"`
	CreatedAt     time.Time
}

// TableName returns the table name for WalletTransactionModel.
func (WalletTransactionModel) TableName() string {
	return "wallet_transactions"
}

// ToDomain converts the GORM model to a domain entity.
func (m *WalletTransactionModel) ToDomain() *domain.WalletTransaction {
	return &domain.WalletTransaction{
		ID:            m.ID,
		SellerID:      m.SellerID,
		Type:          domain.WalletTransactionType(m.Type),
		AmountCents:   m.AmountCents,
		ReferenceType: m.ReferenceType,
		ReferenceID:   m.ReferenceID,
		Description:   m.Description,
		CreatedAt:     m.CreatedAt,
	}
}

// WalletTransactionModelFromDomain creates a GORM model from a domain entity.
func WalletTransactionModelFromDomain(tx *domain.WalletTransaction) *WalletTransactionModel {
	return &WalletTransactionModel{
		ID:            tx.ID,
		SellerID:      tx.SellerID,
		Type:          string(tx.Type),
		AmountCents:   tx.AmountCents,
		ReferenceType: tx.ReferenceType,
		ReferenceID:   tx.ReferenceID,
		Description:   tx.Description,
		CreatedAt:     tx.CreatedAt,
	}
}

// PayoutModel is the GORM model for the payouts table.
type PayoutModel struct {
	ID               string `gorm:"type:varchar(36);primaryKey"`
	SellerID         string `gorm:"type:varchar(36);index;not null"`
	AmountCents      int64  `gorm:"not null"`
	Currency         string `gorm:"type:varchar(3);not null;default:'usd'"`
	Method           string `gorm:"type:varchar(30);not null"`
	StripeTransferID string `gorm:"type:varchar(255)"`
	Status           string `gorm:"type:varchar(20);not null;default:'requested';index"`
	RequestedAt      time.Time
	CompletedAt      *time.Time
}

// TableName returns the table name for PayoutModel.
func (PayoutModel) TableName() string {
	return "payouts"
}

// ToDomain converts the GORM model to a domain entity.
func (m *PayoutModel) ToDomain() *domain.Payout {
	return &domain.Payout{
		ID:               m.ID,
		SellerID:         m.SellerID,
		AmountCents:      m.AmountCents,
		Currency:         m.Currency,
		Method:           m.Method,
		StripeTransferID: m.StripeTransferID,
		Status:           domain.PayoutStatus(m.Status),
		RequestedAt:      m.RequestedAt,
		CompletedAt:      m.CompletedAt,
	}
}

// PayoutModelFromDomain creates a GORM model from a domain entity.
func PayoutModelFromDomain(p *domain.Payout) *PayoutModel {
	return &PayoutModel{
		ID:               p.ID,
		SellerID:         p.SellerID,
		AmountCents:      p.AmountCents,
		Currency:         p.Currency,
		Method:           p.Method,
		StripeTransferID: p.StripeTransferID,
		Status:           string(p.Status),
		RequestedAt:      p.RequestedAt,
		CompletedAt:      p.CompletedAt,
	}
}
