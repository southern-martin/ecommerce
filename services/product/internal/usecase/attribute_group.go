package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/southern-martin/ecommerce/services/product/internal/domain"
)

// AttributeGroupUseCase handles attribute group business logic.
type AttributeGroupUseCase struct {
	groupRepo     domain.AttributeGroupRepository
	attributeRepo domain.AttributeRepository
}

// NewAttributeGroupUseCase creates a new AttributeGroupUseCase.
func NewAttributeGroupUseCase(groupRepo domain.AttributeGroupRepository, attributeRepo domain.AttributeRepository) *AttributeGroupUseCase {
	return &AttributeGroupUseCase{
		groupRepo:     groupRepo,
		attributeRepo: attributeRepo,
	}
}

// CreateAttributeGroupInput holds input for creating an attribute group.
type CreateAttributeGroupInput struct {
	Name        string
	Description string
	SortOrder   int
}

// CreateAttributeGroup creates a new attribute group.
func (uc *AttributeGroupUseCase) CreateAttributeGroup(ctx context.Context, input CreateAttributeGroupInput) (*domain.AttributeGroup, error) {
	if input.Name == "" {
		return nil, fmt.Errorf("attribute group name is required")
	}

	now := time.Now().UTC()
	group := &domain.AttributeGroup{
		ID:          uuid.New().String(),
		Name:        input.Name,
		Slug:        generateSlug(input.Name),
		Description: input.Description,
		SortOrder:   input.SortOrder,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := uc.groupRepo.Create(ctx, group); err != nil {
		return nil, fmt.Errorf("failed to create attribute group: %w", err)
	}

	return group, nil
}

// ListAttributeGroups lists all attribute groups.
func (uc *AttributeGroupUseCase) ListAttributeGroups(ctx context.Context) ([]*domain.AttributeGroup, error) {
	groups, err := uc.groupRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	// Enrich each group with its attribute count by loading attributes
	for _, g := range groups {
		attrs, err := uc.groupRepo.ListAttributes(ctx, g.ID)
		if err == nil {
			g.Attributes = make([]domain.AttributeDefinition, len(attrs))
			for i, a := range attrs {
				g.Attributes[i] = *a
			}
		}
	}

	return groups, nil
}

// GetAttributeGroup retrieves an attribute group by ID.
func (uc *AttributeGroupUseCase) GetAttributeGroup(ctx context.Context, id string) (*domain.AttributeGroup, error) {
	group, err := uc.groupRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("attribute group not found: %w", err)
	}

	// Enrich with attributes
	attrs, err := uc.groupRepo.ListAttributes(ctx, group.ID)
	if err == nil {
		group.Attributes = make([]domain.AttributeDefinition, len(attrs))
		for i, a := range attrs {
			group.Attributes[i] = *a
		}
	}

	return group, nil
}

// UpdateAttributeGroupInput holds input for updating an attribute group.
type UpdateAttributeGroupInput struct {
	Name        *string
	Description *string
	SortOrder   *int
}

// UpdateAttributeGroup updates an attribute group.
func (uc *AttributeGroupUseCase) UpdateAttributeGroup(ctx context.Context, id string, input UpdateAttributeGroupInput) (*domain.AttributeGroup, error) {
	group, err := uc.groupRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("attribute group not found: %w", err)
	}

	if input.Name != nil {
		group.Name = *input.Name
		group.Slug = generateSlug(*input.Name)
	}
	if input.Description != nil {
		group.Description = *input.Description
	}
	if input.SortOrder != nil {
		group.SortOrder = *input.SortOrder
	}
	group.UpdatedAt = time.Now().UTC()

	if err := uc.groupRepo.Update(ctx, group); err != nil {
		return nil, fmt.Errorf("failed to update attribute group: %w", err)
	}

	return group, nil
}

// DeleteAttributeGroup deletes an attribute group.
func (uc *AttributeGroupUseCase) DeleteAttributeGroup(ctx context.Context, id string) error {
	if err := uc.groupRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete attribute group: %w", err)
	}
	return nil
}

// AddAttributeToGroup assigns an attribute to a group.
func (uc *AttributeGroupUseCase) AddAttributeToGroup(ctx context.Context, groupID, attributeID string, sortOrder int) error {
	// Validate group exists
	if _, err := uc.groupRepo.GetByID(ctx, groupID); err != nil {
		return fmt.Errorf("attribute group not found: %w", err)
	}

	// Validate attribute exists
	if _, err := uc.attributeRepo.GetDefinitionByID(ctx, attributeID); err != nil {
		return fmt.Errorf("attribute not found: %w", err)
	}

	if err := uc.groupRepo.AddAttribute(ctx, groupID, attributeID, sortOrder); err != nil {
		return fmt.Errorf("failed to add attribute to group: %w", err)
	}
	return nil
}

// RemoveAttributeFromGroup removes an attribute from a group.
func (uc *AttributeGroupUseCase) RemoveAttributeFromGroup(ctx context.Context, groupID, attributeID string) error {
	if err := uc.groupRepo.RemoveAttribute(ctx, groupID, attributeID); err != nil {
		return fmt.Errorf("failed to remove attribute from group: %w", err)
	}
	return nil
}

// ListGroupAttributes lists all attributes in a group.
func (uc *AttributeGroupUseCase) ListGroupAttributes(ctx context.Context, groupID string) ([]*domain.AttributeDefinition, error) {
	return uc.groupRepo.ListAttributes(ctx, groupID)
}
