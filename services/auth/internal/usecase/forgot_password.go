package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/rs/zerolog"

	pkgerrors "github.com/southern-martin/ecommerce/pkg/errors"
	"github.com/southern-martin/ecommerce/services/auth/internal/domain"
	authnats "github.com/southern-martin/ecommerce/services/auth/internal/infrastructure/nats"
)

// ForgotPasswordInput holds the input data for password reset request.
type ForgotPasswordInput struct {
	Email string `json:"email" binding:"required,email"`
}

// ForgotPasswordUseCase handles forgot password requests.
type ForgotPasswordUseCase struct {
	repo      domain.UserRepository
	publisher *authnats.EventPublisher
	logger    zerolog.Logger
}

// NewForgotPasswordUseCase creates a new ForgotPasswordUseCase.
func NewForgotPasswordUseCase(
	repo domain.UserRepository,
	publisher *authnats.EventPublisher,
	logger zerolog.Logger,
) *ForgotPasswordUseCase {
	return &ForgotPasswordUseCase{
		repo:      repo,
		publisher: publisher,
		logger:    logger,
	}
}

// Execute performs the forgot password flow.
func (uc *ForgotPasswordUseCase) Execute(ctx context.Context, input ForgotPasswordInput) error {
	user, err := uc.repo.GetByEmail(ctx, input.Email)
	if err != nil || user == nil {
		// Return nil to avoid leaking whether the email exists
		return nil
	}

	// Generate random reset token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		uc.logger.Error().Err(err).Msg("failed to generate reset token")
		return pkgerrors.NewInternalError("AUTH_RESET_TOKEN_FAILED", "failed to generate reset token")
	}
	resetToken := hex.EncodeToString(tokenBytes)

	// Store reset token with 1 hour expiry
	expiry := time.Now().Add(1 * time.Hour)
	if err := uc.repo.UpdateResetToken(ctx, user.ID, resetToken, expiry); err != nil {
		uc.logger.Error().Err(err).Msg("failed to store reset token")
		return pkgerrors.NewInternalError("AUTH_RESET_TOKEN_FAILED", "failed to process request")
	}

	// Publish event
	if uc.publisher != nil {
		evt := domain.PasswordResetRequestedEvent{
			UserID:     user.ID,
			Email:      user.Email,
			ResetToken: resetToken,
		}
		if err := uc.publisher.PublishPasswordResetRequested(evt); err != nil {
			uc.logger.Error().Err(err).Msg("failed to publish password.reset.requested event")
		}
	}

	return nil
}
