package nats

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog"

	"github.com/southern-martin/ecommerce/pkg/events"
	"github.com/southern-martin/ecommerce/services/affiliate/internal/usecase"
)

// OrderCompletedEvent matches the order service's event payload for completed orders.
type OrderCompletedEvent struct {
	OrderID       string `json:"order_id"`
	BuyerID       string `json:"buyer_id"`
	TotalCents    int64  `json:"total_cents"`
	AffiliateCode string `json:"affiliate_code,omitempty"`
}

// StartSubscriber subscribes to order events and tracks affiliate conversions.
func StartSubscriber(sub *events.Subscriber, referralUC *usecase.ReferralUseCase, linkUC *usecase.LinkUseCase, logger zerolog.Logger) error {
	// Subscribe to order.completed — track conversion for affiliate referrals
	if err := sub.Subscribe(events.SubjectOrderCompleted, "affiliate-service-order-completed", func(data []byte) {
		var evt OrderCompletedEvent
		if err := json.Unmarshal(data, &evt); err != nil {
			logger.Error().Err(err).Msg("failed to unmarshal order.completed event")
			return
		}

		logger.Info().
			Str("order_id", evt.OrderID).
			Str("buyer_id", evt.BuyerID).
			Int64("total_cents", evt.TotalCents).
			Str("affiliate_code", evt.AffiliateCode).
			Msg("received order.completed event")

		// Only track conversion if affiliate code is present
		if evt.AffiliateCode == "" {
			logger.Debug().Str("order_id", evt.OrderID).Msg("no affiliate code, skipping conversion tracking")
			return
		}

		// Look up affiliate link by code
		link, err := linkUC.GetByCode(context.Background(), evt.AffiliateCode)
		if err != nil {
			logger.Warn().Err(err).
				Str("affiliate_code", evt.AffiliateCode).
				Msg("affiliate link not found for code, skipping conversion")
			return
		}

		_, err = referralUC.TrackConversion(context.Background(), usecase.TrackConversionRequest{
			LinkID:          link.ID,
			ReferredID:      evt.BuyerID,
			OrderID:         evt.OrderID,
			OrderTotalCents: evt.TotalCents,
		})
		if err != nil {
			logger.Error().Err(err).
				Str("order_id", evt.OrderID).
				Str("affiliate_code", evt.AffiliateCode).
				Msg("failed to track affiliate conversion from event")
		}
	}); err != nil {
		return fmt.Errorf("failed to subscribe to order.completed: %w", err)
	}

	logger.Info().Msg("affiliate NATS subscribers started")
	return nil
}
