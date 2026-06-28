ALTER TABLE products DROP CONSTRAINT IF EXISTS fk_products_store;
ALTER TABLE products DROP CONSTRAINT IF EXISTS uq_products_slug;

ALTER TABLE products DROP COLUMN IF EXISTS store_id;
ALTER TABLE products DROP COLUMN IF EXISTS description;
ALTER TABLE products DROP COLUMN IF EXISTS slug;
ALTER TABLE products DROP COLUMN IF EXISTS rating_avg;
ALTER TABLE products DROP COLUMN IF EXISTS total_sold;
ALTER TABLE products DROP COLUMN IF EXISTS is_active;

DROP TRIGGER IF EXISTS trigger_update_rating_after_review ON reviews;
DROP FUNCTION IF EXISTS update_rating_avg();
