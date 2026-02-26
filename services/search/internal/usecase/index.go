package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/southern-martin/ecommerce/services/search/internal/domain"
)

// IndexUseCase handles indexing operations for the search service.
type IndexUseCase struct {
	repo      domain.SearchRepository
	publisher domain.EventPublisher
}

// NewIndexUseCase creates a new IndexUseCase.
func NewIndexUseCase(repo domain.SearchRepository, publisher domain.EventPublisher) *IndexUseCase {
	return &IndexUseCase{
		repo:      repo,
		publisher: publisher,
	}
}

// IndexProduct adds or updates a product in the search index.
func (uc *IndexUseCase) IndexProduct(ctx context.Context, idx *domain.SearchIndex) error {
	if idx.ID == "" {
		idx.ID = uuid.New().String()
	}
	now := time.Now()
	if idx.CreatedAt.IsZero() {
		idx.CreatedAt = now
	}
	idx.UpdatedAt = now

	if err := uc.repo.Index(ctx, idx); err != nil {
		return err
	}

	// Publish event
	event := map[string]interface{}{
		"event":      "product.indexed",
		"product_id": idx.ProductID,
		"name":       idx.Name,
		"timestamp":  now,
	}
	if err := uc.publisher.Publish(ctx, "search.product.indexed", event); err != nil {
		log.Warn().Err(err).Str("product_id", idx.ProductID).Msg("failed to publish index event")
	}

	return nil
}

// RemoveProduct removes a product from the search index.
func (uc *IndexUseCase) RemoveProduct(ctx context.Context, productID string) error {
	if err := uc.repo.Delete(ctx, productID); err != nil {
		return err
	}

	// Publish event
	event := map[string]interface{}{
		"event":      "product.deindexed",
		"product_id": productID,
		"timestamp":  time.Now(),
	}
	if err := uc.publisher.Publish(ctx, "search.product.deindexed", event); err != nil {
		log.Warn().Err(err).Str("product_id", productID).Msg("failed to publish deindex event")
	}

	return nil
}
