CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Coupons
CREATE TABLE coupons (
    id UUID PRIMARY KEY,
    code VARCHAR(100) UNIQUE NOT NULL,
    type VARCHAR(20) NOT NULL,
    discount_value BIGINT NOT NULL DEFAULT 0,
    min_order_cents BIGINT NOT NULL DEFAULT 0,
    max_discount_cents BIGINT NOT NULL DEFAULT 0,
    usage_limit INTEGER NOT NULL DEFAULT 0,
    usage_count INTEGER NOT NULL DEFAULT 0,
    per_user_limit INTEGER NOT NULL DEFAULT 0,
    scope VARCHAR(20) NOT NULL DEFAULT 'all',
    scope_ids TEXT[],
    created_by VARCHAR(255) NOT NULL,
    starts_at TIMESTAMPTZ NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Coupon usages
CREATE TABLE coupon_usages (
    id UUID PRIMARY KEY,
    coupon_id UUID NOT NULL,
    user_id UUID NOT NULL,
    order_id UUID NOT NULL,
    discount_cents BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_coupon_usages_coupon_id ON coupon_usages(coupon_id);
CREATE INDEX idx_coupon_usages_user_id ON coupon_usages(user_id);

-- Flash sales
CREATE TABLE flash_sales (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    starts_at TIMESTAMPTZ NOT NULL,
    ends_at TIMESTAMPTZ NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Flash sale items
CREATE TABLE flash_sale_items (
    id UUID PRIMARY KEY,
    flash_sale_id UUID NOT NULL,
    product_id UUID NOT NULL,
    variant_id UUID,
    sale_price_cents BIGINT NOT NULL DEFAULT 0,
    quantity_limit INTEGER NOT NULL DEFAULT 0,
    sold_count INTEGER NOT NULL DEFAULT 0,
    CONSTRAINT fk_flash_sale_items_flash_sale FOREIGN KEY (flash_sale_id) REFERENCES flash_sales(id) ON DELETE CASCADE
);

CREATE INDEX idx_flash_sale_items_flash_sale_id ON flash_sale_items(flash_sale_id);

-- Bundles
CREATE TABLE bundles (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    seller_id UUID NOT NULL,
    product_ids TEXT[],
    bundle_price_cents BIGINT NOT NULL DEFAULT 0,
    savings_cents BIGINT NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_bundles_seller_id ON bundles(seller_id);
