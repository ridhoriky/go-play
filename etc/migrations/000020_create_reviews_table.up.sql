CREATE TABLE reviews (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    product_id UUID NOT NULL,
    buyer_id UUID NOT NULL,
    order_id UUID NOT NULL,
    rating SMALLINT NOT NULL,
    comment TEXT,
    images JSONB DEFAULT '[]',
    seller_reply TEXT,
    seller_replied_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),

    CONSTRAINT fk_reviews_product
        FOREIGN KEY (product_id)
        REFERENCES products(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_reviews_buyer
        FOREIGN KEY (buyer_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_reviews_order
        FOREIGN KEY (order_id)
        REFERENCES orders(id)
        ON DELETE CASCADE,

    CONSTRAINT chk_reviews_rating
        CHECK (rating >= 1 AND rating <= 5),

    CONSTRAINT uq_reviews_product_buyer_order
        UNIQUE (product_id, buyer_id, order_id)
);
