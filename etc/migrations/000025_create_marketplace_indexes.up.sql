-- Products
CREATE INDEX idx_products_store_id ON products(store_id);
CREATE INDEX idx_products_category_id ON products(category_id);
CREATE INDEX idx_products_price ON products(price);
CREATE INDEX idx_products_rating ON products(rating_avg DESC);
CREATE INDEX idx_products_slug ON products(slug);
CREATE INDEX idx_products_is_active ON products(is_active) WHERE deleted_at IS NULL;

-- Orders
CREATE INDEX idx_orders_buyer_id ON orders(buyer_id);
CREATE INDEX idx_orders_store_id ON orders(store_id);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_created ON orders(created_at DESC);

-- Order Items
CREATE INDEX idx_order_items_order_id ON order_items(order_id);
CREATE INDEX idx_order_items_product_id ON order_items(product_id);

-- Reviews
CREATE INDEX idx_reviews_product_id ON reviews(product_id);
CREATE INDEX idx_reviews_buyer_id ON reviews(buyer_id);
CREATE INDEX idx_reviews_order_id ON reviews(order_id);

-- Wishlists
CREATE INDEX idx_wishlists_buyer_id ON wishlists(buyer_id);
CREATE INDEX idx_wishlists_product_id ON wishlists(product_id);

-- Carts
CREATE INDEX idx_carts_buyer_id ON carts(buyer_id);

-- Stores (slug & user_id already indexed via UNIQUE constraints)

-- Product Images
CREATE INDEX idx_product_images_product_id ON product_images(product_id);
