package nats

import (
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog"

	"github.com/southern-martin/ecommerce/pkg/events"
)

// OrderEvent matches the order service's event payload for confirmed orders.
type OrderEvent struct {
	OrderID    string `json:"order_id"`
	BuyerID    string `json:"buyer_id"`
	TotalCents int64  `json:"total_cents"`
	Status     string `json:"status"`
}

// StartSubscriber subscribes to order events relevant to shipping.
func StartSubscriber(sub *events.Subscriber, logger zerolog.Logger) error {
	// Subscribe to order.confirmed — log that order is ready for shipment preparation
	if err := sub.Subscribe(events.SubjectOrderConfirmed, "shipping-service-order-confirmed", func(data []byte) {
		var evt OrderEvent
		if err := json.Unmarshal(data, &evt); err != nil {
			logger.Error().Err(err).Msg("failed to unmarshal order.confirmed event")
			return
		}

		logger.Info().
			Str("order_id", evt.OrderID).
			Str("buyer_id", evt.BuyerID).
			Int64("total_cents", evt.TotalCents).
			Msg("order confirmed and ready for shipment preparation")
	}); err != nil {
		return fmt.Errorf("failed to subscribe to order.confirmed: %w", err)
	}

	logger.Info().Msg("shipping NATS subscribers started")
	return nil
}
