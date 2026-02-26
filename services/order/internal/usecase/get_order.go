package usecase

import (
	"context"
	"errors"

	"github.com/southern-martin/ecommerce/services/order/internal/domain"
)

// GetOrderUseCase handles retrieving orders.
type GetOrderUseCase struct {
	orderRepo       domain.OrderRepository
	sellerOrderRepo domain.SellerOrderRepository
}

// NewGetOrderUseCase creates a new GetOrderUseCase instance.
func NewGetOrderUseCase(
	orderRepo domain.OrderRepository,
	sellerOrderRepo domain.SellerOrderRepository,
) *GetOrderUseCase {
	return &GetOrderUseCase{
		orderRepo:       orderRepo,
		sellerOrderRepo: sellerOrderRepo,
	}
}

// GetOrder retrieves a single order by ID.
func (uc *GetOrderUseCase) GetOrder(ctx context.Context, id string) (*domain.Order, error) {
	if id == "" {
		return nil, errors.New("order id is required")
	}
	return uc.orderRepo.GetByID(ctx, id)
}

// GetOrderByNumber retrieves a single order by order number.
func (uc *GetOrderUseCase) GetOrderByNumber(ctx context.Context, orderNumber string) (*domain.Order, error) {
	if orderNumber == "" {
		return nil, errors.New("order number is required")
	}
	return uc.orderRepo.GetByOrderNumber(ctx, orderNumber)
}

// ListOrders retrieves a paginated list of orders based on filter criteria.
func (uc *GetOrderUseCase) ListOrders(ctx context.Context, filter domain.OrderFilter) ([]*domain.Order, int64, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}
	if filter.PageSize > 100 {
		filter.PageSize = 100
	}
	return uc.orderRepo.List(ctx, filter)
}

// GetSellerOrder retrieves a single seller order by ID.
func (uc *GetOrderUseCase) GetSellerOrder(ctx context.Context, id string) (*domain.SellerOrder, error) {
	if id == "" {
		return nil, errors.New("seller order id is required")
	}
	return uc.sellerOrderRepo.GetByID(ctx, id)
}

// ListSellerOrders retrieves a paginated list of seller orders.
func (uc *GetOrderUseCase) ListSellerOrders(ctx context.Context, sellerID string, page, pageSize int) ([]*domain.SellerOrder, int64, error) {
	if sellerID == "" {
		return nil, 0, errors.New("seller id is required")
	}
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return uc.sellerOrderRepo.ListBySeller(ctx, sellerID, page, pageSize)
}
