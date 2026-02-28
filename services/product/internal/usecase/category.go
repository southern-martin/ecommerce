package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/southern-martin/ecommerce/services/product/internal/domain"
)

// CategoryUseCase handles category business logic.
type CategoryUseCase struct {
	categoryRepo domain.CategoryRepository
}

// NewCategoryUseCase creates a new CategoryUseCase.
func NewCategoryUseCase(categoryRepo domain.CategoryRepository) *CategoryUseCase {
	return &CategoryUseCase{categoryRepo: categoryRepo}
}

// CreateCategoryInput holds the input for creating a category.
type CreateCategoryInput struct {
	Name      string
	ParentID  string
	SortOrder int
	ImageURL  string
}

// CreateCategory creates a new category (admin only).
func (uc *CategoryUseCase) CreateCategory(ctx context.Context, input CreateCategoryInput) (*domain.Category, error) {
	if input.Name == "" {
		return nil, fmt.Errorf("category name is required")
	}

	// Validate parent exists if provided
	if input.ParentID != "" {
		_, err := uc.categoryRepo.GetByID(ctx, input.ParentID)
		if err != nil {
			return nil, fmt.Errorf("parent category not found: %w", err)
		}
	}

	category := &domain.Category{
		ID:        uuid.New().String(),
		Name:      input.Name,
		Slug:      generateSlug(input.Name),
		ParentID:  input.ParentID,
		SortOrder: input.SortOrder,
		ImageURL:  input.ImageURL,
		IsActive:  true,
		CreatedAt: time.Now().UTC(),
	}

	if err := uc.categoryRepo.Create(ctx, category); err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	return category, nil
}

// UpdateCategoryInput holds the input for updating a category.
type UpdateCategoryInput struct {
	Name      *string
	ParentID  *string
	SortOrder *int
	ImageURL  *string
	IsActive  *bool
}

// UpdateCategory updates an existing category.
func (uc *CategoryUseCase) UpdateCategory(ctx context.Context, id string, input UpdateCategoryInput) (*domain.Category, error) {
	cat, err := uc.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("category not found: %w", err)
	}

	if input.Name != nil {
		cat.Name = *input.Name
		cat.Slug = generateSlug(*input.Name)
	}
	if input.ParentID != nil {
		if *input.ParentID != "" && *input.ParentID != cat.ID {
			if _, err := uc.categoryRepo.GetByID(ctx, *input.ParentID); err != nil {
				return nil, fmt.Errorf("parent category not found: %w", err)
			}
		}
		cat.ParentID = *input.ParentID
	}
	if input.SortOrder != nil {
		cat.SortOrder = *input.SortOrder
	}
	if input.ImageURL != nil {
		cat.ImageURL = *input.ImageURL
	}
	if input.IsActive != nil {
		cat.IsActive = *input.IsActive
	}

	if err := uc.categoryRepo.Update(ctx, cat); err != nil {
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	return cat, nil
}

// DeleteCategory deletes a category by ID.
func (uc *CategoryUseCase) DeleteCategory(ctx context.Context, id string) error {
	if _, err := uc.categoryRepo.GetByID(ctx, id); err != nil {
		return fmt.Errorf("category not found: %w", err)
	}
	return uc.categoryRepo.Delete(ctx, id)
}

// GetCategories lists all categories.
func (uc *CategoryUseCase) GetCategories(ctx context.Context) ([]*domain.Category, error) {
	return uc.categoryRepo.List(ctx)
}
