package postgres

import (
	"time"

	"github.com/lib/pq"

	"github.com/southern-martin/ecommerce/services/product/internal/domain"
)

// ProductModel is the GORM model for the products table.
type ProductModel struct {
	ID             string         `gorm:"type:uuid;primaryKey"`
	SellerID       string         `gorm:"type:uuid;not null;index"`
	CategoryID     *string        `gorm:"type:uuid;index"`
	Name           string         `gorm:"type:varchar(500);not null"`
	Slug           string         `gorm:"type:varchar(600);uniqueIndex;not null"`
	Description    string         `gorm:"type:text"`
	BasePriceCents int64          `gorm:"not null;default:0"`
	Currency       string         `gorm:"type:varchar(3);not null;default:'USD'"`
	Status         string         `gorm:"type:varchar(20);not null;default:'draft';index"`
	HasVariants    bool           `gorm:"not null;default:false"`
	Tags           pq.StringArray `gorm:"type:text[]"`
	ImageURLs      pq.StringArray `gorm:"type:text[];column:image_urls"`
	RatingAvg      float64        `gorm:"not null;default:0"`
	RatingCount    int            `gorm:"not null;default:0"`
	CreatedAt      time.Time      `gorm:"not null"`
	UpdatedAt      time.Time      `gorm:"not null"`
}

func (ProductModel) TableName() string { return "products" }

