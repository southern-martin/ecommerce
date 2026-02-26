package usecase

import (
	"context"
	"fmt"
	"strings"

	"ecommerce/catalog-service/internal/domain"
	"ecommerce/catalog-service/internal/port"

	"github.com/google/uuid"
)

type CategoryService struct {
	categories port.CategoryRepository
}

func NewCategoryService(categories port.CategoryRepository) *CategoryService {
	return &CategoryService{categories: categories}
}

type CreateCategoryInput struct {
	Name      string
	Slug      string
	ParentID  *uuid.UUID
	SortOrder int
	IsActive  bool
}

func (s *CategoryService) CreateCategory(ctx context.Context, in CreateCategoryInput) (domain.Category, error) {
	slug := strings.Trim(strings.ToLower(strings.TrimSpace(in.Slug)), "/")
	if slug == "" {
		return domain.Category{}, fmt.Errorf("%w: slug is required", domain.ErrInvalidCategory)
	}

	exists, err := s.categories.ExistsByParentAndSlug(ctx, in.ParentID, slug)
	if err != nil {
		return domain.Category{}, err
	}
	if exists {
		return domain.Category{}, domain.ErrDuplicateSlugUnderParent
	}

	c := domain.Category{
		ID:        uuid.New(),
		Name:      strings.TrimSpace(in.Name),
		Slug:      slug,
		ParentID:  in.ParentID,
		Level:     0,
		Path:      "/" + slug,
		SortOrder: in.SortOrder,
		IsActive:  in.IsActive,
	}

	if c.ParentID != nil {
		parent, err := s.categories.GetByID(ctx, *c.ParentID)
		if err != nil {
			return domain.Category{}, err
		}
		c.Level = parent.Level + 1
		c.Path = domain.BuildCategoryPath(parent.Path, c.Slug)
	}

	if err := c.Validate(); err != nil {
		return domain.Category{}, err
	}
	if err := s.categories.Create(ctx, c); err != nil {
		return domain.Category{}, err
	}
	return c, nil
}

func (s *CategoryService) ListChildren(ctx context.Context, parentID *uuid.UUID) ([]domain.Category, error) {
	return s.categories.ListChildren(ctx, parentID)
}

type CategoryTreeNode struct {
	Category domain.Category
	Children []CategoryTreeNode
}

func (s *CategoryService) GetCategoryTree(ctx context.Context) ([]CategoryTreeNode, error) {
	categories, err := s.categories.ListAll(ctx)
	if err != nil {
		return nil, err
	}

	childrenByParent := map[string][]domain.Category{}
	for _, category := range categories {
		parentKey := "root"
		if category.ParentID != nil {
			parentKey = category.ParentID.String()
		}
		childrenByParent[parentKey] = append(childrenByParent[parentKey], category)
	}

	var build func(parentKey string) []CategoryTreeNode
	build = func(parentKey string) []CategoryTreeNode {
		children := childrenByParent[parentKey]
		nodes := make([]CategoryTreeNode, 0, len(children))
		for _, category := range children {
			nodes = append(nodes, CategoryTreeNode{
				Category: category,
				Children: build(category.ID.String()),
			})
		}
		return nodes
	}

	return build("root"), nil
}
