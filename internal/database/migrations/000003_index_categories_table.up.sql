CREATE INDEX idx_categories_active
ON categories(id)
WHERE deleted_at IS NULL;


CREATE INDEX idx_categories_name
ON categories(name);