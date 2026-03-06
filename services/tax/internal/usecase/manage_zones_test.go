package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/southern-martin/ecommerce/services/tax/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// Tests for ManageZonesUseCase.ListZones
// ---------------------------------------------------------------------------

func TestManageZones_ListZones(t *testing.T) {
	ctx := context.Background()

	t.Run("success returns all zones", func(t *testing.T) {
		expected := []*domain.TaxZone{
			{ID: "zone-au", CountryCode: "AU", Name: "Australia"},
			{ID: "zone-us-ca", CountryCode: "US", StateCode: "CA", Name: "California"},
			{ID: "zone-de", CountryCode: "DE", Name: "Germany"},
		}
		zoneRepo := &mockTaxZoneRepo{
			listFn: func(_ context.Context) ([]*domain.TaxZone, error) {
				return expected, nil
			},
		}
		uc := NewManageZonesUseCase(zoneRepo)

		zones, err := uc.ListZones(ctx)
		require.NoError(t, err)
		assert.Len(t, zones, 3)
		assert.Equal(t, "zone-au", zones[0].ID)
		assert.Equal(t, "zone-us-ca", zones[1].ID)
		assert.Equal(t, "zone-de", zones[2].ID)
	})

	t.Run("repo error propagates", func(t *testing.T) {
		zoneRepo := &mockTaxZoneRepo{
			listFn: func(_ context.Context) ([]*domain.TaxZone, error) {
				return nil, errors.New("list zones failed")
			},
		}
		uc := NewManageZonesUseCase(zoneRepo)

		zones, err := uc.ListZones(ctx)
		require.Error(t, err)
		assert.Nil(t, zones)
		assert.Contains(t, err.Error(), "list zones failed")
	})
}

// ---------------------------------------------------------------------------
// Tests for ManageZonesUseCase.GetZoneByID
// ---------------------------------------------------------------------------

func TestManageZones_GetZoneByID(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		expected := &domain.TaxZone{
			ID:          "zone-au",
			CountryCode: "AU",
			Name:        "Australia",
		}
		zoneRepo := &mockTaxZoneRepo{
			getByIDFn: func(_ context.Context, id string) (*domain.TaxZone, error) {
				assert.Equal(t, "zone-au", id)
				return expected, nil
			},
		}
		uc := NewManageZonesUseCase(zoneRepo)

		zone, err := uc.GetZoneByID(ctx, "zone-au")
		require.NoError(t, err)
		require.NotNil(t, zone)
		assert.Equal(t, "zone-au", zone.ID)
		assert.Equal(t, "AU", zone.CountryCode)
		assert.Equal(t, "Australia", zone.Name)
	})

	t.Run("not found returns error", func(t *testing.T) {
		zoneRepo := &mockTaxZoneRepo{
			getByIDFn: func(_ context.Context, _ string) (*domain.TaxZone, error) {
				return nil, errors.New("zone not found")
			},
		}
		uc := NewManageZonesUseCase(zoneRepo)

		zone, err := uc.GetZoneByID(ctx, "nonexistent")
		require.Error(t, err)
		assert.Nil(t, zone)
		assert.Contains(t, err.Error(), "zone not found")
	})
}

// ---------------------------------------------------------------------------
// Tests for ManageZonesUseCase.GetZoneByLocation
// ---------------------------------------------------------------------------

func TestManageZones_GetZoneByLocation(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		expected := &domain.TaxZone{
			ID:          "zone-us-ca",
			CountryCode: "US",
			StateCode:   "CA",
			Name:        "California",
		}
		zoneRepo := &mockTaxZoneRepo{
			getByLocationFn: func(_ context.Context, countryCode, stateCode string) (*domain.TaxZone, error) {
				assert.Equal(t, "US", countryCode)
				assert.Equal(t, "CA", stateCode)
				return expected, nil
			},
		}
		uc := NewManageZonesUseCase(zoneRepo)

		zone, err := uc.GetZoneByLocation(ctx, "US", "CA")
		require.NoError(t, err)
		require.NotNil(t, zone)
		assert.Equal(t, "zone-us-ca", zone.ID)
		assert.Equal(t, "US", zone.CountryCode)
		assert.Equal(t, "CA", zone.StateCode)
		assert.Equal(t, "California", zone.Name)
	})

	t.Run("not found returns error", func(t *testing.T) {
		zoneRepo := &mockTaxZoneRepo{
			getByLocationFn: func(_ context.Context, _, _ string) (*domain.TaxZone, error) {
				return nil, errors.New("no zone for location")
			},
		}
		uc := NewManageZonesUseCase(zoneRepo)

		zone, err := uc.GetZoneByLocation(ctx, "ZZ", "")
		require.Error(t, err)
		assert.Nil(t, zone)
		assert.Contains(t, err.Error(), "no zone for location")
	})
}
