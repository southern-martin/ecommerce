CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Payments
CREATE TABLE payments (
    id VARCHAR(36) PRIMARY KEY,
    order_id VARCHAR(36) NOT NULL,
    buyer_id VARCHAR(36) NOT NULL,
    amount_cents BIGINT NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'usd',
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    method VARCHAR(20) NOT NULL DEFAULT 'card',
    stripe_payment_id VARCHAR(255),
    failure_reason TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_payments_order_id ON payments(order_id);
CREATE INDEX idx_payments_buyer_id ON payments(buyer_id);
CREATE INDEX idx_payments_status ON payments(status);
CREATE INDEX idx_payments_stripe_payment_id ON payments(stripe_payment_id);

-- Seller wallets
CREATE TABLE seller_wallets (
    seller_id VARCHAR(36) PRIMARY KEY,
    available_balance BIGINT NOT NULL DEFAULT 0,
    pending_balance BIGINT NOT NULL DEFAULT 0,
    currency VARCHAR(3) NOT NULL DEFAULT 'usd',
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Wallet transactions
CREATE TABLE wallet_transactions (
    id VARCHAR(36) PRIMARY KEY,
    seller_id VARCHAR(36) NOT NULL,
    type VARCHAR(30) NOT NULL,
    amount_cents BIGINT NOT NULL,
    reference_type VARCHAR(20),
    reference_id VARCHAR(36),
    description TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_wallet_transactions_seller_id ON wallet_transactions(seller_id);
CREATE INDEX idx_wallet_transactions_reference_id ON wallet_transactions(reference_id);

-- Payouts
CREATE TABLE payouts (
    id VARCHAR(36) PRIMARY KEY,
    seller_id VARCHAR(36) NOT NULL,
    amount_cents BIGINT NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'usd',
    method VARCHAR(30) NOT NULL,
    stripe_transfer_id VARCHAR(255),
    status VARCHAR(20) NOT NULL DEFAULT 'requested',
    requested_at TIMESTAMPTZ DEFAULT NOW(),
    completed_at TIMESTAMPTZ
);

CREATE INDEX idx_payouts_seller_id ON payouts(seller_id);
CREATE INDEX idx_payouts_status ON payouts(status);
