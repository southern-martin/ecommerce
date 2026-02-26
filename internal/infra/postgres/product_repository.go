package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sort"

	"ecommerce/catalog-service/internal/domain"

	"github.com/google/uuid"
)

type ProductRepo struct {
	db *sql.DB
}

func NewProductRepo(db *sql.DB) *ProductRepo {
	return &ProductRepo{db: db}
}

func (r *ProductRepo) Create(ctx context.Context, product domain.Product) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO products (
			id, name, slug, description, primary_category_id, status, base_price_minor
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, product.ID.String(), product.Name, product.Slug, product.Description, product.PrimaryCategoryID.String(), string(product.Status), product.BasePriceMinor)
	if err != nil {
		return fmt.Errorf("insert product: %w", err)
	}
	return nil
}

func (r *ProductRepo) GetByID(ctx context.Context, id uuid.UUID) (domain.Product, error) {
	var (
		p                 domain.Product
		productID         string
		primaryCategoryID string
		status            string
	)
	err := r.db.QueryRowContext(ctx, `
		SELECT id::text, name, slug, description, primary_category_id::text, status, base_price_minor, created_at, updated_at
		FROM products
		WHERE id = $1
	`, id.String()).Scan(
		&productID, &p.Name, &p.Slug, &p.Description, &primaryCategoryID, &status, &p.BasePriceMinor, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Product{}, fmt.Errorf("%w: product %s", domain.ErrNotFound, id.String())
		}
		return domain.Product{}, fmt.Errorf("get product: %w", err)
	}
	parsedProductID, err := uuid.Parse(productID)
	if err != nil {
		return domain.Product{}, fmt.Errorf("parse product id: %w", err)
	}
	parsedPrimaryCategoryID, err := uuid.Parse(primaryCategoryID)
	if err != nil {
		return domain.Product{}, fmt.Errorf("parse product primary category id: %w", err)
	}
	p.ID = parsedProductID
	p.PrimaryCategoryID = parsedPrimaryCategoryID
	p.Status = domain.ProductStatus(status)

	rows, err := r.db.QueryContext(ctx, `
		SELECT category_id::text
		FROM product_categories
		WHERE product_id = $1
		ORDER BY is_primary DESC, category_id
	`, p.ID.String())
	if err != nil {
		return domain.Product{}, fmt.Errorf("list product categories: %w", err)
	}
	defer rows.Close()

	p.CategoryIDs = make([]uuid.UUID, 0)
	for rows.Next() {
		var categoryID string
		if err := rows.Scan(&categoryID); err != nil {
			return domain.Product{}, fmt.Errorf("scan product category: %w", err)
		}
		id, err := uuid.Parse(categoryID)
		if err != nil {
			return domain.Product{}, fmt.Errorf("parse product category id: %w", err)
		}
		p.CategoryIDs = append(p.CategoryIDs, id)
	}
	return p, rows.Err()
}

func (r *ProductRepo) SetCategories(ctx context.Context, productID uuid.UUID, categoryIDs []uuid.UUID, primaryCategoryID uuid.UUID) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin set categories tx: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `DELETE FROM product_categories WHERE product_id = $1`, productID.String()); err != nil {
		return fmt.Errorf("delete product categories: %w", err)
	}
	for _, categoryID := range categoryIDs {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO product_categories (product_id, category_id, is_primary)
			VALUES ($1, $2, $3)
		`, productID.String(), categoryID.String(), categoryID == primaryCategoryID)
		if err != nil {
			return fmt.Errorf("insert product category: %w", err)
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit set categories tx: %w", err)
	}
	return nil
}

func (r *ProductRepo) UpsertAttributeValues(ctx context.Context, productID uuid.UUID, values []domain.ProductAttributeValue) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin upsert product attributes tx: %w", err)
	}
	defer tx.Rollback()

	for _, value := range values {
		var (
			optionID any
			valText  any
			valNum   any
			valBool  any
			valJSON  any
		)
		if value.OptionID != nil {
			optionID = value.OptionID.String()
		}
		if value.ValueText != nil {
			valText = *value.ValueText
		}
		if value.ValueNumber != nil {
			valNum = *value.ValueNumber
		}
		if value.ValueBoolean != nil {
			valBool = *value.ValueBoolean
		}
		if value.ValueJSON != nil {
			valJSON = string(*value.ValueJSON)
		}

		_, err := tx.ExecContext(ctx, `
			INSERT INTO product_attribute_values (
				id, product_id, attribute_id, option_id, value_text, value_number, value_boolean, value_json
			) VALUES ($1, $2, $3, $4, $5, $6, $7, CAST($8 AS JSONB))
			ON CONFLICT (product_id, attribute_id) DO UPDATE
			SET
				option_id = EXCLUDED.option_id,
				value_text = EXCLUDED.value_text,
				value_number = EXCLUDED.value_number,
				value_boolean = EXCLUDED.value_boolean,
				value_json = EXCLUDED.value_json,
				updated_at = now()
		`, value.ID.String(), productID.String(), value.AttributeID.String(), optionID, valText, valNum, valBool, valJSON)
		if err != nil {
			return fmt.Errorf("upsert product attribute value: %w", err)
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit upsert product attributes tx: %w", err)
	}
	return nil
}

func (r *ProductRepo) ListAttributeValues(ctx context.Context, productID uuid.UUID) ([]domain.ProductAttributeValue, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT
			id::text,
			product_id::text,
			attribute_id::text,
			option_id::text,
			value_text,
			value_number::double precision,
			value_boolean,
			value_json::text,
			created_at,
			updated_at
		FROM product_attribute_values
		WHERE product_id = $1
	`, productID.String())
	if err != nil {
		return nil, fmt.Errorf("list product attribute values: %w", err)
	}
	defer rows.Close()

	out := make([]domain.ProductAttributeValue, 0)
	for rows.Next() {
		value, err := scanProductAttributeValue(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, value)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].AttributeID.String() < out[j].AttributeID.String()
	})
	return out, rows.Err()
}

