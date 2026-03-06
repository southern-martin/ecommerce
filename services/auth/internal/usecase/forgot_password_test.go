package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/southern-martin/ecommerce/services/auth/internal/domain"
)

func TestForgotPassword_ExistingUser(t *testing.T) {
	var storedToken string
	var storedExpiry time.Time

	repo := &mockUserRepo{
		getByEmailFn: func(_ context.Context, _ string) (*domain.AuthUser, error) {
			return &domain.AuthUser{ID: "user-1", Email: "test@example.com"}, nil
		},
		updateResetTokenFn: func(_ context.Context, _ string, token string, exp time.Time) error {
			storedToken = token
			storedExpiry = exp
			return nil
		},
	}
	uc := NewForgotPasswordUseCase(repo, nil, zerolog.Nop())

	err := uc.Execute(context.Background(), ForgotPasswordInput{Email: "test@example.com"})
	require.NoError(t, err)

	// Token should be 64 hex chars (32 bytes)
	assert.Len(t, storedToken, 64)
	// Expiry should be roughly 1 hour from now
	assert.WithinDuration(t, time.Now().Add(1*time.Hour), storedExpiry, 5*time.Second)
}

func TestForgotPassword_NonExistentEmail(t *testing.T) {
	repo := &mockUserRepo{
		getByEmailFn: func(_ context.Context, _ string) (*domain.AuthUser, error) {
			return nil, errors.New("not found")
		},
	}
	uc := NewForgotPasswordUseCase(repo, nil, zerolog.Nop())

	// Should return nil to avoid leaking email existence
	err := uc.Execute(context.Background(), ForgotPasswordInput{Email: "nobody@example.com"})
	assert.NoError(t, err)
}

func TestForgotPassword_RepoError_NoLeak(t *testing.T) {
	repo := &mockUserRepo{
		getByEmailFn: func(_ context.Context, _ string) (*domain.AuthUser, error) {
			return nil, errors.New("db error")
		},
	}
	uc := NewForgotPasswordUseCase(repo, nil, zerolog.Nop())

	// Should still return nil — don't leak that an error occurred
	err := uc.Execute(context.Background(), ForgotPasswordInput{Email: "a@b.com"})
	assert.NoError(t, err)
}

func TestForgotPassword_StoreTokenError(t *testing.T) {
	repo := &mockUserRepo{
		getByEmailFn: func(_ context.Context, _ string) (*domain.AuthUser, error) {
			return &domain.AuthUser{ID: "user-1", Email: "test@example.com"}, nil
		},
		updateResetTokenFn: func(_ context.Context, _ string, _ string, _ time.Time) error {
			return errors.New("db write error")
		},
	}
	uc := NewForgotPasswordUseCase(repo, nil, zerolog.Nop())

	err := uc.Execute(context.Background(), ForgotPasswordInput{Email: "test@example.com"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to process request")
}
