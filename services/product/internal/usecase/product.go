package usecase

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/southern-martin/ecommerce/services/product/internal/domain"
)

// ProductUseCase handles product business logic.
type ProductUseCase struct {
	productRepo   domain.ProductRepository
	categoryRepo  domain.CategoryRepository
	attributeRepo domain.AttributeRepository
	optionRepo    domain.OptionRepository
	variantRepo   domain.VariantRepository
	eventPub      domain.EventPublisher
}

// NewProductUseCase creates a new ProductUseCase.
func NewProductUseCase(
	productRepo domain.ProductRepository,
	categoryRepo domain.CategoryRepository,
	attributeRepo domain.AttributeRepository,
	optionRepo domain.OptionRepository,
	variantRepo domain.VariantRepository,
	eventPub domain.EventPublisher,
) *ProductUseCase {
	return &ProductUseCase{
		productRepo:   productRepo,
		categoryRepo:  categoryRepo,
		attributeRepo: attributeRepo,
		optionRepo:    optionRepo,
		variantRepo:   variantRepo,
		eventPub:      eventPub,
	}
}

// CreateProductInput holds the input for creating a product.
type CreateProductInput struct {
	SellerID       string
	CategoryID     string
	Name           string
	Description    string
	BasePriceCents int64
	Currency       string
	Tags           []string
	ImageURLs      []string
	Attributes     []AttributeValueInput
}

// AttributeValueInput holds attribute values for product creation.
type AttributeValueInput struct {
	AttributeID string
	Value       string
	Values      []string
}

// CreateProduct creates a new product.
func (uc *ProductUseCase) CreateProduct(ctx context.Context, input CreateProductInput) (*domain.Product, error) {
	if input.Name == "" {
		return nil, fmt.Errorf("product name is required")
	}
	if input.SellerID == "" {
		return nil, fmt.Errorf("seller ID is required")
	}
	if input.BasePriceCents < 0 {
		return nil, fmt.Errorf("base price must be non-negative")
	}

	// Validate category exists
	if input.CategoryID != "" {
		_, err := uc.categoryRepo.GetByID(ctx, input.CategoryID)
		if err != nil {
			return nil, fmt.Errorf("category not found: %w", err)
		}
	}

	now := time.Now().UTC()
	product := &domain.Product{
		ID:             uuid.New().String(),
		SellerID:       input.SellerID,
		CategoryID:     input.CategoryID,
		Name:           input.Name,
		Slug:           generateSlug(input.Name),
		Description:    input.Description,
		BasePriceCents: input.BasePriceCents,
		Currency:       input.Currency,
		Status:         domain.ProductStatusDraft,
		HasVariants:    false,
		Tags:           input.Tags,
		ImageURLs:      input.ImageURLs,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if product.Currency == "" {
		product.Currency = "USD"
	}

	if err := uc.productRepo.Create(ctx, product); err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	// Save attribute values
	if len(input.Attributes) > 0 {
		var values []domain.ProductAttributeValue
		for _, av := range input.Attributes {
			values = append(values, domain.ProductAttributeValue{
				ID:          uuid.New().String(),
				ProductID:   product.ID,
				AttributeID: av.AttributeID,
				Value:       av.Value,
				Values:      av.Values,
			})
		}
		if err := uc.attributeRepo.SetProductValues(ctx, product.ID, values); err != nil {
			return nil, fmt.Errorf("failed to set attribute values: %w", err)
		}
	}

	_ = uc.eventPub.PublishProductCreated(ctx, product)

	return product, nil
}

// GetProduct retrieves a product by ID including options and variants.
func (uc *ProductUseCase) GetProduct(ctx context.Context, id string) (*domain.Product, error) {
	product, err := uc.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}
	return product, nil
}

// GetProductBySlug retrieves a product by slug.
func (uc *ProductUseCase) GetProductBySlug(ctx context.Context, slug string) (*domain.Product, error) {
	product, err := uc.productRepo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}
	return product, nil
}

// ListProducts lists products with filtering and pagination.
func (uc *ProductUseCase) ListProducts(ctx context.Context, filter domain.ProductFilter) ([]*domain.Product, int64, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}
	if filter.PageSize > 100 {
		filter.PageSize = 100
	}
	return uc.productRepo.List(ctx, filter)
}

// UpdateProductInput holds the input for updating a product.
type UpdateProductInput struct {
	Name           *string
	Description    *string
	BasePriceCents *int64
	Currency       *string
	Status         *domain.ProductStatus
	Tags           []string
	ImageURLs      []string
	CategoryID     *string
}

// UpdateProduct updates an existing product.
func (uc *ProductUseCase) UpdateProduct(ctx context.Context, id string, sellerID string, input UpdateProductInput) (*domain.Product, error) {
	product, err := uc.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	if product.SellerID != sellerID {
		return nil, fmt.Errorf("unauthorized: product belongs to another seller")
	}

	if input.Name != nil {
		product.Name = *input.Name
		product.Slug = generateSlug(*input.Name)
	}
	if input.Description != nil {
		product.Description = *input.Description
	}
	if input.BasePriceCents != nil {
		product.BasePriceCents = *input.BasePriceCents
	}
	if input.Currency != nil {
		product.Currency = *input.Currency
	}
	if input.Status != nil {
		product.Status = *input.Status
	}
	if input.Tags != nil {
		product.Tags = input.Tags
	}
	if input.ImageURLs != nil {
		product.ImageURLs = input.ImageURLs
	}
	if input.CategoryID != nil {
		product.CategoryID = *input.CategoryID
	}
	product.UpdatedAt = time.Now().UTC()

	if err := uc.productRepo.Update(ctx, product); err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	_ = uc.eventPub.PublishProductUpdated(ctx, product)

	return product, nil
}

// DeleteProduct soft-deletes a product.
func (uc *ProductUseCase) DeleteProduct(ctx context.Context, id string, sellerID string) error {
	product, err := uc.productRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("product not found: %w", err)
	}

	if product.SellerID != sellerID {
		return fmt.Errorf("unauthorized: product belongs to another seller")
	}

	if err := uc.productRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	_ = uc.eventPub.PublishProductDeleted(ctx, id)

	return nil
}

// generateSlug creates a URL-friendly slug from a name with a short UUID suffix.
func generateSlug(name string) string {
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	reg := regexp.MustCompile(`[^a-z0-9\-]`)
	slug = reg.ReplaceAllString(slug, "")
	// Remove consecutive hyphens
	reg = regexp.MustCompile(`-+`)
	slug = reg.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")

	// Append short UUID suffix for uniqueness
	suffix := uuid.New().String()[:8]
	return fmt.Sprintf("%s-%s", slug, suffix)
}
