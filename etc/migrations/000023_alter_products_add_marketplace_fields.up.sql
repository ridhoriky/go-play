ALTER TABLE products ADD COLUMN store_id UUID;
ALTER TABLE products ADD COLUMN description TEXT;
ALTER TABLE products ADD COLUMN slug VARCHAR(255);
ALTER TABLE products ADD COLUMN rating_avg DECIMAL(3,2) DEFAULT 0.00;
ALTER TABLE products ADD COLUMN total_sold INT DEFAULT 0;
ALTER TABLE products ADD COLUMN is_active BOOLEAN DEFAULT TRUE;

ALTER TABLE products
    ADD CONSTRAINT fk_products_store
    FOREIGN KEY (store_id)
    REFERENCES stores(id)
    ON DELETE SET NULL;

ALTER TABLE products
    ADD CONSTRAINT uq_products_slug
    UNIQUE (slug);

-- Review rating triggers
CREATE OR REPLACE FUNCTION update_rating_avg()
RETURNS TRIGGER AS $$
DECLARE
    target_product_id UUID;
    target_store_id UUID;
BEGIN
    -- Skip recalculation if rating hasn't changed on update
    IF TG_OP = 'UPDATE' AND OLD.rating = NEW.rating AND OLD.product_id = NEW.product_id THEN
        RETURN NULL;
    END IF;

    IF TG_OP = 'DELETE' THEN
        target_product_id := OLD.product_id;
    ELSE
        target_product_id := NEW.product_id;
    END IF;

    -- Cache store_id
    SELECT store_id INTO target_store_id 
    FROM products 
    WHERE id = target_product_id;

    -- Update product
    UPDATE products 
    SET rating_avg = (
        SELECT COALESCE(AVG(rating), 0) 
        FROM reviews 
        WHERE product_id = target_product_id
    )
    WHERE id = target_product_id;
    
    -- Update store
    IF target_store_id IS NOT NULL THEN
        UPDATE stores
        SET rating_avg = (
            SELECT COALESCE(AVG(r.rating), 0) 
            FROM reviews r
            JOIN products p ON r.product_id = p.id
            WHERE p.store_id = target_store_id
        )
        WHERE id = target_store_id;
    END IF;

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_rating_after_review
AFTER INSERT OR UPDATE OR DELETE ON reviews
FOR EACH ROW
EXECUTE FUNCTION update_rating_avg();
