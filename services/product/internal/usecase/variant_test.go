package usecase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/southern-martin/ecommerce/services/product/internal/domain"
)

// Mocks are defined in product_test.go (same package).

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func newVariantUseCase(
	productRepo *mockProductRepo,
	optionRepo *mockOptionRepo,
	variantRepo *mockVariantRepo,
	eventPub *mockEventPub,
) *VariantUseCase {
	return NewVariantUseCase(productRepo, optionRepo, variantRepo, eventPub)
}

// configurableProduct returns a configurable product owned by the given seller.
func configurableProduct(id, sellerID string) *domain.Product {
	return &domain.Product{
		ID:             id,
		SellerID:       sellerID,
		Name:           "Test Configurable",
		Slug:           "test-configurable-abcd1234",
		BasePriceCents: 5000,
		ProductType:    domain.ProductTypeConfigurable,
		HasVariants:    false,
	}
}

// simpleProduct returns a simple (non-configurable) product.
func simpleProduct(id, sellerID string) *domain.Product {
	return &domain.Product{
		ID:          id,
		SellerID:    sellerID,
		Name:        "Simple Widget",
		Slug:        "simple-widget-abcd1234",
		ProductType: domain.ProductTypeSimple,
	}
}

// sizeOptions returns a single "Size" option with the given values.
func sizeOptions(values ...string) []domain.ProductOption {
	var pvs []domain.ProductOptionValue
	for i, v := range values {
		pvs = append(pvs, domain.ProductOptionValue{
			ID:       "ov-" + v,
			OptionID: "opt-size",
			Value:    v,
		})
		_ = i
	}
	return []domain.ProductOption{
		{ID: "opt-size", ProductID: "prod-1", Name: "Size", Values: pvs},
	}
}

// colorOptions returns a single "Color" option with the given values.
func colorOptions(values ...string) []domain.ProductOption {
	var pvs []domain.ProductOptionValue
	for _, v := range values {
		pvs = append(pvs, domain.ProductOptionValue{
			ID:       "ov-" + v,
			OptionID: "opt-color",
			Value:    v,
		})
	}
	return []domain.ProductOption{
		{ID: "opt-color", ProductID: "prod-1", Name: "Color", Values: pvs},
	}
}

// ===========================================================================
// GenerateVariants tests
// ===========================================================================

func TestGenerateVariants_Success(t *testing.T) {
	pRepo := &mockProductRepo{
		getByIDFn: func(_ context.Context, _ string) (*domain.Product, error) {
			return configurableProduct("prod-1", "seller-1"), nil
		},
		updateFn: func(_ context.Context, _ *domain.Product) error { return nil },
	}
	oRepo := &mockOptionRepo{
		listByProductFn: func(_ context.Context, _ string) ([]domain.ProductOption, error) {
			return append(
				sizeOptions("S", "M"),
				colorOptions("Red", "Blue", "Green")...,
			), nil
		},
	}
	var bulkCreated []domain.Variant
	vRepo := &mockVariantRepo{
		bulkCreateFn: func(_ context.Context, variants []domain.Variant) error {
			bulkCreated = variants
			return nil
		},
	}

	uc := newVariantUseCase(pRepo, oRepo, vRepo, &mockEventPub{})
	variants, err := uc.GenerateVariants(context.Background(), "prod-1", "seller-1")

	require.NoError(t, err)
	// 2 sizes x 3 colors = 6 variants
	assert.Len(t, variants, 6)
	assert.Len(t, bulkCreated, 6)
}

func TestGenerateVariants_WrongSeller(t *testing.T) {
	pRepo := &mockProductRepo{
		getByIDFn: func(_ context.Context, _ string) (*domain.Product, error) {
			return configurableProduct("prod-1", "seller-1"), nil
		},
	}

	uc := newVariantUseCase(pRepo, &mockOptionRepo{}, &mockVariantRepo{}, &mockEventPub{})
	_, err := uc.GenerateVariants(context.Background(), "prod-1", "seller-WRONG")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "unauthorized")
}

func TestGenerateVariants_SimpleProductRejection(t *testing.T) {
	pRepo := &mockProductRepo{
		getByIDFn: func(_ context.Context, _ string) (*domain.Product, error) {
			return simpleProduct("prod-1", "seller-1"), nil
		},
	}

	uc := newVariantUseCase(pRepo, &mockOptionRepo{}, &mockVariantRepo{}, &mockEventPub{})
	_, err := uc.GenerateVariants(context.Background(), "prod-1", "seller-1")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "cannot generate variants")
	assert.Contains(t, err.Error(), "simple")
}

