package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/southern-martin/ecommerce/services/cart/internal/domain"
)

const (
	cartKeyPrefix = "cart:"
	cartTTL       = 30 * 24 * time.Hour // 30 days
)

// redisCartRepo implements domain.CartRepository using Redis as the primary store.
type redisCartRepo struct {
	client *redis.Client
}

// NewRedisCartRepository creates a new Redis-backed cart repository.
func NewRedisCartRepository(client *redis.Client) domain.CartRepository {
	return &redisCartRepo{client: client}
}

func cartKey(userID string) string {
	return fmt.Sprintf("%s%s", cartKeyPrefix, userID)
}

// GetCart retrieves the cart for a given user from Redis.
// Returns an empty cart if no cart exists.
func (r *redisCartRepo) GetCart(ctx context.Context, userID string) (*domain.Cart, error) {
	data, err := r.client.Get(ctx, cartKey(userID)).Bytes()
	if err == redis.Nil {
		return &domain.Cart{
			UserID:    userID,
			Items:     []domain.CartItem{},
			UpdatedAt: time.Now().UTC(),
		}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("redis get cart: %w", err)
	}

	var cart domain.Cart
	if err := json.Unmarshal(data, &cart); err != nil {
		return nil, fmt.Errorf("unmarshal cart: %w", err)
	}

	if cart.Items == nil {
		cart.Items = []domain.CartItem{}
	}

	return &cart, nil
}

// SaveCart persists the cart to Redis with a 30-day TTL.
func (r *redisCartRepo) SaveCart(ctx context.Context, cart *domain.Cart) error {
	data, err := json.Marshal(cart)
	if err != nil {
		return fmt.Errorf("marshal cart: %w", err)
	}

	if err := r.client.Set(ctx, cartKey(cart.UserID), data, cartTTL).Err(); err != nil {
		return fmt.Errorf("redis set cart: %w", err)
	}

	return nil
}

// DeleteCart removes the cart from Redis.
func (r *redisCartRepo) DeleteCart(ctx context.Context, userID string) error {
	if err := r.client.Del(ctx, cartKey(userID)).Err(); err != nil {
		return fmt.Errorf("redis delete cart: %w", err)
	}

	return nil
}
