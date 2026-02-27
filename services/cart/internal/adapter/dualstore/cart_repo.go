package dualstore

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"

	"github.com/southern-martin/ecommerce/services/cart/internal/domain"
)

// dualStoreCartRepo implements domain.CartRepository using Redis as a read-through
// cache and PostgreSQL as the durable backing store.
//
// Write path: write to Postgres first (source of truth), then update Redis cache.
// Read path: try Redis first; on cache miss, load from Postgres and populate cache.
// Delete path: delete from both stores.
type dualStoreCartRepo struct {
	redis  domain.CartRepository
	pg     domain.CartRepository
	logger zerolog.Logger
}

// NewDualStoreCartRepository creates a repository that writes to both Redis and Postgres.
// Redis serves as a fast cache; Postgres provides permanent persistence.
func NewDualStoreCartRepository(redis, pg domain.CartRepository, logger zerolog.Logger) domain.CartRepository {
	return &dualStoreCartRepo{
		redis:  redis,
		pg:     pg,
		logger: logger.With().Str("component", "dualstore_repo").Logger(),
	}
}

func (d *dualStoreCartRepo) GetCart(ctx context.Context, userID string) (*domain.Cart, error) {
	// Try Redis first (fast path)
	cart, err := d.redis.GetCart(ctx, userID)
	if err != nil {
		d.logger.Warn().Err(err).Str("user_id", userID).Msg("redis read failed, falling back to postgres")
	} else if len(cart.Items) > 0 {
		return cart, nil
	}

	// Cache miss or empty cart in Redis — check Postgres (durable store)
	pgCart, pgErr := d.pg.GetCart(ctx, userID)
	if pgErr != nil {
		return nil, fmt.Errorf("dualstore get cart: %w", pgErr)
	}

	// If Postgres has items, repopulate the Redis cache
	if len(pgCart.Items) > 0 {
		if cacheErr := d.redis.SaveCart(ctx, pgCart); cacheErr != nil {
			d.logger.Warn().Err(cacheErr).Str("user_id", userID).Msg("failed to repopulate redis cache")
		}
	}

	return pgCart, nil
}

func (d *dualStoreCartRepo) SaveCart(ctx context.Context, cart *domain.Cart) error {
	// Write to Postgres first (source of truth)
	if err := d.pg.SaveCart(ctx, cart); err != nil {
		return fmt.Errorf("dualstore save postgres: %w", err)
	}

	// Then update Redis cache (best-effort)
	if err := d.redis.SaveCart(ctx, cart); err != nil {
		d.logger.Warn().Err(err).Str("user_id", cart.UserID).Msg("failed to update redis cache after postgres save")
	}

	return nil
}

func (d *dualStoreCartRepo) DeleteCart(ctx context.Context, userID string) error {
	// Delete from Postgres first
	if err := d.pg.DeleteCart(ctx, userID); err != nil {
		return fmt.Errorf("dualstore delete postgres: %w", err)
	}

	// Then evict from Redis cache (best-effort)
	if err := d.redis.DeleteCart(ctx, userID); err != nil {
		d.logger.Warn().Err(err).Str("user_id", userID).Msg("failed to evict redis cache after postgres delete")
	}

	return nil
}
