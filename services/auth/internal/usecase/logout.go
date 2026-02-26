package usecase

import (
	"context"
	"time"

	"github.com/rs/zerolog"

	pkgauth "github.com/southern-martin/ecommerce/pkg/auth"
	pkgerrors "github.com/southern-martin/ecommerce/pkg/errors"
	"github.com/southern-martin/ecommerce/services/auth/internal/domain"
	authredis "github.com/southern-martin/ecommerce/services/auth/internal/infrastructure/redis"
)

// LogoutInput holds the input data for user logout.
type LogoutInput struct {
	UserID      string `json:"user_id"`
	AccessToken string `json:"access_token" binding:"required"`
}

// LogoutUseCase handles user logout by clearing refresh tokens and blacklisting access tokens.
type LogoutUseCase struct {
	repo      domain.UserRepository
	blacklist *authredis.TokenBlacklist
	jwtSecret string
	logger    zerolog.Logger
}

// NewLogoutUseCase creates a new LogoutUseCase.
func NewLogoutUseCase(
	repo domain.UserRepository,
	blacklist *authredis.TokenBlacklist,
	jwtSecret string,
	logger zerolog.Logger,
) *LogoutUseCase {
	return &LogoutUseCase{
		repo:      repo,
		blacklist: blacklist,
		jwtSecret: jwtSecret,
		logger:    logger,
	}
}

// Execute performs the user logout.
func (uc *LogoutUseCase) Execute(ctx context.Context, input LogoutInput) error {
	// Clear refresh token in DB
	if err := uc.repo.UpdateRefreshToken(ctx, input.UserID, ""); err != nil {
		uc.logger.Error().Err(err).Msg("failed to clear refresh token")
		return pkgerrors.NewInternalError("AUTH_LOGOUT_FAILED", "failed to logout")
	}

	// Blacklist the access token in Redis with TTL = remaining token lifetime
	claims, err := pkgauth.ValidateToken(input.AccessToken, uc.jwtSecret)
	if err == nil && claims.ExpiresAt != nil {
		ttl := time.Until(claims.ExpiresAt.Time)
		if ttl > 0 {
			if err := uc.blacklist.BlacklistToken(ctx, input.AccessToken, ttl); err != nil {
				uc.logger.Error().Err(err).Msg("failed to blacklist access token")
			}
		}
	}

	return nil
}