func (m *ProductModel) ToDomain() *domain.Product {
	catID := ""
	if m.CategoryID != nil {
		catID = *m.CategoryID
	}
	return &domain.Product{
		ID:             m.ID,
		SellerID:       m.SellerID,
		CategoryID:     catID,
		Name:           m.Name,
		Slug:           m.Slug,
		Description:    m.Description,
		BasePriceCents: m.BasePriceCents,
		Currency:       m.Currency,
		Status:         domain.ProductStatus(m.Status),
		HasVariants:    m.HasVariants,
		Tags:           m.Tags,
		ImageURLs:      m.ImageURLs,
		RatingAvg:      m.RatingAvg,
		RatingCount:    m.RatingCount,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
}

func ProductModelFromDomain(p *domain.Product) *ProductModel {
	var catID *string
	if p.CategoryID != "" {
		catID = &p.CategoryID
	}
	return &ProductModel{
		ID:             p.ID,
		SellerID:       p.SellerID,
		CategoryID:     catID,
		Name:           p.Name,
		Slug:           p.Slug,
		Description:    p.Description,
		BasePriceCents: p.BasePriceCents,
		Currency:       p.Currency,
		Status:         string(p.Status),
		HasVariants:    p.HasVariants,
		Tags:           p.Tags,
		ImageURLs:      p.ImageURLs,
		RatingAvg:      p.RatingAvg,
		RatingCount:    p.RatingCount,
		CreatedAt:      p.CreatedAt,
		UpdatedAt:      p.UpdatedAt,
	}
}

// CategoryModel is the GORM model for the categories table.
type CategoryModel struct {
	ID        string    `gorm:"type:uuid;primaryKey"`
	Name      string    `gorm:"type:varchar(255);not null"`
	Slug      string    `gorm:"type:varchar(300);uniqueIndex;not null"`
	ParentID  *string   `gorm:"type:uuid;index"`
	SortOrder int       `gorm:"not null;default:0"`
	ImageURL  string    `gorm:"type:text;column:image_url"`
	IsActive  bool      `gorm:"not null;default:true"`
	CreatedAt time.Time `gorm:"not null"`
}

func (CategoryModel) TableName() string { return "categories" }

func (m *CategoryModel) ToDomain() *domain.Category {
	parentID := ""
	if m.ParentID != nil {
		parentID = *m.ParentID
	}
	return &domain.Category{
		ID:        m.ID,
		Name:      m.Name,
		Slug:      m.Slug,
		ParentID:  parentID,
		SortOrder: m.SortOrder,
		ImageURL:  m.ImageURL,
		IsActive:  m.IsActive,
		CreatedAt: m.CreatedAt,
	}
}

func CategoryModelFromDomain(c *domain.Category) *CategoryModel {
	var parentID *string
	if c.ParentID != "" {
		parentID = &c.ParentID
	}
	return &CategoryModel{
		ID:        c.ID,
		Name:      c.Name,
		Slug:      c.Slug,
		ParentID:  parentID,
		SortOrder: c.SortOrder,
		ImageURL:  c.ImageURL,
		IsActive:  c.IsActive,
		CreatedAt: c.CreatedAt,
	}
}

// AttributeDefinitionModel is the GORM model for attribute_definitions table.
type AttributeDefinitionModel struct {
	ID         string         `gorm:"type:uuid;primaryKey"`
	Name       string         `gorm:"type:varchar(255);not null"`
	Slug       string         `gorm:"type:varchar(300);uniqueIndex;not null"`
	Type       string         `gorm:"type:varchar(20);not null"`
	Required   bool           `gorm:"not null;default:false"`
	Filterable bool           `gorm:"not null;default:false"`
	Options    pq.StringArray `gorm:"type:text[]"`
	Unit       string         `gorm:"type:varchar(50)"`
	SortOrder  int            `gorm:"not null;default:0"`
	CreatedAt  time.Time      `gorm:"not null"`
}

func (AttributeDefinitionModel) TableName() string { return "attribute_definitions" }

func (m *AttributeDefinitionModel) ToDomain() *domain.AttributeDefinition {
	return &domain.AttributeDefinition{
		ID:         m.ID,
		Name:       m.Name,
		Slug:       m.Slug,
		Type:       domain.AttributeType(m.Type),
		Required:   m.Required,
		Filterable: m.Filterable,
		Options:    m.Options,
		Unit:       m.Unit,
		SortOrder:  m.SortOrder,
		CreatedAt:  m.CreatedAt,
	}
}

func AttributeDefinitionModelFromDomain(a *domain.AttributeDefinition) *AttributeDefinitionModel {
	return &AttributeDefinitionModel{
		ID:         a.ID,
		Name:       a.Name,
		Slug:       a.Slug,
		Type:       string(a.Type),
		Required:   a.Required,
		Filterable: a.Filterable,
		Options:    a.Options,
		Unit:       a.Unit,
		SortOrder:  a.SortOrder,
		CreatedAt:  a.CreatedAt,
	}
}

// CategoryAttributeModel is the GORM model for the category_attributes join table.
type CategoryAttributeModel struct {
	CategoryID  string `gorm:"type:uuid;primaryKey"`
	AttributeID string `gorm:"type:uuid;primaryKey"`
	SortOrder   int    `gorm:"not null;default:0"`
}

func (CategoryAttributeModel) TableName() string { return "category_attributes" }

// ProductAttributeValueModel is the GORM model for product_attribute_values table.
type ProductAttributeValueModel struct {
	ID            string         `gorm:"type:uuid;primaryKey"`
	ProductID     string         `gorm:"type:uuid;not null;index"`
	AttributeID   string         `gorm:"type:uuid;not null"`
	AttributeName string         `gorm:"type:varchar(255)"`
	Value         string         `gorm:"type:text"`
	Values        pq.StringArray `gorm:"type:text[]"`
}

func (ProductAttributeValueModel) TableName() string { return "product_attribute_values" }

func (m *ProductAttributeValueModel) ToDomain() domain.ProductAttributeValue {
	return domain.ProductAttributeValue{
		ID:            m.ID,
		ProductID:     m.ProductID,
		AttributeID:   m.AttributeID,
		AttributeName: m.AttributeName,
		Value:         m.Value,
		Values:        m.Values,
	}
}

func ProductAttributeValueModelFromDomain(v domain.ProductAttributeValue) *ProductAttributeValueModel {
	return &ProductAttributeValueModel{
		ID:            v.ID,
		ProductID:     v.ProductID,
		AttributeID:   v.AttributeID,
		AttributeName: v.AttributeName,
		Value:         v.Value,
		Values:        v.Values,
	}
}

// ProductOptionModel is the GORM model for the product_options table.
type ProductOptionModel struct {
	ID        string                    `gorm:"type:uuid;primaryKey"`
	ProductID string                    `gorm:"type:uuid;not null;index"`
	Name      string                    `gorm:"type:varchar(255);not null"`
	SortOrder int                       `gorm:"not null;default:0"`
	Values    []ProductOptionValueModel `gorm:"foreignKey:OptionID;references:ID"`
}

func (ProductOptionModel) TableName() string { return "product_options" }

func (m *ProductOptionModel) ToDomain() domain.ProductOption {
	opt := domain.ProductOption{
		ID:        m.ID,
		ProductID: m.ProductID,
		Name:      m.Name,
		SortOrder: m.SortOrder,
	}
	for _, v := range m.Values {
		opt.Values = append(opt.Values, v.ToDomain())
	}
	return opt
}

// ProductOptionValueModel is the GORM model for product_option_values table.
type ProductOptionValueModel struct {
	ID        string `gorm:"type:uuid;primaryKey"`
	OptionID  string `gorm:"type:uuid;not null;index"`
	Value     string `gorm:"type:varchar(255);not null"`
	ColorHex  string `gorm:"type:varchar(7)"`
	SortOrder int    `gorm:"not null;default:0"`
}

func (ProductOptionValueModel) TableName() string { return "product_option_values" }

func (m *ProductOptionValueModel) ToDomain() domain.ProductOptionValue {
	return domain.ProductOptionValue{
		ID:        m.ID,
		OptionID:  m.OptionID,
		Value:     m.Value,
		ColorHex:  m.ColorHex,
		SortOrder: m.SortOrder,
	}
}

// VariantModel is the GORM model for the product_variants table.
type VariantModel struct {
	ID             string                    `gorm:"type:uuid;primaryKey"`
	ProductID      string                    `gorm:"type:uuid;not null;index"`
	SKU            string                    `gorm:"type:varchar(255);uniqueIndex;not null"`
	Name           string                    `gorm:"type:varchar(500)"`
	PriceCents     int64                     `gorm:"not null;default:0"`
	CompareAtCents int64                     `gorm:"not null;default:0"`
	CostCents      int64                     `gorm:"not null;default:0"`
	Stock          int                       `gorm:"not null;default:0"`
	LowStockAlert  int                       `gorm:"not null;default:0"`
	WeightGrams    int                       `gorm:"not null;default:0"`
	IsDefault      bool                      `gorm:"not null;default:false"`
	IsActive       bool                      `gorm:"not null;default:true"`
	ImageURLs      pq.StringArray            `gorm:"type:text[];column:image_urls"`
	Barcode        string                    `gorm:"type:varchar(255)"`
	OptionValues   []VariantOptionValueModel `gorm:"foreignKey:VariantID;references:ID"`
	CreatedAt      time.Time                 `gorm:"not null"`
	UpdatedAt      time.Time                 `gorm:"not null"`
}

func (VariantModel) TableName() string { return "product_variants" }

func (m *VariantModel) ToDomain() *domain.Variant {
	v := &domain.Variant{
		ID:             m.ID,
		ProductID:      m.ProductID,
		SKU:            m.SKU,
		Name:           m.Name,
		PriceCents:     m.PriceCents,
		CompareAtCents: m.CompareAtCents,
		CostCents:      m.CostCents,
		Stock:          m.Stock,
		LowStockAlert:  m.LowStockAlert,
		WeightGrams:    m.WeightGrams,
		IsDefault:      m.IsDefault,
		IsActive:       m.IsActive,
		ImageURLs:      m.ImageURLs,
		Barcode:        m.Barcode,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
	for _, ov := range m.OptionValues {
		v.OptionValues = append(v.OptionValues, ov.ToDomain())
	}
	return v
}

func VariantModelFromDomain(v *domain.Variant) *VariantModel {
	return &VariantModel{
		ID:             v.ID,
		ProductID:      v.ProductID,
		SKU:            v.SKU,
		Name:           v.Name,
		PriceCents:     v.PriceCents,
		CompareAtCents: v.CompareAtCents,
		CostCents:      v.CostCents,
		Stock:          v.Stock,
		LowStockAlert:  v.LowStockAlert,
		WeightGrams:    v.WeightGrams,
		IsDefault:      v.IsDefault,
		IsActive:       v.IsActive,
		ImageURLs:      v.ImageURLs,
		Barcode:        v.Barcode,
		CreatedAt:      v.CreatedAt,
		UpdatedAt:      v.UpdatedAt,
	}
}

// VariantOptionValueModel is the GORM model for variant_option_values table.
type VariantOptionValueModel struct {
	VariantID     string `gorm:"type:uuid;primaryKey"`
	OptionID      string `gorm:"type:uuid;primaryKey"`
	OptionValueID string `gorm:"type:uuid;not null"`
	OptionName    string `gorm:"type:varchar(255)"`
	Value         string `gorm:"type:varchar(255)"`
}

func (VariantOptionValueModel) TableName() string { return "variant_option_values" }

func (m *VariantOptionValueModel) ToDomain() domain.VariantOptionValue {
	return domain.VariantOptionValue{
		VariantID:     m.VariantID,
		OptionID:      m.OptionID,
		OptionValueID: m.OptionValueID,
		OptionName:    m.OptionName,
		Value:         m.Value,
	}
}
