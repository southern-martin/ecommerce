CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Affiliate programs
CREATE TABLE affiliate_programs (
    id UUID PRIMARY KEY,
    commission_rate DECIMAL(5,4) NOT NULL DEFAULT 0.05,
    min_payout_cents BIGINT NOT NULL DEFAULT 5000,
    cookie_days INTEGER NOT NULL DEFAULT 30,
    referrer_bonus_cents BIGINT NOT NULL DEFAULT 0,
    referred_bonus_cents BIGINT NOT NULL DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Affiliate links
CREATE TABLE affiliate_links (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    code VARCHAR(8) UNIQUE NOT NULL,
    target_url TEXT NOT NULL,
    click_count BIGINT DEFAULT 0,
    conversion_count BIGINT DEFAULT 0,
    total_earnings_cents BIGINT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_affiliate_links_user_id ON affiliate_links(user_id);

-- Referrals
CREATE TABLE referrals (
    id UUID PRIMARY KEY,
    referrer_id UUID NOT NULL,
    referred_id UUID NOT NULL,
    order_id UUID NOT NULL,
    order_total_cents BIGINT NOT NULL DEFAULT 0,
    commission_cents BIGINT NOT NULL DEFAULT 0,
    status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_referrals_referrer_id ON referrals(referrer_id);
CREATE INDEX idx_referrals_referred_id ON referrals(referred_id);
CREATE INDEX idx_referrals_order_id ON referrals(order_id);

-- Affiliate payouts
CREATE TABLE affiliate_payouts (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    amount_cents BIGINT NOT NULL DEFAULT 0,
    status VARCHAR(20) DEFAULT 'requested',
    payout_method VARCHAR(20) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    completed_at TIMESTAMPTZ
);

CREATE INDEX idx_affiliate_payouts_user_id ON affiliate_payouts(user_id);
