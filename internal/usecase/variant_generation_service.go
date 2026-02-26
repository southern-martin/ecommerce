package usecase

import (
	"context"
	"fmt"

	"ecommerce/catalog-service/internal/domain"
	"ecommerce/catalog-service/internal/port"

	"github.com/google/uuid"
)

const defaultMaxGeneratedVariants = 500

type VariantGenerationService struct {
	products    port.ProductRepository
	attributes  port.AttributeRepository
	maxVariants int
}

func NewVariantGenerationService(products port.ProductRepository, attributes port.AttributeRepository) *VariantGenerationService {
	return &VariantGenerationService{
		products:    products,
		attributes:  attributes,
		maxVariants: defaultMaxGeneratedVariants,
	}
}

type VariantAxisInput struct {
	AttributeID uuid.UUID
	OptionIDs   []uuid.UUID
}

type GenerateVariantsInput struct {
	ProductID       uuid.UUID
	Axes            []VariantAxisInput
	BasePriceMinor  int64
	InitialStockQty int64
}

func (s *VariantGenerationService) SetMaxVariants(limit int) {
	if limit > 0 {
		s.maxVariants = limit
	}
}

func (s *VariantGenerationService) GenerateMatrix(in GenerateVariantsInput) ([]domain.ProductVariant, error) {
	if in.ProductID == uuid.Nil {
		return nil, fmt.Errorf("%w: product_id is required", domain.ErrInvalidVariantAxis)
	}
	if len(in.Axes) == 0 {
		return nil, fmt.Errorf("%w: at least one axis is required", domain.ErrInvalidVariantAxis)
	}

	seenAxis := map[uuid.UUID]struct{}{}
	total := 1
	for _, axis := range in.Axes {
		if axis.AttributeID == uuid.Nil || len(axis.OptionIDs) == 0 {
			return nil, fmt.Errorf("%w: axis must have attribute_id and options", domain.ErrInvalidVariantAxis)
		}
		if _, ok := seenAxis[axis.AttributeID]; ok {
			return nil, fmt.Errorf("%w: duplicate axis attribute", domain.ErrInvalidVariantAxis)
		}
		seenAxis[axis.AttributeID] = struct{}{}
		total *= len(axis.OptionIDs)
		if total > s.maxVariants {
			return nil, fmt.Errorf("%w: matrix too large (%d > %d)", domain.ErrInvalidVariantAxis, total, s.maxVariants)
		}
	}

	combinationBuffer := make([]domain.VariantOptionValue, len(in.Axes))
	var generated []domain.ProductVariant
	var walk func(depth int) error
	walk = func(depth int) error {
		if depth == len(in.Axes) {
			options := append([]domain.VariantOptionValue(nil), combinationBuffer...)
			key, err := domain.BuildCombinationKey(options)
			if err != nil {
				return err
			}
			variant := domain.ProductVariant{
				ID:             uuid.New(),
				ProductID:      in.ProductID,
				PriceMinor:     in.BasePriceMinor,
				StockQty:       in.InitialStockQty,
				Status:         domain.ProductVariantStatusActive,
				CombinationKey: key,
				Options:        options,
			}
			if err := variant.Validate(); err != nil {
				return err
			}
			generated = append(generated, variant)
			return nil
		}

		axis := in.Axes[depth]
		for _, optionID := range axis.OptionIDs {
			if optionID == uuid.Nil {
				return fmt.Errorf("%w: option_id is required", domain.ErrInvalidVariantAxis)
			}
			combinationBuffer[depth] = domain.VariantOptionValue{
				AttributeID: axis.AttributeID,
				OptionID:    optionID,
			}
			if err := walk(depth + 1); err != nil {
				return err
			}
		}
		return nil
	}

	if err := walk(0); err != nil {
		return nil, err
	}
	return generated, nil
}

func (s *VariantGenerationService) GenerateAndPersist(
	ctx context.Context,
	in GenerateVariantsInput,
) ([]domain.ProductVariant, error) {
	if err := s.validateBusinessRules(ctx, in); err != nil {
		return nil, err
	}
	variants, err := s.GenerateMatrix(in)
	if err != nil {
		return nil, err
	}
	if s.products == nil {
		return variants, nil
	}
	if err := s.products.CreateVariants(ctx, variants); err != nil {
		return nil, err
	}
	return variants, nil
}

func (s *VariantGenerationService) validateBusinessRules(ctx context.Context, in GenerateVariantsInput) error {
	if s.products == nil || s.attributes == nil {
		return nil
	}

	product, err := s.products.GetByID(ctx, in.ProductID)
	if err != nil {
		return err
	}
	for _, axis := range in.Axes {
		attribute, err := s.attributes.GetByID(ctx, axis.AttributeID)
		if err != nil {
			return err
		}
		if !attribute.IsVariantAxis {
			return fmt.Errorf("%w: attribute %s is not a variant axis", domain.ErrInvalidVariantAxis, attribute.Code)
		}
		if attribute.CategoryID != product.PrimaryCategoryID {
			return fmt.Errorf("%w: attribute %s does not belong to product primary category", domain.ErrInvalidVariantAxis, attribute.Code)
		}
		for _, optionID := range axis.OptionIDs {
			ok, err := s.attributes.OptionBelongsToAttribute(ctx, optionID, axis.AttributeID)
			if err != nil {
				return err
			}
			if !ok {
				return fmt.Errorf("%w: option does not belong to attribute", domain.ErrInvalidVariantAxis)
			}
		}
	}
	return nil
}
