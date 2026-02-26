package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"sort"

	"ecommerce/catalog-service/internal/domain"

	"github.com/google/uuid"
)

type CategoryRepo struct {
	db *sql.DB
}

func NewCategoryRepo(db *sql.DB) *CategoryRepo {
	return &CategoryRepo{db: db}
}

func (r *CategoryRepo) Create(ctx context.Context, category domain.Category) error {
	var parentID any
	if category.ParentID != nil {
		parentID = category.ParentID.String()
	}
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO categories (
			id, name, slug, parent_id, level, path, sort_order, is_active
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, category.ID.String(), category.Name, category.Slug, parentID, category.Level, category.Path, category.SortOrder, category.IsActive)
	if err != nil {
		return fmt.Errorf("insert category: %w", err)
	}
	return nil
}

func (r *CategoryRepo) GetByID(ctx context.Context, id uuid.UUID) (domain.Category, error) {
	var (
		c        domain.Category
		idStr    string
		parentID sql.NullString
	)
	err := r.db.QueryRowContext(ctx, `
		SELECT id::text, name, slug, parent_id::text, level, path, sort_order, is_active, created_at, updated_at
		FROM categories
		WHERE id = $1
	`, id.String()).Scan(
		&idStr, &c.Name, &c.Slug, &parentID, &c.Level, &c.Path, &c.SortOrder, &c.IsActive, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Category{}, fmt.Errorf("%w: category %s", domain.ErrNotFound, id.String())
		}
		return domain.Category{}, fmt.Errorf("get category: %w", err)
	}
	parsedID, err := uuid.Parse(idStr)
	if err != nil {
		return domain.Category{}, fmt.Errorf("parse category id: %w", err)
	}
	c.ID = parsedID
	if parentID.Valid {
		parsed, err := uuid.Parse(parentID.String)
		if err != nil {
			return domain.Category{}, fmt.Errorf("parse parent id: %w", err)
		}
		c.ParentID = &parsed
	}
	return c, nil
}

func (r *CategoryRepo) ExistsByParentAndSlug(ctx context.Context, parentID *uuid.UUID, slug string) (bool, error) {
	var exists bool
	if parentID == nil {
		err := r.db.QueryRowContext(ctx, `
			SELECT EXISTS(
				SELECT 1 FROM categories
				WHERE parent_id IS NULL AND slug = $1
			)
		`, slug).Scan(&exists)
		if err != nil {
			return false, fmt.Errorf("exists category root slug: %w", err)
		}
		return exists, nil
	}
	err := r.db.QueryRowContext(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM categories
			WHERE parent_id = $1 AND slug = $2
		)
	`, parentID.String(), slug).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("exists category child slug: %w", err)
	}
	return exists, nil
}

func (r *CategoryRepo) ListChildren(ctx context.Context, parentID *uuid.UUID) ([]domain.Category, error) {
	var (
		rows *sql.Rows
		err  error
	)
	if parentID == nil {
		rows, err = r.db.QueryContext(ctx, `
			SELECT id::text, name, slug, parent_id::text, level, path, sort_order, is_active, created_at, updated_at
			FROM categories
			WHERE parent_id IS NULL
			ORDER BY sort_order, name
		`)
	} else {
		rows, err = r.db.QueryContext(ctx, `
			SELECT id::text, name, slug, parent_id::text, level, path, sort_order, is_active, created_at, updated_at
			FROM categories
			WHERE parent_id = $1
			ORDER BY sort_order, name
		`, parentID.String())
	}
	if err != nil {
		return nil, fmt.Errorf("list children categories: %w", err)
	}
	defer rows.Close()

	out := make([]domain.Category, 0)
	for rows.Next() {
		c, err := scanCategory(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func (r *CategoryRepo) ListAll(ctx context.Context) ([]domain.Category, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id::text, name, slug, parent_id::text, level, path, sort_order, is_active, created_at, updated_at
		FROM categories
	`)
	if err != nil {
		return nil, fmt.Errorf("list all categories: %w", err)
	}
	defer rows.Close()

	out := make([]domain.Category, 0)
	for rows.Next() {
		c, err := scanCategory(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Path == out[j].Path {
			return out[i].Name < out[j].Name
		}
		return out[i].Path < out[j].Path
	})
	return out, nil
}

func scanCategory(scanner interface {
	Scan(dest ...any) error
}) (domain.Category, error) {
	var (
		c        domain.Category
		id       string
		parentID sql.NullString
	)
	if err := scanner.Scan(
		&id, &c.Name, &c.Slug, &parentID, &c.Level, &c.Path, &c.SortOrder, &c.IsActive, &c.CreatedAt, &c.UpdatedAt,
	); err != nil {
		return domain.Category{}, fmt.Errorf("scan category: %w", err)
	}
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return domain.Category{}, fmt.Errorf("parse category id: %w", err)
	}
	c.ID = parsedID
	if parentID.Valid {
		p, err := uuid.Parse(parentID.String)
		if err != nil {
			return domain.Category{}, fmt.Errorf("parse category parent id: %w", err)
		}
		c.ParentID = &p
	}
	return c, nil
}
