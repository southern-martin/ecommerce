CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 1. Create attribute_option_values table
CREATE TABLE IF NOT EXISTS attribute_option_values (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    attribute_id UUID NOT NULL REFERENCES attribute_definitions(id) ON DELETE CASCADE,
    value VARCHAR(255) NOT NULL,
    color_hex VARCHAR(7),
    sort_order INTEGER NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_attribute_option_values_attribute_id ON attribute_option_values(attribute_id);

-- 2. Migrate existing TEXT[] options into the new table
INSERT INTO attribute_option_values (id, attribute_id, value, sort_order, created_at)
SELECT
    uuid_generate_v4(),
    ad.id,
    opt_val,
    (ord - 1),
    NOW()
FROM attribute_definitions ad,
     unnest(ad.options) WITH ORDINALITY AS t(opt_val, ord)
WHERE ad.options IS NOT NULL AND array_length(ad.options, 1) > 0
ON CONFLICT DO NOTHING;

-- 3. Add option_value_id and option_value_ids to product_attribute_values
ALTER TABLE product_attribute_values ADD COLUMN IF NOT EXISTS option_value_id UUID;
ALTER TABLE product_attribute_values ADD COLUMN IF NOT EXISTS option_value_ids UUID[];

-- 4. Backfill option_value_id for single-select product attribute values
UPDATE product_attribute_values pav
SET option_value_id = aov.id
FROM attribute_option_values aov
WHERE pav.attribute_id = aov.attribute_id
  AND pav.value = aov.value
  AND pav.value IS NOT NULL AND pav.value != ''
  AND pav.option_value_id IS NULL;

-- 5. Backfill option_value_ids for multi-select product attribute values
UPDATE product_attribute_values pav
SET option_value_ids = (
    SELECT array_agg(aov.id ORDER BY aov.sort_order)
    FROM attribute_option_values aov
    WHERE aov.attribute_id = pav.attribute_id
      AND aov.value = ANY(pav."values")
)
WHERE pav."values" IS NOT NULL AND array_length(pav."values", 1) > 0
  AND pav.option_value_ids IS NULL;

-- 6. Drop the old options column from attribute_definitions
ALTER TABLE attribute_definitions DROP COLUMN IF EXISTS options;
