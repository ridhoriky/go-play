CREATE TABLE product_images (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    product_id UUID NOT NULL,
    url VARCHAR(512) NOT NULL,
    alt_text VARCHAR(255),
    sort_order INT DEFAULT 0,
    is_primary BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT now(),

    CONSTRAINT fk_product_images_product
        FOREIGN KEY (product_id)
        REFERENCES products(id)
        ON DELETE CASCADE
);
