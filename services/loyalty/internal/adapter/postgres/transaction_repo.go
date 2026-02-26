package postgres

import (
	"context"

	"github.com/southern-martin/ecommerce/services/loyalty/internal/domain"
	"gorm.io/gorm"
)

// TransactionRepo implements domain.PointsTransactionRepository.
type TransactionRepo struct {
	db *gorm.DB
}

// NewTransactionRepo creates a new TransactionRepo.
func NewTransactionRepo(db *gorm.DB) *TransactionRepo {
	return &TransactionRepo{db: db}
}

func (r *TransactionRepo) GetByID(ctx context.Context, id string) (*domain.PointsTransaction, error) {
	var model PointsTransactionModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *TransactionRepo) ListByUser(ctx context.Context, userID string, page, pageSize int) ([]domain.PointsTransaction, int64, error) {
	var total int64
	r.db.WithContext(ctx).Model(&PointsTransactionModel{}).Where("user_id = ?", userID).Count(&total)

	var models []PointsTransactionModel
	offset := (page - 1) * pageSize
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).
		Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	transactions := make([]domain.PointsTransaction, len(models))
	for i, m := range models {
		transactions[i] = *m.ToDomain()
	}
	return transactions, total, nil
}

func (r *TransactionRepo) Create(ctx context.Context, tx *domain.PointsTransaction) error {
	model := ToPointsTransactionModel(tx)
	return r.db.WithContext(ctx).Create(model).Error
}
