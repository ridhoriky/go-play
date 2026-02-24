CREATE TABLE transaction_details (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

  product_id UUID NOT NULL,
  transaction_id UUID NOT NULL,

  product_name VARCHAR(255) NOT NULL,

  quantity INT NOT NULL CHECK (quantity > 0),

  price DECIMAL(12,2) NOT NULL CHECK (price >= 0),

  subtotal DECIMAL(14,2) NOT NULL CHECK (subtotal >= 0),

  CONSTRAINT fk_td_product
    FOREIGN KEY (product_id)
    REFERENCES products(id)
    ON DELETE RESTRICT,

  CONSTRAINT fk_td_transaction
    FOREIGN KEY (transaction_id)
    REFERENCES transactions(id)
    ON DELETE CASCADE
);