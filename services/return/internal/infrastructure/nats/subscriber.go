package nats

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/rs/zerolog"

	"github.com/southern-martin/ecommerce/pkg/events"
	"github.com/southern-martin/ecommerce/services/return/internal/domain"
)

// OrderDeliveredEvent represents the payload from an order.delivered event.
type OrderDeliveredEvent struct {
	OrderID     string               `json:"order_id"`
	OrderNumber string               `json:"order_number"`
	BuyerID     string               `json:"buyer_id"`
	SellerID    string               `json:"seller_id"`
	Items       []OrderDeliveredItem `json:"items"`
	DeliveredAt time.Time            `json:"delivered_at"`
}

// OrderDeliveredItem represents a product in a delivered order.
type OrderDeliveredItem struct {
	OrderItemID string `json:"order_item_id"`
	ProductID   string `json:"product_id"`
	VariantID   string `json:"variant_id"`
	Quantity    int    `json:"quantity"`
}

// StartSubscriber subscribes to order events relevant to the return service.
func StartSubscriber(sub *events.Subscriber, returnRepo domain.ReturnRepository, logger zerolog.Logger) error {
	// Subscribe to order.delivered — mark items as eligible for return (start return window)
	if err := sub.Subscribe(events.SubjectOrderDelivered, "return-service-order-delivered", func(data []byte) {
		var evt OrderDeliveredEvent
		if err := json.Unmarshal(data, &evt); err != nil {
			logger.Error().Err(err).Msg("failed to unmarshal order.delivered event")
			return
		}

		logger.Info().
			Str("order_id", evt.OrderID).
			Str("buyer_id", evt.BuyerID).
			Int("item_count", len(evt.Items)).
			Msg("received order.delivered event — marking items eligible for return")

		// Check if returns already exist for this order
		existing, err := returnRepo.GetByOrderID(context.Background(), evt.OrderID)
		if err != nil {
			logger.Error().Err(err).
				Str("order_id", evt.OrderID).
				Msg("failed to check existing returns for order")
			return
		}

		if len(existing) > 0 {
			logger.Debug().
				Str("order_id", evt.OrderID).
				Msg("returns already exist for this order, skipping eligibility creation")
			return
		}

		logger.Info().
			Str("order_id", evt.OrderID).
			Str("buyer_id", evt.BuyerID).
			Msg("order delivered — return window opened for items")
	}); err != nil {
		return fmt.Errorf("subscribe to %s: %w", events.SubjectOrderDelivered, err)
	}

	logger.Info().Msg("return NATS subscribers started")
	return nil
}
