ALTER TABLE products DROP CONSTRAINT IF EXISTS fk_products_attribute_group;
DROP INDEX IF EXISTS idx_products_attribute_group_id;
ALTER TABLE products DROP COLUMN IF EXISTS attribute_group_id;
DROP TABLE IF EXISTS attribute_group_items;
DROP TABLE IF EXISTS attribute_groups;
