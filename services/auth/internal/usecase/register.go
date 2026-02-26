package usecase

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"

	pkgauth "github.com/southern-martin/ecommerce/pkg/auth"
	pkgerrors "github.com/southern-martin/ecommerce/pkg/errors"
	"github.com/southern-martin/ecommerce/services/auth/internal/domain"
	authnats "github.com/southern-martin/ecommerce/services/auth/internal/infrastructure/nats"
)

// RegisterInput holds the input data for user registration.
type RegisterInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// RegisterOutput holds the output data after successful registration.
type RegisterOutput struct {
	UserID       string `json:"user_id"`
	Email        string `json:"email"`
	Role         string `json:"role"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// RegisterUseCase handles user registration.
type RegisterUseCase struct {
	repo           domain.UserRepository
	publisher      *authnats.EventPublisher
	jwtSecret      string
	accessExpiry   time.Duration
	refreshExpiry  time.Duration
	logger         zerolog.Logger
}

// NewRegisterUseCase creates a new RegisterUseCase.
func NewRegisterUseCase(
	repo domain.UserRepository,
	publisher *authnats.EventPublisher,
	jwtSecret string,
	accessExpiry, refreshExpiry time.Duration,
	logger zerolog.Logger,
) *RegisterUseCase {
	return &RegisterUseCase{
		repo:          repo,
		publisher:     publisher,
		jwtSecret:     jwtSecret,
		accessExpiry:  accessExpiry,
		refreshExpiry: refreshExpiry,
		logger:        logger,
	}
}

// Execute performs the user registration.
func (uc *RegisterUseCase) Execute(ctx context.Context, input RegisterInput) (*RegisterOutput, error) {
	// Check if email is already taken
	existing, _ := uc.repo.GetByEmail(ctx, input.Email)
	if existing != nil {
		return nil, pkgerrors.NewConflictError("AUTH_EMAIL_EXISTS", "email already registered")
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		uc.logger.Error().Err(err).Msg("failed to hash password")
		return nil, pkgerrors.NewInternalError("AUTH_HASH_FAILED", "failed to process password")
	}

	user := &domain.AuthUser{
		Email:        input.Email,
		PasswordHash: string(hash),
		Role:         "buyer",
	}

	if err := uc.repo.Create(ctx, user); err != nil {
		uc.logger.Error().Err(err).Msg("failed to create user")
		return nil, pkgerrors.NewInternalError("AUTH_CREATE_FAILED", "failed to create user")
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

	// Store refresh token
	if err := uc.repo.UpdateRefreshToken(ctx, user.ID, refreshToken); err != nil {
		uc.logger.Error().Err(err).Msg("failed to store refresh token")
	}

	// Publish event
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

	return &RegisterOutput{
		UserID:       user.ID,
		Email:        user.Email,
		Role:         user.Role,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
