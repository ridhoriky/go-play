DROP INDEX IF EXISTS idx_products_store_id;
DROP INDEX IF EXISTS idx_products_category_id;
DROP INDEX IF EXISTS idx_products_price;
DROP INDEX IF EXISTS idx_products_rating;
DROP INDEX IF EXISTS idx_products_slug;
DROP INDEX IF EXISTS idx_products_is_active;

DROP INDEX IF EXISTS idx_orders_buyer_id;
DROP INDEX IF EXISTS idx_orders_store_id;
DROP INDEX IF EXISTS idx_orders_status;
DROP INDEX IF EXISTS idx_orders_created;

DROP INDEX IF EXISTS idx_order_items_order_id;
DROP INDEX IF EXISTS idx_order_items_product_id;

DROP INDEX IF EXISTS idx_reviews_product_id;
DROP INDEX IF EXISTS idx_reviews_buyer_id;
DROP INDEX IF EXISTS idx_reviews_order_id;

DROP INDEX IF EXISTS idx_wishlists_buyer_id;
DROP INDEX IF EXISTS idx_wishlists_product_id;

DROP INDEX IF EXISTS idx_carts_buyer_id;

-- Stores: slug & user_id indexes removed (redundant with UNIQUE constraints)

DROP INDEX IF EXISTS idx_product_images_product_id;
