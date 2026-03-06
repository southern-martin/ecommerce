package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	pkgerrors "github.com/southern-martin/ecommerce/pkg/errors"
	"github.com/southern-martin/ecommerce/services/user/internal/domain"
)

func TestGetProfile_Success(t *testing.T) {
	expected := &domain.UserProfile{
		ID:        "user-1",
		Email:     "alice@example.com",
		FirstName: "Alice",
		LastName:  "Smith",
		Role:      "buyer",
	}
	repo := &mockUserProfileRepo{
		getByIDFn: func(_ context.Context, id string) (*domain.UserProfile, error) {
			assert.Equal(t, "user-1", id)
			return expected, nil
		},
	}

	uc := NewProfileUseCase(repo, zerolog.Nop())
	profile, err := uc.GetProfile(context.Background(), "user-1")

	require.NoError(t, err)
	assert.Equal(t, "alice@example.com", profile.Email)
	assert.Equal(t, "Alice", profile.FirstName)
}

func TestGetProfile_NotFound(t *testing.T) {
	repo := &mockUserProfileRepo{
		getByIDFn: func(_ context.Context, _ string) (*domain.UserProfile, error) {
			return nil, pkgerrors.NewNotFoundError("NOT_FOUND", "user profile not found")
		},
	}

	uc := NewProfileUseCase(repo, zerolog.Nop())
	profile, err := uc.GetProfile(context.Background(), "nonexistent")

	require.Error(t, err)
	assert.Nil(t, profile)
	var nfErr *pkgerrors.NotFoundError
	assert.ErrorAs(t, err, &nfErr)
}

func TestUpdateProfile_PartialUpdate(t *testing.T) {
	existing := &domain.UserProfile{
		ID:          "user-1",
		Email:       "alice@example.com",
		FirstName:   "Alice",
		LastName:    "Smith",
		DisplayName: "alice_s",
		Phone:       "555-0000",
		AvatarURL:   "https://example.com/old.png",
		Bio:         "Old bio",
		Role:        "buyer",
		CreatedAt:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	repo := &mockUserProfileRepo{
		getByIDFn: func(_ context.Context, _ string) (*domain.UserProfile, error) {
			return existing, nil
		},
		updateFn: func(_ context.Context, _ *domain.UserProfile) error { return nil },
	}

	uc := NewProfileUseCase(repo, zerolog.Nop())
	result, err := uc.UpdateProfile(context.Background(), "user-1", UpdateProfileInput{
		FirstName: strPtr("Alicia"),
		Bio:       strPtr("New bio"),
		// other fields nil -> remain unchanged
	})

	require.NoError(t, err)
	assert.Equal(t, "Alicia", result.FirstName)
	assert.Equal(t, "New bio", result.Bio)
	assert.Equal(t, "Smith", result.LastName, "nil fields should remain unchanged")
	assert.Equal(t, "alice_s", result.DisplayName, "nil fields should remain unchanged")
	assert.Equal(t, "555-0000", result.Phone, "nil fields should remain unchanged")
	assert.Equal(t, "https://example.com/old.png", result.AvatarURL, "nil fields should remain unchanged")
	assert.False(t, result.UpdatedAt.IsZero(), "UpdatedAt should be set")
}

func TestUpdateProfile_FullUpdate(t *testing.T) {
	existing := &domain.UserProfile{
		ID:    "user-1",
		Email: "alice@example.com",
		Role:  "buyer",
	}

	repo := &mockUserProfileRepo{
		getByIDFn: func(_ context.Context, _ string) (*domain.UserProfile, error) {
			return existing, nil
		},
		updateFn: func(_ context.Context, _ *domain.UserProfile) error { return nil },
	}

	uc := NewProfileUseCase(repo, zerolog.Nop())
	result, err := uc.UpdateProfile(context.Background(), "user-1", UpdateProfileInput{
		FirstName:   strPtr("Alice"),
		LastName:    strPtr("Wonder"),
		DisplayName: strPtr("alicew"),
		Phone:       strPtr("555-9999"),
		AvatarURL:   strPtr("https://example.com/new.png"),
		Bio:         strPtr("Full bio"),
	})

	require.NoError(t, err)
	assert.Equal(t, "Alice", result.FirstName)
	assert.Equal(t, "Wonder", result.LastName)
	assert.Equal(t, "alicew", result.DisplayName)
	assert.Equal(t, "555-9999", result.Phone)
	assert.Equal(t, "https://example.com/new.png", result.AvatarURL)
	assert.Equal(t, "Full bio", result.Bio)
}

func TestCreateFromEvent_Success(t *testing.T) {
	var created *domain.UserProfile
	repo := &mockUserProfileRepo{
		createFn: func(_ context.Context, p *domain.UserProfile) error {
			created = p
			return nil
		},
	}

	uc := NewProfileUseCase(repo, zerolog.Nop())
	err := uc.CreateFromEvent("user-99", "bob@example.com", "buyer")

	require.NoError(t, err)
	require.NotNil(t, created)
	assert.Equal(t, "user-99", created.ID)
	assert.Equal(t, "bob@example.com", created.Email)
	assert.Equal(t, "buyer", created.Role)
	assert.False(t, created.CreatedAt.IsZero())
	assert.False(t, created.UpdatedAt.IsZero())
}

func TestCreateFromEvent_RepoError(t *testing.T) {
	repo := &mockUserProfileRepo{
		createFn: func(_ context.Context, _ *domain.UserProfile) error {
			return errors.New("db connection refused")
		},
	}

	uc := NewProfileUseCase(repo, zerolog.Nop())
	err := uc.CreateFromEvent("user-99", "bob@example.com", "buyer")

	require.Error(t, err)
	assert.Equal(t, "db connection refused", err.Error())
}
