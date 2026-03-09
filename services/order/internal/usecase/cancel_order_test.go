package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/southern-martin/ecommerce/services/order/internal/domain"
)

func TestCancelOrder_Success(t *testing.T) {
	var updatedStatus domain.OrderStatus
	oRepo := &mockOrderRepo{
		getByIDFn: func(_ context.Context, _ string) (*domain.Order, error) {
			return &domain.Order{ID: "o-1", BuyerID: "b-1", Status: domain.OrderStatusPending, OrderNumber: "ORD-123"}, nil
		},
		updateStatusFn: func(_ context.Context, _ string, s domain.OrderStatus) error {
			updatedStatus = s
			return nil
		},
	}
	soRepo := &mockSellerOrderRepo{
		listByOrderFn: func(_ context.Context, _ string) ([]*domain.SellerOrder, error) {
			return []*domain.SellerOrder{
				{ID: "so-1", Status: domain.OrderStatusPending},
			}, nil
		},
		updateStatusFn: func(_ context.Context, _ string, _ domain.OrderStatus) error { return nil },
	}

	uc := NewCancelOrderUseCase(oRepo, soRepo, &mockEventPublisher{}, nil)
	order, err := uc.Execute(context.Background(), "o-1", "b-1")
	require.NoError(t, err)
	assert.Equal(t, domain.OrderStatusCancelled, order.Status)
	assert.Equal(t, domain.OrderStatusCancelled, updatedStatus)
}

func TestCancelOrder_WrongBuyer(t *testing.T) {
	oRepo := &mockOrderRepo{
		getByIDFn: func(_ context.Context, _ string) (*domain.Order, error) {
			return &domain.Order{ID: "o-1", BuyerID: "b-1", Status: domain.OrderStatusPending}, nil
		},
	}
	uc := NewCancelOrderUseCase(oRepo, &mockSellerOrderRepo{}, &mockEventPublisher{}, nil)

	_, err := uc.Execute(context.Background(), "o-1", "wrong-buyer")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "does not belong")
}

func TestCancelOrder_InvalidTransition(t *testing.T) {
	oRepo := &mockOrderRepo{
		getByIDFn: func(_ context.Context, _ string) (*domain.Order, error) {
			return &domain.Order{ID: "o-1", BuyerID: "b-1", Status: domain.OrderStatusShipped}, nil
		},
	}
	uc := NewCancelOrderUseCase(oRepo, &mockSellerOrderRepo{}, &mockEventPublisher{}, nil)

	_, err := uc.Execute(context.Background(), "o-1", "b-1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be cancelled")
}

func TestCancelOrder_NotFound(t *testing.T) {
	oRepo := &mockOrderRepo{
		getByIDFn: func(_ context.Context, _ string) (*domain.Order, error) {
			return nil, errors.New("not found")
		},
	}
	uc := NewCancelOrderUseCase(oRepo, &mockSellerOrderRepo{}, &mockEventPublisher{}, nil)

	_, err := uc.Execute(context.Background(), "o-99", "b-1")
	assert.EqualError(t, err, "not found")
}

func TestCancelOrder_CascadesToSellerOrders(t *testing.T) {
	var cancelledIDs []string
	oRepo := &mockOrderRepo{
		getByIDFn: func(_ context.Context, _ string) (*domain.Order, error) {
			return &domain.Order{ID: "o-1", BuyerID: "b-1", Status: domain.OrderStatusPending, OrderNumber: "ORD-123"}, nil
		},
		updateStatusFn: func(_ context.Context, _ string, _ domain.OrderStatus) error { return nil },
	}
	soRepo := &mockSellerOrderRepo{
		listByOrderFn: func(_ context.Context, _ string) ([]*domain.SellerOrder, error) {
			return []*domain.SellerOrder{
				{ID: "so-1", Status: domain.OrderStatusPending},
				{ID: "so-2", Status: domain.OrderStatusConfirmed},
				{ID: "so-3", Status: domain.OrderStatusShipped}, // Cannot be cancelled
			}, nil
		},
		updateStatusFn: func(_ context.Context, id string, _ domain.OrderStatus) error {
			cancelledIDs = append(cancelledIDs, id)
			return nil
		},
	}

	uc := NewCancelOrderUseCase(oRepo, soRepo, &mockEventPublisher{}, nil)
	_, err := uc.Execute(context.Background(), "o-1", "b-1")
	require.NoError(t, err)

	// so-1 (pending) and so-2 (confirmed) can be cancelled; so-3 (shipped) cannot
	assert.Len(t, cancelledIDs, 2)
	assert.Contains(t, cancelledIDs, "so-1")
	assert.Contains(t, cancelledIDs, "so-2")
}
