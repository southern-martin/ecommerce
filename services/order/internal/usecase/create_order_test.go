package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/southern-martin/ecommerce/services/order/internal/domain"
)

// --- Mocks ---

type mockOrderRepo struct {
	createFn        func(ctx context.Context, order *domain.Order) error
	getByIDFn       func(ctx context.Context, id string) (*domain.Order, error)
	getByNumberFn   func(ctx context.Context, num string) (*domain.Order, error)
	listFn          func(ctx context.Context, filter domain.OrderFilter) ([]*domain.Order, int64, error)
	updateStatusFn  func(ctx context.Context, id string, status domain.OrderStatus) error
	updateFn        func(ctx context.Context, order *domain.Order) error
}

func (m *mockOrderRepo) Create(ctx context.Context, order *domain.Order) error { return m.createFn(ctx, order) }
func (m *mockOrderRepo) GetByID(ctx context.Context, id string) (*domain.Order, error) { return m.getByIDFn(ctx, id) }
func (m *mockOrderRepo) GetByOrderNumber(ctx context.Context, num string) (*domain.Order, error) { return m.getByNumberFn(ctx, num) }
func (m *mockOrderRepo) List(ctx context.Context, filter domain.OrderFilter) ([]*domain.Order, int64, error) { return m.listFn(ctx, filter) }
func (m *mockOrderRepo) UpdateStatus(ctx context.Context, id string, status domain.OrderStatus) error { return m.updateStatusFn(ctx, id, status) }
func (m *mockOrderRepo) Update(ctx context.Context, order *domain.Order) error { return m.updateFn(ctx, order) }

type mockSellerOrderRepo struct {
	createFn       func(ctx context.Context, so *domain.SellerOrder) error
	getByIDFn      func(ctx context.Context, id string) (*domain.SellerOrder, error)
	listByOrderFn  func(ctx context.Context, orderID string) ([]*domain.SellerOrder, error)
	listBySellerFn func(ctx context.Context, sellerID string, page, pageSize int) ([]*domain.SellerOrder, int64, error)
	updateStatusFn func(ctx context.Context, id string, status domain.OrderStatus) error
}

func (m *mockSellerOrderRepo) Create(ctx context.Context, so *domain.SellerOrder) error { return m.createFn(ctx, so) }
func (m *mockSellerOrderRepo) GetByID(ctx context.Context, id string) (*domain.SellerOrder, error) { return m.getByIDFn(ctx, id) }
func (m *mockSellerOrderRepo) ListByOrder(ctx context.Context, orderID string) ([]*domain.SellerOrder, error) { return m.listByOrderFn(ctx, orderID) }
func (m *mockSellerOrderRepo) ListBySeller(ctx context.Context, sellerID string, page, pageSize int) ([]*domain.SellerOrder, int64, error) { return m.listBySellerFn(ctx, sellerID, page, pageSize) }
func (m *mockSellerOrderRepo) UpdateStatus(ctx context.Context, id string, status domain.OrderStatus) error { return m.updateStatusFn(ctx, id, status) }

type mockEventPublisher struct {
	publishFn func(ctx context.Context, subject string, data interface{}) error
}

func (m *mockEventPublisher) Publish(ctx context.Context, subject string, data interface{}) error {
	if m.publishFn != nil {
		return m.publishFn(ctx, subject, data)
	}
	return nil
}

// --- CreateOrder tests ---

func validInput() CreateOrderInput {
	return CreateOrderInput{
		BuyerID:  "buyer-1",
		Currency: "USD",
		Items: []CreateOrderItemInput{
			{ProductID: "p1", Quantity: 2, UnitPriceCents: 1000, SellerID: "s1"},
		},
	}
}

