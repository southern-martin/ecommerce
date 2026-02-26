package domain

import "context"

// EventPublisher defines the interface for publishing domain events.
type EventPublisher interface {
	Publish(ctx context.Context, subject string, data interface{}) error
}

// Event subjects for order domain events.
const (
	EventOrderCreated   = "order.created"
	EventOrderConfirmed = "order.confirmed"
	EventOrderCancelled = "order.cancelled"
	EventOrderShipped   = "order.shipped"
	EventOrderDelivered = "order.delivered"
	EventOrderCompleted = "order.completed"
)

// OrderCreatedEvent is the payload published when an order is created.
type OrderCreatedEvent struct {
	OrderID     string      `json:"order_id"`
	OrderNumber string      `json:"order_number"`
	BuyerID     string      `json:"buyer_id"`
	TotalCents  int64       `json:"total_cents"`
	Currency    string      `json:"currency"`
	Items       []ItemEvent `json:"items"`
}

// ItemEvent represents an order item in an event payload.
type ItemEvent struct {
	ProductID      string `json:"product_id"`
	VariantID      string `json:"variant_id"`
	Quantity       int    `json:"quantity"`
	UnitPriceCents int64  `json:"unit_price_cents"`
	SellerID       string `json:"seller_id"`
}

// OrderStatusEvent is the payload published when an order status changes.
type OrderStatusEvent struct {
	OrderID     string      `json:"order_id"`
	OrderNumber string      `json:"order_number"`
	BuyerID     string      `json:"buyer_id"`
	Status      OrderStatus `json:"status"`
}
