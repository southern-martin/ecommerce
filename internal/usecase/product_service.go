package usecase

import (
	"context"
	"fmt"
	"strings"

	"ecommerce/catalog-service/internal/domain"
	"ecommerce/catalog-service/internal/port"

	"github.com/google/uuid"
)

type ProductService struct {
	products   port.ProductRepository
	categories port.CategoryRepository
	attributes port.AttributeRepository
}

func NewProductService(
	products port.ProductRepository,
	categories port.CategoryRepository,
	attributes port.AttributeRepository,
) *ProductService {
	return &ProductService{
		products:   products,
		categories: categories,
		attributes: attributes,
	}
}

type CreateProductInput struct {
	Name                  string
	Slug                  string
	Description           string
	PrimaryCategoryID     uuid.UUID
	AdditionalCategoryIDs []uuid.UUID
	BasePriceMinor        int64
}

func (s *ProductService) CreateProduct(ctx context.Context, in CreateProductInput) (domain.Product, error) {
	if in.PrimaryCategoryID == uuid.Nil {
		return domain.Product{}, fmt.Errorf("%w: primary category is required", domain.ErrInvalidProduct)
	}
	if _, err := s.categories.GetByID(ctx, in.PrimaryCategoryID); err != nil {
		return domain.Product{}, err
	}
	for _, categoryID := range in.AdditionalCategoryIDs {
		if categoryID == uuid.Nil || categoryID == in.PrimaryCategoryID {
			continue
		}
		if _, err := s.categories.GetByID(ctx, categoryID); err != nil {
			return domain.Product{}, err
		}
	}

	categoryIDs := domain.NormalizeCategoryIDs(in.PrimaryCategoryID, in.AdditionalCategoryIDs)
	product := domain.Product{
		ID:                uuid.New(),
		Name:              strings.TrimSpace(in.Name),
		Slug:              strings.Trim(strings.ToLower(strings.TrimSpace(in.Slug)), "/"),
		Description:       strings.TrimSpace(in.Description),
		PrimaryCategoryID: in.PrimaryCategoryID,
		CategoryIDs:       categoryIDs,
		Status:            domain.ProductStatusDraft,
		BasePriceMinor:    in.BasePriceMinor,
	}
	if err := product.Validate(); err != nil {
		return domain.Product{}, err
	}

	if err := s.products.Create(ctx, product); err != nil {
		return domain.Product{}, err
	}
	if err := s.products.SetCategories(ctx, product.ID, categoryIDs, product.PrimaryCategoryID); err != nil {
		return domain.Product{}, err
	}
	return product, nil
}

type SetProductAttributesInput struct {
	ProductID uuid.UUID
	Values    []domain.ProductAttributeValue
}

func (s *ProductService) SetProductAttributes(ctx context.Context, in SetProductAttributesInput) error {
	if in.ProductID == uuid.Nil {
		return fmt.Errorf("%w: product_id is required", domain.ErrInvalidProduct)
	}
	if _, err := s.products.GetByID(ctx, in.ProductID); err != nil {
		return err
	}

	for _, value := range in.Values {
		if value.ProductID == uuid.Nil {
			value.ProductID = in.ProductID
		}
		if err := value.Validate(); err != nil {
			return err
		}
		attribute, err := s.attributes.GetByID(ctx, value.AttributeID)
		if err != nil {
			return err
		}
		if attribute.IsVariantAxis {
			return fmt.Errorf("%w: attribute %s is a variant axis", domain.ErrInvalidAttributeValue, attribute.Code)
		}
		if value.OptionID != nil {
			ok, err := s.attributes.OptionBelongsToAttribute(ctx, *value.OptionID, value.AttributeID)
			if err != nil {
				return err
			}
			if !ok {
				return fmt.Errorf("%w: option does not belong to attribute", domain.ErrInvalidAttributeValue)
			}
		}
	}

	return s.products.UpsertAttributeValues(ctx, in.ProductID, in.Values)
}

func (s *ProductService) GetProduct(ctx context.Context, productID uuid.UUID) (domain.Product, error) {
	if productID == uuid.Nil {
		return domain.Product{}, fmt.Errorf("%w: product_id is required", domain.ErrInvalidProduct)
	}
	return s.products.GetByID(ctx, productID)
}

func (s *ProductService) ListProductsByCategory(ctx context.Context, categoryID uuid.UUID) ([]domain.Product, error) {
	if categoryID == uuid.Nil {
		return nil, fmt.Errorf("%w: category_id is required", domain.ErrInvalidCategory)
	}
	if _, err := s.categories.GetByID(ctx, categoryID); err != nil {
		return nil, err
	}
	return s.products.ListByCategory(ctx, categoryID)
}

func (s *ProductService) ListProductAttributeValues(
	ctx context.Context,
	productID uuid.UUID,
) ([]domain.ProductAttributeValue, error) {
	if productID == uuid.Nil {
		return nil, fmt.Errorf("%w: product_id is required", domain.ErrInvalidProduct)
	}
	if _, err := s.products.GetByID(ctx, productID); err != nil {
		return nil, err
	}
	return s.products.ListAttributeValues(ctx, productID)
}