func TestCreateOrder_Success(t *testing.T) {
	var createdOrder *domain.Order
	var createdSellerOrders []*domain.SellerOrder

	repo := &mockOrderRepo{
		createFn: func(_ context.Context, o *domain.Order) error { createdOrder = o; return nil },
	}
	soRepo := &mockSellerOrderRepo{
		createFn: func(_ context.Context, so *domain.SellerOrder) error {
			createdSellerOrders = append(createdSellerOrders, so)
			return nil
		},
	}
	pub := &mockEventPublisher{}
	uc := NewCreateOrderUseCase(repo, soRepo, pub, nil)

	order, err := uc.Execute(context.Background(), validInput())
	require.NoError(t, err)
	assert.NotEmpty(t, order.ID)
	assert.Equal(t, "buyer-1", order.BuyerID)
	assert.Equal(t, domain.OrderStatusPending, order.Status)
	assert.Equal(t, int64(2000), order.TotalCents)
	assert.NotNil(t, createdOrder)
	assert.Len(t, createdSellerOrders, 1)
}

func TestCreateOrder_EmptyBuyerID(t *testing.T) {
	uc := NewCreateOrderUseCase(&mockOrderRepo{}, &mockSellerOrderRepo{}, &mockEventPublisher{}, nil)
	input := validInput()
	input.BuyerID = ""
	_, err := uc.Execute(context.Background(), input)
	assert.EqualError(t, err, "buyer_id is required")
}

func TestCreateOrder_NoItems(t *testing.T) {
	uc := NewCreateOrderUseCase(&mockOrderRepo{}, &mockSellerOrderRepo{}, &mockEventPublisher{}, nil)
	input := validInput()
	input.Items = nil
	_, err := uc.Execute(context.Background(), input)
	assert.EqualError(t, err, "at least one item is required")
}

func TestCreateOrder_InvalidQuantity(t *testing.T) {
	uc := NewCreateOrderUseCase(&mockOrderRepo{}, &mockSellerOrderRepo{}, &mockEventPublisher{}, nil)
	input := validInput()
	input.Items[0].Quantity = 0
	_, err := uc.Execute(context.Background(), input)
	assert.EqualError(t, err, "item quantity must be greater than 0")
}

func TestCreateOrder_InvalidPrice(t *testing.T) {
	uc := NewCreateOrderUseCase(&mockOrderRepo{}, &mockSellerOrderRepo{}, &mockEventPublisher{}, nil)
	input := validInput()
	input.Items[0].UnitPriceCents = -100
	_, err := uc.Execute(context.Background(), input)
	assert.EqualError(t, err, "item unit price must be greater than 0")
}

func TestCreateOrder_MissingSellerID(t *testing.T) {
	uc := NewCreateOrderUseCase(&mockOrderRepo{}, &mockSellerOrderRepo{}, &mockEventPublisher{}, nil)
	input := validInput()
	input.Items[0].SellerID = ""
	_, err := uc.Execute(context.Background(), input)
	assert.EqualError(t, err, "seller_id is required for each item")
}

func TestCreateOrder_DefaultCurrency(t *testing.T) {
	repo := &mockOrderRepo{createFn: func(_ context.Context, _ *domain.Order) error { return nil }}
	soRepo := &mockSellerOrderRepo{createFn: func(_ context.Context, _ *domain.SellerOrder) error { return nil }}
	uc := NewCreateOrderUseCase(repo, soRepo, &mockEventPublisher{}, nil)

	input := validInput()
	input.Currency = ""
	order, err := uc.Execute(context.Background(), input)
	require.NoError(t, err)
	assert.Equal(t, "USD", order.Currency)
}

func TestCreateOrder_RepoError(t *testing.T) {
	repo := &mockOrderRepo{
		createFn: func(_ context.Context, _ *domain.Order) error { return errors.New("db error") },
	}
	uc := NewCreateOrderUseCase(repo, &mockSellerOrderRepo{}, &mockEventPublisher{}, nil)
	_, err := uc.Execute(context.Background(), validInput())
	assert.EqualError(t, err, "db error")
}

func TestCreateOrder_SellerRepoError(t *testing.T) {
	repo := &mockOrderRepo{createFn: func(_ context.Context, _ *domain.Order) error { return nil }}
	soRepo := &mockSellerOrderRepo{
		createFn: func(_ context.Context, _ *domain.SellerOrder) error { return errors.New("seller db error") },
	}
	uc := NewCreateOrderUseCase(repo, soRepo, &mockEventPublisher{}, nil)
	_, err := uc.Execute(context.Background(), validInput())
	assert.EqualError(t, err, "seller db error")
}
