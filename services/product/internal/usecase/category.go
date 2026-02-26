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

// GetCategories lists all categories.
func (uc *CategoryUseCase) GetCategories(ctx context.Context) ([]*domain.Category, error) {
	return uc.categoryRepo.List(ctx)
}
