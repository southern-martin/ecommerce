package postgres

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/lib/pq"
	"github.com/southern-martin/ecommerce/services/search/internal/domain"
)

// AttributesJSON is a custom type for storing map[string]string as JSONB.
type AttributesJSON map[string]string

// Value implements the driver.Valuer interface.
func (a AttributesJSON) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}
	return json.Marshal(a)
}

// Scan implements the sql.Scanner interface.
func (a *AttributesJSON) Scan(value interface{}) error {
	if value == nil {
		*a = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan AttributesJSON: invalid type")
	}
	return json.Unmarshal(bytes, a)
}

// SearchIndexModel is the GORM model for the search_indices table.
type SearchIndexModel struct {
	ID          string         `gorm:"type:uuid;primaryKey"`
	ProductID   string         `gorm:"type:varchar(255);uniqueIndex;not null"`
	Name        string         `gorm:"type:varchar(500);not null"`
	Slug        string         `gorm:"type:varchar(500)"`
	Description string         `gorm:"type:text"`
	PriceCents  int64          `gorm:"not null;default:0"`
	Currency    string         `gorm:"type:varchar(3);default:'USD'"`
	CategoryID  string         `gorm:"type:varchar(255);index"`
	SellerID    string         `gorm:"type:varchar(255);index"`
	ImageURL    string         `gorm:"type:text"`
	Rating      float64        `gorm:"type:decimal(3,2);default:0"`
	ReviewCount int            `gorm:"default:0"`
	InStock     bool           `gorm:"default:true;index"`
	Tags        pq.StringArray `gorm:"type:text[]"`
	Attributes  AttributesJSON `gorm:"type:jsonb"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// TableName returns the table name for SearchIndexModel.
func (SearchIndexModel) TableName() string {
	return "search_indices"
}

// ToDomain converts the GORM model to a domain SearchIndex entity.
func (m *SearchIndexModel) ToDomain() *domain.SearchIndex {
	attrs := make(map[string]string)
	for k, v := range m.Attributes {
		attrs[k] = v
	}

	tags := make([]string, len(m.Tags))
	copy(tags, m.Tags)

	return &domain.SearchIndex{
		ID:          m.ID,
		ProductID:   m.ProductID,
		Name:        m.Name,
		Slug:        m.Slug,
		Description: m.Description,
		PriceCents:  m.PriceCents,
		Currency:    m.Currency,
		CategoryID:  m.CategoryID,
		SellerID:    m.SellerID,
		ImageURL:    m.ImageURL,
		Rating:      m.Rating,
		ReviewCount: m.ReviewCount,
		InStock:     m.InStock,
		Tags:        tags,
		Attributes:  attrs,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

// ToModel converts a domain SearchIndex entity to a GORM model.
func ToModel(idx *domain.SearchIndex) *SearchIndexModel {
	attrs := make(AttributesJSON)
	for k, v := range idx.Attributes {
		attrs[k] = v
	}

	tags := make(pq.StringArray, len(idx.Tags))
	copy(tags, idx.Tags)

	return &SearchIndexModel{
		ID:          idx.ID,
		ProductID:   idx.ProductID,
		Name:        idx.Name,
		Slug:        idx.Slug,
		Description: idx.Description,
		PriceCents:  idx.PriceCents,
		Currency:    idx.Currency,
		CategoryID:  idx.CategoryID,
		SellerID:    idx.SellerID,
		ImageURL:    idx.ImageURL,
		Rating:      idx.Rating,
		ReviewCount: idx.ReviewCount,
		InStock:     idx.InStock,
		Tags:        tags,
		Attributes:  attrs,
		CreatedAt:   idx.CreatedAt,
		UpdatedAt:   idx.UpdatedAt,
	}
}
