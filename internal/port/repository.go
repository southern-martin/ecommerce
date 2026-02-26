package port

import (
	"context"

	"ecommerce/catalog-service/internal/domain"

	"github.com/google/uuid"
)

type CategoryRepository interface {
	Create(ctx context.Context, category domain.Category) error
	GetByID(ctx context.Context, id uuid.UUID) (domain.Category, error)
	ExistsByParentAndSlug(ctx context.Context, parentID *uuid.UUID, slug string) (bool, error)
	ListChildren(ctx context.Context, parentID *uuid.UUID) ([]domain.Category, error)
	ListAll(ctx context.Context) ([]domain.Category, error)
}

type AttributeRepository interface {
	Create(ctx context.Context, attribute domain.CategoryAttribute) error
	GetByID(ctx context.Context, id uuid.UUID) (domain.CategoryAttribute, error)
	GetByCategory(ctx context.Context, categoryID uuid.UUID) ([]domain.CategoryAttribute, error)
	ExistsByCategoryAndCode(ctx context.Context, categoryID uuid.UUID, code string) (bool, error)
	CreateOption(ctx context.Context, option domain.AttributeOption) error
	ListOptionsByAttribute(ctx context.Context, attributeID uuid.UUID) ([]domain.AttributeOption, error)
	OptionBelongsToAttribute(ctx context.Context, optionID, attributeID uuid.UUID) (bool, error)
}

type ProductRepository interface {
	Create(ctx context.Context, product domain.Product) error
	GetByID(ctx context.Context, id uuid.UUID) (domain.Product, error)
	SetCategories(ctx context.Context, productID uuid.UUID, categoryIDs []uuid.UUID, primaryCategoryID uuid.UUID) error
	UpsertAttributeValues(ctx context.Context, productID uuid.UUID, values []domain.ProductAttributeValue) error
	ListAttributeValues(ctx context.Context, productID uuid.UUID) ([]domain.ProductAttributeValue, error)
	ListByCategory(ctx context.Context, categoryID uuid.UUID) ([]domain.Product, error)
	CreateVariants(ctx context.Context, variants []domain.ProductVariant) error
	ListVariantsByProduct(ctx context.Context, productID uuid.UUID) ([]domain.ProductVariant, error)
}
