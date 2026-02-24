CREATE INDEX idx_td_product
ON transaction_details(product_id);

CREATE INDEX idx_td_transaction
ON transaction_details(transaction_id);

CREATE INDEX idx_td_transaction_product
ON transaction_details(transaction_id, product_id);