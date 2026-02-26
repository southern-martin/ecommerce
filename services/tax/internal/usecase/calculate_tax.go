package usecase

import (
	"context"
	"math"

	"github.com/southern-martin/ecommerce/services/tax/internal/domain"
)

// CalculateTaxUseCase handles tax calculation logic.
type CalculateTaxUseCase struct {
	zoneRepo domain.TaxZoneRepository
	ruleRepo domain.TaxRuleRepository
}

// NewCalculateTaxUseCase creates a new CalculateTaxUseCase.
func NewCalculateTaxUseCase(zoneRepo domain.TaxZoneRepository, ruleRepo domain.TaxRuleRepository) *CalculateTaxUseCase {
	return &CalculateTaxUseCase{
		zoneRepo: zoneRepo,
		ruleRepo: ruleRepo,
	}
}

// Execute performs the tax calculation for the given request.
func (uc *CalculateTaxUseCase) Execute(ctx context.Context, req *domain.TaxCalculationRequest) (*domain.TaxCalculation, error) {
	// 1. Look up TaxZone by shipping address
	zone, err := uc.zoneRepo.GetByLocation(ctx, req.ShippingAddress.CountryCode, req.ShippingAddress.StateCode)
	if err != nil {
		// No tax zone found - return zero tax
		return uc.zeroTaxResult(req), nil
	}

	// 2. Calculate subtotal
	var subtotalCents int64
	for _, item := range req.Items {
		subtotalCents += item.PriceCents * int64(item.Quantity)
	}

	// 3. For each item, find matching rules and calculate tax
	breakdownMap := make(map[string]*domain.TaxBreakdown)
	var totalTaxCents int64

	for _, item := range req.Items {
		// Find matching rules for this item's category
		rules, err := uc.ruleRepo.GetByZoneAndCategory(ctx, zone.ID, item.Category)
		if err != nil {
			return nil, err
		}

		for _, rule := range rules {
			itemTotal := item.PriceCents * int64(item.Quantity)
			var taxCents int64

			if rule.Inclusive {
				// Tax is included in price: tax = price - (price / (1 + rate))
				priceExclTax := float64(itemTotal) / (1.0 + rule.Rate)
				taxCents = int64(math.Round(float64(itemTotal) - priceExclTax))
			} else {
				// Tax is added on top: tax = price * rate
				taxCents = int64(math.Round(float64(itemTotal) * rule.Rate))
			}

			key := rule.TaxName
			if existing, ok := breakdownMap[key]; ok {
				existing.AmountCents += taxCents
			} else {
				breakdownMap[key] = &domain.TaxBreakdown{
					TaxName:      rule.TaxName,
					Rate:         rule.Rate,
					AmountCents:  taxCents,
					Jurisdiction: zone.Name,
				}
			}
			totalTaxCents += taxCents
		}
	}

	// 4. Build breakdown slice
	breakdown := make([]domain.TaxBreakdown, 0, len(breakdownMap))
	for _, b := range breakdownMap {
		breakdown = append(breakdown, *b)
	}

	return &domain.TaxCalculation{
		SubtotalCents:  subtotalCents,
		TaxAmountCents: totalTaxCents,
		Breakdown:      breakdown,
	}, nil
}

func (uc *CalculateTaxUseCase) zeroTaxResult(req *domain.TaxCalculationRequest) *domain.TaxCalculation {
	var subtotalCents int64
	for _, item := range req.Items {
		subtotalCents += item.PriceCents * int64(item.Quantity)
	}
	return &domain.TaxCalculation{
		SubtotalCents:  subtotalCents,
		TaxAmountCents: 0,
		Breakdown:      []domain.TaxBreakdown{},
	}
}
