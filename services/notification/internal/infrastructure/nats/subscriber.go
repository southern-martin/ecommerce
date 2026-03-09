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
	BuyerEmail  string `json:"buyer_email"`
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

// PaymentEvent matches the payment service's event payloads.
type PaymentEvent struct {
	PaymentID   string `json:"payment_id"`
	OrderID     string `json:"order_id"`
	BuyerID     string `json:"buyer_id"`
	AmountCents int64  `json:"amount_cents"`
	Currency    string `json:"currency"`
	Method      string `json:"method"`
	Reason      string `json:"reason"`
}

// ReturnEvent matches the return service's event payloads.
type ReturnEvent struct {
	ReturnID    string `json:"return_id"`
	OrderID     string `json:"order_id"`
	BuyerID     string `json:"buyer_id"`
	Reason      string `json:"reason"`
	AmountCents int64  `json:"amount_cents"`
	Currency    string `json:"currency"`
}

// ReviewEvent matches the review service's event payloads.
type ReviewEvent struct {
	ReviewID  string `json:"review_id"`
	ProductID string `json:"product_id"`
	UserID    string `json:"user_id"`
	Rating    int    `json:"rating"`
}

// withRetry wraps a handler so that if it returns an error, the message is
// published to the retry topic for later redelivery with exponential backoff.
// If retryPub is nil, errors are only logged (no retry).
func withRetry(retryPub *events.RetryPublisher, subject string, logger zerolog.Logger, handler func(data []byte) error) func(data []byte) {
	return func(data []byte) {
		if err := handler(data); err != nil {
			if retryPub != nil {
				if retryErr := retryPub.PublishToRetry(subject, data, err); retryErr != nil {
					logger.Error().Err(retryErr).Str("subject", subject).Msg("failed to publish to retry topic")
				}
			} else {
				logger.Error().Err(err).Str("subject", subject).Msg("handler failed (no retry publisher configured)")
			}
		}
	}
}

