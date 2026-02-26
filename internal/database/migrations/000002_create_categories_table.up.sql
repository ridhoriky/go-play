CREATE TABLE categories (
  id UUID PRIMARY KEY DEFAULT uuidv7(),

  name VARCHAR(255) UNIQUE NOT NULL,
  description TEXT,

  created_at TIMESTAMP NOT NULL DEFAULT now(),
  updated_at TIMESTAMP NOT NULL DEFAULT now(),
  deleted_at TIMESTAMP
);