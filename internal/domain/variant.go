package domain

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
)

type ProductVariantStatus string

const (
	ProductVariantStatusActive   ProductVariantStatus = "active"
	ProductVariantStatusInactive ProductVariantStatus = "inactive"
)

type VariantOptionValue struct {
	AttributeID uuid.UUID
	OptionID    uuid.UUID
}

type ProductVariant struct {
	ID             uuid.UUID
	ProductID      uuid.UUID
	SKU            string
	PriceMinor     int64
	StockQty       int64
	ImageURL       string
	Status         ProductVariantStatus
	CombinationKey string
	Options        []VariantOptionValue
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (v ProductVariant) Validate() error {
	if v.ID == uuid.Nil || v.ProductID == uuid.Nil {
		return fmt.Errorf("%w: variant id and product id are required", ErrInvalidVariantAxis)
	}
	if v.PriceMinor < 0 || v.StockQty < 0 {
		return fmt.Errorf("%w: negative price/stock is not allowed", ErrInvalidVariantAxis)
	}
	switch v.Status {
	case ProductVariantStatusActive, ProductVariantStatusInactive:
	default:
		return fmt.Errorf("%w: unsupported variant status %q", ErrInvalidVariantAxis, v.Status)
	}
	if len(v.Options) == 0 {
		return fmt.Errorf("%w: variant requires at least one option", ErrInvalidVariantAxis)
	}
	key, err := BuildCombinationKey(v.Options)
	if err != nil {
		return err
	}
	if v.CombinationKey != "" && v.CombinationKey != key {
		return fmt.Errorf("%w: combination key mismatch", ErrInvalidVariantAxis)
	}
	return nil
}

func BuildCombinationKey(options []VariantOptionValue) (string, error) {
	if len(options) == 0 {
		return "", fmt.Errorf("%w: no option values provided", ErrInvalidVariantAxis)
	}
	sorted := append([]VariantOptionValue(nil), options...)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].AttributeID.String() < sorted[j].AttributeID.String()
	})

	seenAttribute := map[uuid.UUID]struct{}{}
	parts := make([]string, 0, len(sorted))
	for _, ov := range sorted {
		if ov.AttributeID == uuid.Nil || ov.OptionID == uuid.Nil {
			return "", fmt.Errorf("%w: attribute_id and option_id are required", ErrInvalidVariantAxis)
		}
		if _, ok := seenAttribute[ov.AttributeID]; ok {
			return "", fmt.Errorf("%w: duplicate attribute in combination", ErrDuplicateVariantCombination)
		}
		seenAttribute[ov.AttributeID] = struct{}{}
		parts = append(parts, ov.AttributeID.String()+":"+ov.OptionID.String())
	}
	return strings.Join(parts, "|"), nil
}
