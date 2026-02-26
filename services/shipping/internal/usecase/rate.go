package usecase

import (
	"context"

	"github.com/southern-martin/ecommerce/services/shipping/internal/domain"
)

// RateUseCase handles shipping rate calculations.
type RateUseCase struct {
	carrierRepo domain.CarrierRepository
}

// NewRateUseCase creates a new RateUseCase.
func NewRateUseCase(carrierRepo domain.CarrierRepository) *RateUseCase {
	return &RateUseCase{carrierRepo: carrierRepo}
}

// GetShippingRatesRequest is the input for GetShippingRates.
type GetShippingRatesRequest struct {
	Origin      domain.Address
	Destination domain.Address
	WeightGrams int
	Currency    string
}

// GetShippingRates returns mock shipping rates from all active carriers.
func (uc *RateUseCase) GetShippingRates(ctx context.Context, req GetShippingRatesRequest) ([]domain.ShippingRate, error) {
	carriers, err := uc.carrierRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	currency := req.Currency
	if currency == "" {
		currency = "USD"
	}

	var rates []domain.ShippingRate
	for _, carrier := range carriers {
		// Mock rate calculation based on weight
		baseRate := int64(500) // $5.00 base
		weightRate := int64(req.WeightGrams) / 100 * 50 // $0.50 per 100g

		// Standard service
		rates = append(rates, domain.ShippingRate{
			CarrierCode:      carrier.Code,
			ServiceName:      carrier.Name + " Standard",
			RateCents:        baseRate + weightRate,
			Currency:         currency,
			EstimatedDaysMin: 5,
			EstimatedDaysMax: 10,
		})

		// Express service
		rates = append(rates, domain.ShippingRate{
			CarrierCode:      carrier.Code,
			ServiceName:      carrier.Name + " Express",
			RateCents:        (baseRate + weightRate) * 2,
			Currency:         currency,
			EstimatedDaysMin: 2,
			EstimatedDaysMax: 4,
		})
	}

	return rates, nil
}
