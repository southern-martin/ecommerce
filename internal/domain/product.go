package domain

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type ProductStatus string

const (
	ProductStatusDraft    ProductStatus = "draft"
	ProductStatusActive   ProductStatus = "active"
	ProductStatusArchived ProductStatus = "archived"
)

type Product struct {
	ID                uuid.UUID
	Name              string
	Slug              string
	Description       string
	PrimaryCategoryID uuid.UUID
	CategoryIDs       []uuid.UUID
	Status            ProductStatus
	BasePriceMinor    int64
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (p Product) Validate() error {
	if p.ID == uuid.Nil || p.PrimaryCategoryID == uuid.Nil {
		return fmt.Errorf("%w: id and primary_category_id are required", ErrInvalidProduct)
	}
	if strings.TrimSpace(p.Name) == "" || strings.TrimSpace(p.Slug) == "" {
		return fmt.Errorf("%w: name and slug are required", ErrInvalidProduct)
	}
	if p.BasePriceMinor < 0 {
		return fmt.Errorf("%w: base price cannot be negative", ErrInvalidProduct)
	}
	switch p.Status {
	case ProductStatusDraft, ProductStatusActive, ProductStatusArchived:
	default:
		return fmt.Errorf("%w: unsupported status %q", ErrInvalidProduct, p.Status)
	}
	return nil
}

func NormalizeCategoryIDs(primaryCategoryID uuid.UUID, additional []uuid.UUID) []uuid.UUID {
	result := make([]uuid.UUID, 0, 1+len(additional))
	seen := map[uuid.UUID]struct{}{}

	if primaryCategoryID != uuid.Nil {
		seen[primaryCategoryID] = struct{}{}
		result = append(result, primaryCategoryID)
	}
	for _, id := range additional {
		if id == uuid.Nil {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		result = append(result, id)
	}
	return result
}

type ProductAttributeValue struct {
	ID           uuid.UUID
	ProductID    uuid.UUID
	AttributeID  uuid.UUID
	OptionID     *uuid.UUID
	ValueText    *string
	ValueNumber  *float64
	ValueBoolean *bool
	ValueJSON    *json.RawMessage
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (v ProductAttributeValue) Validate() error {
	if v.ProductID == uuid.Nil || v.AttributeID == uuid.Nil {
		return fmt.Errorf("%w: product_id and attribute_id are required", ErrInvalidAttributeValue)
	}
	count := 0
	if v.OptionID != nil {
		count++
	}
	if v.ValueText != nil {
		count++
	}
	if v.ValueNumber != nil {
		count++
	}
	if v.ValueBoolean != nil {
		count++
	}
	if v.ValueJSON != nil {
		count++
	}
	if count != 1 {
		return fmt.Errorf("%w: exactly one value must be set", ErrInvalidAttributeValue)
	}
	return nil
}
