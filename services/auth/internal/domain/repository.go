package domain

import (
	"context"
	"time"
)

// UserRepository defines the persistence operations for AuthUser.
type UserRepository interface {
	Create(ctx context.Context, user *AuthUser) error
	GetByID(ctx context.Context, id string) (*AuthUser, error)
	GetByEmail(ctx context.Context, email string) (*AuthUser, error)
	GetByOAuthProvider(ctx context.Context, provider, providerID string) (*AuthUser, error)
	UpdateRefreshToken(ctx context.Context, id, token string) error
	UpdatePassword(ctx context.Context, id, passwordHash string) error
	UpdateResetToken(ctx context.Context, id, token string, exp time.Time) error
	ClearResetToken(ctx context.Context, id string) error
	UpdateRole(ctx context.Context, id, role string) error
}