func (r *ProductRepo) ListByCategory(ctx context.Context, categoryID uuid.UUID) ([]domain.Product, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT p.id::text, p.name, p.slug, p.description, p.primary_category_id::text, p.status, p.base_price_minor, p.created_at, p.updated_at
		FROM products p
		INNER JOIN product_categories pc ON pc.product_id = p.id
		WHERE pc.category_id = $1
		ORDER BY p.name
	`, categoryID.String())
	if err != nil {
		return nil, fmt.Errorf("list products by category: %w", err)
	}
	defer rows.Close()

	out := make([]domain.Product, 0)
	for rows.Next() {
		p, err := scanProduct(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func (r *ProductRepo) CreateVariants(ctx context.Context, variants []domain.ProductVariant) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin create variants tx: %w", err)
	}
	defer tx.Rollback()

	for _, variant := range variants {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO product_variants (
				id, product_id, sku, price_minor, stock_qty, image_url, status, combination_key
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		`,
			variant.ID.String(),
			variant.ProductID.String(),
			variant.SKU,
			variant.PriceMinor,
			variant.StockQty,
			variant.ImageURL,
			string(variant.Status),
			variant.CombinationKey,
		)
		if err != nil {
			return fmt.Errorf("insert product variant: %w", err)
		}
		for _, option := range variant.Options {
			_, err := tx.ExecContext(ctx, `
				INSERT INTO variant_option_values (variant_id, attribute_id, option_id)
				VALUES ($1, $2, $3)
			`, variant.ID.String(), option.AttributeID.String(), option.OptionID.String())
			if err != nil {
				return fmt.Errorf("insert variant option value: %w", err)
			}
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit create variants tx: %w", err)
	}
	return nil
}

