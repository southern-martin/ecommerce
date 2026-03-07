package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/southern-martin/ecommerce/services/notification/internal/domain"
)

// Mocks are defined in notification_test.go (shared within same package).

// ===========================================================================
// GetPreferences tests
// ===========================================================================

func TestGetPreferences_Success(t *testing.T) {
	prefRepo := &mockPreferenceRepo{}
	prefs := []domain.NotificationPreference{
		{ID: "pref-1", UserID: "user-1", Channel: domain.ChannelEmail, Enabled: true},
		{ID: "pref-2", UserID: "user-1", Channel: domain.ChannelPush, Enabled: false},
		{ID: "pref-3", UserID: "user-1", Channel: domain.ChannelInApp, Enabled: true},
	}
	prefRepo.getByUserFn = func(_ context.Context, userID string) ([]domain.NotificationPreference, error) {
		assert.Equal(t, "user-1", userID)
		return prefs, nil
	}

	uc := NewPreferenceUseCase(prefRepo)
	result, err := uc.GetPreferences(context.Background(), "user-1")

	require.NoError(t, err)
	assert.Len(t, result, 3)
	assert.Equal(t, domain.ChannelEmail, result[0].Channel)
	assert.True(t, result[0].Enabled)
	assert.Equal(t, domain.ChannelPush, result[1].Channel)
	assert.False(t, result[1].Enabled)
}

func TestGetPreferences_Empty(t *testing.T) {
	prefRepo := &mockPreferenceRepo{}
	// Default getByUserFn returns nil, nil

	uc := NewPreferenceUseCase(prefRepo)
	result, err := uc.GetPreferences(context.Background(), "user-1")

	require.NoError(t, err)
	assert.Empty(t, result)
}

func TestGetPreferences_RepoError(t *testing.T) {
	prefRepo := &mockPreferenceRepo{}
	prefRepo.getByUserFn = func(_ context.Context, _ string) ([]domain.NotificationPreference, error) {
		return nil, errors.New("query failed")
	}

	uc := NewPreferenceUseCase(prefRepo)
	_, err := uc.GetPreferences(context.Background(), "user-1")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "query failed")
}

// ===========================================================================
// UpdatePreference tests
// ===========================================================================

func TestUpdatePreference_Success(t *testing.T) {
	prefRepo := &mockPreferenceRepo{}
	var upserted *domain.NotificationPreference
	prefRepo.upsertFn = func(_ context.Context, pref *domain.NotificationPreference) error {
		upserted = pref
		return nil
	}

	uc := NewPreferenceUseCase(prefRepo)
	err := uc.UpdatePreference(context.Background(), UpdatePreferenceRequest{
		UserID:  "user-1",
		Channel: "email",
		Enabled: true,
	})

	require.NoError(t, err)
	require.NotNil(t, upserted)
	assert.Equal(t, "user-1", upserted.UserID)
	assert.Equal(t, domain.ChannelEmail, upserted.Channel)
	assert.True(t, upserted.Enabled)
}

func TestUpdatePreference_DisableChannel(t *testing.T) {
	prefRepo := &mockPreferenceRepo{}
	var upserted *domain.NotificationPreference
	prefRepo.upsertFn = func(_ context.Context, pref *domain.NotificationPreference) error {
		upserted = pref
		return nil
	}

	uc := NewPreferenceUseCase(prefRepo)
	err := uc.UpdatePreference(context.Background(), UpdatePreferenceRequest{
		UserID:  "user-1",
		Channel: "push",
		Enabled: false,
	})

	require.NoError(t, err)
	require.NotNil(t, upserted)
	assert.Equal(t, domain.NotificationChannel("push"), upserted.Channel)
	assert.False(t, upserted.Enabled)
}

func TestUpdatePreference_RepoError(t *testing.T) {
	prefRepo := &mockPreferenceRepo{}
	prefRepo.upsertFn = func(_ context.Context, _ *domain.NotificationPreference) error {
		return errors.New("upsert failed")
	}

	uc := NewPreferenceUseCase(prefRepo)
	err := uc.UpdatePreference(context.Background(), UpdatePreferenceRequest{
		UserID:  "user-1",
		Channel: "email",
		Enabled: true,
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "upsert failed")
}

func TestUpdatePreference_InAppChannel(t *testing.T) {
	prefRepo := &mockPreferenceRepo{}
	var upserted *domain.NotificationPreference
	prefRepo.upsertFn = func(_ context.Context, pref *domain.NotificationPreference) error {
		upserted = pref
		return nil
	}

	uc := NewPreferenceUseCase(prefRepo)
	err := uc.UpdatePreference(context.Background(), UpdatePreferenceRequest{
		UserID:  "user-2",
		Channel: "in_app",
		Enabled: true,
	})

	require.NoError(t, err)
	require.NotNil(t, upserted)
	assert.Equal(t, "user-2", upserted.UserID)
	assert.Equal(t, domain.ChannelInApp, upserted.Channel)
	assert.True(t, upserted.Enabled)
}
