CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Carriers
CREATE TABLE carriers (
    code VARCHAR(20) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    is_active BOOLEAN DEFAULT true,
    supported_countries TEXT[],
    api_base_url VARCHAR(500),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Carrier credentials (per seller)
CREATE TABLE carrier_credentials (
    id UUID PRIMARY KEY,
    seller_id UUID NOT NULL,
    carrier_code VARCHAR(20) NOT NULL,
    credentials JSONB NOT NULL,
    is_active BOOLEAN DEFAULT true,
    CONSTRAINT fk_carrier_credentials_carrier FOREIGN KEY (carrier_code) REFERENCES carriers(code) ON DELETE CASCADE,
    CONSTRAINT uq_carrier_credentials_seller_carrier UNIQUE (seller_id, carrier_code)
);

-- Shipments
CREATE TABLE shipments (
    id UUID PRIMARY KEY,
    order_id UUID NOT NULL,
    seller_id UUID NOT NULL,
    carrier_code VARCHAR(20),
    service_code VARCHAR(50),
    tracking_number VARCHAR(100),
    label_url TEXT,
    status VARCHAR(20) DEFAULT 'pending',
    origin JSONB NOT NULL,
    destination JSONB NOT NULL,
    weight_grams INTEGER DEFAULT 0,
    rate_cents BIGINT DEFAULT 0,
    currency VARCHAR(3) DEFAULT 'USD',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_shipments_order_id ON shipments(order_id);
CREATE INDEX idx_shipments_seller_id ON shipments(seller_id);
CREATE INDEX idx_shipments_tracking_number ON shipments(tracking_number);

-- Shipment items
CREATE TABLE shipment_items (
    id UUID PRIMARY KEY,
    shipment_id UUID NOT NULL,
    product_id UUID NOT NULL,
    variant_id UUID,
    product_name VARCHAR(255),
    quantity INTEGER NOT NULL DEFAULT 1,
    CONSTRAINT fk_shipment_items_shipment FOREIGN KEY (shipment_id) REFERENCES shipments(id) ON DELETE CASCADE
);

CREATE INDEX idx_shipment_items_shipment_id ON shipment_items(shipment_id);

-- Tracking events
CREATE TABLE tracking_events (
    id UUID PRIMARY KEY,
    shipment_id UUID NOT NULL,
    status VARCHAR(50) NOT NULL,
    description TEXT,
    location VARCHAR(200),
    event_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT fk_tracking_events_shipment FOREIGN KEY (shipment_id) REFERENCES shipments(id) ON DELETE CASCADE
);

CREATE INDEX idx_tracking_events_shipment_id ON tracking_events(shipment_id);
