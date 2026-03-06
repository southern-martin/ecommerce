package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/southern-martin/ecommerce/services/order/internal/domain"
)

func TestUpdateSellerOrderStatus_ValidTransition(t *testing.T) {
	soRepo := &mockSellerOrderRepo{
		getByIDFn: func(_ context.Context, _ string) (*domain.SellerOrder, error) {
			return &domain.SellerOrder{ID: "so-1", OrderID: "o-1", Status: domain.OrderStatusPending}, nil
		},
		updateStatusFn: func(_ context.Context, _ string, _ domain.OrderStatus) error { return nil },
	}
	oRepo := &mockOrderRepo{
		getByIDFn: func(_ context.Context, _ string) (*domain.Order, error) {
			return &domain.Order{ID: "o-1", OrderNumber: "ORD-123", BuyerID: "b-1"}, nil
		},
	}
	uc := NewUpdateOrderStatusUseCase(oRepo, soRepo, &mockEventPublisher{})

	so, err := uc.Execute(context.Background(), "so-1", domain.OrderStatusConfirmed)
	require.NoError(t, err)
	assert.Equal(t, domain.OrderStatusConfirmed, so.Status)
}

func TestUpdateSellerOrderStatus_InvalidTransition(t *testing.T) {
	soRepo := &mockSellerOrderRepo{
		getByIDFn: func(_ context.Context, _ string) (*domain.SellerOrder, error) {
			return &domain.SellerOrder{ID: "so-1", Status: domain.OrderStatusPending}, nil
		},
	}
	uc := NewUpdateOrderStatusUseCase(&mockOrderRepo{}, soRepo, &mockEventPublisher{})

	_, err := uc.Execute(context.Background(), "so-1", domain.OrderStatusShipped)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid status transition")
}

func TestUpdateSellerOrderStatus_NotFound(t *testing.T) {
	soRepo := &mockSellerOrderRepo{
		getByIDFn: func(_ context.Context, _ string) (*domain.SellerOrder, error) {
			return nil, errors.New("not found")
		},
	}
	uc := NewUpdateOrderStatusUseCase(&mockOrderRepo{}, soRepo, &mockEventPublisher{})

	_, err := uc.Execute(context.Background(), "so-99", domain.OrderStatusConfirmed)
	assert.EqualError(t, err, "not found")
}

func TestUpdateOrderStatus_ValidTransition(t *testing.T) {
	oRepo := &mockOrderRepo{
		getByIDFn: func(_ context.Context, _ string) (*domain.Order, error) {
			return &domain.Order{ID: "o-1", OrderNumber: "ORD-123", BuyerID: "b-1", Status: domain.OrderStatusPending}, nil
		},
		updateStatusFn: func(_ context.Context, _ string, _ domain.OrderStatus) error { return nil },
	}
	uc := NewUpdateOrderStatusUseCase(oRepo, &mockSellerOrderRepo{}, &mockEventPublisher{})

	order, err := uc.UpdateOrderStatus(context.Background(), "o-1", domain.OrderStatusConfirmed)
	require.NoError(t, err)
	assert.Equal(t, domain.OrderStatusConfirmed, order.Status)
}

func TestUpdateOrderStatus_InvalidTransition(t *testing.T) {
	oRepo := &mockOrderRepo{
		getByIDFn: func(_ context.Context, _ string) (*domain.Order, error) {
			return &domain.Order{ID: "o-1", Status: domain.OrderStatusShipped}, nil
		},
	}
	uc := NewUpdateOrderStatusUseCase(oRepo, &mockSellerOrderRepo{}, &mockEventPublisher{})

	_, err := uc.UpdateOrderStatus(context.Background(), "o-1", domain.OrderStatusPending)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid status transition")
}

func TestStatusToEventSubject(t *testing.T) {
	tests := []struct {
		status   domain.OrderStatus
		expected string
	}{
		{domain.OrderStatusConfirmed, domain.EventOrderConfirmed},
		{domain.OrderStatusCancelled, domain.EventOrderCancelled},
		{domain.OrderStatusShipped, domain.EventOrderShipped},
		{domain.OrderStatusDelivered, domain.EventOrderDelivered},
		{domain.OrderStatusCompleted, domain.EventOrderCompleted},
		{domain.OrderStatusPending, ""},
		{"unknown", ""},
	}
	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			assert.Equal(t, tt.expected, statusToEventSubject(tt.status))
		})
	}
}
