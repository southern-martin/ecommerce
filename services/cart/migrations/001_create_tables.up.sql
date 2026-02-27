-- Cart service: durable cart persistence
-- Redis remains the fast-read cache; this table is the permanent backing store.

CREATE TABLE IF NOT EXISTS carts (
    user_id     VARCHAR(36) PRIMARY KEY,
    cart_data   JSONB       NOT NULL DEFAULT '{}',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index for finding recently-updated carts (cleanup jobs, analytics)
CREATE INDEX IF NOT EXISTS idx_carts_updated_at ON carts (updated_at);
