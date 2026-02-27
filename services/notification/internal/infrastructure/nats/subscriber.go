package nats

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog"

	"github.com/southern-martin/ecommerce/pkg/events"
	"github.com/southern-martin/ecommerce/services/notification/internal/usecase"
)

// OrderCreatedEvent matches the order service's event payload.
type OrderCreatedEvent struct {
	OrderID     string `json:"order_id"`
	OrderNumber string `json:"order_number"`
	BuyerID     string `json:"buyer_id"`
	TotalCents  int64  `json:"total_cents"`
	Currency    string `json:"currency"`
	Items       []struct {
		ProductID      string `json:"product_id"`
		VariantID      string `json:"variant_id"`
		Quantity       int    `json:"quantity"`
		UnitPriceCents int64  `json:"unit_price_cents"`
		SellerID       string `json:"seller_id"`
	} `json:"items"`
}

// OrderStatusEvent matches the order service's status change event payload.
type OrderStatusEvent struct {
	OrderID     string `json:"order_id"`
	OrderNumber string `json:"order_number"`
	BuyerID     string `json:"buyer_id"`
	Status      string `json:"status"`
}

// StartSubscriber subscribes to order events and creates notifications.
func StartSubscriber(sub *events.Subscriber, notificationUC *usecase.NotificationUseCase, logger zerolog.Logger) error {
	// Subscribe to order.created
	if err := sub.Subscribe(events.SubjectOrderCreated, "notification-service-order-created", func(data []byte) {
		var evt OrderCreatedEvent
		if err := json.Unmarshal(data, &evt); err != nil {
			logger.Error().Err(err).Msg("failed to unmarshal order.created event")
			return
		}

		logger.Info().
			Str("order_id", evt.OrderID).
			Str("order_number", evt.OrderNumber).
			Str("buyer_id", evt.BuyerID).
			Msg("received order.created event")

		itemCount := 0
		for _, item := range evt.Items {
			itemCount += item.Quantity
		}

		totalFormatted := fmt.Sprintf("%s %.2f", evt.Currency, float64(evt.TotalCents)/100)
		subject := fmt.Sprintf("Order #%s Confirmed", evt.OrderNumber)
		body := fmt.Sprintf("Your order #%s with %d item(s) totaling %s has been confirmed.", evt.OrderNumber, itemCount, totalFormatted)

		// Create in_app notification
		_, err := notificationUC.SendNotification(context.Background(), usecase.SendNotificationRequest{
			UserID:  evt.BuyerID,
			Type:    "order_update",
			Channel: "in_app",
			Subject: subject,
			Body:    body,
			Data:    evt.BuyerID,
		})
		if err != nil {
			logger.Error().Err(err).Str("order_id", evt.OrderID).Msg("failed to create order confirmation notification")
		}
	}); err != nil {
		return fmt.Errorf("failed to subscribe to order.created: %w", err)
	}

	// Subscribe to order.shipped
	if err := sub.Subscribe(events.SubjectOrderShipped, "notification-service-order-shipped", func(data []byte) {
		var evt OrderStatusEvent
		if err := json.Unmarshal(data, &evt); err != nil {
			logger.Error().Err(err).Msg("failed to unmarshal order.shipped event")
			return
		}

		logger.Info().
			Str("order_id", evt.OrderID).
			Str("order_number", evt.OrderNumber).
			Str("buyer_id", evt.BuyerID).
			Msg("received order.shipped event")

		subject := fmt.Sprintf("Order #%s Shipped", evt.OrderNumber)
		body := fmt.Sprintf("Your order #%s has been shipped and is on its way!", evt.OrderNumber)

		_, err := notificationUC.SendNotification(context.Background(), usecase.SendNotificationRequest{
			UserID:  evt.BuyerID,
			Type:    "shipment_update",
			Channel: "in_app",
			Subject: subject,
			Body:    body,
			Data:    evt.BuyerID,
		})
		if err != nil {
			logger.Error().Err(err).Str("order_id", evt.OrderID).Msg("failed to create order shipped notification")
		}
	}); err != nil {
		return fmt.Errorf("failed to subscribe to order.shipped: %w", err)
	}

	// Subscribe to order.delivered
	if err := sub.Subscribe(events.SubjectOrderDelivered, "notification-service-order-delivered", func(data []byte) {
		var evt OrderStatusEvent
		if err := json.Unmarshal(data, &evt); err != nil {
			logger.Error().Err(err).Msg("failed to unmarshal order.delivered event")
			return
		}

		logger.Info().
			Str("order_id", evt.OrderID).
			Str("order_number", evt.OrderNumber).
			Str("buyer_id", evt.BuyerID).
			Msg("received order.delivered event")

		subject := fmt.Sprintf("Order #%s Delivered", evt.OrderNumber)
		body := fmt.Sprintf("Your order #%s has been delivered. We hope you enjoy your purchase!", evt.OrderNumber)

		_, err := notificationUC.SendNotification(context.Background(), usecase.SendNotificationRequest{
			UserID:  evt.BuyerID,
			Type:    "shipment_update",
			Channel: "in_app",
			Subject: subject,
			Body:    body,
			Data:    evt.BuyerID,
		})
		if err != nil {
			logger.Error().Err(err).Str("order_id", evt.OrderID).Msg("failed to create order delivered notification")
		}
	}); err != nil {
		return fmt.Errorf("failed to subscribe to order.delivered: %w", err)
	}

	logger.Info().Msg("notification NATS subscribers started")
	return nil
}
