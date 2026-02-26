package postgres

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/southern-martin/ecommerce/services/payment/internal/domain"
)

// PayoutRepo implements domain.PayoutRepository using PostgreSQL via GORM.
type PayoutRepo struct {
	db *gorm.DB
}

// NewPayoutRepo creates a new PayoutRepo.
func NewPayoutRepo(db *gorm.DB) *PayoutRepo {
	return &PayoutRepo{db: db}
}

// Create persists a new payout record.
func (r *PayoutRepo) Create(ctx context.Context, payout *domain.Payout) error {
	model := PayoutModelFromDomain(payout)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return fmt.Errorf("failed to create payout: %w", err)
	}
	return nil
}

// GetByID retrieves a payout by its ID.
func (r *PayoutRepo) GetByID(ctx context.Context, id string) (*domain.Payout, error) {
	var model PayoutModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		return nil, fmt.Errorf("payout not found: %w", err)
	}
	return model.ToDomain(), nil
}

// ListBySeller retrieves a paginated list of payouts for a seller.
func (r *PayoutRepo) ListBySeller(ctx context.Context, sellerID string, page, pageSize int) ([]*domain.Payout, int64, error) {
	var total int64
	query := r.db.WithContext(ctx).Model(&PayoutModel{}).Where("seller_id = ?", sellerID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count payouts: %w", err)
	}

	var models []PayoutModel
	offset := (page - 1) * pageSize
	if err := query.Order("requested_at DESC").Offset(offset).Limit(pageSize).Find(&models).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list payouts: %w", err)
	}

	payouts := make([]*domain.Payout, len(models))
	for i, m := range models {
		payouts[i] = m.ToDomain()
	}
	return payouts, total, nil
}

// UpdateStatus updates the status of a payout.
func (r *PayoutRepo) UpdateStatus(ctx context.Context, id string, status domain.PayoutStatus) error {
	result := r.db.WithContext(ctx).Model(&PayoutModel{}).Where("id = ?", id).Update("status", string(status))
	if result.Error != nil {
		return fmt.Errorf("failed to update payout status: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("payout %s not found", id)
	}
	return nil
}

// Ensure PayoutRepo implements domain.PayoutRepository.
var _ domain.PayoutRepository = (*PayoutRepo)(nil)
