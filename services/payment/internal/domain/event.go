package domain

import "context"

// Event subjects for payment service.
const (
	EventPaymentInitiated = "payment.initiated"
	EventPaymentCompleted = "payment.completed"
	EventPaymentFailed    = "payment.failed"
	EventPaymentRefunded  = "payment.refunded"
	EventOrderCreated     = "order.created"
)

// EventPublisher defines the interface for publishing domain events.
type EventPublisher interface {
	Publish(ctx context.Context, subject string, data interface{}) error
}

// PaymentEvent represents a payment-related event payload.
type PaymentEvent struct {
	PaymentID   string `json:"payment_id"`
	OrderID     string `json:"order_id"`
	BuyerID     string `json:"buyer_id"`
	AmountCents int64  `json:"amount_cents"`
	Currency    string `json:"currency"`
	Status      string `json:"status"`
}

// OrderCreatedEvent represents the payload from an order.created event.
type OrderCreatedEvent struct {
	OrderID     string             `json:"order_id"`
	BuyerID     string             `json:"buyer_id"`
	AmountCents int64              `json:"amount_cents"`
	Currency    string             `json:"currency"`
	SellerItems []OrderSellerItem  `json:"seller_items"`
}

// OrderSellerItem represents a seller's portion of an order.
type OrderSellerItem struct {
	SellerID    string `json:"seller_id"`
	AmountCents int64  `json:"amount_cents"`
}
