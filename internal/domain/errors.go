package domain

import "errors"

var (
	ErrInvalidCategory             = errors.New("invalid category")
	ErrInvalidCategoryHierarchy    = errors.New("invalid category hierarchy")
	ErrDuplicateSlugUnderParent    = errors.New("duplicate category slug under parent")
	ErrInvalidAttribute            = errors.New("invalid category attribute")
	ErrInvalidAttributeValue       = errors.New("invalid product attribute value")
	ErrInvalidProduct              = errors.New("invalid product")
	ErrInvalidVariantAxis          = errors.New("invalid variant axis")
	ErrDuplicateVariantCombination = errors.New("duplicate variant combination")
	ErrNotFound                    = errors.New("not found")
)
