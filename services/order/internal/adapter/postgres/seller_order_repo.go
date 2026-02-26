package postgres

import (
	"context"
	"errors"

	"github.com/southern-martin/ecommerce/services/order/internal/domain"
	"gorm.io/gorm"
)

// SellerOrderRepo implements domain.SellerOrderRepository using GORM/Postgres.
type SellerOrderRepo struct {
	db *gorm.DB
}

// NewSellerOrderRepo creates a new SellerOrderRepo.
func NewSellerOrderRepo(db *gorm.DB) *SellerOrderRepo {
	return &SellerOrderRepo{db: db}
}

// Create persists a new seller order.
func (r *SellerOrderRepo) Create(ctx context.Context, sellerOrder *domain.SellerOrder) error {
	model := ToSellerOrderModel(sellerOrder)
	return r.db.WithContext(ctx).Create(model).Error
}

// GetByID retrieves a seller order by its UUID.
func (r *SellerOrderRepo) GetByID(ctx context.Context, id string) (*domain.SellerOrder, error) {
	var model SellerOrderModel
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("seller order not found")
		}
		return nil, err
	}
	return model.ToDomain(), nil
}

// ListByOrder retrieves all seller orders for a given order.
func (r *SellerOrderRepo) ListByOrder(ctx context.Context, orderID string) ([]*domain.SellerOrder, error) {
	var models []SellerOrderModel
	err := r.db.WithContext(ctx).
		Where("order_id = ?", orderID).
		Order("created_at ASC").
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	var sellerOrders []*domain.SellerOrder
	for i := range models {
		sellerOrders = append(sellerOrders, models[i].ToDomain())
	}
	return sellerOrders, nil
}

// ListBySeller retrieves a paginated list of seller orders for a given seller.
func (r *SellerOrderRepo) ListBySeller(ctx context.Context, sellerID string, page, pageSize int) ([]*domain.SellerOrder, int64, error) {
	var models []SellerOrderModel
	var total int64

	query := r.db.WithContext(ctx).Model(&SellerOrderModel{}).Where("seller_id = ?", sellerID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&models).Error
	if err != nil {
		return nil, 0, err
	}

	var sellerOrders []*domain.SellerOrder
	for i := range models {
		sellerOrders = append(sellerOrders, models[i].ToDomain())
	}
	return sellerOrders, total, nil
}

// UpdateStatus updates only the status field of a seller order.
func (r *SellerOrderRepo) UpdateStatus(ctx context.Context, id string, status domain.OrderStatus) error {
	result := r.db.WithContext(ctx).
		Model(&SellerOrderModel{}).
		Where("id = ?", id).
		Update("status", string(status))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("seller order not found")
	}
	return nil
}
