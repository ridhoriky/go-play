DROP INDEX IF EXISTS products_search_vector_idx;
DROP TRIGGER IF EXISTS products_search_vector_update_trigger ON products;
DROP FUNCTION IF EXISTS products_search_vector_update();
ALTER TABLE products DROP COLUMN IF EXISTS search_vector;
