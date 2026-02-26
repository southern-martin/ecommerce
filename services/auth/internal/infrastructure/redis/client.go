package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// TokenBlacklist manages blacklisted JWT tokens in Redis.
type TokenBlacklist struct {
	client *redis.Client
}

// NewTokenBlacklist creates a new TokenBlacklist connected to the given Redis URL.
func NewTokenBlacklist(redisURL string) (*TokenBlacklist, error) {
	opts, err := redis.ParseURL("redis://" + redisURL)
	if err != nil {
		// Fallback: treat as host:port
		opts = &redis.Options{
			Addr: redisURL,
		}
	}

	client := redis.NewClient(opts)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &TokenBlacklist{client: client}, nil
}

// BlacklistToken adds a token to the blacklist with the given TTL.
func (tb *TokenBlacklist) BlacklistToken(ctx context.Context, token string, ttl time.Duration) error {
	return tb.client.Set(ctx, "blacklist:"+token, "1", ttl).Err()
}

// IsBlacklisted checks whether a token has been blacklisted.
func (tb *TokenBlacklist) IsBlacklisted(ctx context.Context, token string) (bool, error) {
	result, err := tb.client.Exists(ctx, "blacklist:"+token).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

// Close closes the Redis client connection.
func (tb *TokenBlacklist) Close() error {
	return tb.client.Close()
}