func TestGenerateVariants_NoOptionsError(t *testing.T) {
	pRepo := &mockProductRepo{
		getByIDFn: func(_ context.Context, _ string) (*domain.Product, error) {
			return configurableProduct("prod-1", "seller-1"), nil
		},
	}
	oRepo := &mockOptionRepo{
		listByProductFn: func(_ context.Context, _ string) ([]domain.ProductOption, error) {
			return []domain.ProductOption{}, nil
		},
	}

	uc := newVariantUseCase(pRepo, oRepo, &mockVariantRepo{}, &mockEventPub{})
	_, err := uc.GenerateVariants(context.Background(), "prod-1", "seller-1")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "no options defined")
}

func TestGenerateVariants_SingleOption_3Values(t *testing.T) {
	pRepo := &mockProductRepo{
		getByIDFn: func(_ context.Context, _ string) (*domain.Product, error) {
			return configurableProduct("prod-1", "seller-1"), nil
		},
		updateFn: func(_ context.Context, _ *domain.Product) error { return nil },
	}
	oRepo := &mockOptionRepo{
		listByProductFn: func(_ context.Context, _ string) ([]domain.ProductOption, error) {
			return sizeOptions("S", "M", "L"), nil
		},
	}
	vRepo := &mockVariantRepo{
		bulkCreateFn: func(_ context.Context, _ []domain.Variant) error { return nil },
	}

	uc := newVariantUseCase(pRepo, oRepo, vRepo, &mockEventPub{})
	variants, err := uc.GenerateVariants(context.Background(), "prod-1", "seller-1")

	require.NoError(t, err)
	assert.Len(t, variants, 3)
}

func TestGenerateVariants_TwoOptions_2x3(t *testing.T) {
	pRepo := &mockProductRepo{
		getByIDFn: func(_ context.Context, _ string) (*domain.Product, error) {
			return configurableProduct("prod-1", "seller-1"), nil
		},
		updateFn: func(_ context.Context, _ *domain.Product) error { return nil },
	}
	oRepo := &mockOptionRepo{
		listByProductFn: func(_ context.Context, _ string) ([]domain.ProductOption, error) {
			return append(sizeOptions("S", "M"), colorOptions("Red", "Blue", "Green")...), nil
		},
	}
	vRepo := &mockVariantRepo{
		bulkCreateFn: func(_ context.Context, _ []domain.Variant) error { return nil },
	}

	uc := newVariantUseCase(pRepo, oRepo, vRepo, &mockEventPub{})
	variants, err := uc.GenerateVariants(context.Background(), "prod-1", "seller-1")

	require.NoError(t, err)
	assert.Len(t, variants, 6) // 2 x 3
}

func TestGenerateVariants_VariantNamesJoinedWithSlash(t *testing.T) {
	pRepo := &mockProductRepo{
		getByIDFn: func(_ context.Context, _ string) (*domain.Product, error) {
			return configurableProduct("prod-1", "seller-1"), nil
		},
		updateFn: func(_ context.Context, _ *domain.Product) error { return nil },
	}
	oRepo := &mockOptionRepo{
		listByProductFn: func(_ context.Context, _ string) ([]domain.ProductOption, error) {
			return append(sizeOptions("S"), colorOptions("Red")...), nil
		},
	}
	vRepo := &mockVariantRepo{
		bulkCreateFn: func(_ context.Context, _ []domain.Variant) error { return nil },
	}

	uc := newVariantUseCase(pRepo, oRepo, vRepo, &mockEventPub{})
	variants, err := uc.GenerateVariants(context.Background(), "prod-1", "seller-1")

	require.NoError(t, err)
	require.Len(t, variants, 1)
	assert.Equal(t, "S / Red", variants[0].Name)
}

func TestGenerateVariants_FirstVariantIsDefault(t *testing.T) {
	pRepo := &mockProductRepo{
		getByIDFn: func(_ context.Context, _ string) (*domain.Product, error) {
			return configurableProduct("prod-1", "seller-1"), nil
		},
		updateFn: func(_ context.Context, _ *domain.Product) error { return nil },
	}
	oRepo := &mockOptionRepo{
		listByProductFn: func(_ context.Context, _ string) ([]domain.ProductOption, error) {
			return sizeOptions("S", "M", "L"), nil
		},
	}
	vRepo := &mockVariantRepo{
		bulkCreateFn: func(_ context.Context, _ []domain.Variant) error { return nil },
	}

	uc := newVariantUseCase(pRepo, oRepo, vRepo, &mockEventPub{})
	variants, err := uc.GenerateVariants(context.Background(), "prod-1", "seller-1")

	require.NoError(t, err)
	require.Len(t, variants, 3)
	assert.True(t, variants[0].IsDefault, "first variant should be default")
	assert.False(t, variants[1].IsDefault, "second variant should not be default")
	assert.False(t, variants[2].IsDefault, "third variant should not be default")
}

