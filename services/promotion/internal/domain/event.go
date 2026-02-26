package domain

import "context"

// EventPublisher defines the interface for publishing domain events.
type EventPublisher interface {
	Publish(ctx context.Context, subject string, data interface{}) error
}

// Event subjects for promotion domain events.
const (
	EventCouponRedeemed   = "coupon.redeemed"
	EventFlashSaleStarted = "flash_sale.started"
	EventFlashSaleEnded   = "flash_sale.ended"
)

// CouponRedeemedEvent is the payload published when a coupon is redeemed.
type CouponRedeemedEvent struct {
	CouponID      string `json:"coupon_id"`
	CouponCode    string `json:"coupon_code"`
	UserID        string `json:"user_id"`
	OrderID       string `json:"order_id"`
	DiscountCents int64  `json:"discount_cents"`
}

// FlashSaleEvent is the payload published when a flash sale starts or ends.
type FlashSaleEvent struct {
	FlashSaleID string `json:"flash_sale_id"`
	Name        string `json:"name"`
}
