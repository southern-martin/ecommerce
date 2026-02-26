package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/southern-martin/ecommerce/services/product/internal/domain"
)

// AttributeUseCase handles attribute definition business logic.
type AttributeUseCase struct {
	attributeRepo domain.AttributeRepository
	categoryRepo  domain.CategoryRepository
}

// NewAttributeUseCase creates a new AttributeUseCase.
func NewAttributeUseCase(attributeRepo domain.AttributeRepository, categoryRepo domain.CategoryRepository) *AttributeUseCase {
	return &AttributeUseCase{
		attributeRepo: attributeRepo,
		categoryRepo:  categoryRepo,
	}
}

// CreateAttributeInput holds input for creating an attribute definition.
type CreateAttributeInput struct {
	Name       string
	Type       domain.AttributeType
	Required   bool
	Filterable bool
	Options    []string
	Unit       string
	SortOrder  int
}

// CreateAttributeDefinition creates a new attribute definition (admin only).
func (uc *AttributeUseCase) CreateAttributeDefinition(ctx context.Context, input CreateAttributeInput) (*domain.AttributeDefinition, error) {
	if input.Name == "" {
		return nil, fmt.Errorf("attribute name is required")
	}
	if input.Type == "" {
		return nil, fmt.Errorf("attribute type is required")
	}

	attr := &domain.AttributeDefinition{
		ID:         uuid.New().String(),
		Name:       input.Name,
		Slug:       generateSlug(input.Name),
		Type:       input.Type,
		Required:   input.Required,
		Filterable: input.Filterable,
		Options:    input.Options,
		Unit:       input.Unit,
		SortOrder:  input.SortOrder,
		CreatedAt:  time.Now().UTC(),
	}

	if err := uc.attributeRepo.CreateDefinition(ctx, attr); err != nil {
		return nil, fmt.Errorf("failed to create attribute definition: %w", err)
	}

	return attr, nil
}

// ListAttributeDefinitions lists all attribute definitions.
func (uc *AttributeUseCase) ListAttributeDefinitions(ctx context.Context) ([]*domain.AttributeDefinition, error) {
	return uc.attributeRepo.ListDefinitions(ctx)
}

// UpdateAttributeInput holds input for updating an attribute definition.
type UpdateAttributeInput struct {
	Name       *string
	Required   *bool
	Filterable *bool
	Options    []string
	Unit       *string
	SortOrder  *int
}

// UpdateAttributeDefinition updates an attribute definition.
func (uc *AttributeUseCase) UpdateAttributeDefinition(ctx context.Context, id string, input UpdateAttributeInput) (*domain.AttributeDefinition, error) {
	attr, err := uc.attributeRepo.GetDefinitionByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("attribute definition not found: %w", err)
	}

	if input.Name != nil {
		attr.Name = *input.Name
		attr.Slug = generateSlug(*input.Name)
	}
	if input.Required != nil {
		attr.Required = *input.Required
	}
	if input.Filterable != nil {
		attr.Filterable = *input.Filterable
	}
	if input.Options != nil {
		attr.Options = input.Options
	}
	if input.Unit != nil {
		attr.Unit = *input.Unit
	}
	if input.SortOrder != nil {
		attr.SortOrder = *input.SortOrder
	}

	if err := uc.attributeRepo.UpdateDefinition(ctx, attr); err != nil {
		return nil, fmt.Errorf("failed to update attribute definition: %w", err)
	}

	return attr, nil
}

// DeleteAttributeDefinition deletes an attribute definition.
func (uc *AttributeUseCase) DeleteAttributeDefinition(ctx context.Context, id string) error {
	if err := uc.attributeRepo.DeleteDefinition(ctx, id); err != nil {
		return fmt.Errorf("failed to delete attribute definition: %w", err)
	}
	return nil
}

// AssignAttributeToCategory assigns an attribute definition to a category.
func (uc *AttributeUseCase) AssignAttributeToCategory(ctx context.Context, categoryID, attributeID string, sortOrder int) error {
	// Validate category
	if _, err := uc.categoryRepo.GetByID(ctx, categoryID); err != nil {
		return fmt.Errorf("category not found: %w", err)
	}

	// Validate attribute
	if _, err := uc.attributeRepo.GetDefinitionByID(ctx, attributeID); err != nil {
		return fmt.Errorf("attribute not found: %w", err)
	}

	if err := uc.attributeRepo.AssignToCategory(ctx, categoryID, attributeID, sortOrder); err != nil {
		return fmt.Errorf("failed to assign attribute to category: %w", err)
	}
	return nil
}

// RemoveAttributeFromCategory removes an attribute from a category.
func (uc *AttributeUseCase) RemoveAttributeFromCategory(ctx context.Context, categoryID, attributeID string) error {
	if err := uc.attributeRepo.RemoveFromCategory(ctx, categoryID, attributeID); err != nil {
		return fmt.Errorf("failed to remove attribute from category: %w", err)
	}
	return nil
}

// ListCategoryAttributes lists all attribute definitions for a category.
func (uc *AttributeUseCase) ListCategoryAttributes(ctx context.Context, categoryID string) ([]*domain.AttributeDefinition, error) {
	return uc.attributeRepo.ListByCategory(ctx, categoryID)
}

// SetProductAttributeValues sets attribute values on a product.
func (uc *AttributeUseCase) SetProductAttributeValues(ctx context.Context, productID string, values []domain.ProductAttributeValue) error {
	for i := range values {
		if values[i].ID == "" {
			values[i].ID = uuid.New().String()
		}
		values[i].ProductID = productID
	}
	return uc.attributeRepo.SetProductValues(ctx, productID, values)
}

// GetProductAttributeValues retrieves attribute values for a product.
func (uc *AttributeUseCase) GetProductAttributeValues(ctx context.Context, productID string) ([]domain.ProductAttributeValue, error) {
	return uc.attributeRepo.GetProductValues(ctx, productID)
}
