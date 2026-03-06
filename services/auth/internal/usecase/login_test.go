package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"github.com/southern-martin/ecommerce/services/auth/internal/domain"
)

func hashPassword(t *testing.T, password string) string {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	require.NoError(t, err)
	return string(hash)
}

func TestLogin_Success(t *testing.T) {
	repo := &mockUserRepo{
		getByEmailFn: func(_ context.Context, _ string) (*domain.AuthUser, error) {
			return &domain.AuthUser{
				ID:           "user-1",
				Email:        "test@example.com",
				PasswordHash: hashPassword(t, "correctpassword"),
				Role:         "buyer",
			}, nil
		},
	}
	uc := NewLoginUseCase(repo, testJWTSecret, 15*time.Minute, 7*24*time.Hour, zerolog.Nop())

	out, err := uc.Execute(context.Background(), LoginInput{Email: "test@example.com", Password: "correctpassword"})
	require.NoError(t, err)
	assert.Equal(t, "user-1", out.UserID)
	assert.Equal(t, "test@example.com", out.Email)
	assert.NotEmpty(t, out.AccessToken)
	assert.NotEmpty(t, out.RefreshToken)
}

func TestLogin_UserNotFound(t *testing.T) {
	repo := &mockUserRepo{
		getByEmailFn: func(_ context.Context, _ string) (*domain.AuthUser, error) {
			return nil, errors.New("not found")
		},
	}
	uc := NewLoginUseCase(repo, testJWTSecret, 15*time.Minute, 7*24*time.Hour, zerolog.Nop())

	_, err := uc.Execute(context.Background(), LoginInput{Email: "nobody@example.com", Password: "pass"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid email or password")
}

func TestLogin_WrongPassword(t *testing.T) {
	repo := &mockUserRepo{
		getByEmailFn: func(_ context.Context, _ string) (*domain.AuthUser, error) {
			return &domain.AuthUser{
				ID:           "user-1",
				Email:        "test@example.com",
				PasswordHash: hashPassword(t, "correctpassword"),
				Role:         "buyer",
			}, nil
		},
	}
	uc := NewLoginUseCase(repo, testJWTSecret, 15*time.Minute, 7*24*time.Hour, zerolog.Nop())

	_, err := uc.Execute(context.Background(), LoginInput{Email: "test@example.com", Password: "wrongpassword"})
	assert.Error(t, err)
	// Same error message as "not found" to prevent info leaking
	assert.Contains(t, err.Error(), "invalid email or password")
}

func TestLogin_RepoError(t *testing.T) {
	repo := &mockUserRepo{
		getByEmailFn: func(_ context.Context, _ string) (*domain.AuthUser, error) {
			return nil, errors.New("connection refused")
		},
	}
	uc := NewLoginUseCase(repo, testJWTSecret, 15*time.Minute, 7*24*time.Hour, zerolog.Nop())

	_, err := uc.Execute(context.Background(), LoginInput{Email: "a@b.com", Password: "pass"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid email or password")
}
