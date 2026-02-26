package usecase

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"

	pkgauth "github.com/southern-martin/ecommerce/pkg/auth"
	pkgerrors "github.com/southern-martin/ecommerce/pkg/errors"
	"github.com/southern-martin/ecommerce/services/auth/internal/domain"
)

// LoginInput holds the input data for user login.
type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginOutput holds the output data after successful login.
type LoginOutput struct {
	UserID       string `json:"user_id"`
	Email        string `json:"email"`
	Role         string `json:"role"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// LoginUseCase handles user login.
type LoginUseCase struct {
	repo          domain.UserRepository
	jwtSecret     string
	accessExpiry  time.Duration
	refreshExpiry time.Duration
	logger        zerolog.Logger
}

// NewLoginUseCase creates a new LoginUseCase.
func NewLoginUseCase(
	repo domain.UserRepository,
	jwtSecret string,
	accessExpiry, refreshExpiry time.Duration,
	logger zerolog.Logger,
) *LoginUseCase {
	return &LoginUseCase{
		repo:          repo,
		jwtSecret:     jwtSecret,
		accessExpiry:  accessExpiry,
		refreshExpiry: refreshExpiry,
		logger:        logger,
	}
}

// Execute performs the user login.
func (uc *LoginUseCase) Execute(ctx context.Context, input LoginInput) (*LoginOutput, error) {
	user, err := uc.repo.GetByEmail(ctx, input.Email)
	if err != nil || user == nil {
		return nil, pkgerrors.NewUnauthorizedError("AUTH_INVALID_CREDENTIALS", "invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, pkgerrors.NewUnauthorizedError("AUTH_INVALID_CREDENTIALS", "invalid email or password")
	}

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

	return &LoginOutput{
		UserID:       user.ID,
		Email:        user.Email,
		Role:         user.Role,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
