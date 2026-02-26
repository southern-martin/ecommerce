package postgres

import (
	"context"

	"github.com/southern-martin/ecommerce/services/return/internal/domain"
	"gorm.io/gorm"
)

// DisputeMessageRepo implements domain.DisputeMessageRepository.
type DisputeMessageRepo struct {
	db *gorm.DB
}

// NewDisputeMessageRepo creates a new DisputeMessageRepo.
func NewDisputeMessageRepo(db *gorm.DB) *DisputeMessageRepo {
	return &DisputeMessageRepo{db: db}
}

func (r *DisputeMessageRepo) GetByDisputeID(ctx context.Context, disputeID string) ([]domain.DisputeMessage, error) {
	var models []DisputeMessageModel
	if err := r.db.WithContext(ctx).Where("dispute_id = ?", disputeID).
		Order("created_at ASC").Find(&models).Error; err != nil {
		return nil, err
	}
	messages := make([]domain.DisputeMessage, len(models))
	for i, m := range models {
		messages[i] = *m.ToDomain()
	}
	return messages, nil
}

func (r *DisputeMessageRepo) Create(ctx context.Context, msg *domain.DisputeMessage) error {
	model := ToDisputeMessageModel(msg)
	return r.db.WithContext(ctx).Create(model).Error
}
