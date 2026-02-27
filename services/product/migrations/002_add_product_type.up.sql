ALTER TABLE products
  ADD COLUMN product_type VARCHAR(20) NOT NULL DEFAULT 'simple',
  ADD COLUMN stock_quantity INTEGER NOT NULL DEFAULT 0;

-- Set existing products with variants to 'configurable'
UPDATE products SET product_type = 'configurable' WHERE has_variants = true;
