CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Tax zones
CREATE TABLE tax_zones (
    id UUID PRIMARY KEY,
    country_code VARCHAR(2) NOT NULL,
    state_code VARCHAR(10),
    name VARCHAR(100) NOT NULL
);

CREATE INDEX idx_tax_zones_country_state ON tax_zones(country_code, state_code);

-- Tax rules
CREATE TABLE tax_rules (
    id UUID PRIMARY KEY,
    zone_id UUID NOT NULL,
    tax_name VARCHAR(50) NOT NULL,
    rate DECIMAL(10,6) NOT NULL,
    category VARCHAR(100),
    inclusive BOOLEAN DEFAULT false,
    starts_at TIMESTAMPTZ NOT NULL,
    expires_at TIMESTAMPTZ,
    is_active BOOLEAN DEFAULT true,
    CONSTRAINT fk_tax_rules_zone FOREIGN KEY (zone_id) REFERENCES tax_zones(id) ON DELETE CASCADE
);

CREATE INDEX idx_tax_rules_zone_id ON tax_rules(zone_id);
CREATE INDEX idx_tax_rules_is_active ON tax_rules(is_active);
CREATE INDEX idx_tax_rules_zone_category ON tax_rules(zone_id, category);
