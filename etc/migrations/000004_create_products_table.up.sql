CREATE TABLE products (
  id UUID PRIMARY KEY DEFAULT uuidv7(),

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