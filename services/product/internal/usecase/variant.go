package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/southern-martin/ecommerce/services/product/internal/domain"
)

// VariantUseCase handles variant and option business logic.
type VariantUseCase struct {
	productRepo domain.ProductRepository
	optionRepo  domain.OptionRepository
	variantRepo domain.VariantRepository
	eventPub    domain.EventPublisher
}

// NewVariantUseCase creates a new VariantUseCase.
func NewVariantUseCase(
	productRepo domain.ProductRepository,
	optionRepo domain.OptionRepository,
	variantRepo domain.VariantRepository,
	eventPub domain.EventPublisher,
) *VariantUseCase {
	return &VariantUseCase{
		productRepo: productRepo,
		optionRepo:  optionRepo,
		variantRepo: variantRepo,
		eventPub:    eventPub,
	}
}

// AddOptionInput holds input for adding an option to a product.
type AddOptionInput struct {
	Name      string
	SortOrder int
	Values    []OptionValueInput
}

// OptionValueInput holds input for an option value.
type OptionValueInput struct {
	Value     string
	ColorHex  string
	SortOrder int
}

// AddOption adds an option group to a product.
func (uc *VariantUseCase) AddOption(ctx context.Context, productID string, sellerID string, input AddOptionInput) (*domain.ProductOption, error) {
	product, err := uc.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}
	if product.SellerID != sellerID {
		return nil, fmt.Errorf("unauthorized: product belongs to another seller")
	}

	if input.Name == "" {
		return nil, fmt.Errorf("option name is required")
	}

	option := &domain.ProductOption{
		ID:        uuid.New().String(),
		ProductID: productID,
		Name:      input.Name,
		SortOrder: input.SortOrder,
	}

	for _, v := range input.Values {
		option.Values = append(option.Values, domain.ProductOptionValue{
			ID:        uuid.New().String(),
			OptionID:  option.ID,
			Value:     v.Value,
			ColorHex:  v.ColorHex,
			SortOrder: v.SortOrder,
		})
	}

	if err := uc.optionRepo.CreateOption(ctx, option); err != nil {
		return nil, fmt.Errorf("failed to create option: %w", err)
	}

	// Mark product as having variants
	if !product.HasVariants {
		product.HasVariants = true
		product.UpdatedAt = time.Now().UTC()
		_ = uc.productRepo.Update(ctx, product)
	}

	return option, nil
}

// RemoveOption removes an option from a product.
func (uc *VariantUseCase) RemoveOption(ctx context.Context, productID string, optionID string, sellerID string) error {
	product, err := uc.productRepo.GetByID(ctx, productID)
	if err != nil {
		return fmt.Errorf("product not found: %w", err)
	}
	if product.SellerID != sellerID {
		return fmt.Errorf("unauthorized: product belongs to another seller")
	}

	if err := uc.optionRepo.DeleteOption(ctx, optionID); err != nil {
		return fmt.Errorf("failed to delete option: %w", err)
	}

	return nil
}

// GenerateVariants creates variants as the cartesian product of all option values.
func (uc *VariantUseCase) GenerateVariants(ctx context.Context, productID string, sellerID string) ([]domain.Variant, error) {
	product, err := uc.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}
	if product.SellerID != sellerID {
		return nil, fmt.Errorf("unauthorized: product belongs to another seller")
	}

	options, err := uc.optionRepo.ListByProduct(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to list options: %w", err)
	}
	if len(options) == 0 {
		return nil, fmt.Errorf("no options defined for this product")
	}

	// Build cartesian product of option values
	combos := cartesianProduct(options)

	now := time.Now().UTC()
	var variants []domain.Variant
	for i, combo := range combos {
		sku := fmt.Sprintf("%s-%d", product.Slug, i+1)
		variantID := uuid.New().String()

		// Build name from combo values
		var names []string
		var optValues []domain.VariantOptionValue
		for _, ov := range combo {
			names = append(names, ov.Value)
			optValues = append(optValues, domain.VariantOptionValue{
				VariantID:     variantID,
				OptionID:      ov.OptionID,
				OptionValueID: ov.ID,
				OptionName:    ov.OptionName,
				Value:         ov.Value,
			})
		}

		variant := domain.Variant{
			ID:           variantID,
			ProductID:    productID,
			SKU:          sku,
			Name:         fmt.Sprintf("%s", joinNames(names)),
			PriceCents:   product.BasePriceCents,
			Stock:        0,
			IsDefault:    i == 0,
			IsActive:     true,
			OptionValues: optValues,
			CreatedAt:    now,
			UpdatedAt:    now,
		}
		variants = append(variants, variant)
	}

	if err := uc.variantRepo.BulkCreate(ctx, variants); err != nil {
		return nil, fmt.Errorf("failed to bulk create variants: %w", err)
	}

	// Mark product as having variants
	product.HasVariants = true
	product.UpdatedAt = now
	_ = uc.productRepo.Update(ctx, product)

	return variants, nil
}

// optionValueWithMeta holds an option value with its parent option name for building combos.
type optionValueWithMeta struct {
	domain.ProductOptionValue
	OptionName string
}

