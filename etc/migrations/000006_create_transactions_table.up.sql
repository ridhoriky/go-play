CREATE TABLE transactions (
  id UUID PRIMARY KEY DEFAULT uuidv7(),

  total_amount DECIMAL(14,2) NOT NULL CHECK (total_amount >= 0),

  status VARCHAR(20) NOT NULL,

  created_at TIMESTAMP NOT NULL DEFAULT now(),

  CONSTRAINT chk_transaction_status
    CHECK (status IN ('pending','paid','canceled','completed'))
);