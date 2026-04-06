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
        INSERT INTO transaction_details 
                (id, transaction_id, product_id, product_name, quantity, price, subtotal)
        VALUES `

	getProductForCheckoutQuery = `
        SELECT id, name, price, stock FROM products
        WHERE id = $1 AND deleted_at IS NULL`

	getTransactionByIDQuery = `
        SELECT 
            t.id, 
            t.total_amount, 
            t.status,
            t.created_at,
            td.id,
            td.transaction_id,
            td.product_id,
            td.product_name, 
            td.quantity,
            td.price, 
            td.subtotal 
        FROM 
            transactions t 
        LEFT JOIN transaction_details td ON t.id = td.transaction_id
        WHERE 
            t.id = $1
        `

	updateTransactionStatusQuery = `
        UPDATE transactions
        SET status = $1
        WHERE id = $2 AND deleted_at IS NULL
        RETURNING id, total_amount, status, created_at
        `

	addStockQuery = `
	UPDATE products
	SET stock = stock + $1, updated_at = NOW()
	WHERE id = $2 AND deleted_at IS NULL
	`
)
