package usecase

import (
	"context"
	"fmt"
	"testing"

	"github.com/southern-martin/ecommerce/services/shipping/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCarrierUseCase_ListCarriers(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expected := []domain.Carrier{
			{Code: "fedex", Name: "FedEx", IsActive: true},
			{Code: "ups", Name: "UPS", IsActive: true},
		}
		repo := &mockCarrierRepo{
			getAllFn: func(ctx context.Context) ([]domain.Carrier, error) {
				return expected, nil
			},
		}
		uc := NewCarrierUseCase(repo, nil)

		result, err := uc.ListCarriers(context.Background())

		require.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("repo error", func(t *testing.T) {
		repo := &mockCarrierRepo{
			getAllFn: func(ctx context.Context) ([]domain.Carrier, error) {
				return nil, fmt.Errorf("db error")
			},
		}
		uc := NewCarrierUseCase(repo, nil)

		result, err := uc.ListCarriers(context.Background())

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestCarrierUseCase_CreateCarrier(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var created *domain.Carrier
		repo := &mockCarrierRepo{
			createFn: func(ctx context.Context, carrier *domain.Carrier) error {
				created = carrier
				return nil
			},
		}
		uc := NewCarrierUseCase(repo, nil)

		carrier := &domain.Carrier{Code: "fedex", Name: "FedEx", IsActive: true}
		err := uc.CreateCarrier(context.Background(), carrier)

		require.NoError(t, err)
		assert.Equal(t, carrier, created)
	})

	t.Run("empty code error", func(t *testing.T) {
		repo := &mockCarrierRepo{}
		uc := NewCarrierUseCase(repo, nil)

		carrier := &domain.Carrier{Code: "", Name: "FedEx"}
		err := uc.CreateCarrier(context.Background(), carrier)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "carrier code and name are required")
	})

	t.Run("empty name error", func(t *testing.T) {
		repo := &mockCarrierRepo{}
		uc := NewCarrierUseCase(repo, nil)

		carrier := &domain.Carrier{Code: "fedex", Name: ""}
		err := uc.CreateCarrier(context.Background(), carrier)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "carrier code and name are required")
	})

	t.Run("both empty error", func(t *testing.T) {
		repo := &mockCarrierRepo{}
		uc := NewCarrierUseCase(repo, nil)

		carrier := &domain.Carrier{Code: "", Name: ""}
		err := uc.CreateCarrier(context.Background(), carrier)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "carrier code and name are required")
	})
}

func TestCarrierUseCase_UpdateCarrier(t *testing.T) {
	t.Run("merges non-empty fields", func(t *testing.T) {
		existing := &domain.Carrier{
			Code:               "fedex",
			Name:               "FedEx",
			IsActive:           true,
			SupportedCountries: []string{"US"},
			APIBaseURL:         "https://api.fedex.com",
		}
		var updated *domain.Carrier
		repo := &mockCarrierRepo{
			getByCodeFn: func(ctx context.Context, code string) (*domain.Carrier, error) {
				return existing, nil
			},
			updateFn: func(ctx context.Context, carrier *domain.Carrier) error {
				updated = carrier
				return nil
			},
		}
		uc := NewCarrierUseCase(repo, nil)

		input := &domain.Carrier{
			Code:               "fedex",
			Name:               "FedEx Express",
			IsActive:           false,
			SupportedCountries: []string{"US", "CA"},
			APIBaseURL:         "https://api-v2.fedex.com",
		}
		err := uc.UpdateCarrier(context.Background(), input)

		require.NoError(t, err)
		assert.Equal(t, "FedEx Express", updated.Name)
		assert.False(t, updated.IsActive)
		assert.Equal(t, []string{"US", "CA"}, updated.SupportedCountries)
		assert.Equal(t, "https://api-v2.fedex.com", updated.APIBaseURL)
	})

	t.Run("empty fields preserve existing values", func(t *testing.T) {
		existing := &domain.Carrier{
			Code:               "fedex",
			Name:               "FedEx",
			IsActive:           true,
			SupportedCountries: []string{"US"},
			APIBaseURL:         "https://api.fedex.com",
		}
		var updated *domain.Carrier
		repo := &mockCarrierRepo{
			getByCodeFn: func(ctx context.Context, code string) (*domain.Carrier, error) {
				return existing, nil
			},
			updateFn: func(ctx context.Context, carrier *domain.Carrier) error {
				updated = carrier
				return nil
			},
		}
		uc := NewCarrierUseCase(repo, nil)

		// Name is empty, SupportedCountries is nil, APIBaseURL is empty
		// IsActive is always overwritten (set to false here)
		input := &domain.Carrier{
			Code:     "fedex",
			IsActive: false,
		}
		err := uc.UpdateCarrier(context.Background(), input)

		require.NoError(t, err)
		assert.Equal(t, "FedEx", updated.Name)           // preserved
		assert.False(t, updated.IsActive)                  // overwritten
		assert.Equal(t, []string{"US"}, updated.SupportedCountries) // preserved
		assert.Equal(t, "https://api.fedex.com", updated.APIBaseURL) // preserved
	})

	t.Run("not found error", func(t *testing.T) {
		repo := &mockCarrierRepo{
			getByCodeFn: func(ctx context.Context, code string) (*domain.Carrier, error) {
				return nil, fmt.Errorf("not found")
			},
		}
		uc := NewCarrierUseCase(repo, nil)

		err := uc.UpdateCarrier(context.Background(), &domain.Carrier{Code: "nope"})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "carrier not found")
	})
}

