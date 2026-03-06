package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	pkgauth "github.com/southern-martin/ecommerce/pkg/auth"
	"github.com/southern-martin/ecommerce/services/auth/internal/domain"
)

func generateTestRefreshToken(t *testing.T, userID string) string {
	t.Helper()
	token, err := pkgauth.GenerateRefreshToken(userID, testJWTSecret, 24*time.Hour)
	require.NoError(t, err)
	return token
}

func TestRefreshToken_Success(t *testing.T) {
	validToken := generateTestRefreshToken(t, "user-1")

	repo := &mockUserRepo{
		getByIDFn: func(_ context.Context, id string) (*domain.AuthUser, error) {
			return &domain.AuthUser{
				ID:           "user-1",
				Email:        "test@example.com",
				Role:         "buyer",
				RefreshToken: validToken,
			}, nil
		},
	}
	uc := NewRefreshTokenUseCase(repo, testJWTSecret, 15*time.Minute, 7*24*time.Hour, zerolog.Nop())

	out, err := uc.Execute(context.Background(), RefreshTokenInput{RefreshToken: validToken})
	require.NoError(t, err)
	assert.NotEmpty(t, out.AccessToken)
	assert.NotEmpty(t, out.RefreshToken)
	// New tokens should differ from old
	assert.NotEqual(t, validToken, out.RefreshToken)
}

func TestRefreshToken_InvalidToken(t *testing.T) {
	uc := NewRefreshTokenUseCase(&mockUserRepo{}, testJWTSecret, 15*time.Minute, 7*24*time.Hour, zerolog.Nop())

	_, err := uc.Execute(context.Background(), RefreshTokenInput{RefreshToken: "garbage-token"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid refresh token")
}

func TestRefreshToken_TokenMismatch(t *testing.T) {
	validToken := generateTestRefreshToken(t, "user-1")

	repo := &mockUserRepo{
		getByIDFn: func(_ context.Context, _ string) (*domain.AuthUser, error) {
			return &domain.AuthUser{
				ID:           "user-1",
				Email:        "test@example.com",
				RefreshToken: "different-stored-token",
			}, nil
		},
	}
	uc := NewRefreshTokenUseCase(repo, testJWTSecret, 15*time.Minute, 7*24*time.Hour, zerolog.Nop())

	_, err := uc.Execute(context.Background(), RefreshTokenInput{RefreshToken: validToken})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "refresh token mismatch")
}

func TestRefreshToken_UserNotFound(t *testing.T) {
	validToken := generateTestRefreshToken(t, "user-1")

	repo := &mockUserRepo{
		getByIDFn: func(_ context.Context, _ string) (*domain.AuthUser, error) {
			return nil, errors.New("not found")
		},
	}
	uc := NewRefreshTokenUseCase(repo, testJWTSecret, 15*time.Minute, 7*24*time.Hour, zerolog.Nop())

	_, err := uc.Execute(context.Background(), RefreshTokenInput{RefreshToken: validToken})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
}

func TestRefreshToken_ExpiredToken(t *testing.T) {
	// Generate a token with -1 hour expiry (already expired)
	token, err := pkgauth.GenerateRefreshToken("user-1", testJWTSecret, -1*time.Hour)
	require.NoError(t, err)

	uc := NewRefreshTokenUseCase(&mockUserRepo{}, testJWTSecret, 15*time.Minute, 7*24*time.Hour, zerolog.Nop())

	_, err = uc.Execute(context.Background(), RefreshTokenInput{RefreshToken: token})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid refresh token")
}
