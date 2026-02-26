package usecase

import (
	"context"
	"time"

	"github.com/rs/zerolog"

	pkgauth "github.com/southern-martin/ecommerce/pkg/auth"
	pkgerrors "github.com/southern-martin/ecommerce/pkg/errors"
	"github.com/southern-martin/ecommerce/services/auth/internal/domain"
)

// RefreshTokenInput holds the input data for token refresh.
type RefreshTokenInput struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshTokenOutput holds the output data after successful token refresh.
type RefreshTokenOutput struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// RefreshTokenUseCase handles JWT token refresh.
type RefreshTokenUseCase struct {
	repo          domain.UserRepository
	jwtSecret     string
	accessExpiry  time.Duration
	refreshExpiry time.Duration
	logger        zerolog.Logger
}

// NewRefreshTokenUseCase creates a new RefreshTokenUseCase.
func NewRefreshTokenUseCase(
	repo domain.UserRepository,
	jwtSecret string,
	accessExpiry, refreshExpiry time.Duration,
	logger zerolog.Logger,
) *RefreshTokenUseCase {
	return &RefreshTokenUseCase{
		repo:          repo,
		jwtSecret:     jwtSecret,
		accessExpiry:  accessExpiry,
		refreshExpiry: refreshExpiry,
		logger:        logger,
	}
}

// Execute performs the token refresh.
func (uc *RefreshTokenUseCase) Execute(ctx context.Context, input RefreshTokenInput) (*RefreshTokenOutput, error) {
	// Validate old refresh token
	claims, err := pkgauth.ValidateToken(input.RefreshToken, uc.jwtSecret)
	if err != nil {
		return nil, pkgerrors.NewUnauthorizedError("AUTH_INVALID_REFRESH", "invalid refresh token")
	}

	// Verify refresh token matches what is stored in DB
	user, err := uc.repo.GetByID(ctx, claims.UserID)
	if err != nil || user == nil {
		return nil, pkgerrors.NewUnauthorizedError("AUTH_INVALID_REFRESH", "user not found")
	}

	if user.RefreshToken != input.RefreshToken {
		return nil, pkgerrors.NewUnauthorizedError("AUTH_INVALID_REFRESH", "refresh token mismatch")
	}

	// Generate new token pair
	accessToken, err := pkgauth.GenerateAccessToken(user.ID, user.Email, user.Role, uc.jwtSecret, uc.accessExpiry)
	if err != nil {
		uc.logger.Error().Err(err).Msg("failed to generate access token")
		return nil, pkgerrors.NewInternalError("AUTH_TOKEN_FAILED", "failed to generate tokens")
	}

	newRefreshToken, err := pkgauth.GenerateRefreshToken(user.ID, uc.jwtSecret, uc.refreshExpiry)
	if err != nil {
		uc.logger.Error().Err(err).Msg("failed to generate refresh token")
		return nil, pkgerrors.NewInternalError("AUTH_TOKEN_FAILED", "failed to generate tokens")
	}

	// Rotate refresh token in DB
	if err := uc.repo.UpdateRefreshToken(ctx, user.ID, newRefreshToken); err != nil {
		uc.logger.Error().Err(err).Msg("failed to rotate refresh token")
	}

	return &RefreshTokenOutput{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
