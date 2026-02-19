-- EXTENSION
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- CATEGORIES
CREATE TABLE categories (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

  name VARCHAR(255) UNIQUE NOT NULL,
  description TEXT,

  created_at TIMESTAMP NOT NULL DEFAULT now(),
  updated_at TIMESTAMP NOT NULL DEFAULT now(),
  deleted_at TIMESTAMP
);

-- PRODUCTS
CREATE TABLE products (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

  name VARCHAR(255) UNIQUE NOT NULL,
  price DECIMAL(12,2) NOT NULL CHECK (price >= 0),
  stock INT NOT NULL CHECK (stock >= 0),

  category_id UUID NOT NULL,

  created_at TIMESTAMP NOT NULL DEFAULT now(),
  updated_at TIMESTAMP NOT NULL DEFAULT now(),
  deleted_at TIMESTAMP,

  CONSTRAINT fk_products_category
    FOREIGN KEY (category_id)
    REFERENCES categories(id)
    ON DELETE RESTRICT
);

-- TRANSACTIONS
CREATE TABLE transactions (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

  total_amount DECIMAL(14,2) NOT NULL CHECK (total_amount >= 0),

  status VARCHAR(20) NOT NULL,

  created_at TIMESTAMP NOT NULL DEFAULT now(),

  CONSTRAINT chk_transaction_status
    CHECK (status IN ('pending','paid','cancelled','completed'))
);

-- TRANSACTION DETAILS
CREATE TABLE transaction_details (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

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

-- INDEXES

CREATE INDEX idx_products_category
ON products(category_id);

CREATE INDEX idx_td_product
ON transaction_details(product_id);

CREATE INDEX idx_td_transaction
ON transaction_details(transaction_id);

CREATE INDEX idx_transactions_created
ON transactions(created_at);

-- Soft delete optimization
CREATE INDEX idx_products_active
ON products(id)
WHERE deleted_at IS NULL;

CREATE INDEX idx_categories_active
ON categories(id)
WHERE deleted_at IS NULL;

-- Composite index
CREATE INDEX idx_td_transaction_product
ON transaction_details(transaction_id, product_id);
