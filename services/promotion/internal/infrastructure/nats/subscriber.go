package nats

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog"

	"github.com/southern-martin/ecommerce/pkg/events"
	"github.com/southern-martin/ecommerce/services/promotion/internal/usecase"
)

// OrderCompletedEvent represents the payload from an order.completed event.
type OrderCompletedEvent struct {
	OrderID    string `json:"order_id"`
	BuyerID    string `json:"buyer_id"`
	CouponCode string `json:"coupon_code"`
	CouponID   string `json:"coupon_id"`
}

// StartSubscriber subscribes to order events relevant to the promotion service.
func StartSubscriber(sub *events.Subscriber, couponUC *usecase.CouponUseCase, logger zerolog.Logger) error {
	// Subscribe to order.completed — increment coupon usage count when an order completes
	if err := sub.Subscribe(events.SubjectOrderCompleted, "promotion-service-order-completed", func(data []byte) {
		var evt OrderCompletedEvent
		if err := json.Unmarshal(data, &evt); err != nil {
			logger.Error().Err(err).Msg("failed to unmarshal order.completed event")
			return
		}

		logger.Info().
			Str("order_id", evt.OrderID).
			Str("buyer_id", evt.BuyerID).
			Str("coupon_code", evt.CouponCode).
			Msg("received order.completed event")

		// Only process if a coupon was used in this order
		if evt.CouponCode == "" {
			logger.Debug().Str("order_id", evt.OrderID).Msg("no coupon used in order, skipping")
			return
		}

		// Look up the coupon to verify it exists
		coupon, err := couponUC.GetCoupon(context.Background(), evt.CouponID)
		if err != nil {
			logger.Error().Err(err).
				Str("coupon_id", evt.CouponID).
				Str("order_id", evt.OrderID).
				Msg("failed to get coupon for usage tracking")
			return
		}

		// Increment the usage count
		coupon.UsageCount++
		if err := couponUC.UpdateCoupon(context.Background(), coupon); err != nil {
			logger.Error().Err(err).
				Str("coupon_id", evt.CouponID).
				Str("order_id", evt.OrderID).
				Msg("failed to increment coupon usage count")
			return
		}

		logger.Info().
			Str("coupon_id", evt.CouponID).
			Str("coupon_code", evt.CouponCode).
			Str("order_id", evt.OrderID).
			Int("new_usage_count", coupon.UsageCount).
			Msg("coupon usage count incremented after order completed")
	}); err != nil {
		return fmt.Errorf("subscribe to %s: %w", events.SubjectOrderCompleted, err)
	}

	logger.Info().Msg("promotion NATS subscribers started")
	return nil
}
