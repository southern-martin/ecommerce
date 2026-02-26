package postgres

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/southern-martin/ecommerce/services/payment/internal/domain"
)

// PaymentRepo implements domain.PaymentRepository using PostgreSQL via GORM.
type PaymentRepo struct {
	db *gorm.DB
}

// NewPaymentRepo creates a new PaymentRepo.
func NewPaymentRepo(db *gorm.DB) *PaymentRepo {
	return &PaymentRepo{db: db}
}

// Create persists a new payment record.
func (r *PaymentRepo) Create(ctx context.Context, payment *domain.Payment) error {
	model := PaymentModelFromDomain(payment)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return fmt.Errorf("failed to create payment: %w", err)
	}
	return nil
}

// GetByID retrieves a payment by its ID.
func (r *PaymentRepo) GetByID(ctx context.Context, id string) (*domain.Payment, error) {
	var model PaymentModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		return nil, fmt.Errorf("payment not found: %w", err)
	}
	return model.ToDomain(), nil
}

// GetByOrderID retrieves a payment by its order ID.
func (r *PaymentRepo) GetByOrderID(ctx context.Context, orderID string) (*domain.Payment, error) {
	var model PaymentModel
	if err := r.db.WithContext(ctx).Where("order_id = ?", orderID).First(&model).Error; err != nil {
		return nil, fmt.Errorf("payment not found for order %s: %w", orderID, err)
	}
	return model.ToDomain(), nil
}

// GetByStripeID retrieves a payment by its Stripe PaymentIntent ID.
func (r *PaymentRepo) GetByStripeID(ctx context.Context, stripePaymentID string) (*domain.Payment, error) {
	var model PaymentModel
	if err := r.db.WithContext(ctx).Where("stripe_payment_id = ?", stripePaymentID).First(&model).Error; err != nil {
		return nil, fmt.Errorf("payment not found for stripe ID %s: %w", stripePaymentID, err)
	}
	return model.ToDomain(), nil
}

// UpdateStatus updates the status and optional failure reason of a payment.
func (r *PaymentRepo) UpdateStatus(ctx context.Context, id string, status domain.PaymentStatus, failureReason string) error {
	updates := map[string]interface{}{
		"status":     string(status),
		"updated_at": time.Now(),
	}
	if failureReason != "" {
		updates["failure_reason"] = failureReason
	}

	result := r.db.WithContext(ctx).Model(&PaymentModel{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to update payment status: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("payment %s not found", id)
	}
	return nil
}

// List retrieves a paginated list of payments for a buyer.
func (r *PaymentRepo) List(ctx context.Context, buyerID string, page, pageSize int) ([]*domain.Payment, int64, error) {
	var total int64
	query := r.db.WithContext(ctx).Model(&PaymentModel{})

	if buyerID != "" {
		query = query.Where("buyer_id = ?", buyerID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count payments: %w", err)
	}

	var models []PaymentModel
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&models).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list payments: %w", err)
	}

	payments := make([]*domain.Payment, len(models))
	for i, m := range models {
		payments[i] = m.ToDomain()
	}
	return payments, total, nil
}

// Ensure PaymentRepo implements domain.PaymentRepository.
var _ domain.PaymentRepository = (*PaymentRepo)(nil)
