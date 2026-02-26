package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"sort"

	"ecommerce/catalog-service/internal/domain"

	"github.com/google/uuid"
)

type AttributeRepo struct {
	db *sql.DB
}

func NewAttributeRepo(db *sql.DB) *AttributeRepo {
	return &AttributeRepo{db: db}
}

func (r *AttributeRepo) Create(ctx context.Context, attribute domain.CategoryAttribute) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO category_attributes (
			id, category_id, name, code, attribute_type, required, is_variant_axis, is_filterable, sort_order
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`,
		attribute.ID.String(),
		attribute.CategoryID.String(),
		attribute.Name,
		attribute.Code,
		string(attribute.Type),
		attribute.Required,
		attribute.IsVariantAxis,
		attribute.IsFilterable,
		attribute.SortOrder,
	)
	if err != nil {
		return fmt.Errorf("insert category attribute: %w", err)
	}
	return nil
}

func (r *AttributeRepo) GetByID(ctx context.Context, id uuid.UUID) (domain.CategoryAttribute, error) {
	var (
		a          domain.CategoryAttribute
		attrID     string
		categoryID string
		attrType   string
	)
	err := r.db.QueryRowContext(ctx, `
		SELECT id::text, category_id::text, name, code, attribute_type, required, is_variant_axis, is_filterable, sort_order, created_at, updated_at
		FROM category_attributes
		WHERE id = $1
	`, id.String()).Scan(
		&attrID, &categoryID, &a.Name, &a.Code, &attrType, &a.Required, &a.IsVariantAxis, &a.IsFilterable, &a.SortOrder, &a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.CategoryAttribute{}, fmt.Errorf("%w: attribute %s", domain.ErrNotFound, id.String())
		}
		return domain.CategoryAttribute{}, fmt.Errorf("get attribute: %w", err)
	}
	parsedID, err := uuid.Parse(attrID)
	if err != nil {
		return domain.CategoryAttribute{}, fmt.Errorf("parse attribute id: %w", err)
	}
	parsedCategoryID, err := uuid.Parse(categoryID)
	if err != nil {
		return domain.CategoryAttribute{}, fmt.Errorf("parse attribute category id: %w", err)
	}
	a.ID = parsedID
	a.CategoryID = parsedCategoryID
	a.Type = domain.AttributeType(attrType)
	return a, nil
}

func (r *AttributeRepo) GetByCategory(ctx context.Context, categoryID uuid.UUID) ([]domain.CategoryAttribute, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id::text, category_id::text, name, code, attribute_type, required, is_variant_axis, is_filterable, sort_order, created_at, updated_at
		FROM category_attributes
		WHERE category_id = $1
		ORDER BY sort_order, code
	`, categoryID.String())
	if err != nil {
		return nil, fmt.Errorf("list attributes by category: %w", err)
	}
	defer rows.Close()

	out := make([]domain.CategoryAttribute, 0)
	for rows.Next() {
		a, err := scanAttribute(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

func (r *AttributeRepo) ExistsByCategoryAndCode(ctx context.Context, categoryID uuid.UUID, code string) (bool, error) {
	var exists bool
	if err := r.db.QueryRowContext(ctx, `
		SELECT EXISTS(
			SELECT 1
			FROM category_attributes
			WHERE category_id = $1 AND code = $2
		)
	`, categoryID.String(), code).Scan(&exists); err != nil {
		return false, fmt.Errorf("exists category attribute code: %w", err)
	}
	return exists, nil
}

func (r *AttributeRepo) CreateOption(ctx context.Context, option domain.AttributeOption) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO attribute_options (
			id, attribute_id, value, label, sort_order
		) VALUES ($1, $2, $3, $4, $5)
	`, option.ID.String(), option.AttributeID.String(), option.Value, option.Label, option.SortOrder)
	if err != nil {
		return fmt.Errorf("insert attribute option: %w", err)
	}
	return nil
}

func (r *AttributeRepo) ListOptionsByAttribute(ctx context.Context, attributeID uuid.UUID) ([]domain.AttributeOption, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id::text, attribute_id::text, value, label, sort_order, created_at, updated_at
		FROM attribute_options
		WHERE attribute_id = $1
		ORDER BY sort_order, value
	`, attributeID.String())
	if err != nil {
		return nil, fmt.Errorf("list options by attribute: %w", err)
	}
	defer rows.Close()

	out := make([]domain.AttributeOption, 0)
	for rows.Next() {
		opt, err := scanAttributeOption(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, opt)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].SortOrder == out[j].SortOrder {
			return out[i].Value < out[j].Value
		}
		return out[i].SortOrder < out[j].SortOrder
	})
	return out, rows.Err()
}

func (r *AttributeRepo) OptionBelongsToAttribute(ctx context.Context, optionID, attributeID uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, `
		SELECT EXISTS(
			SELECT 1
			FROM attribute_options
			WHERE id = $1 AND attribute_id = $2
		)
	`, optionID.String(), attributeID.String()).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check option belongs to attribute: %w", err)
	}
	return exists, nil
}

func scanAttribute(scanner interface {
	Scan(dest ...any) error
}) (domain.CategoryAttribute, error) {
	var (
		a          domain.CategoryAttribute
		id         string
		categoryID string
		attrType   string
	)
	if err := scanner.Scan(
		&id, &categoryID, &a.Name, &a.Code, &attrType, &a.Required, &a.IsVariantAxis, &a.IsFilterable, &a.SortOrder, &a.CreatedAt, &a.UpdatedAt,
	); err != nil {
		return domain.CategoryAttribute{}, fmt.Errorf("scan category attribute: %w", err)
	}
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return domain.CategoryAttribute{}, fmt.Errorf("parse category attribute id: %w", err)
	}
	parsedCategoryID, err := uuid.Parse(categoryID)
	if err != nil {
		return domain.CategoryAttribute{}, fmt.Errorf("parse category attribute category id: %w", err)
	}
	a.ID = parsedID
	a.CategoryID = parsedCategoryID
	a.Type = domain.AttributeType(attrType)
	return a, nil
}

func scanAttributeOption(scanner interface {
	Scan(dest ...any) error
}) (domain.AttributeOption, error) {
	var (
		o           domain.AttributeOption
		id          string
		attributeID string
	)
	if err := scanner.Scan(
		&id, &attributeID, &o.Value, &o.Label, &o.SortOrder, &o.CreatedAt, &o.UpdatedAt,
	); err != nil {
		return domain.AttributeOption{}, fmt.Errorf("scan attribute option: %w", err)
	}
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return domain.AttributeOption{}, fmt.Errorf("parse attribute option id: %w", err)
	}
	parsedAttributeID, err := uuid.Parse(attributeID)
	if err != nil {
		return domain.AttributeOption{}, fmt.Errorf("parse attribute option attribute id: %w", err)
	}
	o.ID = parsedID
	o.AttributeID = parsedAttributeID
	return o, nil
}
