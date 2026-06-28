CREATE TABLE wishlists (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    buyer_id UUID NOT NULL,
    product_id UUID NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now(),

    CONSTRAINT fk_wishlists_buyer
        FOREIGN KEY (buyer_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_wishlists_product
        FOREIGN KEY (product_id)
        REFERENCES products(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_wishlists_buyer_product
        UNIQUE (buyer_id, product_id)
);
