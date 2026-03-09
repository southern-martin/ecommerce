package nats

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog"

	"github.com/southern-martin/ecommerce/pkg/events"
	"github.com/southern-martin/ecommerce/services/review/internal/usecase"
)

// OrderDeliveredEvent represents the payload from an order.delivered event.
type OrderDeliveredEvent struct {
	OrderID     string               `json:"order_id"`
	OrderNumber string               `json:"order_number"`
	BuyerID     string               `json:"buyer_id"`
	Items       []OrderDeliveredItem `json:"items"`
}

// OrderDeliveredItem represents a product in a delivered order.
type OrderDeliveredItem struct {
	ProductID string `json:"product_id"`
}

// StartSubscriber subscribes to order events relevant to the review service.
func StartSubscriber(sub *events.Subscriber, reviewUC *usecase.ReviewUseCase, logger zerolog.Logger) error {
	// Subscribe to order.delivered — create review eligibility prompts for the buyer
	if err := sub.Subscribe(events.SubjectOrderDelivered, "review-service-order-delivered", func(data []byte) {
		var evt OrderDeliveredEvent
		if err := json.Unmarshal(data, &evt); err != nil {
			logger.Error().Err(err).Msg("failed to unmarshal order.delivered event")
			return
		}

		logger.Info().
			Str("order_id", evt.OrderID).
			Str("buyer_id", evt.BuyerID).
			Int("item_count", len(evt.Items)).
			Msg("received order.delivered event")

		// Create a placeholder verified-purchase review prompt for each product in the order.
		// This marks the buyer as eligible to leave a verified purchase review.
		for _, item := range evt.Items {
			_, err := reviewUC.CreateReview(context.Background(), usecase.CreateReviewRequest{
				ProductID:          item.ProductID,
				UserID:             evt.BuyerID,
				Rating:             5,
				Title:              fmt.Sprintf("Review prompt for order %s", evt.OrderNumber),
				Content:            "",
				IsVerifiedPurchase: true,
			})
			if err != nil {
				logger.Warn().Err(err).
					Str("order_id", evt.OrderID).
					Str("buyer_id", evt.BuyerID).
					Str("product_id", item.ProductID).
					Msg("failed to create review prompt for delivered order item")
				continue
			}

			logger.Info().
				Str("order_id", evt.OrderID).
				Str("buyer_id", evt.BuyerID).
				Str("product_id", item.ProductID).
				Msg("review eligibility created for delivered order item")
		}
	}); err != nil {
		return fmt.Errorf("subscribe to %s: %w", events.SubjectOrderDelivered, err)
	}

	logger.Info().Msg("review NATS subscribers started")
	return nil
}
