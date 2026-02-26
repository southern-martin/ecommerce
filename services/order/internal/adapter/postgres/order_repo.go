package postgres

import (
	"context"
	"errors"

	"github.com/southern-martin/ecommerce/services/order/internal/domain"
	"gorm.io/gorm"
)

// OrderRepo implements domain.OrderRepository using GORM/Postgres.
type OrderRepo struct {
	db *gorm.DB
}

// NewOrderRepo creates a new OrderRepo.
func NewOrderRepo(db *gorm.DB) *OrderRepo {
	return &OrderRepo{db: db}
}

// Create persists a new order with its items.
func (r *OrderRepo) Create(ctx context.Context, order *domain.Order) error {
	model := ToOrderModel(order)
	return r.db.WithContext(ctx).Create(model).Error
}

// GetByID retrieves an order by its UUID, including items and seller orders.
func (r *OrderRepo) GetByID(ctx context.Context, id string) (*domain.Order, error) {
	var model OrderModel
	err := r.db.WithContext(ctx).
		Preload("Items").
		Preload("SellerOrders").
		Where("id = ?", id).
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order not found")
		}
		return nil, err
	}
	return model.ToDomain(), nil
}

// GetByOrderNumber retrieves an order by its human-readable order number.
func (r *OrderRepo) GetByOrderNumber(ctx context.Context, orderNumber string) (*domain.Order, error) {
	var model OrderModel
	err := r.db.WithContext(ctx).
		Preload("Items").
		Preload("SellerOrders").
		Where("order_number = ?", orderNumber).
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order not found")
		}
		return nil, err
	}
	return model.ToDomain(), nil
}

// List retrieves a paginated list of orders based on filter criteria.
func (r *OrderRepo) List(ctx context.Context, filter domain.OrderFilter) ([]*domain.Order, int64, error) {
	var models []OrderModel
	var total int64

	query := r.db.WithContext(ctx).Model(&OrderModel{})

	if filter.BuyerID != "" {
		query = query.Where("buyer_id = ?", filter.BuyerID)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.PageSize
	err := query.
		Preload("Items").
		Preload("SellerOrders").
		Order("created_at DESC").
		Offset(offset).
		Limit(filter.PageSize).
		Find(&models).Error
	if err != nil {
		return nil, 0, err
	}

	var orders []*domain.Order
	for i := range models {
		orders = append(orders, models[i].ToDomain())
	}
	return orders, total, nil
}

// UpdateStatus updates only the status field of an order.
func (r *OrderRepo) UpdateStatus(ctx context.Context, id string, status domain.OrderStatus) error {
	result := r.db.WithContext(ctx).
		Model(&OrderModel{}).
		Where("id = ?", id).
		Update("status", string(status))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("order not found")
	}
	return nil
}

// Update persists all changes to an existing order.
func (r *OrderRepo) Update(ctx context.Context, order *domain.Order) error {
	model := ToOrderModel(order)
	return r.db.WithContext(ctx).Save(model).Error
}
