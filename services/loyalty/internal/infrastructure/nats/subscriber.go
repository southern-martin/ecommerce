package nats

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog"

	"github.com/southern-martin/ecommerce/pkg/events"
	"github.com/southern-martin/ecommerce/services/loyalty/internal/domain"
	"github.com/southern-martin/ecommerce/services/loyalty/internal/usecase"
)

// OrderEvent matches the order service's event payload for completed/delivered orders.
type OrderEvent struct {
	OrderID     string `json:"order_id"`
	OrderNumber string `json:"order_number"`
	BuyerID     string `json:"buyer_id"`
	TotalCents  int64  `json:"total_cents"`
	Currency    string `json:"currency"`
}

// StartSubscriber subscribes to order events and awards loyalty points.
func StartSubscriber(sub *events.Subscriber, pointsUC *usecase.PointsUseCase, logger zerolog.Logger) error {
	// Subscribe to order.delivered — award 1 point per dollar spent
	if err := sub.Subscribe(events.SubjectOrderDelivered, "loyalty-service-order-delivered", func(data []byte) {
		var evt OrderEvent
		if err := json.Unmarshal(data, &evt); err != nil {
			logger.Error().Err(err).Msg("failed to unmarshal order.delivered event")
			return
		}

		logger.Info().
			Str("order_id", evt.OrderID).
			Str("buyer_id", evt.BuyerID).
			Int64("total_cents", evt.TotalCents).
			Msg("received order.delivered event")

		points := evt.TotalCents / 100 // 1 point per dollar
		if points < 1 {
			points = 1
		}

		description := fmt.Sprintf("Points earned for order #%s", evt.OrderNumber)
		_, err := pointsUC.EarnPoints(context.Background(), usecase.EarnPointsRequest{
			UserID:      evt.BuyerID,
			Points:      points,
			Source:      domain.SourceOrder,
			ReferenceID: evt.OrderID,
			Description: description,
		})
		if err != nil {
			logger.Error().Err(err).
				Str("order_id", evt.OrderID).
				Str("buyer_id", evt.BuyerID).
				Msg("failed to award loyalty points from event")
		}
	}); err != nil {
		return fmt.Errorf("failed to subscribe to order.delivered: %w", err)
	}

	// Subscribe to order.completed — also award points (covers both statuses)
	if err := sub.Subscribe(events.SubjectOrderCompleted, "loyalty-service-order-completed", func(data []byte) {
		var evt OrderEvent
		if err := json.Unmarshal(data, &evt); err != nil {
			logger.Error().Err(err).Msg("failed to unmarshal order.completed event")
			return
		}

		logger.Info().
			Str("order_id", evt.OrderID).
			Str("buyer_id", evt.BuyerID).
			Int64("total_cents", evt.TotalCents).
			Msg("received order.completed event")

		points := evt.TotalCents / 100 // 1 point per dollar
		if points < 1 {
			points = 1
		}

		description := fmt.Sprintf("Points earned for completed order #%s", evt.OrderNumber)
		_, err := pointsUC.EarnPoints(context.Background(), usecase.EarnPointsRequest{
			UserID:      evt.BuyerID,
			Points:      points,
			Source:      domain.SourceOrder,
			ReferenceID: evt.OrderID,
			Description: description,
		})
		if err != nil {
			logger.Error().Err(err).
				Str("order_id", evt.OrderID).
				Str("buyer_id", evt.BuyerID).
				Msg("failed to award loyalty points from event")
		}
	}); err != nil {
		return fmt.Errorf("failed to subscribe to order.completed: %w", err)
	}

	logger.Info().Msg("loyalty NATS subscribers started")
	return nil
}
