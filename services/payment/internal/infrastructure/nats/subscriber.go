package nats

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog"

	"github.com/southern-martin/ecommerce/pkg/events"
	"github.com/southern-martin/ecommerce/services/payment/internal/usecase"
)

// ReturnApprovedEvent represents the payload from a return.approved event.
type ReturnApprovedEvent struct {
	ReturnID          string `json:"return_id"`
	OrderID           string `json:"order_id"`
	RefundAmountCents int64  `json:"refund_amount_cents"`
}

// StartSubscriber subscribes to events relevant to the payment service.
func StartSubscriber(sub *events.Subscriber, refundUC *usecase.RefundUseCase, logger zerolog.Logger) error {
	// Subscribe to return.approved — initiate a refund when a return is approved
	if err := sub.Subscribe(events.SubjectReturnApproved, "payment-service-return-approved", func(data []byte) {
		var evt ReturnApprovedEvent
		if err := json.Unmarshal(data, &evt); err != nil {
			logger.Error().Err(err).Msg("failed to unmarshal return.approved event")
			return
		}

		logger.Info().
			Str("return_id", evt.ReturnID).
			Str("order_id", evt.OrderID).
			Int64("refund_amount_cents", evt.RefundAmountCents).
			Msg("received return.approved event")

		input := usecase.RefundInput{
			OrderID:     evt.OrderID,
			AmountCents: evt.RefundAmountCents,
		}

		if err := refundUC.ProcessRefund(context.Background(), input); err != nil {
			logger.Error().Err(err).
				Str("return_id", evt.ReturnID).
				Str("order_id", evt.OrderID).
				Msg("failed to process refund from return.approved event")
			return
		}

		logger.Info().
			Str("return_id", evt.ReturnID).
			Str("order_id", evt.OrderID).
			Msg("refund initiated from return approval")
	}); err != nil {
		return fmt.Errorf("subscribe to %s: %w", events.SubjectReturnApproved, err)
	}

	logger.Info().Msg("payment NATS subscribers started")
	return nil
}
