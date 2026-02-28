package domain

import "context"

// ProductFilter defines filtering and pagination for product listing.
type ProductFilter struct {
	SellerID   string
	CategoryID string
	Status     string
	Query      string
	MinPrice   int64
	MaxPrice   int64
	SortBy     string
	Page       int
	PageSize   int
}

// ProductRepository defines persistence operations for products.
type ProductRepository interface {
	Create(ctx context.Context, p *Product) error
	GetByID(ctx context.Context, id string) (*Product, error)
	GetBySlug(ctx context.Context, slug string) (*Product, error)
	List(ctx context.Context, filter ProductFilter) ([]*Product, int64, error)
	Update(ctx context.Context, p *Product) error
	Delete(ctx context.Context, id string) error
}

// CategoryRepository defines persistence operations for categories.
type CategoryRepository interface {
	Create(ctx context.Context, c *Category) error
	GetByID(ctx context.Context, id string) (*Category, error)
	List(ctx context.Context) ([]*Category, error)
	Update(ctx context.Context, c *Category) error
	Delete(ctx context.Context, id string) error
}

// AttributeRepository defines persistence operations for attribute definitions and values.
type AttributeRepository interface {
	CreateDefinition(ctx context.Context, attr *AttributeDefinition) error
	GetDefinitionByID(ctx context.Context, id string) (*AttributeDefinition, error)
	ListDefinitions(ctx context.Context) ([]*AttributeDefinition, error)
	UpdateDefinition(ctx context.Context, attr *AttributeDefinition) error
	DeleteDefinition(ctx context.Context, id string) error
	AssignToCategory(ctx context.Context, categoryID, attributeID string, sortOrder int) error
	RemoveFromCategory(ctx context.Context, categoryID, attributeID string) error
	ListByCategory(ctx context.Context, categoryID string) ([]*AttributeDefinition, error)
	SetProductValues(ctx context.Context, productID string, values []ProductAttributeValue) error
	GetProductValues(ctx context.Context, productID string) ([]ProductAttributeValue, error)
	// Attribute option values
	SetOptionValues(ctx context.Context, attributeID string, opts []AttributeOptionValue) error
	ListOptionValues(ctx context.Context, attributeID string) ([]AttributeOptionValue, error)
}

// AttributeGroupRepository defines persistence operations for attribute groups.
type AttributeGroupRepository interface {
	Create(ctx context.Context, group *AttributeGroup) error
	GetByID(ctx context.Context, id string) (*AttributeGroup, error)
	List(ctx context.Context) ([]*AttributeGroup, error)
	Update(ctx context.Context, group *AttributeGroup) error
	Delete(ctx context.Context, id string) error
	AddAttribute(ctx context.Context, groupID, attributeID string, sortOrder int) error
	RemoveAttribute(ctx context.Context, groupID, attributeID string) error
	ListAttributes(ctx context.Context, groupID string) ([]*AttributeDefinition, error)
}

// OptionRepository defines persistence operations for product options.
type OptionRepository interface {
	CreateOption(ctx context.Context, option *ProductOption) error
	UpdateOption(ctx context.Context, option *ProductOption) error
	DeleteOption(ctx context.Context, optionID string) error
	ListByProduct(ctx context.Context, productID string) ([]ProductOption, error)
	CreateOptionValue(ctx context.Context, value *ProductOptionValue) error
	UpdateOptionValue(ctx context.Context, value *ProductOptionValue) error
	DeleteOptionValue(ctx context.Context, valueID string) error
}

// VariantRepository defines persistence operations for product variants.
type VariantRepository interface {
	Create(ctx context.Context, v *Variant) error
	GetByID(ctx context.Context, id string) (*Variant, error)
	GetBySKU(ctx context.Context, sku string) (*Variant, error)
	ListByProduct(ctx context.Context, productID string) ([]Variant, error)
	Update(ctx context.Context, v *Variant) error
	Delete(ctx context.Context, id string) error
	BulkCreate(ctx context.Context, variants []Variant) error
	UpdateStock(ctx context.Context, variantID string, delta int) error
	SetOptionValues(ctx context.Context, variantID string, values []VariantOptionValue) error
}
