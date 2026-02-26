CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    slug TEXT NOT NULL,
    parent_id UUID NULL REFERENCES categories(id) ON DELETE RESTRICT,
    level INT NOT NULL DEFAULT 0 CHECK (level >= 0),
    path TEXT NOT NULL CHECK (path LIKE '/%'),
    sort_order INT NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (parent_id, slug),
    CHECK (parent_id IS NULL OR parent_id <> id)
);

CREATE INDEX idx_categories_parent_id ON categories(parent_id);
CREATE INDEX idx_categories_path ON categories(path);

CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL DEFAULT '',
    primary_category_id UUID NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
    status TEXT NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'active', 'archived')),
    base_price_minor BIGINT NOT NULL DEFAULT 0 CHECK (base_price_minor >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_products_primary_category_id ON products(primary_category_id);

CREATE TABLE product_categories (
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    category_id UUID NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
    is_primary BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (product_id, category_id)
);

CREATE UNIQUE INDEX ux_product_primary_category
    ON product_categories(product_id)
    WHERE is_primary = true;

CREATE INDEX idx_product_categories_category_id ON product_categories(category_id);

CREATE TABLE category_attributes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    category_id UUID NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    code TEXT NOT NULL,
    attribute_type TEXT NOT NULL CHECK (attribute_type IN ('select', 'text', 'number', 'bool')),
    required BOOLEAN NOT NULL DEFAULT false,
    is_variant_axis BOOLEAN NOT NULL DEFAULT false,
    is_filterable BOOLEAN NOT NULL DEFAULT false,
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (category_id, code)
);

CREATE INDEX idx_category_attributes_category_id ON category_attributes(category_id);

CREATE TABLE attribute_options (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    attribute_id UUID NOT NULL REFERENCES category_attributes(id) ON DELETE CASCADE,
    value TEXT NOT NULL,
    label TEXT NOT NULL,
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (attribute_id, value),
    UNIQUE (id, attribute_id)
);

CREATE INDEX idx_attribute_options_attribute_id ON attribute_options(attribute_id);

CREATE TABLE product_attribute_values (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    attribute_id UUID NOT NULL REFERENCES category_attributes(id) ON DELETE RESTRICT,
    option_id UUID NULL,
    value_text TEXT NULL,
    value_number NUMERIC(18, 6) NULL,
    value_boolean BOOLEAN NULL,
    value_json JSONB NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (product_id, attribute_id),
    CHECK (
        ((option_id IS NOT NULL)::INT +
         (value_text IS NOT NULL)::INT +
         (value_number IS NOT NULL)::INT +
         (value_boolean IS NOT NULL)::INT +
         (value_json IS NOT NULL)::INT) = 1
    ),
    FOREIGN KEY (option_id, attribute_id)
        REFERENCES attribute_options(id, attribute_id)
        ON DELETE RESTRICT
);

CREATE INDEX idx_pav_product_id ON product_attribute_values(product_id);
CREATE INDEX idx_pav_attribute_id ON product_attribute_values(attribute_id);

CREATE TABLE product_variants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    sku TEXT NOT NULL,
    price_minor BIGINT NOT NULL CHECK (price_minor >= 0),
    stock_qty BIGINT NOT NULL DEFAULT 0 CHECK (stock_qty >= 0),
    image_url TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive')),
    combination_key TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (product_id, sku),
    UNIQUE (product_id, combination_key)
);

CREATE INDEX idx_product_variants_product_id ON product_variants(product_id);

CREATE TABLE variant_option_values (
    variant_id UUID NOT NULL REFERENCES product_variants(id) ON DELETE CASCADE,
    attribute_id UUID NOT NULL REFERENCES category_attributes(id) ON DELETE RESTRICT,
    option_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (variant_id, attribute_id),
    FOREIGN KEY (option_id, attribute_id)
        REFERENCES attribute_options(id, attribute_id)
        ON DELETE RESTRICT
);

CREATE INDEX idx_variant_option_values_option_id ON variant_option_values(option_id);
