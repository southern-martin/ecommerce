package domain

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

// OrderStatus represents the current state of an order.
type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusConfirmed  OrderStatus = "confirmed"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusShipped    OrderStatus = "shipped"
	OrderStatusDelivered  OrderStatus = "delivered"
	OrderStatusCancelled  OrderStatus = "cancelled"
	OrderStatusRefunded   OrderStatus = "refunded"
	OrderStatusCompleted  OrderStatus = "completed"
)

// AllowedTransitions defines the state machine for order status transitions.
var AllowedTransitions = map[OrderStatus][]OrderStatus{
	OrderStatusPending:    {OrderStatusConfirmed, OrderStatusCancelled},
	OrderStatusConfirmed:  {OrderStatusProcessing, OrderStatusCancelled},
	OrderStatusProcessing: {OrderStatusShipped, OrderStatusCancelled},
	OrderStatusShipped:    {OrderStatusDelivered},
	OrderStatusDelivered:  {OrderStatusCompleted, OrderStatusRefunded},
	OrderStatusCompleted:  {OrderStatusRefunded},
}

// CanTransition checks whether a transition from one status to another is allowed.
func CanTransition(from, to OrderStatus) bool {
	allowed, ok := AllowedTransitions[from]
	if !ok {
		return false
	}
	for _, s := range allowed {
		if s == to {
			return true
		}
	}
	return false
}

// Address represents a shipping address.
type Address struct {
	FullName    string `json:"full_name"`
	Line1       string `json:"line1"`
	Line2       string `json:"line2"`
	City        string `json:"city"`
	State       string `json:"state"`
	PostalCode  string `json:"postal_code"`
	CountryCode string `json:"country_code"`
	Phone       string `json:"phone"`
}

// Order represents a buyer's order which may contain items from multiple sellers.
type Order struct {
	ID              string
	OrderNumber     string
	BuyerID         string
	Status          OrderStatus
	SubtotalCents   int64
	ShippingCents   int64
	TaxCents        int64
	DiscountCents   int64
	TotalCents      int64
	Currency        string
	ShippingAddress Address
	Items           []OrderItem
	SellerOrders    []SellerOrder
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// OrderItem represents a single line item in an order.
type OrderItem struct {
	ID             string
	OrderID        string
	ProductID      string
	VariantID      string
	ProductName    string
	VariantName    string
	SKU            string
	Quantity       int
	UnitPriceCents int64
	TotalCents     int64
	SellerID       string
	ImageURL       string
}

// SellerOrder groups items by seller for multi-seller marketplace orders.
type SellerOrder struct {
	ID            string
	OrderID       string
	SellerID      string
	Status        OrderStatus
	SubtotalCents int64
	Items         []OrderItem
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// NewOrder creates a new Order with generated ID and order number.
func NewOrder(buyerID string, currency string, shippingAddress Address, items []OrderItem) *Order {
	now := time.Now()
	orderID := uuid.New().String()

	var subtotal int64
	for i := range items {
		items[i].ID = uuid.New().String()
		items[i].OrderID = orderID
		items[i].TotalCents = items[i].UnitPriceCents * int64(items[i].Quantity)
		subtotal += items[i].TotalCents
	}

	order := &Order{
		ID:              orderID,
		OrderNumber:     generateOrderNumber(now),
		BuyerID:         buyerID,
		Status:          OrderStatusPending,
		SubtotalCents:   subtotal,
		ShippingCents:   0,
		TaxCents:        0,
		DiscountCents:   0,
		TotalCents:      subtotal,
		Currency:        currency,
		ShippingAddress: shippingAddress,
		Items:           items,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	// Split items by seller to create seller orders
	order.SellerOrders = splitBySeller(order)

	return order
}

// generateOrderNumber creates a human-readable order number like "ORD-20240101-ABCD".
func generateOrderNumber(t time.Time) string {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	suffix := make([]byte, 4)
	for i := range suffix {
		suffix[i] = chars[rand.Intn(len(chars))]
	}
	return fmt.Sprintf("ORD-%s-%s", t.Format("20060102"), string(suffix))
}

// splitBySeller groups order items by seller and creates SellerOrder entries.
func splitBySeller(order *Order) []SellerOrder {
	sellerItemsMap := make(map[string][]OrderItem)
	for _, item := range order.Items {
		sellerItemsMap[item.SellerID] = append(sellerItemsMap[item.SellerID], item)
	}

	var sellerOrders []SellerOrder
	for sellerID, items := range sellerItemsMap {
		var subtotal int64
		for _, item := range items {
			subtotal += item.TotalCents
		}
		sellerOrders = append(sellerOrders, SellerOrder{
			ID:            uuid.New().String(),
			OrderID:       order.ID,
			SellerID:      sellerID,
			Status:        OrderStatusPending,
			SubtotalCents: subtotal,
			Items:         items,
			CreatedAt:     order.CreatedAt,
			UpdatedAt:     order.UpdatedAt,
		})
	}
	return sellerOrders
}
