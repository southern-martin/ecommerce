CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Returns
CREATE TABLE returns (
    id UUID PRIMARY KEY,
    order_id UUID NOT NULL,
    buyer_id UUID NOT NULL,
    seller_id UUID NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'requested',
    reason VARCHAR(50) NOT NULL,
    description TEXT,
    image_urls TEXT[],
    refund_amount_cents BIGINT DEFAULT 0,
    refund_method VARCHAR(20),
    return_tracking VARCHAR(100),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_returns_order_id ON returns(order_id);
CREATE INDEX idx_returns_buyer_id ON returns(buyer_id);
CREATE INDEX idx_returns_seller_id ON returns(seller_id);

-- Return items
CREATE TABLE return_items (
    id UUID PRIMARY KEY,
    return_id UUID NOT NULL,
    order_item_id UUID NOT NULL,
    product_id UUID NOT NULL,
    variant_id UUID,
    quantity INTEGER NOT NULL DEFAULT 1,
    reason VARCHAR(50),
    CONSTRAINT fk_return_items_return FOREIGN KEY (return_id) REFERENCES returns(id) ON DELETE CASCADE
);

CREATE INDEX idx_return_items_return_id ON return_items(return_id);

-- Disputes
CREATE TABLE disputes (
    id UUID PRIMARY KEY,
    order_id UUID NOT NULL,
    return_id UUID,
    buyer_id UUID NOT NULL,
    seller_id UUID NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'open',
    type VARCHAR(30) NOT NULL,
    description TEXT NOT NULL,
    resolution TEXT,
    resolved_by UUID,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    resolved_at TIMESTAMPTZ
);

CREATE INDEX idx_disputes_order_id ON disputes(order_id);
CREATE INDEX idx_disputes_return_id ON disputes(return_id);
CREATE INDEX idx_disputes_buyer_id ON disputes(buyer_id);
CREATE INDEX idx_disputes_seller_id ON disputes(seller_id);

-- Dispute messages
CREATE TABLE dispute_messages (
    id UUID PRIMARY KEY,
    dispute_id UUID NOT NULL,
    sender_id UUID NOT NULL,
    sender_role VARCHAR(10) NOT NULL,
    message TEXT NOT NULL,
    attachments TEXT[],
    created_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT fk_dispute_messages_dispute FOREIGN KEY (dispute_id) REFERENCES disputes(id) ON DELETE CASCADE
);

CREATE INDEX idx_dispute_messages_dispute_id ON dispute_messages(dispute_id);
