package transaction

const (
	lockProductStockQuery = `
        SELECT stock FROM products
        WHERE id = $1 AND deleted_at IS NULL
        FOR UPDATE NOWAIT`

	deductStockQuery = `
        UPDATE products
        SET stock = stock - $1, updated_at = NOW()
        WHERE id = $2`

	insertTransactionQuery = `
        INSERT INTO transactions (id, total_amount, status, created_at)
        VALUES ($1, $2, $3, NOW())
        RETURNING id, total_amount, status, created_at`

	insertTransactionDetailQuery = `
        INSERT INTO transaction_details (id, transaction_id, product_id, product_name, quantity, price, subtotal)
        VALUES ($1, $2, $3, $4, $5, $6, $7)`

	getProductForCheckoutQuery = `
        SELECT id, name, price, stock FROM products
        WHERE id = $1 AND deleted_at IS NULL`
)