func TestCarrierUseCase_SetupSellerCarrier(t *testing.T) {
	t.Run("creates new credential", func(t *testing.T) {
		var created *domain.CarrierCredential
		repo := &mockCarrierRepo{
			getByCodeFn: func(ctx context.Context, code string) (*domain.Carrier, error) {
				return &domain.Carrier{Code: "fedex", Name: "FedEx"}, nil
			},
		}
		credRepo := &mockCarrierCredentialRepo{
			getBySellerAndCarrierFn: func(ctx context.Context, sellerID, carrierCode string) (*domain.CarrierCredential, error) {
				return nil, fmt.Errorf("not found")
			},
			createFn: func(ctx context.Context, cred *domain.CarrierCredential) error {
				created = cred
				return nil
			},
		}
		uc := NewCarrierUseCase(repo, credRepo)

		result, err := uc.SetupSellerCarrier(context.Background(), "seller-1", "fedex", `{"api_key":"xyz"}`)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.NotEmpty(t, result.ID)
		assert.Equal(t, "seller-1", result.SellerID)
		assert.Equal(t, "fedex", result.CarrierCode)
		assert.Equal(t, `{"api_key":"xyz"}`, result.Credentials)
		assert.True(t, result.IsActive)
		assert.Equal(t, result, created)
	})

	t.Run("updates existing credential", func(t *testing.T) {
		existing := &domain.CarrierCredential{
			ID:          "cred-1",
			SellerID:    "seller-1",
			CarrierCode: "fedex",
			Credentials: `{"api_key":"old"}`,
			IsActive:    false,
		}
		var updated *domain.CarrierCredential
		repo := &mockCarrierRepo{
			getByCodeFn: func(ctx context.Context, code string) (*domain.Carrier, error) {
				return &domain.Carrier{Code: "fedex", Name: "FedEx"}, nil
			},
		}
		credRepo := &mockCarrierCredentialRepo{
			getBySellerAndCarrierFn: func(ctx context.Context, sellerID, carrierCode string) (*domain.CarrierCredential, error) {
				return existing, nil
			},
			updateFn: func(ctx context.Context, cred *domain.CarrierCredential) error {
				updated = cred
				return nil
			},
		}
		uc := NewCarrierUseCase(repo, credRepo)

		result, err := uc.SetupSellerCarrier(context.Background(), "seller-1", "fedex", `{"api_key":"new"}`)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "cred-1", result.ID)
		assert.Equal(t, `{"api_key":"new"}`, result.Credentials)
		assert.True(t, result.IsActive)
		assert.Equal(t, result, updated)
	})

	t.Run("carrier not found error", func(t *testing.T) {
		repo := &mockCarrierRepo{
			getByCodeFn: func(ctx context.Context, code string) (*domain.Carrier, error) {
				return nil, fmt.Errorf("not found")
			},
		}
		uc := NewCarrierUseCase(repo, nil)

		result, err := uc.SetupSellerCarrier(context.Background(), "seller-1", "nope", `{}`)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "carrier not found")
	})
}

func TestCarrierUseCase_GetSellerCarriers(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expected := []domain.CarrierCredential{
			{ID: "1", SellerID: "seller-1", CarrierCode: "fedex"},
			{ID: "2", SellerID: "seller-1", CarrierCode: "ups"},
		}
		credRepo := &mockCarrierCredentialRepo{
			listBySellerFn: func(ctx context.Context, sellerID string) ([]domain.CarrierCredential, error) {
				assert.Equal(t, "seller-1", sellerID)
				return expected, nil
			},
		}
		uc := NewCarrierUseCase(nil, credRepo)

		result, err := uc.GetSellerCarriers(context.Background(), "seller-1")

		require.NoError(t, err)
		assert.Equal(t, expected, result)
	})
}
