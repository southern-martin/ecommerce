package usecase

import (
	"context"
	"fmt"
	"testing"

	"github.com/southern-martin/ecommerce/services/shipping/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRateUseCase_GetShippingRates(t *testing.T) {
	t.Run("returns correct rates for carriers", func(t *testing.T) {
		carriers := []domain.Carrier{
			{Code: "fedex", Name: "FedEx", IsActive: true},
			{Code: "ups", Name: "UPS", IsActive: true},
		}
		repo := &mockCarrierRepo{
			getAllFn: func(ctx context.Context) ([]domain.Carrier, error) {
				return carriers, nil
			},
		}
		uc := NewRateUseCase(repo)

		req := GetShippingRatesRequest{
			Origin:      domain.Address{Country: "US"},
			Destination: domain.Address{Country: "US"},
			WeightGrams: 1000,
			Currency:    "USD",
		}

		rates, err := uc.GetShippingRates(context.Background(), req)

		require.NoError(t, err)
		require.Len(t, rates, 4) // 2 carriers * 2 services

		// Verify rate formula: base=500 + weight/100*50
		// 1000g: 1000/100*50 = 500, total standard = 500 + 500 = 1000
		assert.Equal(t, "fedex", rates[0].CarrierCode)
		assert.Equal(t, "FedEx Standard", rates[0].ServiceName)
		assert.Equal(t, int64(1000), rates[0].RateCents)
		assert.Equal(t, "USD", rates[0].Currency)
		assert.Equal(t, 5, rates[0].EstimatedDaysMin)
		assert.Equal(t, 10, rates[0].EstimatedDaysMax)

		// Express = 2x standard
		assert.Equal(t, "fedex", rates[1].CarrierCode)
		assert.Equal(t, "FedEx Express", rates[1].ServiceName)
		assert.Equal(t, int64(2000), rates[1].RateCents)
		assert.Equal(t, 2, rates[1].EstimatedDaysMin)
		assert.Equal(t, 4, rates[1].EstimatedDaysMax)

		// Second carrier
		assert.Equal(t, "ups", rates[2].CarrierCode)
		assert.Equal(t, "UPS Standard", rates[2].ServiceName)
		assert.Equal(t, int64(1000), rates[2].RateCents)

		assert.Equal(t, "ups", rates[3].CarrierCode)
		assert.Equal(t, "UPS Express", rates[3].ServiceName)
		assert.Equal(t, int64(2000), rates[3].RateCents)
	})

	t.Run("defaults currency to USD", func(t *testing.T) {
		carriers := []domain.Carrier{
			{Code: "fedex", Name: "FedEx"},
		}
		repo := &mockCarrierRepo{
			getAllFn: func(ctx context.Context) ([]domain.Carrier, error) {
				return carriers, nil
			},
		}
		uc := NewRateUseCase(repo)

		req := GetShippingRatesRequest{
			WeightGrams: 500,
			Currency:    "", // empty should default
		}

		rates, err := uc.GetShippingRates(context.Background(), req)

		require.NoError(t, err)
		for _, rate := range rates {
			assert.Equal(t, "USD", rate.Currency)
		}
	})

	t.Run("no carriers returns empty", func(t *testing.T) {
		repo := &mockCarrierRepo{
			getAllFn: func(ctx context.Context) ([]domain.Carrier, error) {
				return []domain.Carrier{}, nil
			},
		}
		uc := NewRateUseCase(repo)

		rates, err := uc.GetShippingRates(context.Background(), GetShippingRatesRequest{WeightGrams: 100})

		require.NoError(t, err)
		assert.Empty(t, rates)
	})

	t.Run("weight affects rate calculation", func(t *testing.T) {
		carriers := []domain.Carrier{
			{Code: "dhl", Name: "DHL"},
		}
		repo := &mockCarrierRepo{
			getAllFn: func(ctx context.Context) ([]domain.Carrier, error) {
				return carriers, nil
			},
		}
		uc := NewRateUseCase(repo)

		tests := []struct {
			name         string
			weightGrams  int
			expectedRate int64
		}{
			// base=500 + weightGrams/100*50 (integer division)
			{"zero weight", 0, 500},
			{"50g (under 100g bucket)", 50, 500},   // 50/100=0, 0*50=0
			{"100g", 100, 550},                       // 100/100=1, 1*50=50
			{"500g", 500, 750},                       // 500/100=5, 5*50=250
			{"2000g", 2000, 1500},                    // 2000/100=20, 20*50=1000
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				req := GetShippingRatesRequest{
					WeightGrams: tc.weightGrams,
					Currency:    "USD",
				}
				rates, err := uc.GetShippingRates(context.Background(), req)

				require.NoError(t, err)
				require.Len(t, rates, 2)
				assert.Equal(t, tc.expectedRate, rates[0].RateCents,
					"standard rate for %dg", tc.weightGrams)
				assert.Equal(t, tc.expectedRate*2, rates[1].RateCents,
					"express rate for %dg", tc.weightGrams)
			})
		}
	})

	t.Run("repo error", func(t *testing.T) {
		repo := &mockCarrierRepo{
			getAllFn: func(ctx context.Context) ([]domain.Carrier, error) {
				return nil, fmt.Errorf("db error")
			},
		}
		uc := NewRateUseCase(repo)

		rates, err := uc.GetShippingRates(context.Background(), GetShippingRatesRequest{})

		assert.Error(t, err)
		assert.Nil(t, rates)
	})
}
