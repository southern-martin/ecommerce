-- Enable uuid_generate_v4 if not already available
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Attribute groups: reusable bundles of attributes (e.g., "Clothing Specs", "Electronics")
CREATE TABLE IF NOT EXISTS attribute_groups (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(300) UNIQUE NOT NULL,
    description TEXT,
    sort_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Join table: which attributes belong to which group
CREATE TABLE IF NOT EXISTS attribute_group_items (
    group_id UUID NOT NULL,
    attribute_id UUID NOT NULL,
    sort_order INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (group_id, attribute_id),
    CONSTRAINT fk_agi_group FOREIGN KEY (group_id) REFERENCES attribute_groups(id) ON DELETE CASCADE,
    CONSTRAINT fk_agi_attribute FOREIGN KEY (attribute_id) REFERENCES attribute_definitions(id) ON DELETE CASCADE
);

-- Products get an optional attribute_group_id
ALTER TABLE products ADD COLUMN IF NOT EXISTS attribute_group_id UUID;
ALTER TABLE products ADD CONSTRAINT fk_products_attribute_group
    FOREIGN KEY (attribute_group_id) REFERENCES attribute_groups(id) ON DELETE SET NULL;
CREATE INDEX IF NOT EXISTS idx_products_attribute_group_id ON products(attribute_group_id);

-- Data migration: create attribute groups from existing category-attribute bindings
INSERT INTO attribute_groups (id, name, slug, description, sort_order, created_at, updated_at)
SELECT
    uuid_generate_v4(),
    c.name || ' Attributes',
    c.slug || '-attributes',
    'Auto-migrated from category: ' || c.name,
    0,
    NOW(),
    NOW()
FROM categories c
WHERE EXISTS (SELECT 1 FROM category_attributes ca WHERE ca.category_id = c.id)
ON CONFLICT (slug) DO NOTHING;

-- Copy category-attribute bindings into attribute_group_items
INSERT INTO attribute_group_items (group_id, attribute_id, sort_order)
SELECT ag.id, ca.attribute_id, ca.sort_order
FROM category_attributes ca
JOIN categories c ON c.id = ca.category_id
JOIN attribute_groups ag ON ag.slug = c.slug || '-attributes'
ON CONFLICT DO NOTHING;

-- Set attribute_group_id on products that had a category with attributes
UPDATE products p
SET attribute_group_id = ag.id
FROM categories c
JOIN attribute_groups ag ON ag.slug = c.slug || '-attributes'
WHERE p.category_id = c.id
AND p.attribute_group_id IS NULL
AND EXISTS (SELECT 1 FROM category_attributes ca WHERE ca.category_id = c.id);
