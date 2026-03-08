CREATE INDEX idx_products_category
ON products(category_id);

CREATE INDEX idx_products_active
ON products(id)
WHERE deleted_at IS NULL;

CREATE INDEX idx_products_name
ON products(name);