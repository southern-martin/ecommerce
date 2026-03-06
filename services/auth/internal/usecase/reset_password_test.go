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

func validResetUser() *domain.AuthUser {
	exp := time.Now().Add(1 * time.Hour)
	return &domain.AuthUser{
		ID:            "user-1",
		Email:         "test@example.com",
		ResetToken:    "valid-reset-token",
		ResetTokenExp: &exp,
	}
}

func TestResetPassword_Success(t *testing.T) {
	var updatedHash string
	var clearedID string
	repo := &mockUserRepo{
		getByEmailFn: func(_ context.Context, _ string) (*domain.AuthUser, error) {
			return validResetUser(), nil
		},
		updatePasswordFn: func(_ context.Context, _ string, hash string) error {
			updatedHash = hash
			return nil
		},
		clearResetTokenFn: func(_ context.Context, id string) error {
			clearedID = id
			return nil
		},
	}
	uc := NewResetPasswordUseCase(repo, zerolog.Nop())

	err := uc.Execute(context.Background(), ResetPasswordInput{
		Email:       "test@example.com",
		ResetToken:  "valid-reset-token",
		NewPassword: "newpassword123",
	})
	require.NoError(t, err)
	assert.NotEmpty(t, updatedHash)
	assert.Equal(t, "user-1", clearedID)
	// Verify hash is valid bcrypt
	assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(updatedHash), []byte("newpassword123")))
}

func TestResetPassword_UserNotFound(t *testing.T) {
	repo := &mockUserRepo{
		getByEmailFn: func(_ context.Context, _ string) (*domain.AuthUser, error) {
			return nil, errors.New("not found")
		},
	}
	uc := NewResetPasswordUseCase(repo, zerolog.Nop())

	err := uc.Execute(context.Background(), ResetPasswordInput{
		Email: "nobody@example.com", ResetToken: "abc", NewPassword: "newpass",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
}

func TestResetPassword_WrongToken(t *testing.T) {
	repo := &mockUserRepo{
		getByEmailFn: func(_ context.Context, _ string) (*domain.AuthUser, error) {
			return validResetUser(), nil
		},
	}
	uc := NewResetPasswordUseCase(repo, zerolog.Nop())

	err := uc.Execute(context.Background(), ResetPasswordInput{
		Email: "test@example.com", ResetToken: "wrong-token", NewPassword: "newpass",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid reset token")
}

func TestResetPassword_ExpiredToken(t *testing.T) {
	pastExp := time.Now().Add(-1 * time.Hour)
	repo := &mockUserRepo{
		getByEmailFn: func(_ context.Context, _ string) (*domain.AuthUser, error) {
			return &domain.AuthUser{
				ID:            "user-1",
				Email:         "test@example.com",
				ResetToken:    "valid-reset-token",
				ResetTokenExp: &pastExp,
			}, nil
		},
	}
	uc := NewResetPasswordUseCase(repo, zerolog.Nop())

	err := uc.Execute(context.Background(), ResetPasswordInput{
		Email: "test@example.com", ResetToken: "valid-reset-token", NewPassword: "newpass",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "reset token has expired")
}

func TestResetPassword_NilExpiry(t *testing.T) {
	repo := &mockUserRepo{
		getByEmailFn: func(_ context.Context, _ string) (*domain.AuthUser, error) {
			return &domain.AuthUser{
				ID:            "user-1",
				Email:         "test@example.com",
				ResetToken:    "valid-reset-token",
				ResetTokenExp: nil,
			}, nil
		},
	}
	uc := NewResetPasswordUseCase(repo, zerolog.Nop())

	err := uc.Execute(context.Background(), ResetPasswordInput{
		Email: "test@example.com", ResetToken: "valid-reset-token", NewPassword: "newpass",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "reset token has expired")
}

func TestResetPassword_NewPasswordHashed(t *testing.T) {
	var savedHash string
	repo := &mockUserRepo{
		getByEmailFn: func(_ context.Context, _ string) (*domain.AuthUser, error) {
			return validResetUser(), nil
		},
		updatePasswordFn: func(_ context.Context, _ string, hash string) error {
			savedHash = hash
			return nil
		},
	}
	uc := NewResetPasswordUseCase(repo, zerolog.Nop())

	err := uc.Execute(context.Background(), ResetPasswordInput{
		Email: "test@example.com", ResetToken: "valid-reset-token", NewPassword: "secretpass",
	})
	require.NoError(t, err)
	assert.NotEqual(t, "secretpass", savedHash)
	assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(savedHash), []byte("secretpass")))
}