// ===========================================================================
// cartesianProduct tests
// ===========================================================================

func TestCartesianProduct_EmptyOptions(t *testing.T) {
	result := cartesianProduct(nil)
	assert.Nil(t, result)

	result = cartesianProduct([]domain.ProductOption{})
	assert.Nil(t, result)
}

func TestCartesianProduct_SingleOption(t *testing.T) {
	options := sizeOptions("S", "M", "L")
	result := cartesianProduct(options)

	assert.Len(t, result, 3)
	for _, combo := range result {
		assert.Len(t, combo, 1)
	}
	assert.Equal(t, "S", result[0][0].Value)
	assert.Equal(t, "M", result[1][0].Value)
	assert.Equal(t, "L", result[2][0].Value)
}

func TestCartesianProduct_TwoOptions(t *testing.T) {
	options := append(sizeOptions("S", "M"), colorOptions("Red", "Blue", "Green")...)
	result := cartesianProduct(options)

	assert.Len(t, result, 6) // 2 x 3
	for _, combo := range result {
		assert.Len(t, combo, 2, "each combo should have values from both options")
	}
}

func TestCartesianProduct_ThreeOptions_2x2x2(t *testing.T) {
	options := []domain.ProductOption{
		{
			ID: "opt-a", Name: "A",
			Values: []domain.ProductOptionValue{
				{ID: "a1", OptionID: "opt-a", Value: "A1"},
				{ID: "a2", OptionID: "opt-a", Value: "A2"},
			},
		},
		{
			ID: "opt-b", Name: "B",
			Values: []domain.ProductOptionValue{
				{ID: "b1", OptionID: "opt-b", Value: "B1"},
				{ID: "b2", OptionID: "opt-b", Value: "B2"},
			},
		},
		{
			ID: "opt-c", Name: "C",
			Values: []domain.ProductOptionValue{
				{ID: "c1", OptionID: "opt-c", Value: "C1"},
				{ID: "c2", OptionID: "opt-c", Value: "C2"},
			},
		},
	}
	result := cartesianProduct(options)

	assert.Len(t, result, 8) // 2 x 2 x 2
	for _, combo := range result {
		assert.Len(t, combo, 3, "each combo should have 3 values")
	}
}

// ===========================================================================
// joinNames tests
// ===========================================================================

func TestJoinNames_Empty(t *testing.T) {
	assert.Equal(t, "", joinNames(nil))
	assert.Equal(t, "", joinNames([]string{}))
}

func TestJoinNames_Single(t *testing.T) {
	assert.Equal(t, "Large", joinNames([]string{"Large"}))
}

func TestJoinNames_Multiple(t *testing.T) {
	assert.Equal(t, "Large / Red / Cotton", joinNames([]string{"Large", "Red", "Cotton"}))
}

// ===========================================================================
// AddOption tests
// ===========================================================================

func TestAddOption_Success(t *testing.T) {
	pRepo := &mockProductRepo{
		getByIDFn: func(_ context.Context, _ string) (*domain.Product, error) {
			return configurableProduct("prod-1", "seller-1"), nil
		},
		updateFn: func(_ context.Context, _ *domain.Product) error { return nil },
	}
	var createdOption *domain.ProductOption
	oRepo := &mockOptionRepo{
		createOptionFn: func(_ context.Context, opt *domain.ProductOption) error {
			createdOption = opt
			return nil
		},
	}

	uc := newVariantUseCase(pRepo, oRepo, &mockVariantRepo{}, &mockEventPub{})
	option, err := uc.AddOption(context.Background(), "prod-1", "seller-1", AddOptionInput{
		Name: "Size",
		Values: []OptionValueInput{
			{Value: "S"},
			{Value: "M"},
			{Value: "L"},
		},
	})

	require.NoError(t, err)
	require.NotNil(t, option)
	assert.Equal(t, "Size", option.Name)
	assert.Len(t, option.Values, 3)
	assert.NotNil(t, createdOption)
}

func TestAddOption_WrongSeller(t *testing.T) {
	pRepo := &mockProductRepo{
		getByIDFn: func(_ context.Context, _ string) (*domain.Product, error) {
			return configurableProduct("prod-1", "seller-1"), nil
		},
	}

	uc := newVariantUseCase(pRepo, &mockOptionRepo{}, &mockVariantRepo{}, &mockEventPub{})
	_, err := uc.AddOption(context.Background(), "prod-1", "seller-WRONG", AddOptionInput{Name: "Size"})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "unauthorized")
}