func (r *ProductRepo) ListVariantsByProduct(ctx context.Context, productID uuid.UUID) ([]domain.ProductVariant, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id::text, product_id::text, sku, price_minor, stock_qty, image_url, status, combination_key, created_at, updated_at
		FROM product_variants
		WHERE product_id = $1
		ORDER BY combination_key
	`, productID.String())
	if err != nil {
		return nil, fmt.Errorf("list variants by product: %w", err)
	}
	defer rows.Close()

	variants := make([]domain.ProductVariant, 0)
	variantByID := map[string]*domain.ProductVariant{}
	for rows.Next() {
		variant, err := scanVariant(rows)
		if err != nil {
			return nil, err
		}
		variants = append(variants, variant)
		variantByID[variant.ID.String()] = &variants[len(variants)-1]
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	optionRows, err := r.db.QueryContext(ctx, `
		SELECT vov.variant_id::text, vov.attribute_id::text, vov.option_id::text
		FROM variant_option_values vov
		INNER JOIN product_variants pv ON pv.id = vov.variant_id
		WHERE pv.product_id = $1
		ORDER BY vov.variant_id, vov.attribute_id
	`, productID.String())
	if err != nil {
		return nil, fmt.Errorf("list variant options by product: %w", err)
	}
	defer optionRows.Close()

	for optionRows.Next() {
		var (
			variantID   string
			attributeID string
			optionID    string
		)
		if err := optionRows.Scan(&variantID, &attributeID, &optionID); err != nil {
			return nil, fmt.Errorf("scan variant option value: %w", err)
		}
		variant := variantByID[variantID]
		if variant == nil {
			continue
		}
		attr, err := uuid.Parse(attributeID)
		if err != nil {
			return nil, fmt.Errorf("parse variant option attribute id: %w", err)
		}
		opt, err := uuid.Parse(optionID)
		if err != nil {
			return nil, fmt.Errorf("parse variant option id: %w", err)
		}
		variant.Options = append(variant.Options, domain.VariantOptionValue{
			AttributeID: attr,
			OptionID:    opt,
		})
	}
	return variants, optionRows.Err()
}

func scanProduct(scanner interface {
	Scan(dest ...any) error
}) (domain.Product, error) {
	var (
		p                 domain.Product
		id                string
		primaryCategoryID string
		status            string
	)
	if err := scanner.Scan(
		&id, &p.Name, &p.Slug, &p.Description, &primaryCategoryID, &status, &p.BasePriceMinor, &p.CreatedAt, &p.UpdatedAt,
	); err != nil {
		return domain.Product{}, fmt.Errorf("scan product: %w", err)
	}
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return domain.Product{}, fmt.Errorf("parse product id: %w", err)
	}
	parsedPrimaryCategoryID, err := uuid.Parse(primaryCategoryID)
	if err != nil {
		return domain.Product{}, fmt.Errorf("parse product primary category id: %w", err)
	}
	p.ID = parsedID
	p.PrimaryCategoryID = parsedPrimaryCategoryID
	p.Status = domain.ProductStatus(status)
	p.CategoryIDs = []uuid.UUID{parsedPrimaryCategoryID}
	return p, nil
}

func scanProductAttributeValue(scanner interface {
	Scan(dest ...any) error
}) (domain.ProductAttributeValue, error) {
	var (
		v           domain.ProductAttributeValue
		id          string
		productID   string
		attributeID string
		optionID    sql.NullString
		valueText   sql.NullString
		valueNum    sql.NullFloat64
		valueBool   sql.NullBool
		valueJSON   sql.NullString
	)
	if err := scanner.Scan(
		&id,
		&productID,
		&attributeID,
		&optionID,
		&valueText,
		&valueNum,
		&valueBool,
		&valueJSON,
		&v.CreatedAt,
		&v.UpdatedAt,
	); err != nil {
		return domain.ProductAttributeValue{}, fmt.Errorf("scan product attribute value: %w", err)
	}
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return domain.ProductAttributeValue{}, fmt.Errorf("parse product attribute value id: %w", err)
	}
	parsedProductID, err := uuid.Parse(productID)
	if err != nil {
		return domain.ProductAttributeValue{}, fmt.Errorf("parse product attribute value product id: %w", err)
	}
	parsedAttributeID, err := uuid.Parse(attributeID)
	if err != nil {
		return domain.ProductAttributeValue{}, fmt.Errorf("parse product attribute value attribute id: %w", err)
	}
	v.ID = parsedID
	v.ProductID = parsedProductID
	v.AttributeID = parsedAttributeID
	if optionID.Valid {
		parsedOptionID, err := uuid.Parse(optionID.String)
		if err != nil {
			return domain.ProductAttributeValue{}, fmt.Errorf("parse product attribute value option id: %w", err)
		}
		v.OptionID = &parsedOptionID
	}
	if valueText.Valid {
		v.ValueText = &valueText.String
	}
	if valueNum.Valid {
		v.ValueNumber = &valueNum.Float64
	}
	if valueBool.Valid {
		v.ValueBoolean = &valueBool.Bool
	}
	if valueJSON.Valid {
		raw := []byte(valueJSON.String)
		parsed := json.RawMessage(raw)
		v.ValueJSON = &parsed
	}
	return v, nil
}

func scanVariant(scanner interface {
	Scan(dest ...any) error
}) (domain.ProductVariant, error) {
	var (
		v         domain.ProductVariant
		id        string
		productID string
		status    string
	)
	if err := scanner.Scan(
		&id,
		&productID,
		&v.SKU,
		&v.PriceMinor,
		&v.StockQty,
		&v.ImageURL,
		&status,
		&v.CombinationKey,
		&v.CreatedAt,
		&v.UpdatedAt,
	); err != nil {
		return domain.ProductVariant{}, fmt.Errorf("scan variant: %w", err)
	}
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return domain.ProductVariant{}, fmt.Errorf("parse variant id: %w", err)
	}
	parsedProductID, err := uuid.Parse(productID)
	if err != nil {
		return domain.ProductVariant{}, fmt.Errorf("parse variant product id: %w", err)
	}
	v.ID = parsedID
	v.ProductID = parsedProductID
	v.Status = domain.ProductVariantStatus(status)
	v.Options = make([]domain.VariantOptionValue, 0)
	return v, nil
}
