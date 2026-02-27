CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Tiers
CREATE TABLE tiers (
    name VARCHAR(20) PRIMARY KEY,
    min_points BIGINT NOT NULL DEFAULT 0,
    cashback_rate DOUBLE PRECISION NOT NULL DEFAULT 0,
    points_multiplier DOUBLE PRECISION NOT NULL DEFAULT 1,
    free_shipping BOOLEAN NOT NULL DEFAULT false,
    priority_support_hours INTEGER NOT NULL DEFAULT 48
);

-- Memberships
CREATE TABLE memberships (
    user_id VARCHAR(36) PRIMARY KEY,
    tier VARCHAR(20) NOT NULL DEFAULT 'bronze',
    points_balance BIGINT NOT NULL DEFAULT 0,
    lifetime_points BIGINT NOT NULL DEFAULT 0,
    tier_expires_at TIMESTAMPTZ,
    joined_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_memberships_tier_expires_at ON memberships(tier_expires_at);

-- Points transactions
CREATE TABLE points_transactions (
    id UUID PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    type VARCHAR(20) NOT NULL,
    points BIGINT NOT NULL,
    source VARCHAR(20) NOT NULL,
    reference_id VARCHAR(255),
    description TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_points_transactions_user_id ON points_transactions(user_id);

-- Insert default tiers
INSERT INTO tiers (name, min_points, cashback_rate, points_multiplier, free_shipping, priority_support_hours) VALUES
('bronze', 0, 0, 1, false, 48),
('silver', 1000, 0.02, 1.5, false, 24),
('gold', 5000, 0.05, 2, true, 12),
('platinum', 20000, 0.1, 3, true, 4);
