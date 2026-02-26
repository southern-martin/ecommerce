package usecase

import (
	"context"
	"fmt"
	"strings"

	"ecommerce/catalog-service/internal/domain"
	"ecommerce/catalog-service/internal/port"

	"github.com/google/uuid"
)

type AttributeService struct {
	attributes port.AttributeRepository
	categories port.CategoryRepository
}

func NewAttributeService(attributes port.AttributeRepository, categories port.CategoryRepository) *AttributeService {
	return &AttributeService{
		attributes: attributes,
		categories: categories,
	}
}

type CreateCategoryAttributeInput struct {
	CategoryID    uuid.UUID
	Name          string
	Code          string
	Type          string
	Required      bool
	IsVariantAxis bool
	IsFilterable  bool
	SortOrder     int
}

func (s *AttributeService) CreateCategoryAttribute(
	ctx context.Context,
	in CreateCategoryAttributeInput,
) (domain.CategoryAttribute, error) {
	if in.CategoryID == uuid.Nil {
		return domain.CategoryAttribute{}, fmt.Errorf("%w: category_id is required", domain.ErrInvalidAttribute)
	}
	if _, err := s.categories.GetByID(ctx, in.CategoryID); err != nil {
		return domain.CategoryAttribute{}, err
	}

	code := strings.Trim(strings.ToLower(strings.TrimSpace(in.Code)), "/")
	if code == "" {
		return domain.CategoryAttribute{}, fmt.Errorf("%w: code is required", domain.ErrInvalidAttribute)
	}
	exists, err := s.attributes.ExistsByCategoryAndCode(ctx, in.CategoryID, code)
	if err != nil {
		return domain.CategoryAttribute{}, err
	}
	if exists {
		return domain.CategoryAttribute{}, fmt.Errorf("%w: duplicate code in category", domain.ErrInvalidAttribute)
	}

	attribute := domain.CategoryAttribute{
		ID:            uuid.New(),
		CategoryID:    in.CategoryID,
		Name:          strings.TrimSpace(in.Name),
		Code:          code,
		Type:          domain.AttributeType(strings.TrimSpace(in.Type)),
		Required:      in.Required,
		IsVariantAxis: in.IsVariantAxis,
		IsFilterable:  in.IsFilterable,
		SortOrder:     in.SortOrder,
	}
	if attribute.IsVariantAxis && attribute.Type != domain.AttributeTypeSelect {
		return domain.CategoryAttribute{}, fmt.Errorf("%w: variant axis must be select type", domain.ErrInvalidAttribute)
	}
	if err := attribute.Validate(); err != nil {
		return domain.CategoryAttribute{}, err
	}
	if err := s.attributes.Create(ctx, attribute); err != nil {
		return domain.CategoryAttribute{}, err
	}
	return attribute, nil
}

type AddAttributeOptionInput struct {
	AttributeID uuid.UUID
	Value       string
	Label       string
	SortOrder   int
}

func (s *AttributeService) AddAttributeOption(ctx context.Context, in AddAttributeOptionInput) (domain.AttributeOption, error) {
	attribute, err := s.attributes.GetByID(ctx, in.AttributeID)
	if err != nil {
		return domain.AttributeOption{}, err
	}
	if attribute.Type != domain.AttributeTypeSelect {
		return domain.AttributeOption{}, fmt.Errorf("%w: options only allowed for select attributes", domain.ErrInvalidAttribute)
	}

	option := domain.AttributeOption{
		ID:          uuid.New(),
		AttributeID: in.AttributeID,
		Value:       strings.TrimSpace(in.Value),
		Label:       strings.TrimSpace(in.Label),
		SortOrder:   in.SortOrder,
	}
	if err := option.Validate(); err != nil {
		return domain.AttributeOption{}, err
	}
	if err := s.attributes.CreateOption(ctx, option); err != nil {
		return domain.AttributeOption{}, err
	}
	return option, nil
}

func (s *AttributeService) ListCategoryAttributes(
	ctx context.Context,
	categoryID uuid.UUID,
) ([]domain.CategoryAttribute, error) {
	return s.attributes.GetByCategory(ctx, categoryID)
}

func (s *AttributeService) ListAttributeOptions(
	ctx context.Context,
	attributeID uuid.UUID,
) ([]domain.AttributeOption, error) {
	return s.attributes.ListOptionsByAttribute(ctx, attributeID)
}
