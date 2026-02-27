CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Search indices
CREATE TABLE search_indices (
    id UUID PRIMARY KEY,
    product_id VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(500) NOT NULL,
    slug VARCHAR(500),
    description TEXT,
    price_cents BIGINT NOT NULL DEFAULT 0,
    currency VARCHAR(3) DEFAULT 'USD',
    category_id VARCHAR(255),
    seller_id VARCHAR(255),
    image_url TEXT,
    rating DECIMAL(3,2) DEFAULT 0,
    review_count INTEGER DEFAULT 0,
    in_stock BOOLEAN DEFAULT true,
    tags TEXT[],
    attributes JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_search_indices_category_id ON search_indices(category_id);
CREATE INDEX idx_search_indices_seller_id ON search_indices(seller_id);
CREATE INDEX idx_search_indices_in_stock ON search_indices(in_stock);
