package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	pkgerrors "github.com/southern-martin/ecommerce/pkg/errors"
	"github.com/southern-martin/ecommerce/services/user/internal/domain"
)

func TestCreateSeller_Success(t *testing.T) {
	var created *domain.SellerProfile
	repo := &mockSellerProfileRepo{
		getByUserFn: func(_ context.Context, _ string) (*domain.SellerProfile, error) {
			return nil, errors.New("not found")
		},
		createFn: func(_ context.Context, s *domain.SellerProfile) error {
			created = s
			return nil
		},
	}
	pub := &mockEventPublisher{}

	uc := NewSellerUseCase(repo, pub, zerolog.Nop())
	seller, err := uc.CreateSeller(context.Background(), "user-1", CreateSellerInput{
		StoreName:   "My Store",
		Description: "Sells things",
		LogoURL:     "https://example.com/logo.png",
	})

	require.NoError(t, err)
	require.NotNil(t, seller)
	assert.Equal(t, "pending", seller.Status)
	assert.Equal(t, "user-1", created.UserID)
	assert.Equal(t, "My Store", created.StoreName)
	assert.Equal(t, "Sells things", created.Description)
	assert.Equal(t, "https://example.com/logo.png", created.LogoURL)
}

func TestCreateSeller_AlreadyExists(t *testing.T) {
	repo := &mockSellerProfileRepo{
		getByUserFn: func(_ context.Context, _ string) (*domain.SellerProfile, error) {
			return &domain.SellerProfile{ID: "existing"}, nil
		},
	}
	pub := &mockEventPublisher{}

	uc := NewSellerUseCase(repo, pub, zerolog.Nop())
	_, err := uc.CreateSeller(context.Background(), "user-1", CreateSellerInput{
		StoreName: "My Store",
	})

	require.Error(t, err)
	var conflictErr *pkgerrors.ConflictError
	assert.ErrorAs(t, err, &conflictErr)
	assert.Equal(t, "SELLER_EXISTS", conflictErr.Code)
}

func TestGetSeller_Success(t *testing.T) {
	expected := &domain.SellerProfile{ID: "seller-1", StoreName: "Cool Store"}
	repo := &mockSellerProfileRepo{
		getByIDFn: func(_ context.Context, id string) (*domain.SellerProfile, error) {
			assert.Equal(t, "seller-1", id)
			return expected, nil
		},
	}
	pub := &mockEventPublisher{}

	uc := NewSellerUseCase(repo, pub, zerolog.Nop())
	seller, err := uc.GetSeller(context.Background(), "seller-1")

	require.NoError(t, err)
	assert.Equal(t, "Cool Store", seller.StoreName)
}

func TestGetSellerByUserID_Success(t *testing.T) {
	expected := &domain.SellerProfile{ID: "seller-1", UserID: "user-1", StoreName: "My Shop"}
	repo := &mockSellerProfileRepo{
		getByUserFn: func(_ context.Context, userID string) (*domain.SellerProfile, error) {
			assert.Equal(t, "user-1", userID)
			return expected, nil
		},
	}
	pub := &mockEventPublisher{}

	uc := NewSellerUseCase(repo, pub, zerolog.Nop())
	seller, err := uc.GetSellerByUserID(context.Background(), "user-1")

	require.NoError(t, err)
	assert.Equal(t, "My Shop", seller.StoreName)
}

func TestUpdateSeller_Success(t *testing.T) {
	existing := &domain.SellerProfile{
		ID:          "seller-1",
		UserID:      "user-1",
		StoreName:   "Old Name",
		Description: "Old desc",
		LogoURL:     "old-logo.png",
		Status:      "approved",
	}
	repo := &mockSellerProfileRepo{
		getByUserFn: func(_ context.Context, _ string) (*domain.SellerProfile, error) {
			return existing, nil
		},
		updateFn: func(_ context.Context, _ *domain.SellerProfile) error { return nil },
	}
	pub := &mockEventPublisher{}

	uc := NewSellerUseCase(repo, pub, zerolog.Nop())
	result, err := uc.UpdateSeller(context.Background(), "user-1", UpdateSellerInput{
		StoreName:   strPtr("New Name"),
		Description: strPtr("New desc"),
		// LogoURL nil -> unchanged
	})

	require.NoError(t, err)
	assert.Equal(t, "New Name", result.StoreName)
	assert.Equal(t, "New desc", result.Description)
	assert.Equal(t, "old-logo.png", result.LogoURL, "nil fields should remain unchanged")
}

func TestApproveSeller_Success(t *testing.T) {
	existing := &domain.SellerProfile{
		ID:     "seller-1",
		UserID: "user-1",
		Status: "pending",
	}

	var publishedSubject string
	var publishedData interface{}

	repo := &mockSellerProfileRepo{
		getByIDFn: func(_ context.Context, _ string) (*domain.SellerProfile, error) {
			return existing, nil
		},
		updateFn: func(_ context.Context, _ *domain.SellerProfile) error { return nil },
	}
	pub := &mockEventPublisher{
		publishFn: func(subject string, data interface{}) error {
			publishedSubject = subject
			publishedData = data
			return nil
		},
	}

	uc := NewSellerUseCase(repo, pub, zerolog.Nop())
	seller, err := uc.ApproveSeller(context.Background(), "seller-1")

	require.NoError(t, err)
	assert.Equal(t, "approved", seller.Status)
	assert.Equal(t, "seller.approved", publishedSubject)

	evt, ok := publishedData.(SellerApprovedEvent)
	require.True(t, ok)
	assert.Equal(t, "seller-1", evt.SellerID)
	assert.Equal(t, "user-1", evt.UserID)
}

func TestApproveSeller_PublishFailureDoesNotFail(t *testing.T) {
	existing := &domain.SellerProfile{
		ID:     "seller-2",
		UserID: "user-2",
		Status: "pending",
	}

	repo := &mockSellerProfileRepo{
		getByIDFn: func(_ context.Context, _ string) (*domain.SellerProfile, error) {
			return existing, nil
		},
		updateFn: func(_ context.Context, _ *domain.SellerProfile) error { return nil },
	}
	pub := &mockEventPublisher{
		publishFn: func(_ string, _ interface{}) error {
			return errors.New("nats connection lost")
		},
	}

	uc := NewSellerUseCase(repo, pub, zerolog.Nop())
	seller, err := uc.ApproveSeller(context.Background(), "seller-2")

	require.NoError(t, err, "publish failure should be swallowed")
	assert.Equal(t, "approved", seller.Status)
}

func TestUpdateSeller_RepoError(t *testing.T) {
	repo := &mockSellerProfileRepo{
		getByUserFn: func(_ context.Context, _ string) (*domain.SellerProfile, error) {
			return nil, errors.New("db error")
		},
	}
	pub := &mockEventPublisher{}

	uc := NewSellerUseCase(repo, pub, zerolog.Nop())
	_, err := uc.UpdateSeller(context.Background(), "user-1", UpdateSellerInput{
		StoreName: strPtr("X"),
	})

	require.Error(t, err)
	assert.Equal(t, "db error", err.Error())
}
