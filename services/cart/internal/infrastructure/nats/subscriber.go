package nats

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog"

	"github.com/southern-martin/ecommerce/pkg/events"
	"github.com/southern-martin/ecommerce/services/cart/internal/usecase"
)

// OrderCreatedEvent represents the payload from an order.created event.
type OrderCreatedEvent struct {
	OrderID string `json:"order_id"`
	BuyerID string `json:"buyer_id"`
}

// StartSubscriber subscribes to order events relevant to the cart service.
func StartSubscriber(sub *events.Subscriber, cartUC *usecase.CartUseCase, logger zerolog.Logger) error {
	// Subscribe to order.created — clear the buyer's cart after an order is placed
	if err := sub.Subscribe(events.SubjectOrderCreated, "cart-service-order-created", func(data []byte) {
		var evt OrderCreatedEvent
		if err := json.Unmarshal(data, &evt); err != nil {
			logger.Error().Err(err).Msg("failed to unmarshal order.created event")
			return
		}

		logger.Info().
			Str("order_id", evt.OrderID).
			Str("buyer_id", evt.BuyerID).
			Msg("received order.created event")

		if err := cartUC.ClearCart(context.Background(), evt.BuyerID); err != nil {
			logger.Error().Err(err).
				Str("order_id", evt.OrderID).
				Str("buyer_id", evt.BuyerID).
				Msg("failed to clear cart after order created")
			return
		}

		logger.Info().
			Str("order_id", evt.OrderID).
			Str("buyer_id", evt.BuyerID).
			Msg("cart cleared after order created")
	}); err != nil {
		return fmt.Errorf("subscribe to %s: %w", events.SubjectOrderCreated, err)
	}

	logger.Info().Msg("cart NATS subscribers started")
	return nil
}