// StartSubscriber subscribes to order events and creates notifications.
// The retryPub parameter is optional; pass nil to disable retry behaviour.
func StartSubscriber(sub *events.Subscriber, notificationUC *usecase.NotificationUseCase, logger zerolog.Logger, retryPub ...*events.RetryPublisher) error {
	var rp *events.RetryPublisher
	if len(retryPub) > 0 {
		rp = retryPub[0]
	}
	_ = rp // used by withRetry calls below
	// Subscribe to order.created — uses withRetry so that failures are retried
	// with exponential backoff before landing in the dead-letter queue.
	if err := sub.Subscribe(events.SubjectOrderCreated, "notification-service-order-created",
		withRetry(rp, events.SubjectOrderCreated, logger, func(data []byte) error {
			var evt OrderCreatedEvent
			if err := json.Unmarshal(data, &evt); err != nil {
				logger.Error().Err(err).Msg("failed to unmarshal order.created event")
				return nil // bad payload, don't retry
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
				return fmt.Errorf("failed to create order confirmation notification: %w", err)
			}

			// Send confirmation email if buyer email is available
			if evt.BuyerEmail != "" {
				emailHTML := buildOrderConfirmationEmail(evt.OrderNumber, itemCount, totalFormatted)
				_, emailErr := notificationUC.SendNotification(context.Background(), usecase.SendNotificationRequest{
					UserID:  evt.BuyerID,
					Type:    "order_update",
					Channel: "email",
					Subject: subject,
					Body:    emailHTML,
					Data:    evt.BuyerEmail,
				})
				if emailErr != nil {
					return fmt.Errorf("failed to send order confirmation email: %w", emailErr)
				}
			}

			return nil
		})); err != nil {
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

	// Subscribe to payment.completed
	if err := sub.Subscribe(events.SubjectPaymentCompleted, "notification-service-payment-completed", func(data []byte) {
		var evt PaymentEvent
		if err := json.Unmarshal(data, &evt); err != nil {
			logger.Error().Err(err).Msg("failed to unmarshal payment.completed event")
			return
		}

		logger.Info().
			Str("payment_id", evt.PaymentID).
			Str("order_id", evt.OrderID).
			Str("buyer_id", evt.BuyerID).
			Msg("received payment.completed event")

		totalFormatted := fmt.Sprintf("%s %.2f", evt.Currency, float64(evt.AmountCents)/100)
		subject := "Payment Confirmed"
		body := fmt.Sprintf("Your payment of %s for order %s has been successfully processed.", totalFormatted, evt.OrderID)

		_, err := notificationUC.SendNotification(context.Background(), usecase.SendNotificationRequest{
			UserID:  evt.BuyerID,
			Type:    "payment_update",
			Channel: "in_app",
			Subject: subject,
			Body:    body,
			Data:    evt.BuyerID,
		})
		if err != nil {
			logger.Error().Err(err).Str("payment_id", evt.PaymentID).Msg("failed to create payment confirmation notification")
		}
	}); err != nil {
		return fmt.Errorf("failed to subscribe to payment.completed: %w", err)
	}

	// Subscribe to payment.failed
	if err := sub.Subscribe(events.SubjectPaymentFailed, "notification-service-payment-failed", func(data []byte) {
		var evt PaymentEvent
		if err := json.Unmarshal(data, &evt); err != nil {
			logger.Error().Err(err).Msg("failed to unmarshal payment.failed event")
			return
		}

		logger.Info().
			Str("payment_id", evt.PaymentID).
			Str("order_id", evt.OrderID).
			Str("buyer_id", evt.BuyerID).
			Msg("received payment.failed event")

		totalFormatted := fmt.Sprintf("%s %.2f", evt.Currency, float64(evt.AmountCents)/100)
		subject := "Payment Failed"
		body := fmt.Sprintf("Your payment of %s for order %s could not be processed. Please try again or use a different payment method.", totalFormatted, evt.OrderID)

		_, err := notificationUC.SendNotification(context.Background(), usecase.SendNotificationRequest{
			UserID:  evt.BuyerID,
			Type:    "payment_update",
			Channel: "in_app",
			Subject: subject,
			Body:    body,
			Data:    evt.BuyerID,
		})
		if err != nil {
			logger.Error().Err(err).Str("payment_id", evt.PaymentID).Msg("failed to create payment failed notification")
		}
	}); err != nil {
		return fmt.Errorf("failed to subscribe to payment.failed: %w", err)
	}

	// Subscribe to return.approved
	if err := sub.Subscribe(events.SubjectReturnApproved, "notification-service-return-approved", func(data []byte) {
		var evt ReturnEvent
		if err := json.Unmarshal(data, &evt); err != nil {
			logger.Error().Err(err).Msg("failed to unmarshal return.approved event")
			return
		}

		logger.Info().
			Str("return_id", evt.ReturnID).
			Str("order_id", evt.OrderID).
			Str("buyer_id", evt.BuyerID).
			Msg("received return.approved event")

		subject := "Return Approved"
		body := fmt.Sprintf("Your return request %s for order %s has been approved. Please follow the return instructions.", evt.ReturnID, evt.OrderID)

		_, err := notificationUC.SendNotification(context.Background(), usecase.SendNotificationRequest{
			UserID:  evt.BuyerID,
			Type:    "return_update",
			Channel: "in_app",
			Subject: subject,
			Body:    body,
			Data:    evt.BuyerID,
		})
		if err != nil {
			logger.Error().Err(err).Str("return_id", evt.ReturnID).Msg("failed to create return approved notification")
		}
	}); err != nil {
		return fmt.Errorf("failed to subscribe to return.approved: %w", err)
	}

	// Subscribe to return.rejected
	if err := sub.Subscribe(events.SubjectReturnRejected, "notification-service-return-rejected", func(data []byte) {
		var evt ReturnEvent
		if err := json.Unmarshal(data, &evt); err != nil {
			logger.Error().Err(err).Msg("failed to unmarshal return.rejected event")
			return
		}

		logger.Info().
			Str("return_id", evt.ReturnID).
			Str("order_id", evt.OrderID).
			Str("buyer_id", evt.BuyerID).
			Msg("received return.rejected event")

		subject := "Return Request Declined"
		body := fmt.Sprintf("Your return request %s for order %s has been declined. Reason: %s", evt.ReturnID, evt.OrderID, evt.Reason)

		_, err := notificationUC.SendNotification(context.Background(), usecase.SendNotificationRequest{
			UserID:  evt.BuyerID,
			Type:    "return_update",
			Channel: "in_app",
			Subject: subject,
			Body:    body,
			Data:    evt.BuyerID,
		})
		if err != nil {
			logger.Error().Err(err).Str("return_id", evt.ReturnID).Msg("failed to create return rejected notification")
		}
	}); err != nil {
		return fmt.Errorf("failed to subscribe to return.rejected: %w", err)
	}

	// Subscribe to review.approved
	if err := sub.Subscribe(events.SubjectReviewApproved, "notification-service-review-approved", func(data []byte) {
		var evt ReviewEvent
		if err := json.Unmarshal(data, &evt); err != nil {
			logger.Error().Err(err).Msg("failed to unmarshal review.approved event")
			return
		}

		logger.Info().
			Str("review_id", evt.ReviewID).
			Str("product_id", evt.ProductID).
			Str("user_id", evt.UserID).
			Msg("received review.approved event")

		subject := "Your Review Has Been Published"
		body := fmt.Sprintf("Your review for product %s has been approved and is now visible to other shoppers. Thank you for your feedback!", evt.ProductID)

		_, err := notificationUC.SendNotification(context.Background(), usecase.SendNotificationRequest{
			UserID:  evt.UserID,
			Type:    "review_update",
			Channel: "in_app",
			Subject: subject,
			Body:    body,
			Data:    evt.UserID,
		})
		if err != nil {
			logger.Error().Err(err).Str("review_id", evt.ReviewID).Msg("failed to create review approved notification")
		}
	}); err != nil {
		return fmt.Errorf("failed to subscribe to review.approved: %w", err)
	}

	logger.Info().Msg("notification NATS subscribers started")
	return nil
}

// buildOrderConfirmationEmail creates an HTML email body for order confirmation.
func buildOrderConfirmationEmail(orderNumber string, itemCount int, totalFormatted string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"></head>
<body style="font-family:Arial,sans-serif;max-width:600px;margin:0 auto;padding:20px;">
  <div style="text-align:center;padding:20px 0;">
    <h1 style="color:#16a34a;">Order Confirmed!</h1>
  </div>
  <div style="background:#f9fafb;border-radius:12px;padding:24px;margin:16px 0;">
    <p style="margin:0 0 8px;font-size:16px;">Order <strong>#%s</strong></p>
    <p style="margin:0 0 8px;color:#6b7280;">%d item(s)</p>
    <p style="margin:0;font-size:20px;font-weight:bold;">Total: %s</p>
  </div>
  <p style="color:#6b7280;font-size:14px;">Thank you for your purchase! You can view your order details in your account.</p>
  <hr style="border:none;border-top:1px solid #e5e7eb;margin:24px 0;">
  <p style="color:#9ca3af;font-size:12px;text-align:center;">This is an automated message from the e-commerce platform.</p>
</body>
</html>`, orderNumber, itemCount, totalFormatted)
}
