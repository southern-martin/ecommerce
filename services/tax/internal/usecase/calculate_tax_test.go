package usecase

import (
	"context"
	"errors"
	"math"
	"testing"

	"github.com/southern-martin/ecommerce/services/tax/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// Hand-written function-field mocks
// ---------------------------------------------------------------------------

type mockTaxZoneRepo struct {
	createFn        func(ctx context.Context, zone *domain.TaxZone) error
	getByIDFn       func(ctx context.Context, id string) (*domain.TaxZone, error)
	getByLocationFn func(ctx context.Context, countryCode, stateCode string) (*domain.TaxZone, error)
	listFn          func(ctx context.Context) ([]*domain.TaxZone, error)
}

func (m *mockTaxZoneRepo) Create(ctx context.Context, zone *domain.TaxZone) error {
	if m.createFn != nil {
		return m.createFn(ctx, zone)
	}
	return nil
}

func (m *mockTaxZoneRepo) GetByID(ctx context.Context, id string) (*domain.TaxZone, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, errors.New("not found")
}

func (m *mockTaxZoneRepo) GetByLocation(ctx context.Context, countryCode, stateCode string) (*domain.TaxZone, error) {
	if m.getByLocationFn != nil {
		return m.getByLocationFn(ctx, countryCode, stateCode)
	}
	return nil, errors.New("not found")
}

func (m *mockTaxZoneRepo) List(ctx context.Context) ([]*domain.TaxZone, error) {
	if m.listFn != nil {
		return m.listFn(ctx)
	}
	return nil, nil
}

type mockTaxRuleRepo struct {
	createFn               func(ctx context.Context, rule *domain.TaxRule) error
	getByIDFn              func(ctx context.Context, id string) (*domain.TaxRule, error)
	listByZoneFn           func(ctx context.Context, zoneID string) ([]*domain.TaxRule, error)
	listActiveFn           func(ctx context.Context) ([]*domain.TaxRule, error)
	updateFn               func(ctx context.Context, rule *domain.TaxRule) error
	deleteFn               func(ctx context.Context, id string) error
	getByZoneAndCategoryFn func(ctx context.Context, zoneID, category string) ([]*domain.TaxRule, error)
}

func (m *mockTaxRuleRepo) Create(ctx context.Context, rule *domain.TaxRule) error {
	if m.createFn != nil {
		return m.createFn(ctx, rule)
	}
	return nil
}

func (m *mockTaxRuleRepo) GetByID(ctx context.Context, id string) (*domain.TaxRule, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, errors.New("not found")
}

func (m *mockTaxRuleRepo) ListByZone(ctx context.Context, zoneID string) ([]*domain.TaxRule, error) {
	if m.listByZoneFn != nil {
		return m.listByZoneFn(ctx, zoneID)
	}
	return nil, nil
}

func (m *mockTaxRuleRepo) ListActive(ctx context.Context) ([]*domain.TaxRule, error) {
	if m.listActiveFn != nil {
		return m.listActiveFn(ctx)
	}
	return nil, nil
}

func (m *mockTaxRuleRepo) Update(ctx context.Context, rule *domain.TaxRule) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, rule)
	}
	return nil
}

func (m *mockTaxRuleRepo) Delete(ctx context.Context, id string) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id)
	}
	return nil
}

