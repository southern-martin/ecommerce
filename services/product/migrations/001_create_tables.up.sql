CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Categories
CREATE TABLE categories (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(300) UNIQUE NOT NULL,
    parent_id UUID,
    sort_order INTEGER NOT NULL DEFAULT 0,
    image_url TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_categories_parent FOREIGN KEY (parent_id) REFERENCES categories(id) ON DELETE SET NULL
);

CREATE INDEX idx_categories_parent_id ON categories(parent_id);

-- Attribute definitions
CREATE TABLE attribute_definitions (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(300) UNIQUE NOT NULL,
    type VARCHAR(20) NOT NULL,
    required BOOLEAN NOT NULL DEFAULT false,
    filterable BOOLEAN NOT NULL DEFAULT false,
    options TEXT[],
    unit VARCHAR(50),
    sort_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Category-attribute join table
CREATE TABLE category_attributes (
    category_id UUID NOT NULL,
    attribute_id UUID NOT NULL,
    sort_order INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (category_id, attribute_id),
    CONSTRAINT fk_category_attributes_category FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE,
    CONSTRAINT fk_category_attributes_attribute FOREIGN KEY (attribute_id) REFERENCES attribute_definitions(id) ON DELETE CASCADE
);

-- Products
CREATE TABLE products (
    id UUID PRIMARY KEY,
    seller_id UUID NOT NULL,
    category_id UUID,
    name VARCHAR(500) NOT NULL,
    slug VARCHAR(600) UNIQUE NOT NULL,
    description TEXT,
    base_price_cents BIGINT NOT NULL DEFAULT 0,
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    status VARCHAR(20) NOT NULL DEFAULT 'draft',
    has_variants BOOLEAN NOT NULL DEFAULT false,
    tags TEXT[],
    image_urls TEXT[],
    rating_avg DOUBLE PRECISION NOT NULL DEFAULT 0,
    rating_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_products_category FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL
);

CREATE INDEX idx_products_seller_id ON products(seller_id);
CREATE INDEX idx_products_category_id ON products(category_id);
CREATE INDEX idx_products_status ON products(status);

-- Product attribute values
CREATE TABLE product_attribute_values (
    id UUID PRIMARY KEY,
    product_id UUID NOT NULL,
    attribute_id UUID NOT NULL,
    attribute_name VARCHAR(255),
    value TEXT,
    "values" TEXT[],
    CONSTRAINT fk_product_attribute_values_product FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
);

CREATE INDEX idx_product_attribute_values_product_id ON product_attribute_values(product_id);

-- Product options
CREATE TABLE product_options (
    id UUID PRIMARY KEY,
    product_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    sort_order INTEGER NOT NULL DEFAULT 0,
    CONSTRAINT fk_product_options_product FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
);

CREATE INDEX idx_product_options_product_id ON product_options(product_id);

-- Product option values
CREATE TABLE product_option_values (
    id UUID PRIMARY KEY,
    option_id UUID NOT NULL,
    value VARCHAR(255) NOT NULL,
    color_hex VARCHAR(7),
    sort_order INTEGER NOT NULL DEFAULT 0,
    CONSTRAINT fk_product_option_values_option FOREIGN KEY (option_id) REFERENCES product_options(id) ON DELETE CASCADE
);

CREATE INDEX idx_product_option_values_option_id ON product_option_values(option_id);

-- Product variants
CREATE TABLE product_variants (
    id UUID PRIMARY KEY,
    product_id UUID NOT NULL,
    sku VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(500),
    price_cents BIGINT NOT NULL DEFAULT 0,
    compare_at_cents BIGINT NOT NULL DEFAULT 0,
    cost_cents BIGINT NOT NULL DEFAULT 0,
    stock INTEGER NOT NULL DEFAULT 0,
    low_stock_alert INTEGER NOT NULL DEFAULT 0,
    weight_grams INTEGER NOT NULL DEFAULT 0,
    is_default BOOLEAN NOT NULL DEFAULT false,
    is_active BOOLEAN NOT NULL DEFAULT true,
    image_urls TEXT[],
    barcode VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_product_variants_product FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
);

CREATE INDEX idx_product_variants_product_id ON product_variants(product_id);

-- Variant option values (which options a variant has selected)
CREATE TABLE variant_option_values (
    variant_id UUID NOT NULL,
    option_id UUID NOT NULL,
    option_value_id UUID NOT NULL,
    option_name VARCHAR(255),
    value VARCHAR(255),
    PRIMARY KEY (variant_id, option_id),
    CONSTRAINT fk_variant_option_values_variant FOREIGN KEY (variant_id) REFERENCES product_variants(id) ON DELETE CASCADE
);
