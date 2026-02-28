-- Restore options TEXT[] column on attribute_definitions
ALTER TABLE attribute_definitions ADD COLUMN IF NOT EXISTS options TEXT[];

-- Backfill options from attribute_option_values
UPDATE attribute_definitions ad
SET options = (
    SELECT array_agg(aov.value ORDER BY aov.sort_order)
    FROM attribute_option_values aov
    WHERE aov.attribute_id = ad.id
);

-- Remove new columns from product_attribute_values
ALTER TABLE product_attribute_values DROP COLUMN IF EXISTS option_value_id;
ALTER TABLE product_attribute_values DROP COLUMN IF EXISTS option_value_ids;

-- Drop the new table
DROP TABLE IF EXISTS attribute_option_values;
