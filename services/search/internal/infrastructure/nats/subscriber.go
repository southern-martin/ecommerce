package nats

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog"

	"github.com/southern-martin/ecommerce/pkg/events"
	"github.com/southern-martin/ecommerce/services/search/internal/domain"
	"github.com/southern-martin/ecommerce/services/search/internal/usecase"
)

// ProductEvent matches the product service's event payloads.
type ProductEvent struct {
	ID          string  `json:"id"`
	SellerID    string  `json:"seller_id"`
	Name        string  `json:"name"`
	Slug        string  `json:"slug"`
	Description string  `json:"description"`
	CategoryID  string  `json:"category_id"`
	Status      string  `json:"status"`
	PriceCents  int64   `json:"price_cents"`
	Currency    string  `json:"currency"`
	ImageURL    string  `json:"image_url"`
	Rating      float64 `json:"rating"`
	InStock     bool    `json:"in_stock"`
}

// ProductDeletedEvent matches the product service's delete event payload.
type ProductDeletedEvent struct {
	ID string `json:"id"`
}

// StartSubscriber subscribes to product events and indexes/removes products.
func StartSubscriber(sub *events.Subscriber, indexUC *usecase.IndexUseCase, logger zerolog.Logger) error {
	// Subscribe to product.created
	if err := sub.Subscribe(events.SubjectProductCreated, "search-service-product-created", func(data []byte) {
		var evt ProductEvent
		if err := json.Unmarshal(data, &evt); err != nil {
			logger.Error().Err(err).Msg("failed to unmarshal product.created event")
			return
		}

		logger.Info().
			Str("product_id", evt.ID).
			Str("name", evt.Name).
			Msg("received product.created event")

		idx := &domain.SearchIndex{
			ProductID:   evt.ID,
			Name:        evt.Name,
			Slug:        evt.Slug,
			Description: evt.Description,
			CategoryID:  evt.CategoryID,
			SellerID:    evt.SellerID,
			PriceCents:  evt.PriceCents,
			Currency:    evt.Currency,
			ImageURL:    evt.ImageURL,
			Rating:      evt.Rating,
			InStock:     evt.InStock,
		}
		if err := indexUC.IndexProduct(context.Background(), idx); err != nil {
			logger.Error().Err(err).Str("product_id", evt.ID).Msg("failed to index product from event")
		}
	}); err != nil {
		return fmt.Errorf("failed to subscribe to product.created: %w", err)
	}

	// Subscribe to product.updated — upsert into index
	if err := sub.Subscribe(events.SubjectProductUpdated, "search-service-product-updated", func(data []byte) {
		var evt ProductEvent
		if err := json.Unmarshal(data, &evt); err != nil {
			logger.Error().Err(err).Msg("failed to unmarshal product.updated event")
			return
		}

		logger.Info().
			Str("product_id", evt.ID).
			Str("name", evt.Name).
			Msg("received product.updated event")

		idx := &domain.SearchIndex{
			ProductID:   evt.ID,
			Name:        evt.Name,
			Slug:        evt.Slug,
			Description: evt.Description,
			CategoryID:  evt.CategoryID,
			SellerID:    evt.SellerID,
			PriceCents:  evt.PriceCents,
			Currency:    evt.Currency,
			ImageURL:    evt.ImageURL,
			Rating:      evt.Rating,
			InStock:     evt.InStock,
		}
		if err := indexUC.IndexProduct(context.Background(), idx); err != nil {
			logger.Error().Err(err).Str("product_id", evt.ID).Msg("failed to index product from update event")
		}
	}); err != nil {
		return fmt.Errorf("failed to subscribe to product.updated: %w", err)
	}

	// Subscribe to product.deleted — remove from index
	if err := sub.Subscribe(events.SubjectProductDeleted, "search-service-product-deleted", func(data []byte) {
		var evt ProductDeletedEvent
		if err := json.Unmarshal(data, &evt); err != nil {
			logger.Error().Err(err).Msg("failed to unmarshal product.deleted event")
			return
		}

		logger.Info().
			Str("product_id", evt.ID).
			Msg("received product.deleted event")

		if err := indexUC.RemoveProduct(context.Background(), evt.ID); err != nil {
			logger.Error().Err(err).Str("product_id", evt.ID).Msg("failed to remove product from index")
		}
	}); err != nil {
		return fmt.Errorf("failed to subscribe to product.deleted: %w", err)
	}

	logger.Info().Msg("search NATS subscribers started")
	return nil
}
