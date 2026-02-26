package usecase

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"

	pkgerrors "github.com/southern-martin/ecommerce/pkg/errors"
	"github.com/southern-martin/ecommerce/services/auth/internal/domain"
)

// ResetPasswordInput holds the input data for password reset.
type ResetPasswordInput struct {
	Email       string `json:"email" binding:"required,email"`
	ResetToken  string `json:"reset_token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// ResetPasswordUseCase handles password reset.
type ResetPasswordUseCase struct {
	repo   domain.UserRepository
	logger zerolog.Logger
}

// NewResetPasswordUseCase creates a new ResetPasswordUseCase.
func NewResetPasswordUseCase(
	repo domain.UserRepository,
	logger zerolog.Logger,
) *ResetPasswordUseCase {
	return &ResetPasswordUseCase{
		repo:   repo,
		logger: logger,
	}
}

// Execute performs the password reset.
func (uc *ResetPasswordUseCase) Execute(ctx context.Context, input ResetPasswordInput) error {
	user, err := uc.repo.GetByEmail(ctx, input.Email)
	if err != nil || user == nil {
		return pkgerrors.NewNotFoundError("AUTH_USER_NOT_FOUND", "user not found")
	}

	// Validate reset token
	if user.ResetToken != input.ResetToken {
		return pkgerrors.NewUnauthorizedError("AUTH_INVALID_RESET_TOKEN", "invalid reset token")
	}

	// Check expiration
	if user.ResetTokenExp == nil || user.ResetTokenExp.Before(time.Now()) {
		return pkgerrors.NewUnauthorizedError("AUTH_RESET_TOKEN_EXPIRED", "reset token has expired")
	}

	// Hash new password
	hash, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		uc.logger.Error().Err(err).Msg("failed to hash new password")
		return pkgerrors.NewInternalError("AUTH_HASH_FAILED", "failed to process password")
	}

	// Update password
	if err := uc.repo.UpdatePassword(ctx, user.ID, string(hash)); err != nil {
		uc.logger.Error().Err(err).Msg("failed to update password")
		return pkgerrors.NewInternalError("AUTH_UPDATE_FAILED", "failed to update password")
	}

	// Clear reset token
	if err := uc.repo.ClearResetToken(ctx, user.ID); err != nil {
		uc.logger.Error().Err(err).Msg("failed to clear reset token")
	}

	return nil
}