func (m *mockTaxRuleRepo) GetByZoneAndCategory(ctx context.Context, zoneID, category string) ([]*domain.TaxRule, error) {
	if m.getByZoneAndCategoryFn != nil {
		return m.getByZoneAndCategoryFn(ctx, zoneID, category)
	}
	return nil, nil
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func auZone() *domain.TaxZone {
	return &domain.TaxZone{
		ID:          "zone-au",
		CountryCode: "AU",
		StateCode:   "",
		Name:        "Australia",
	}
}

func usCAZone() *domain.TaxZone {
	return &domain.TaxZone{
		ID:          "zone-us-ca",
		CountryCode: "US",
		StateCode:   "CA",
		Name:        "California",
	}
}

func exclusiveRule(taxName string, rate float64, category string) *domain.TaxRule {
	return &domain.TaxRule{
		ID:       "rule-" + taxName,
		ZoneID:   "zone-us-ca",
		TaxName:  taxName,
		Rate:     rate,
		Category: category,
		Inclusive: false,
		IsActive: true,
	}
}

func inclusiveRule(taxName string, rate float64, category string) *domain.TaxRule {
	return &domain.TaxRule{
		ID:       "rule-" + taxName,
		ZoneID:   "zone-au",
		TaxName:  taxName,
		Rate:     rate,
		Category: category,
		Inclusive: true,
		IsActive: true,
	}
}

// findBreakdown finds a TaxBreakdown by TaxName in the slice.
func findBreakdown(breakdowns []domain.TaxBreakdown, taxName string) *domain.TaxBreakdown {
	for i := range breakdowns {
		if breakdowns[i].TaxName == taxName {
			return &breakdowns[i]
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// Tests for CalculateTaxUseCase.Execute
// ---------------------------------------------------------------------------

func TestCalculateTax_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("zone not found returns zero tax with subtotal", func(t *testing.T) {
		zoneRepo := &mockTaxZoneRepo{
			getByLocationFn: func(_ context.Context, _, _ string) (*domain.TaxZone, error) {
				return nil, errors.New("zone not found")
			},
		}
		ruleRepo := &mockTaxRuleRepo{}
		uc := NewCalculateTaxUseCase(zoneRepo, ruleRepo)

		req := &domain.TaxCalculationRequest{
			Items: []domain.TaxItem{
				{ProductID: "p1", Category: "electronics", PriceCents: 10000, Quantity: 2},
				{ProductID: "p2", Category: "clothing", PriceCents: 5000, Quantity: 1},
			},
			ShippingAddress: domain.TaxAddress{CountryCode: "XX", StateCode: ""},
		}

		result, err := uc.Execute(ctx, req)
		require.NoError(t, err)
		assert.Equal(t, int64(25000), result.SubtotalCents)
		assert.Equal(t, int64(0), result.TaxAmountCents)
		assert.Empty(t, result.Breakdown)
	})

	t.Run("single item exclusive tax 10% GST", func(t *testing.T) {
		zone := usCAZone()
		zoneRepo := &mockTaxZoneRepo{
			getByLocationFn: func(_ context.Context, _, _ string) (*domain.TaxZone, error) {
				return zone, nil
			},
		}
		ruleRepo := &mockTaxRuleRepo{
			getByZoneAndCategoryFn: func(_ context.Context, _, _ string) ([]*domain.TaxRule, error) {
				return []*domain.TaxRule{
					{ID: "r1", ZoneID: zone.ID, TaxName: "GST", Rate: 0.10, Inclusive: false, IsActive: true},
				}, nil
			},
		}
		uc := NewCalculateTaxUseCase(zoneRepo, ruleRepo)

		req := &domain.TaxCalculationRequest{
			Items: []domain.TaxItem{
				{ProductID: "p1", Category: "electronics", PriceCents: 10000, Quantity: 1},
			},
			ShippingAddress: domain.TaxAddress{CountryCode: "US", StateCode: "CA"},
		}

		result, err := uc.Execute(ctx, req)
		require.NoError(t, err)
		assert.Equal(t, int64(10000), result.SubtotalCents)
		// exclusive: tax = 10000 * 0.10 = 1000
		assert.Equal(t, int64(1000), result.TaxAmountCents)
		require.Len(t, result.Breakdown, 1)
		assert.Equal(t, "GST", result.Breakdown[0].TaxName)
		assert.Equal(t, int64(1000), result.Breakdown[0].AmountCents)
		assert.Equal(t, 0.10, result.Breakdown[0].Rate)
		assert.Equal(t, zone.Name, result.Breakdown[0].Jurisdiction)
	})

	t.Run("single item inclusive tax 10% GST Australian style", func(t *testing.T) {
		zone := auZone()
		zoneRepo := &mockTaxZoneRepo{
			getByLocationFn: func(_ context.Context, _, _ string) (*domain.TaxZone, error) {
				return zone, nil
			},
		}
		ruleRepo := &mockTaxRuleRepo{
			getByZoneAndCategoryFn: func(_ context.Context, _, _ string) ([]*domain.TaxRule, error) {
				return []*domain.TaxRule{
					{ID: "r1", ZoneID: zone.ID, TaxName: "GST", Rate: 0.10, Inclusive: true, IsActive: true},
				}, nil
			},
		}
		uc := NewCalculateTaxUseCase(zoneRepo, ruleRepo)

		req := &domain.TaxCalculationRequest{
			Items: []domain.TaxItem{
				{ProductID: "p1", Category: "electronics", PriceCents: 11000, Quantity: 1},
			},
			ShippingAddress: domain.TaxAddress{CountryCode: "AU", StateCode: ""},
		}

		result, err := uc.Execute(ctx, req)
		require.NoError(t, err)
		assert.Equal(t, int64(11000), result.SubtotalCents)
		// inclusive: tax = 11000 - 11000 / 1.10 = 11000 - 10000 = 1000
		expectedTax := int64(math.Round(float64(11000) - float64(11000)/1.10))
		assert.Equal(t, expectedTax, result.TaxAmountCents)
		assert.Equal(t, int64(1000), result.TaxAmountCents)
		require.Len(t, result.Breakdown, 1)
		assert.Equal(t, "GST", result.Breakdown[0].TaxName)
	})

	t.Run("multiple items same category taxes aggregate", func(t *testing.T) {
		zone := usCAZone()
		zoneRepo := &mockTaxZoneRepo{
			getByLocationFn: func(_ context.Context, _, _ string) (*domain.TaxZone, error) {
				return zone, nil
			},
		}
		ruleRepo := &mockTaxRuleRepo{
			getByZoneAndCategoryFn: func(_ context.Context, _, _ string) ([]*domain.TaxRule, error) {
				return []*domain.TaxRule{
					{ID: "r1", ZoneID: zone.ID, TaxName: "Sales Tax", Rate: 0.08, Inclusive: false, IsActive: true},
				}, nil
			},
		}
		uc := NewCalculateTaxUseCase(zoneRepo, ruleRepo)

		req := &domain.TaxCalculationRequest{
			Items: []domain.TaxItem{
				{ProductID: "p1", Category: "electronics", PriceCents: 10000, Quantity: 1},
				{ProductID: "p2", Category: "electronics", PriceCents: 20000, Quantity: 1},
			},
			ShippingAddress: domain.TaxAddress{CountryCode: "US", StateCode: "CA"},
		}

		result, err := uc.Execute(ctx, req)
		require.NoError(t, err)
		assert.Equal(t, int64(30000), result.SubtotalCents)
		// tax = 10000*0.08 + 20000*0.08 = 800 + 1600 = 2400
		assert.Equal(t, int64(2400), result.TaxAmountCents)
		require.Len(t, result.Breakdown, 1)
		assert.Equal(t, "Sales Tax", result.Breakdown[0].TaxName)
		assert.Equal(t, int64(2400), result.Breakdown[0].AmountCents)
	})

	t.Run("multiple items different categories category-specific rules", func(t *testing.T) {
		zone := usCAZone()
		zoneRepo := &mockTaxZoneRepo{
			getByLocationFn: func(_ context.Context, _, _ string) (*domain.TaxZone, error) {
				return zone, nil
			},
		}
		ruleRepo := &mockTaxRuleRepo{
			getByZoneAndCategoryFn: func(_ context.Context, _, category string) ([]*domain.TaxRule, error) {
				switch category {
				case "electronics":
					return []*domain.TaxRule{
						{ID: "r1", ZoneID: zone.ID, TaxName: "Electronics Tax", Rate: 0.10, Inclusive: false, IsActive: true},
					}, nil
				case "food":
					return []*domain.TaxRule{
						{ID: "r2", ZoneID: zone.ID, TaxName: "Food Tax", Rate: 0.02, Inclusive: false, IsActive: true},
					}, nil
				default:
					return nil, nil
				}
			},
		}
		uc := NewCalculateTaxUseCase(zoneRepo, ruleRepo)

		req := &domain.TaxCalculationRequest{
			Items: []domain.TaxItem{
				{ProductID: "p1", Category: "electronics", PriceCents: 10000, Quantity: 1},
				{ProductID: "p2", Category: "food", PriceCents: 5000, Quantity: 1},
			},
			ShippingAddress: domain.TaxAddress{CountryCode: "US", StateCode: "CA"},
		}

		result, err := uc.Execute(ctx, req)
		require.NoError(t, err)
		assert.Equal(t, int64(15000), result.SubtotalCents)
		// electronics: 10000*0.10 = 1000; food: 5000*0.02 = 100
		assert.Equal(t, int64(1100), result.TaxAmountCents)
		require.Len(t, result.Breakdown, 2)

		elecBreakdown := findBreakdown(result.Breakdown, "Electronics Tax")
		require.NotNil(t, elecBreakdown)
		assert.Equal(t, int64(1000), elecBreakdown.AmountCents)

		foodBreakdown := findBreakdown(result.Breakdown, "Food Tax")
		require.NotNil(t, foodBreakdown)
		assert.Equal(t, int64(100), foodBreakdown.AmountCents)
	})

	t.Run("multiple tax rules on same item state plus federal", func(t *testing.T) {
		zone := usCAZone()
		zoneRepo := &mockTaxZoneRepo{
			getByLocationFn: func(_ context.Context, _, _ string) (*domain.TaxZone, error) {
				return zone, nil
			},
		}
		ruleRepo := &mockTaxRuleRepo{
			getByZoneAndCategoryFn: func(_ context.Context, _, _ string) ([]*domain.TaxRule, error) {
				return []*domain.TaxRule{
					{ID: "r1", ZoneID: zone.ID, TaxName: "State Tax", Rate: 0.06, Inclusive: false, IsActive: true},
					{ID: "r2", ZoneID: zone.ID, TaxName: "Federal Tax", Rate: 0.04, Inclusive: false, IsActive: true},
				}, nil
			},
		}
		uc := NewCalculateTaxUseCase(zoneRepo, ruleRepo)

		req := &domain.TaxCalculationRequest{
			Items: []domain.TaxItem{
				{ProductID: "p1", Category: "electronics", PriceCents: 20000, Quantity: 1},
			},
			ShippingAddress: domain.TaxAddress{CountryCode: "US", StateCode: "CA"},
		}

		result, err := uc.Execute(ctx, req)
		require.NoError(t, err)
		assert.Equal(t, int64(20000), result.SubtotalCents)
		// state: 20000*0.06 = 1200; federal: 20000*0.04 = 800
		assert.Equal(t, int64(2000), result.TaxAmountCents)
		require.Len(t, result.Breakdown, 2)

		stateTax := findBreakdown(result.Breakdown, "State Tax")
		require.NotNil(t, stateTax)
		assert.Equal(t, int64(1200), stateTax.AmountCents)

		fedTax := findBreakdown(result.Breakdown, "Federal Tax")
		require.NotNil(t, fedTax)
		assert.Equal(t, int64(800), fedTax.AmountCents)
	})

	t.Run("items with quantity greater than 1 multiplied correctly", func(t *testing.T) {
		zone := usCAZone()
		zoneRepo := &mockTaxZoneRepo{
			getByLocationFn: func(_ context.Context, _, _ string) (*domain.TaxZone, error) {
				return zone, nil
			},
		}
		ruleRepo := &mockTaxRuleRepo{
			getByZoneAndCategoryFn: func(_ context.Context, _, _ string) ([]*domain.TaxRule, error) {
				return []*domain.TaxRule{
					{ID: "r1", ZoneID: zone.ID, TaxName: "Sales Tax", Rate: 0.10, Inclusive: false, IsActive: true},
				}, nil
			},
		}
		uc := NewCalculateTaxUseCase(zoneRepo, ruleRepo)

		req := &domain.TaxCalculationRequest{
			Items: []domain.TaxItem{
				{ProductID: "p1", Category: "electronics", PriceCents: 5000, Quantity: 3},
			},
			ShippingAddress: domain.TaxAddress{CountryCode: "US", StateCode: "CA"},
		}

		result, err := uc.Execute(ctx, req)
		require.NoError(t, err)
		// subtotal: 5000*3 = 15000
		assert.Equal(t, int64(15000), result.SubtotalCents)
		// tax: 15000*0.10 = 1500
		assert.Equal(t, int64(1500), result.TaxAmountCents)
	})

	t.Run("no matching rules for category returns zero tax for item", func(t *testing.T) {
		zone := usCAZone()
		zoneRepo := &mockTaxZoneRepo{
			getByLocationFn: func(_ context.Context, _, _ string) (*domain.TaxZone, error) {
				return zone, nil
			},
		}
		ruleRepo := &mockTaxRuleRepo{
			getByZoneAndCategoryFn: func(_ context.Context, _, _ string) ([]*domain.TaxRule, error) {
				// no rules at all
				return []*domain.TaxRule{}, nil
			},
		}
		uc := NewCalculateTaxUseCase(zoneRepo, ruleRepo)

		req := &domain.TaxCalculationRequest{
			Items: []domain.TaxItem{
				{ProductID: "p1", Category: "tax-exempt", PriceCents: 50000, Quantity: 1},
			},
			ShippingAddress: domain.TaxAddress{CountryCode: "US", StateCode: "CA"},
		}

		result, err := uc.Execute(ctx, req)
		require.NoError(t, err)
		assert.Equal(t, int64(50000), result.SubtotalCents)
		assert.Equal(t, int64(0), result.TaxAmountCents)
		assert.Empty(t, result.Breakdown)
	})

	t.Run("mixed inclusive and exclusive rules", func(t *testing.T) {
		zone := &domain.TaxZone{ID: "zone-mix", CountryCode: "XX", Name: "MixedZone"}
		zoneRepo := &mockTaxZoneRepo{
			getByLocationFn: func(_ context.Context, _, _ string) (*domain.TaxZone, error) {
				return zone, nil
			},
		}
		ruleRepo := &mockTaxRuleRepo{
			getByZoneAndCategoryFn: func(_ context.Context, _, _ string) ([]*domain.TaxRule, error) {
				return []*domain.TaxRule{
					{ID: "r1", ZoneID: zone.ID, TaxName: "Inclusive GST", Rate: 0.10, Inclusive: true, IsActive: true},
					{ID: "r2", ZoneID: zone.ID, TaxName: "Exclusive Levy", Rate: 0.05, Inclusive: false, IsActive: true},
				}, nil
			},
		}
		uc := NewCalculateTaxUseCase(zoneRepo, ruleRepo)

		req := &domain.TaxCalculationRequest{
			Items: []domain.TaxItem{
				{ProductID: "p1", Category: "general", PriceCents: 11000, Quantity: 1},
			},
			ShippingAddress: domain.TaxAddress{CountryCode: "XX"},
		}

		result, err := uc.Execute(ctx, req)
		require.NoError(t, err)
		assert.Equal(t, int64(11000), result.SubtotalCents)

		// inclusive: 11000 - 11000/1.10 = 11000 - 10000 = 1000
		inclusiveTax := int64(math.Round(float64(11000) - float64(11000)/1.10))
		// exclusive: 11000 * 0.05 = 550
		exclusiveTax := int64(math.Round(float64(11000) * 0.05))
		expectedTotal := inclusiveTax + exclusiveTax
		assert.Equal(t, expectedTotal, result.TaxAmountCents)
		require.Len(t, result.Breakdown, 2)

		gstBreakdown := findBreakdown(result.Breakdown, "Inclusive GST")
		require.NotNil(t, gstBreakdown)
		assert.Equal(t, inclusiveTax, gstBreakdown.AmountCents)

		levyBreakdown := findBreakdown(result.Breakdown, "Exclusive Levy")
		require.NotNil(t, levyBreakdown)
		assert.Equal(t, exclusiveTax, levyBreakdown.AmountCents)
	})

	t.Run("breakdown groups by TaxName across items", func(t *testing.T) {
		zone := auZone()
		zoneRepo := &mockTaxZoneRepo{
			getByLocationFn: func(_ context.Context, _, _ string) (*domain.TaxZone, error) {
				return zone, nil
			},
		}
		ruleRepo := &mockTaxRuleRepo{
			getByZoneAndCategoryFn: func(_ context.Context, _, _ string) ([]*domain.TaxRule, error) {
				return []*domain.TaxRule{
					{ID: "r1", ZoneID: zone.ID, TaxName: "GST", Rate: 0.10, Inclusive: false, IsActive: true},
				}, nil
			},
		}
		uc := NewCalculateTaxUseCase(zoneRepo, ruleRepo)

		req := &domain.TaxCalculationRequest{
			Items: []domain.TaxItem{
				{ProductID: "p1", Category: "books", PriceCents: 2000, Quantity: 1},
				{ProductID: "p2", Category: "toys", PriceCents: 3000, Quantity: 1},
			},
			ShippingAddress: domain.TaxAddress{CountryCode: "AU"},
		}

		result, err := uc.Execute(ctx, req)
		require.NoError(t, err)

		// Both items fall under same "GST" rule, so only one breakdown entry
		require.Len(t, result.Breakdown, 1)
		assert.Equal(t, "GST", result.Breakdown[0].TaxName)
		// 2000*0.10 + 3000*0.10 = 200 + 300 = 500
		assert.Equal(t, int64(500), result.Breakdown[0].AmountCents)
		assert.Equal(t, int64(500), result.TaxAmountCents)
	})

	t.Run("empty items list returns zero subtotal and zero tax", func(t *testing.T) {
		zone := usCAZone()
		zoneRepo := &mockTaxZoneRepo{
			getByLocationFn: func(_ context.Context, _, _ string) (*domain.TaxZone, error) {
				return zone, nil
			},
		}
		ruleRepo := &mockTaxRuleRepo{}
		uc := NewCalculateTaxUseCase(zoneRepo, ruleRepo)

		req := &domain.TaxCalculationRequest{
			Items:           []domain.TaxItem{},
			ShippingAddress: domain.TaxAddress{CountryCode: "US", StateCode: "CA"},
		}

		result, err := uc.Execute(ctx, req)
		require.NoError(t, err)
		assert.Equal(t, int64(0), result.SubtotalCents)
		assert.Equal(t, int64(0), result.TaxAmountCents)
		assert.Empty(t, result.Breakdown)
	})

	t.Run("GetByZoneAndCategory repo error propagates", func(t *testing.T) {
		zone := usCAZone()
		zoneRepo := &mockTaxZoneRepo{
			getByLocationFn: func(_ context.Context, _, _ string) (*domain.TaxZone, error) {
				return zone, nil
			},
		}
		ruleRepo := &mockTaxRuleRepo{
			getByZoneAndCategoryFn: func(_ context.Context, _, _ string) ([]*domain.TaxRule, error) {
				return nil, errors.New("database error")
			},
		}
		uc := NewCalculateTaxUseCase(zoneRepo, ruleRepo)

		req := &domain.TaxCalculationRequest{
			Items: []domain.TaxItem{
				{ProductID: "p1", Category: "electronics", PriceCents: 10000, Quantity: 1},
			},
			ShippingAddress: domain.TaxAddress{CountryCode: "US", StateCode: "CA"},
		}

		result, err := uc.Execute(ctx, req)
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "database error")
	})

	t.Run("inclusive tax with non-round result rounds correctly", func(t *testing.T) {
		zone := auZone()
		zoneRepo := &mockTaxZoneRepo{
			getByLocationFn: func(_ context.Context, _, _ string) (*domain.TaxZone, error) {
				return zone, nil
			},
		}
		ruleRepo := &mockTaxRuleRepo{
			getByZoneAndCategoryFn: func(_ context.Context, _, _ string) ([]*domain.TaxRule, error) {
				return []*domain.TaxRule{
					{ID: "r1", ZoneID: zone.ID, TaxName: "GST", Rate: 0.10, Inclusive: true, IsActive: true},
				}, nil
			},
		}
		uc := NewCalculateTaxUseCase(zoneRepo, ruleRepo)

		// Use a price that does not divide evenly: 9999 cents
		req := &domain.TaxCalculationRequest{
			Items: []domain.TaxItem{
				{ProductID: "p1", Category: "general", PriceCents: 9999, Quantity: 1},
			},
			ShippingAddress: domain.TaxAddress{CountryCode: "AU"},
		}

		result, err := uc.Execute(ctx, req)
		require.NoError(t, err)
		// inclusive: 9999 - 9999/1.10 = 9999 - 9090 = 909 (rounded)
		expected := int64(math.Round(float64(9999) - float64(9999)/1.10))
		assert.Equal(t, expected, result.TaxAmountCents)
	})

	t.Run("exclusive tax with non-round result rounds correctly", func(t *testing.T) {
		zone := usCAZone()
		zoneRepo := &mockTaxZoneRepo{
			getByLocationFn: func(_ context.Context, _, _ string) (*domain.TaxZone, error) {
				return zone, nil
			},
		}
		ruleRepo := &mockTaxRuleRepo{
			getByZoneAndCategoryFn: func(_ context.Context, _, _ string) ([]*domain.TaxRule, error) {
				return []*domain.TaxRule{
					{ID: "r1", ZoneID: zone.ID, TaxName: "Sales Tax", Rate: 0.0725, Inclusive: false, IsActive: true},
				}, nil
			},
		}
		uc := NewCalculateTaxUseCase(zoneRepo, ruleRepo)

		// 3333 * 0.0725 = 241.6425 => 242 rounded
		req := &domain.TaxCalculationRequest{
			Items: []domain.TaxItem{
				{ProductID: "p1", Category: "general", PriceCents: 3333, Quantity: 1},
			},
			ShippingAddress: domain.TaxAddress{CountryCode: "US", StateCode: "CA"},
		}

		result, err := uc.Execute(ctx, req)
		require.NoError(t, err)
		expected := int64(math.Round(float64(3333) * 0.0725))
		assert.Equal(t, expected, result.TaxAmountCents)
	})

	t.Run("multiple items with quantity greater than 1 and multiple rules", func(t *testing.T) {
		zone := usCAZone()
		zoneRepo := &mockTaxZoneRepo{
			getByLocationFn: func(_ context.Context, _, _ string) (*domain.TaxZone, error) {
				return zone, nil
			},
		}
		ruleRepo := &mockTaxRuleRepo{
			getByZoneAndCategoryFn: func(_ context.Context, _, _ string) ([]*domain.TaxRule, error) {
				return []*domain.TaxRule{
					{ID: "r1", ZoneID: zone.ID, TaxName: "State Tax", Rate: 0.05, Inclusive: false, IsActive: true},
					{ID: "r2", ZoneID: zone.ID, TaxName: "County Tax", Rate: 0.02, Inclusive: false, IsActive: true},
				}, nil
			},
		}
		uc := NewCalculateTaxUseCase(zoneRepo, ruleRepo)

		req := &domain.TaxCalculationRequest{
			Items: []domain.TaxItem{
				{ProductID: "p1", Category: "general", PriceCents: 1000, Quantity: 5},
				{ProductID: "p2", Category: "general", PriceCents: 2000, Quantity: 3},
			},
			ShippingAddress: domain.TaxAddress{CountryCode: "US", StateCode: "CA"},
		}

		result, err := uc.Execute(ctx, req)
		require.NoError(t, err)
		// subtotal: 1000*5 + 2000*3 = 5000 + 6000 = 11000
		assert.Equal(t, int64(11000), result.SubtotalCents)

		// State Tax: (5000*0.05) + (6000*0.05) = 250 + 300 = 550
		stateTax := findBreakdown(result.Breakdown, "State Tax")
		require.NotNil(t, stateTax)
		assert.Equal(t, int64(550), stateTax.AmountCents)

		// County Tax: (5000*0.02) + (6000*0.02) = 100 + 120 = 220
		countyTax := findBreakdown(result.Breakdown, "County Tax")
		require.NotNil(t, countyTax)
		assert.Equal(t, int64(220), countyTax.AmountCents)

		assert.Equal(t, int64(770), result.TaxAmountCents)
	})

	t.Run("breakdown stores jurisdiction from zone name", func(t *testing.T) {
		zone := &domain.TaxZone{ID: "zone-de", CountryCode: "DE", Name: "Germany"}
		zoneRepo := &mockTaxZoneRepo{
			getByLocationFn: func(_ context.Context, _, _ string) (*domain.TaxZone, error) {
				return zone, nil
			},
		}
		ruleRepo := &mockTaxRuleRepo{
			getByZoneAndCategoryFn: func(_ context.Context, _, _ string) ([]*domain.TaxRule, error) {
				return []*domain.TaxRule{
					{ID: "r1", ZoneID: zone.ID, TaxName: "VAT", Rate: 0.19, Inclusive: false, IsActive: true},
				}, nil
			},
		}
		uc := NewCalculateTaxUseCase(zoneRepo, ruleRepo)

		req := &domain.TaxCalculationRequest{
			Items: []domain.TaxItem{
				{ProductID: "p1", Category: "general", PriceCents: 10000, Quantity: 1},
			},
			ShippingAddress: domain.TaxAddress{CountryCode: "DE"},
		}

		result, err := uc.Execute(ctx, req)
		require.NoError(t, err)
		require.Len(t, result.Breakdown, 1)
		assert.Equal(t, "Germany", result.Breakdown[0].Jurisdiction)
	})

	t.Run("some items taxed and some not", func(t *testing.T) {
		zone := usCAZone()
		zoneRepo := &mockTaxZoneRepo{
			getByLocationFn: func(_ context.Context, _, _ string) (*domain.TaxZone, error) {
				return zone, nil
			},
		}
		ruleRepo := &mockTaxRuleRepo{
			getByZoneAndCategoryFn: func(_ context.Context, _, category string) ([]*domain.TaxRule, error) {
				if category == "taxable" {
					return []*domain.TaxRule{
						{ID: "r1", ZoneID: zone.ID, TaxName: "Sales Tax", Rate: 0.08, Inclusive: false, IsActive: true},
					}, nil
				}
				// "exempt" category gets no rules
				return []*domain.TaxRule{}, nil
			},
		}
		uc := NewCalculateTaxUseCase(zoneRepo, ruleRepo)

		req := &domain.TaxCalculationRequest{
			Items: []domain.TaxItem{
				{ProductID: "p1", Category: "taxable", PriceCents: 10000, Quantity: 1},
				{ProductID: "p2", Category: "exempt", PriceCents: 5000, Quantity: 2},
			},
			ShippingAddress: domain.TaxAddress{CountryCode: "US", StateCode: "CA"},
		}

		result, err := uc.Execute(ctx, req)
		require.NoError(t, err)
		// subtotal = 10000 + 5000*2 = 20000
		assert.Equal(t, int64(20000), result.SubtotalCents)
		// tax only on taxable: 10000 * 0.08 = 800
		assert.Equal(t, int64(800), result.TaxAmountCents)
		require.Len(t, result.Breakdown, 1)
		assert.Equal(t, "Sales Tax", result.Breakdown[0].TaxName)
	})

	t.Run("breakdown stores rate from rule", func(t *testing.T) {
		zone := usCAZone()
		zoneRepo := &mockTaxZoneRepo{
			getByLocationFn: func(_ context.Context, _, _ string) (*domain.TaxZone, error) {
				return zone, nil
			},
		}
		ruleRepo := &mockTaxRuleRepo{
			getByZoneAndCategoryFn: func(_ context.Context, _, _ string) ([]*domain.TaxRule, error) {
				return []*domain.TaxRule{
					{ID: "r1", ZoneID: zone.ID, TaxName: "Sales Tax", Rate: 0.0725, Inclusive: false, IsActive: true},
				}, nil
			},
		}
		uc := NewCalculateTaxUseCase(zoneRepo, ruleRepo)

		req := &domain.TaxCalculationRequest{
			Items: []domain.TaxItem{
				{ProductID: "p1", Category: "general", PriceCents: 10000, Quantity: 1},
			},
			ShippingAddress: domain.TaxAddress{CountryCode: "US", StateCode: "CA"},
		}

		result, err := uc.Execute(ctx, req)
		require.NoError(t, err)
		require.Len(t, result.Breakdown, 1)
		assert.Equal(t, 0.0725, result.Breakdown[0].Rate)
	})
}