// cartesianProduct computes the cartesian product of all option value sets.
func cartesianProduct(options []domain.ProductOption) [][]optionValueWithMeta {
	if len(options) == 0 {
		return nil
	}

	// Start with first option's values
	var result [][]optionValueWithMeta
	for _, v := range options[0].Values {
		result = append(result, []optionValueWithMeta{{ProductOptionValue: v, OptionName: options[0].Name}})
	}

	// Cross with remaining options
	for i := 1; i < len(options); i++ {
		var newResult [][]optionValueWithMeta
		for _, combo := range result {
			for _, v := range options[i].Values {
				newCombo := make([]optionValueWithMeta, len(combo))
				copy(newCombo, combo)
				newCombo = append(newCombo, optionValueWithMeta{ProductOptionValue: v, OptionName: options[i].Name})
				newResult = append(newResult, newCombo)
			}
		}
		result = newResult
	}

	return result
}

func joinNames(names []string) string {
	result := ""
	for i, n := range names {
		if i > 0 {
			result += " / "
		}
		result += n
	}
	return result
}

// UpdateVariantInput holds input for updating a variant.
type UpdateVariantInput struct {
	Name           *string
	PriceCents     *int64
	CompareAtCents *int64
	CostCents      *int64
	WeightGrams    *int
	IsActive       *bool
	ImageURLs      []string
	Barcode        *string
	LowStockAlert  *int
}

// UpdateVariant updates a variant's details.
func (uc *VariantUseCase) UpdateVariant(ctx context.Context, productID string, variantID string, sellerID string, input UpdateVariantInput) (*domain.Variant, error) {
	product, err := uc.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}
	if product.SellerID != sellerID {
		return nil, fmt.Errorf("unauthorized: product belongs to another seller")
	}

	variant, err := uc.variantRepo.GetByID(ctx, variantID)
	if err != nil {
		return nil, fmt.Errorf("variant not found: %w", err)
	}
	if variant.ProductID != productID {
		return nil, fmt.Errorf("variant does not belong to this product")
	}

	if input.Name != nil {
		variant.Name = *input.Name
	}
	if input.PriceCents != nil {
		variant.PriceCents = *input.PriceCents
	}
	if input.CompareAtCents != nil {
		variant.CompareAtCents = *input.CompareAtCents
	}
	if input.CostCents != nil {
		variant.CostCents = *input.CostCents
	}
	if input.WeightGrams != nil {
		variant.WeightGrams = *input.WeightGrams
	}
	if input.IsActive != nil {
		variant.IsActive = *input.IsActive
	}
	if input.ImageURLs != nil {
		variant.ImageURLs = input.ImageURLs
	}
	if input.Barcode != nil {
		variant.Barcode = *input.Barcode
	}
	if input.LowStockAlert != nil {
		variant.LowStockAlert = *input.LowStockAlert
	}
	variant.UpdatedAt = time.Now().UTC()

	if err := uc.variantRepo.Update(ctx, variant); err != nil {
		return nil, fmt.Errorf("failed to update variant: %w", err)
	}

	return variant, nil
}

// UpdateStock atomically adjusts the stock of a variant.
func (uc *VariantUseCase) UpdateStock(ctx context.Context, productID string, variantID string, sellerID string, delta int) error {
	if sellerID != "" {
		product, err := uc.productRepo.GetByID(ctx, productID)
		if err != nil {
			return fmt.Errorf("product not found: %w", err)
		}
		if product.SellerID != sellerID {
			return fmt.Errorf("unauthorized: product belongs to another seller")
		}
	}

	if err := uc.variantRepo.UpdateStock(ctx, variantID, delta); err != nil {
		return fmt.Errorf("failed to update stock: %w", err)
	}

	// Get updated variant to publish event
	variant, err := uc.variantRepo.GetByID(ctx, variantID)
	if err == nil {
		_ = uc.eventPub.PublishStockUpdated(ctx, variantID, variant.Stock, delta)
	}

	return nil
}

// UpdateStockDirect atomically adjusts stock without seller validation (for gRPC inter-service calls).
func (uc *VariantUseCase) UpdateStockDirect(ctx context.Context, variantID string, delta int) error {
	if err := uc.variantRepo.UpdateStock(ctx, variantID, delta); err != nil {
		return fmt.Errorf("failed to update stock: %w", err)
	}

	variant, err := uc.variantRepo.GetByID(ctx, variantID)
	if err == nil {
		_ = uc.eventPub.PublishStockUpdated(ctx, variantID, variant.Stock, delta)
	}

	return nil
}

// GetVariant retrieves a variant by ID.
func (uc *VariantUseCase) GetVariant(ctx context.Context, variantID string) (*domain.Variant, error) {
	variant, err := uc.variantRepo.GetByID(ctx, variantID)
	if err != nil {
		return nil, fmt.Errorf("variant not found: %w", err)
	}
	return variant, nil
}

// ListVariantsByProduct lists all variants for a product.
func (uc *VariantUseCase) ListVariantsByProduct(ctx context.Context, productID string) ([]domain.Variant, error) {
	return uc.variantRepo.ListByProduct(ctx, productID)
}
