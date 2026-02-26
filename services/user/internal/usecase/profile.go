package usecase

import (
	"context"
	"time"

	"github.com/rs/zerolog"

	"github.com/southern-martin/ecommerce/services/user/internal/domain"
)

// UpdateProfileInput holds the fields that can be updated on a user profile.
type UpdateProfileInput struct {
	FirstName   *string `json:"first_name"`
	LastName    *string `json:"last_name"`
	DisplayName *string `json:"display_name"`
	Phone       *string `json:"phone"`
	AvatarURL   *string `json:"avatar_url"`
	Bio         *string `json:"bio"`
}

// ProfileUseCase handles user profile business logic.
type ProfileUseCase struct {
	repo   domain.UserProfileRepository
	logger zerolog.Logger
}

// NewProfileUseCase creates a new ProfileUseCase.
func NewProfileUseCase(repo domain.UserProfileRepository, logger zerolog.Logger) *ProfileUseCase {
	return &ProfileUseCase{
		repo:   repo,
		logger: logger,
	}
}

// GetProfile retrieves a user profile by ID.
func (uc *ProfileUseCase) GetProfile(ctx context.Context, userID string) (*domain.UserProfile, error) {
	profile, err := uc.repo.GetByID(ctx, userID)
	if err != nil {
		uc.logger.Error().Err(err).Str("user_id", userID).Msg("failed to get profile")
		return nil, err
	}
	return profile, nil
}

// UpdateProfile updates a user profile with the provided input fields.
func (uc *ProfileUseCase) UpdateProfile(ctx context.Context, userID string, input UpdateProfileInput) (*domain.UserProfile, error) {
	profile, err := uc.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if input.FirstName != nil {
		profile.FirstName = *input.FirstName
	}
	if input.LastName != nil {
		profile.LastName = *input.LastName
	}
	if input.DisplayName != nil {
		profile.DisplayName = *input.DisplayName
	}
	if input.Phone != nil {
		profile.Phone = *input.Phone
	}
	if input.AvatarURL != nil {
		profile.AvatarURL = *input.AvatarURL
	}
	if input.Bio != nil {
		profile.Bio = *input.Bio
	}

	profile.UpdatedAt = time.Now()

	if err := uc.repo.Update(ctx, profile); err != nil {
		uc.logger.Error().Err(err).Str("user_id", userID).Msg("failed to update profile")
		return nil, err
	}

	return profile, nil
}

// CreateFromEvent creates a user profile from a NATS user.registered event.
func (uc *ProfileUseCase) CreateFromEvent(userID, email, role string) error {
	now := time.Now()
	profile := &domain.UserProfile{
		ID:        userID,
		Email:     email,
		Role:      role,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := uc.repo.Create(context.Background(), profile); err != nil {
		uc.logger.Error().Err(err).Str("user_id", userID).Msg("failed to create profile from event")
		return err
	}

	uc.logger.Info().Str("user_id", userID).Str("email", email).Msg("profile created from registration event")
	return nil
}
