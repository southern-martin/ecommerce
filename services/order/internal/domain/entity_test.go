package domain

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCanTransition_ValidTransitions(t *testing.T) {
	tests := []struct {
		from, to OrderStatus
	}{
		{OrderStatusPending, OrderStatusConfirmed},
		{OrderStatusPending, OrderStatusCancelled},
		{OrderStatusConfirmed, OrderStatusProcessing},
		{OrderStatusConfirmed, OrderStatusCancelled},
		{OrderStatusProcessing, OrderStatusShipped},
		{OrderStatusProcessing, OrderStatusCancelled},
		{OrderStatusShipped, OrderStatusDelivered},
		{OrderStatusDelivered, OrderStatusCompleted},
		{OrderStatusDelivered, OrderStatusRefunded},
		{OrderStatusCompleted, OrderStatusRefunded},
	}
	for _, tt := range tests {
		t.Run(string(tt.from)+"->"+string(tt.to), func(t *testing.T) {
			assert.True(t, CanTransition(tt.from, tt.to))
		})
	}
}

func TestCanTransition_InvalidTransitions(t *testing.T) {
	tests := []struct {
		from, to OrderStatus
	}{
		{OrderStatusPending, OrderStatusShipped},
		{OrderStatusShipped, OrderStatusPending},
		{OrderStatusCancelled, OrderStatusConfirmed},
		{OrderStatusCancelled, OrderStatusShipped},
		{OrderStatusRefunded, OrderStatusPending},
		{OrderStatusDelivered, OrderStatusProcessing},
		{OrderStatusCompleted, OrderStatusPending},
	}
	for _, tt := range tests {
		t.Run(string(tt.from)+"->"+string(tt.to), func(t *testing.T) {
			assert.False(t, CanTransition(tt.from, tt.to))
		})
	}
}

func TestCanTransition_UnknownStatus(t *testing.T) {
	assert.False(t, CanTransition("unknown", OrderStatusConfirmed))
}

func TestNewOrder_BasicFields(t *testing.T) {
	items := []OrderItem{
		{ProductID: "p1", Quantity: 2, UnitPriceCents: 1000, SellerID: "s1"},
	}
	order := NewOrder("buyer-1", "USD", Address{FullName: "John"}, items)

	assert.NotEmpty(t, order.ID)
	assert.Equal(t, "buyer-1", order.BuyerID)
	assert.Equal(t, "USD", order.Currency)
	assert.Equal(t, OrderStatusPending, order.Status)
	assert.Equal(t, "John", order.ShippingAddress.FullName)
}

func TestNewOrder_OrderNumberFormat(t *testing.T) {
	items := []OrderItem{
		{ProductID: "p1", Quantity: 1, UnitPriceCents: 100, SellerID: "s1"},
	}
	order := NewOrder("buyer-1", "USD", Address{}, items)

	re := regexp.MustCompile(`^ORD-\d{8}-[A-Z0-9]{4}$`)
	assert.Regexp(t, re, order.OrderNumber)
}

func TestNewOrder_SubtotalCalculation(t *testing.T) {
	items := []OrderItem{
		{ProductID: "p1", Quantity: 2, UnitPriceCents: 1000, SellerID: "s1"},
		{ProductID: "p2", Quantity: 3, UnitPriceCents: 500, SellerID: "s1"},
	}
	order := NewOrder("buyer-1", "USD", Address{}, items)

	// 2*1000 + 3*500 = 3500
	assert.Equal(t, int64(3500), order.SubtotalCents)
	assert.Equal(t, int64(3500), order.TotalCents)
}

func TestNewOrder_ItemIDsGenerated(t *testing.T) {
	items := []OrderItem{
		{ProductID: "p1", Quantity: 1, UnitPriceCents: 100, SellerID: "s1"},
		{ProductID: "p2", Quantity: 1, UnitPriceCents: 200, SellerID: "s1"},
	}
	order := NewOrder("buyer-1", "USD", Address{}, items)

	for _, item := range order.Items {
		assert.NotEmpty(t, item.ID)
		assert.Equal(t, order.ID, item.OrderID)
	}
	// IDs should be unique
	assert.NotEqual(t, order.Items[0].ID, order.Items[1].ID)
}

func TestNewOrder_ItemTotalCentsCalculated(t *testing.T) {
	items := []OrderItem{
		{ProductID: "p1", Quantity: 3, UnitPriceCents: 1500, SellerID: "s1"},
	}
	order := NewOrder("buyer-1", "USD", Address{}, items)

	assert.Equal(t, int64(4500), order.Items[0].TotalCents)
}

func TestNewOrder_SingleSellerSingleSellerOrder(t *testing.T) {
	items := []OrderItem{
		{ProductID: "p1", Quantity: 1, UnitPriceCents: 1000, SellerID: "s1"},
		{ProductID: "p2", Quantity: 2, UnitPriceCents: 500, SellerID: "s1"},
	}
	order := NewOrder("buyer-1", "USD", Address{}, items)

	require.Len(t, order.SellerOrders, 1)
	assert.Equal(t, "s1", order.SellerOrders[0].SellerID)
	assert.Equal(t, int64(2000), order.SellerOrders[0].SubtotalCents)
	assert.Equal(t, order.ID, order.SellerOrders[0].OrderID)
	assert.Equal(t, OrderStatusPending, order.SellerOrders[0].Status)
}

func TestNewOrder_MultipleSellersSplit(t *testing.T) {
	items := []OrderItem{
		{ProductID: "p1", Quantity: 1, UnitPriceCents: 1000, SellerID: "seller-a"},
		{ProductID: "p2", Quantity: 2, UnitPriceCents: 500, SellerID: "seller-b"},
		{ProductID: "p3", Quantity: 1, UnitPriceCents: 2000, SellerID: "seller-a"},
	}
	order := NewOrder("buyer-1", "USD", Address{}, items)

	require.Len(t, order.SellerOrders, 2)

	// Build a map of sellerID -> sellerOrder for order-independent assertion
	soMap := make(map[string]*SellerOrder)
	for i := range order.SellerOrders {
		soMap[order.SellerOrders[i].SellerID] = &order.SellerOrders[i]
	}

	// seller-a: 1*1000 + 1*2000 = 3000
	require.Contains(t, soMap, "seller-a")
	assert.Equal(t, int64(3000), soMap["seller-a"].SubtotalCents)
	assert.Len(t, soMap["seller-a"].Items, 2)

	// seller-b: 2*500 = 1000
	require.Contains(t, soMap, "seller-b")
	assert.Equal(t, int64(1000), soMap["seller-b"].SubtotalCents)
	assert.Len(t, soMap["seller-b"].Items, 1)
}
