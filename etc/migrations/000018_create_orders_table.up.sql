CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    buyer_id UUID NOT NULL,
    store_id UUID NOT NULL,
    order_number VARCHAR(50) NOT NULL,
    total_amount DECIMAL(14,2) NOT NULL CHECK (total_amount >= 0),
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    shipping_address JSONB NOT NULL,
    shipping_cost DECIMAL(12,2) DEFAULT 0.00 CHECK (shipping_cost >= 0),
    payment_method VARCHAR(50),
    payment_ref VARCHAR(255),
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),

    CONSTRAINT fk_orders_buyer
        FOREIGN KEY (buyer_id)
        REFERENCES users(id)
        ON DELETE RESTRICT,

    CONSTRAINT fk_orders_store
        FOREIGN KEY (store_id)
        REFERENCES stores(id)
        ON DELETE RESTRICT,

    CONSTRAINT uq_orders_order_number
        UNIQUE (order_number),

    CONSTRAINT chk_orders_status
        CHECK (status IN ('pending','paid','processing','shipped','delivered','canceled','refunded'))
);
