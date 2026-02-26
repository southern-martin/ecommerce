package usecase

import (
	"context"
	"time"

	"github.com/rs/zerolog"

	pkgauth "github.com/southern-martin/ecommerce/pkg/auth"
	pkgerrors "github.com/southern-martin/ecommerce/pkg/errors"
	"github.com/southern-martin/ecommerce/services/auth/internal/domain"
	authnats "github.com/southern-martin/ecommerce/services/auth/internal/infrastructure/nats"
)

// OAuthLoginInput holds the input data for OAuth login.
type OAuthLoginInput struct {
	Provider   string `json:"provider" binding:"required"`
	ProviderID string `json:"provider_id" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
}

// OAuthLoginOutput holds the output data after successful OAuth login.
type OAuthLoginOutput struct {
	UserID       string `json:"user_id"`
	Email        string `json:"email"`
	Role         string `json:"role"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IsNew        bool   `json:"is_new"`
}

// OAuthLoginUseCase handles OAuth-based authentication.
type OAuthLoginUseCase struct {
	repo          domain.UserRepository
	publisher     *authnats.EventPublisher
	jwtSecret     string
	accessExpiry  time.Duration
	refreshExpiry time.Duration
	logger        zerolog.Logger
}

// NewOAuthLoginUseCase creates a new OAuthLoginUseCase.
func NewOAuthLoginUseCase(
	repo domain.UserRepository,
	publisher *authnats.EventPublisher,
	jwtSecret string,
	accessExpiry, refreshExpiry time.Duration,
	logger zerolog.Logger,
) *OAuthLoginUseCase {
	return &OAuthLoginUseCase{
		repo:          repo,
		publisher:     publisher,
		jwtSecret:     jwtSecret,
		accessExpiry:  accessExpiry,
		refreshExpiry: refreshExpiry,
		logger:        logger,
	}
}

// Execute performs the OAuth login or registration.
func (uc *OAuthLoginUseCase) Execute(ctx context.Context, input OAuthLoginInput) (*OAuthLoginOutput, error) {
	isNew := false

	user, _ := uc.repo.GetByOAuthProvider(ctx, input.Provider, input.ProviderID)
	if user == nil {
		// Create new user without password
		user = &domain.AuthUser{
			Email:           input.Email,
			Role:            "buyer",
			OAuthProvider:   input.Provider,
			OAuthProviderID: input.ProviderID,
		}
		if err := uc.repo.Create(ctx, user); err != nil {
			uc.logger.Error().Err(err).Msg("failed to create oauth user")
			return nil, pkgerrors.NewInternalError("AUTH_CREATE_FAILED", "failed to create user")
		}
		isNew = true

		// Publish event for new user
		if uc.publisher != nil {
			evt := domain.UserRegisteredEvent{
				UserID:    user.ID,
				Email:     user.Email,
				Role:      user.Role,
				CreatedAt: user.CreatedAt.Format(time.RFC3339),
			}
			if err := uc.publisher.PublishUserRegistered(evt); err != nil {
				uc.logger.Error().Err(err).Msg("failed to publish user.registered event")
			}
		}
	}

	// Generate tokens
	accessToken, err := pkgauth.GenerateAccessToken(user.ID, user.Email, user.Role, uc.jwtSecret, uc.accessExpiry)
	if err != nil {
		uc.logger.Error().Err(err).Msg("failed to generate access token")
		return nil, pkgerrors.NewInternalError("AUTH_TOKEN_FAILED", "failed to generate tokens")
	}

	refreshToken, err := pkgauth.GenerateRefreshToken(user.ID, uc.jwtSecret, uc.refreshExpiry)
	if err != nil {
		uc.logger.Error().Err(err).Msg("failed to generate refresh token")
		return nil, pkgerrors.NewInternalError("AUTH_TOKEN_FAILED", "failed to generate tokens")
	}

	if err := uc.repo.UpdateRefreshToken(ctx, user.ID, refreshToken); err != nil {
		uc.logger.Error().Err(err).Msg("failed to store refresh token")
	}

	return &OAuthLoginOutput{
		UserID:       user.ID,
		Email:        user.Email,
		Role:         user.Role,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		IsNew:        isNew,
	}, nil
}
