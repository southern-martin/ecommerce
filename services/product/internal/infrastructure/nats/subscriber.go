package nats

import (
	"context"
	"encoding/json"

	"github.com/rs/zerolog"

	"github.com/southern-martin/ecommerce/pkg/events"
	"github.com/southern-martin/ecommerce/services/product/internal/usecase"
)

// OrderCreatedEvent matches the order service's event payload.
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

// StartSubscriber subscribes to the order.created subject and decrements variant stock.
func StartSubscriber(sub *events.Subscriber, variantUC *usecase.VariantUseCase, logger zerolog.Logger) error {
	return sub.Subscribe(events.SubjectOrderCreated, "product-service-order-created", func(data []byte) {
		var evt OrderCreatedEvent
		if err := json.Unmarshal(data, &evt); err != nil {
			logger.Error().Err(err).Msg("failed to unmarshal order.created event")
			return
		}

		logger.Info().
			Str("order_id", evt.OrderID).
			Str("order_number", evt.OrderNumber).
			Int("item_count", len(evt.Items)).
			Msg("received order.created event, decrementing stock")

		for _, item := range evt.Items {
			delta := -item.Quantity
			if err := variantUC.UpdateStockDirect(context.Background(), item.VariantID, delta); err != nil {
				logger.Error().
					Err(err).
					Str("order_id", evt.OrderID).
					Str("variant_id", item.VariantID).
					Int("quantity", item.Quantity).
					Msg("failed to decrement stock for variant")
				continue
			}

			logger.Info().
				Str("variant_id", item.VariantID).
				Int("decremented_by", item.Quantity).
				Msg("stock decremented successfully")
		}
	})
}
