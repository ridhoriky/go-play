CREATE INDEX idx_transactions_created
ON transactions(created_at);

CREATE INDEX idx_transactions_status
ON transactions(status);