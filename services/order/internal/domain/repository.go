package domain

import "context"

// OrderFilter provides filtering and pagination for order queries.
type OrderFilter struct {
	BuyerID  string
	SellerID string
	Status   string
	Page     int
	PageSize int
}

// OrderRepository defines the interface for order persistence.
type OrderRepository interface {
	Create(ctx context.Context, order *Order) error
	GetByID(ctx context.Context, id string) (*Order, error)
	GetByOrderNumber(ctx context.Context, orderNumber string) (*Order, error)
	List(ctx context.Context, filter OrderFilter) ([]*Order, int64, error)
	UpdateStatus(ctx context.Context, id string, status OrderStatus) error
	Update(ctx context.Context, order *Order) error
}

// SellerOrderRepository defines the interface for seller order persistence.
type SellerOrderRepository interface {
	Create(ctx context.Context, sellerOrder *SellerOrder) error
	GetByID(ctx context.Context, id string) (*SellerOrder, error)
	ListByOrder(ctx context.Context, orderID string) ([]*SellerOrder, error)
	ListBySeller(ctx context.Context, sellerID string, page, pageSize int) ([]*SellerOrder, int64, error)
	UpdateStatus(ctx context.Context, id string, status OrderStatus) error
}
