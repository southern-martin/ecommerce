ALTER TABLE products
  DROP COLUMN IF EXISTS product_type,
  DROP COLUMN IF EXISTS stock_quantity;
