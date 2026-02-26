package domain

import "time"

// ProductStatus represents the lifecycle status of a product.
type ProductStatus string

const (
	ProductStatusDraft    ProductStatus = "draft"
	ProductStatusActive   ProductStatus = "active"
	ProductStatusInactive ProductStatus = "inactive"
	ProductStatusArchived ProductStatus = "archived"
)

// Product represents a product in the catalog.
type Product struct {
	ID             string        `json:"id"`
	SellerID       string        `json:"seller_id"`
	CategoryID     string        `json:"category_id"`
	Name           string        `json:"name"`
	Slug           string        `json:"slug"`
	Description    string        `json:"description"`
	BasePriceCents int64         `json:"base_price_cents"`
	Currency       string        `json:"currency"`
	Status         ProductStatus `json:"status"`
	HasVariants    bool          `json:"has_variants"`
	Tags           []string      `json:"tags"`
	ImageURLs      []string      `json:"image_urls"`
	RatingAvg      float64       `json:"rating_avg"`
	RatingCount    int           `json:"rating_count"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
}

// Category represents a product category with optional parent for hierarchy.
type Category struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	ParentID  string    `json:"parent_id,omitempty"`
	SortOrder int       `json:"sort_order"`
	ImageURL  string    `json:"image_url,omitempty"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

// AttributeType defines the type of an attribute definition.
type AttributeType string

const (
	AttributeTypeText        AttributeType = "text"
	AttributeTypeNumber      AttributeType = "number"
	AttributeTypeSelect      AttributeType = "select"
	AttributeTypeMultiSelect AttributeType = "multi_select"
	AttributeTypeColor       AttributeType = "color"
	AttributeTypeBool        AttributeType = "bool"
)

// AttributeDefinition defines a product attribute schema.
type AttributeDefinition struct {
	ID         string        `json:"id"`
	Name       string        `json:"name"`
	Slug       string        `json:"slug"`
	Type       AttributeType `json:"type"`
	Required   bool          `json:"required"`
	Filterable bool          `json:"filterable"`
	Options    []string      `json:"options,omitempty"`
	Unit       string        `json:"unit,omitempty"`
	SortOrder  int           `json:"sort_order"`
	CreatedAt  time.Time     `json:"created_at"`
}

// ProductAttributeValue holds the value of an attribute for a specific product.
type ProductAttributeValue struct {
	ID            string   `json:"id"`
	ProductID     string   `json:"product_id"`
	AttributeID   string   `json:"attribute_id"`
	AttributeName string   `json:"attribute_name"`
	Value         string   `json:"value"`
	Values        []string `json:"values,omitempty"`
}

// ProductOption defines an option group for a product (e.g., Size, Color).
type ProductOption struct {
	ID        string               `json:"id"`
	ProductID string               `json:"product_id"`
	Name      string               `json:"name"`
	SortOrder int                  `json:"sort_order"`
	Values    []ProductOptionValue `json:"values,omitempty"`
}

// ProductOptionValue defines a specific option value (e.g., "Large", "Red").
type ProductOptionValue struct {
	ID        string `json:"id"`
	OptionID  string `json:"option_id"`
	Value     string `json:"value"`
	ColorHex  string `json:"color_hex,omitempty"`
	SortOrder int    `json:"sort_order"`
}

// Variant represents a specific purchasable variant of a product.
type Variant struct {
	ID             string               `json:"id"`
	ProductID      string               `json:"product_id"`
	SKU            string               `json:"sku"`
	Name           string               `json:"name"`
	PriceCents     int64                `json:"price_cents"`
	CompareAtCents int64                `json:"compare_at_cents"`
	CostCents      int64                `json:"cost_cents"`
	Stock          int                  `json:"stock"`
	LowStockAlert  int                  `json:"low_stock_alert"`
	WeightGrams    int                  `json:"weight_grams"`
	IsDefault      bool                 `json:"is_default"`
	IsActive       bool                 `json:"is_active"`
	ImageURLs      []string             `json:"image_urls,omitempty"`
	Barcode        string               `json:"barcode,omitempty"`
	OptionValues   []VariantOptionValue `json:"option_values,omitempty"`
	CreatedAt      time.Time            `json:"created_at"`
	UpdatedAt      time.Time            `json:"updated_at"`
}

// VariantOptionValue links a variant to specific option values.
type VariantOptionValue struct {
	VariantID     string `json:"variant_id"`
	OptionID      string `json:"option_id"`
	OptionValueID string `json:"option_value_id"`
	OptionName    string `json:"option_name"`
	Value         string `json:"value"`
}