func TestAddOption_SimpleProductRejection(t *testing.T) {
	pRepo := &mockProductRepo{
		getByIDFn: func(_ context.Context, _ string) (*domain.Product, error) {
			return simpleProduct("prod-1", "seller-1"), nil
		},
	}

	uc := newVariantUseCase(pRepo, &mockOptionRepo{}, &mockVariantRepo{}, &mockEventPub{})
	_, err := uc.AddOption(context.Background(), "prod-1", "seller-1", AddOptionInput{Name: "Size"})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "cannot add options")
	assert.Contains(t, err.Error(), "simple")
}

func TestAddOption_EmptyName(t *testing.T) {
	pRepo := &mockProductRepo{
		getByIDFn: func(_ context.Context, _ string) (*domain.Product, error) {
			return configurableProduct("prod-1", "seller-1"), nil
		},
	}

	uc := newVariantUseCase(pRepo, &mockOptionRepo{}, &mockVariantRepo{}, &mockEventPub{})
	_, err := uc.AddOption(context.Background(), "prod-1", "seller-1", AddOptionInput{Name: ""})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "option name is required")
}

// ===========================================================================
// UpdateVariant tests
// ===========================================================================

func TestUpdateVariant_Success(t *testing.T) {
	pRepo := &mockProductRepo{
		getByIDFn: func(_ context.Context, _ string) (*domain.Product, error) {
			return configurableProduct("prod-1", "seller-1"), nil
		},
	}
	vRepo := &mockVariantRepo{
		getByIDFn: func(_ context.Context, _ string) (*domain.Variant, error) {
			return &domain.Variant{
				ID:         "var-1",
				ProductID:  "prod-1",
				Name:       "S / Red",
				PriceCents: 5000,
				Stock:      10,
				IsActive:   true,
			}, nil
		},
		updateFn: func(_ context.Context, _ *domain.Variant) error { return nil },
	}

	uc := newVariantUseCase(pRepo, &mockOptionRepo{}, vRepo, &mockEventPub{})
	newName := "Updated Variant"
	newPrice := int64(6000)
	variant, err := uc.UpdateVariant(context.Background(), "prod-1", "var-1", "seller-1", UpdateVariantInput{
		Name:       &newName,
		PriceCents: &newPrice,
	})

	require.NoError(t, err)
	assert.Equal(t, "Updated Variant", variant.Name)
	assert.Equal(t, int64(6000), variant.PriceCents)
}

func TestUpdateVariant_WrongProduct(t *testing.T) {
	pRepo := &mockProductRepo{
		getByIDFn: func(_ context.Context, _ string) (*domain.Product, error) {
			return configurableProduct("prod-1", "seller-1"), nil
		},
	}
	vRepo := &mockVariantRepo{
		getByIDFn: func(_ context.Context, _ string) (*domain.Variant, error) {
			return &domain.Variant{
				ID:        "var-1",
				ProductID: "prod-OTHER", // belongs to a different product
			}, nil
		},
	}

	uc := newVariantUseCase(pRepo, &mockOptionRepo{}, vRepo, &mockEventPub{})
	_, err := uc.UpdateVariant(context.Background(), "prod-1", "var-1", "seller-1", UpdateVariantInput{})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "variant does not belong to this product")
}

// ===========================================================================
// UpdateStock tests
// ===========================================================================

func TestUpdateStock_Success(t *testing.T) {
	pRepo := &mockProductRepo{
		getByIDFn: func(_ context.Context, _ string) (*domain.Product, error) {
			return configurableProduct("prod-1", "seller-1"), nil
		},
	}
	var updatedVariantID string
	var updatedDelta int
	vRepo := &mockVariantRepo{
		updateStockFn: func(_ context.Context, variantID string, delta int) error {
			updatedVariantID = variantID
			updatedDelta = delta
			return nil
		},
		getByIDFn: func(_ context.Context, _ string) (*domain.Variant, error) {
			return &domain.Variant{ID: "var-1", ProductID: "prod-1", Stock: 15}, nil
		},
	}

	uc := newVariantUseCase(pRepo, &mockOptionRepo{}, vRepo, &mockEventPub{})
	err := uc.UpdateStock(context.Background(), "prod-1", "var-1", "seller-1", 5)

	require.NoError(t, err)
	assert.Equal(t, "var-1", updatedVariantID)
	assert.Equal(t, 5, updatedDelta)
}

func TestUpdateStock_WrongSeller(t *testing.T) {
	pRepo := &mockProductRepo{
		getByIDFn: func(_ context.Context, _ string) (*domain.Product, error) {
			return configurableProduct("prod-1", "seller-1"), nil
		},
	}

	uc := newVariantUseCase(pRepo, &mockOptionRepo{}, &mockVariantRepo{}, &mockEventPub{})
	err := uc.UpdateStock(context.Background(), "prod-1", "var-1", "seller-WRONG", 5)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "unauthorized")
}
