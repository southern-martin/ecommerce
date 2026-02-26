package usecase

import (
	"context"
	"errors"
	"testing"

	"ecommerce/catalog-service/internal/domain"
	"ecommerce/catalog-service/internal/infra/memory"
)

func TestAttributeServiceCreateAndList(t *testing.T) {
	ctx := context.Background()
	store := memory.NewStore()
	categoryRepo := memory.NewCategoryRepo(store)
	attributeRepo := memory.NewAttributeRepo(store)

	categoryService := NewCategoryService(categoryRepo)
	attributeService := NewAttributeService(attributeRepo, categoryRepo)

	category, err := categoryService.CreateCategory(ctx, CreateCategoryInput{
		Name:     "Shoes",
		Slug:     "shoes",
		IsActive: true,
	})
	if err != nil {
		t.Fatalf("CreateCategory() error = %v", err)
	}

	attribute, err := attributeService.CreateCategoryAttribute(ctx, CreateCategoryAttributeInput{
		CategoryID:    category.ID,
		Name:          "Size",
		Code:          "size",
		Type:          "select",
		IsVariantAxis: true,
	})
	if err != nil {
		t.Fatalf("CreateCategoryAttribute() error = %v", err)
	}

	_, err = attributeService.AddAttributeOption(ctx, AddAttributeOptionInput{
		AttributeID: attribute.ID,
		Value:       "42",
		Label:       "EU 42",
	})
	if err != nil {
		t.Fatalf("AddAttributeOption() error = %v", err)
	}

	attributes, err := attributeService.ListCategoryAttributes(ctx, category.ID)
	if err != nil {
		t.Fatalf("ListCategoryAttributes() error = %v", err)
	}
	if len(attributes) != 1 {
		t.Fatalf("ListCategoryAttributes() len = %d, want 1", len(attributes))
	}

	options, err := attributeService.ListAttributeOptions(ctx, attribute.ID)
	if err != nil {
		t.Fatalf("ListAttributeOptions() error = %v", err)
	}
	if len(options) != 1 {
		t.Fatalf("ListAttributeOptions() len = %d, want 1", len(options))
	}
}

func TestAttributeServiceRejectsDuplicateCode(t *testing.T) {
	ctx := context.Background()
	store := memory.NewStore()
	categoryRepo := memory.NewCategoryRepo(store)
	attributeRepo := memory.NewAttributeRepo(store)

	categoryService := NewCategoryService(categoryRepo)
	attributeService := NewAttributeService(attributeRepo, categoryRepo)

	category, err := categoryService.CreateCategory(ctx, CreateCategoryInput{
		Name:     "Clothes",
		Slug:     "clothes",
		IsActive: true,
	})
	if err != nil {
		t.Fatalf("CreateCategory() error = %v", err)
	}

	_, err = attributeService.CreateCategoryAttribute(ctx, CreateCategoryAttributeInput{
		CategoryID: category.ID,
		Name:       "Color",
		Code:       "color",
		Type:       "select",
	})
	if err != nil {
		t.Fatalf("CreateCategoryAttribute() error = %v", err)
	}

	_, err = attributeService.CreateCategoryAttribute(ctx, CreateCategoryAttributeInput{
		CategoryID: category.ID,
		Name:       "Color 2",
		Code:       "color",
		Type:       "select",
	})
	if err == nil {
		t.Fatal("CreateCategoryAttribute() error = nil, want duplicate code error")
	}
}

func TestAttributeServiceRejectsOptionForNonSelectAttribute(t *testing.T) {
	ctx := context.Background()
	store := memory.NewStore()
	categoryRepo := memory.NewCategoryRepo(store)
	attributeRepo := memory.NewAttributeRepo(store)

	categoryService := NewCategoryService(categoryRepo)
	attributeService := NewAttributeService(attributeRepo, categoryRepo)

	category, err := categoryService.CreateCategory(ctx, CreateCategoryInput{
		Name:     "Electronics",
		Slug:     "electronics",
		IsActive: true,
	})
	if err != nil {
		t.Fatalf("CreateCategory() error = %v", err)
	}

	attribute, err := attributeService.CreateCategoryAttribute(ctx, CreateCategoryAttributeInput{
		CategoryID: category.ID,
		Name:       "Warranty",
		Code:       "warranty",
		Type:       "text",
	})
	if err != nil {
		t.Fatalf("CreateCategoryAttribute() error = %v", err)
	}

	_, err = attributeService.AddAttributeOption(ctx, AddAttributeOptionInput{
		AttributeID: attribute.ID,
		Value:       "1-year",
		Label:       "1 year",
	})
	if err == nil {
		t.Fatal("AddAttributeOption() error = nil, want non-select rejection")
	}
	if !errors.Is(err, domain.ErrInvalidAttribute) {
		t.Fatalf("AddAttributeOption() error = %v, want ErrInvalidAttribute", err)
	}
}
