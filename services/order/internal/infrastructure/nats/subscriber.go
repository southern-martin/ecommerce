package nats

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog"

	"github.com/southern-martin/ecommerce/pkg/events"
	"github.com/southern-martin/ecommerce/services/order/internal/domain"
	"github.com/southern-martin/ecommerce/services/order/internal/usecase"
)

// PaymentEvent matches the payment service's event payloads.
type PaymentEvent struct {
	OrderID     string `json:"order_id"`
	PaymentID   string `json:"payment_id"`
	AmountCents int64  `json:"amount_cents"`
	Currency    string `json:"currency"`
	Status      string `json:"status"`
}

// ShipmentEvent matches the shipping service's shipment event payloads.
type ShipmentEvent struct {
	ShipmentID     string `json:"shipment_id"`
	OrderID        string `json:"order_id"`
	TrackingNumber string `json:"tracking_number"`
	Status         string `json:"status"`
}

// StartSubscriber subscribes to payment and shipping events and updates order statuses.
func StartSubscriber(sub *events.Subscriber, updateStatusUC *usecase.UpdateOrderStatusUseCase, logger zerolog.Logger) error {
	// Subscribe to payment.completed — update order status to "confirmed"
	if err := sub.Subscribe(events.SubjectPaymentCompleted, "order-service-payment-completed", func(data []byte) {
		var evt PaymentEvent
		if err := json.Unmarshal(data, &evt); err != nil {
			logger.Error().Err(err).Msg("failed to unmarshal payment.completed event")
			return
		}

		logger.Info().
			Str("order_id", evt.OrderID).
			Str("payment_id", evt.PaymentID).
			Int64("amount_cents", evt.AmountCents).
			Msg("received payment.completed event")

		_, err := updateStatusUC.UpdateOrderStatus(context.Background(), evt.OrderID, domain.OrderStatusConfirmed)
		if err != nil {
			logger.Error().Err(err).
				Str("order_id", evt.OrderID).
				Msg("failed to update order status to confirmed")
		}
	}); err != nil {
		return fmt.Errorf("failed to subscribe to payment.completed: %w", err)
	}

	// Subscribe to payment.failed — update order status to "cancelled"
	if err := sub.Subscribe(events.SubjectPaymentFailed, "order-service-payment-failed", func(data []byte) {
		var evt PaymentEvent
		if err := json.Unmarshal(data, &evt); err != nil {
			logger.Error().Err(err).Msg("failed to unmarshal payment.failed event")
			return
		}

		logger.Info().
			Str("order_id", evt.OrderID).
			Str("payment_id", evt.PaymentID).
			Msg("received payment.failed event")

		_, err := updateStatusUC.UpdateOrderStatus(context.Background(), evt.OrderID, domain.OrderStatusCancelled)
		if err != nil {
			logger.Error().Err(err).
				Str("order_id", evt.OrderID).
				Msg("failed to update order status to cancelled")
		}
	}); err != nil {
		return fmt.Errorf("failed to subscribe to payment.failed: %w", err)
	}

	// Subscribe to shipping.shipment.delivered — update order status to "delivered"
	if err := sub.Subscribe(events.SubjectShipmentDelivered, "order-service-shipment-delivered", func(data []byte) {
		var evt ShipmentEvent
		if err := json.Unmarshal(data, &evt); err != nil {
			logger.Error().Err(err).Msg("failed to unmarshal shipping.shipment.delivered event")
			return
		}

		logger.Info().
			Str("order_id", evt.OrderID).
			Str("shipment_id", evt.ShipmentID).
			Str("tracking_number", evt.TrackingNumber).
			Msg("received shipping.shipment.delivered event")

		_, err := updateStatusUC.UpdateOrderStatus(context.Background(), evt.OrderID, domain.OrderStatusDelivered)
		if err != nil {
			logger.Error().Err(err).
				Str("order_id", evt.OrderID).
				Msg("failed to update order status to delivered")
		}
	}); err != nil {
		return fmt.Errorf("failed to subscribe to shipping.shipment.delivered: %w", err)
	}

	logger.Info().Msg("order NATS subscribers started")
	return nil
}
