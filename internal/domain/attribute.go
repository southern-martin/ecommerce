package domain

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type AttributeType string

const (
	AttributeTypeSelect AttributeType = "select"
	AttributeTypeText   AttributeType = "text"
	AttributeTypeNumber AttributeType = "number"
	AttributeTypeBool   AttributeType = "bool"
)

func (t AttributeType) Validate() error {
	switch t {
	case AttributeTypeSelect, AttributeTypeText, AttributeTypeNumber, AttributeTypeBool:
		return nil
	default:
		return fmt.Errorf("%w: unsupported attribute type %q", ErrInvalidAttribute, t)
	}
}

type CategoryAttribute struct {
	ID            uuid.UUID
	CategoryID    uuid.UUID
	Name          string
	Code          string
	Type          AttributeType
	Required      bool
	IsVariantAxis bool
	IsFilterable  bool
	SortOrder     int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (a CategoryAttribute) Validate() error {
	if a.ID == uuid.Nil || a.CategoryID == uuid.Nil {
		return fmt.Errorf("%w: id and category_id are required", ErrInvalidAttribute)
	}
	if strings.TrimSpace(a.Name) == "" || strings.TrimSpace(a.Code) == "" {
		return fmt.Errorf("%w: name and code are required", ErrInvalidAttribute)
	}
	if err := a.Type.Validate(); err != nil {
		return err
	}
	return nil
}

type AttributeOption struct {
	ID          uuid.UUID
	AttributeID uuid.UUID
	Value       string
	Label       string
	SortOrder   int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (o AttributeOption) Validate() error {
	if o.ID == uuid.Nil || o.AttributeID == uuid.Nil {
		return fmt.Errorf("%w: option id and attribute id are required", ErrInvalidAttribute)
	}
	if strings.TrimSpace(o.Value) == "" || strings.TrimSpace(o.Label) == "" {
		return fmt.Errorf("%w: option value and label are required", ErrInvalidAttribute)
	}
	return nil
}
