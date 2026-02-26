package domain

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID        uuid.UUID
	Name      string
	Slug      string
	ParentID  *uuid.UUID
	Level     int
	Path      string
	SortOrder int
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (c Category) Validate() error {
	if c.ID == uuid.Nil || strings.TrimSpace(c.Name) == "" || strings.TrimSpace(c.Slug) == "" {
		return fmt.Errorf("%w: id, name and slug are required", ErrInvalidCategory)
	}
	if strings.Contains(c.Slug, "/") || strings.ContainsAny(c.Slug, " \t\n") {
		return fmt.Errorf("%w: slug must not contain slash or spaces", ErrInvalidCategory)
	}
	if c.ParentID != nil && *c.ParentID == c.ID {
		return fmt.Errorf("%w: category cannot parent itself", ErrInvalidCategoryHierarchy)
	}
	if c.Level < 0 || !strings.HasPrefix(c.Path, "/") {
		return fmt.Errorf("%w: invalid level/path", ErrInvalidCategory)
	}
	return nil
}

func BuildCategoryPath(parentPath, slug string) string {
	parentPath = strings.TrimSpace(parentPath)
	slug = strings.Trim(strings.TrimSpace(slug), "/")
	if parentPath == "" || parentPath == "/" {
		return "/" + slug
	}
	return strings.TrimSuffix(parentPath, "/") + "/" + slug
}
