package usecase

import (
	"context"

	"github.com/southern-martin/ecommerce/services/notification/internal/domain"
)

// PreferenceUseCase handles notification preference business logic.
type PreferenceUseCase struct {
	repo domain.PreferenceRepository
}

// NewPreferenceUseCase creates a new PreferenceUseCase.
func NewPreferenceUseCase(repo domain.PreferenceRepository) *PreferenceUseCase {
	return &PreferenceUseCase{repo: repo}
}

// UpdatePreferenceRequest holds the data needed to update a preference.
type UpdatePreferenceRequest struct {
	UserID  string `json:"user_id"`
	Channel string `json:"channel" binding:"required"`
	Enabled bool   `json:"enabled"`
}

// GetPreferences retrieves all notification preferences for a user.
func (uc *PreferenceUseCase) GetPreferences(ctx context.Context, userID string) ([]domain.NotificationPreference, error) {
	return uc.repo.GetByUser(ctx, userID)
}

// UpdatePreference creates or updates a notification preference.
func (uc *PreferenceUseCase) UpdatePreference(ctx context.Context, req UpdatePreferenceRequest) error {
	pref := &domain.NotificationPreference{
		UserID:  req.UserID,
		Channel: domain.NotificationChannel(req.Channel),
		Enabled: req.Enabled,
	}
	return uc.repo.Upsert(ctx, pref)
}
