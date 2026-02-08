package transaction

const (
	lockProductQuery = `
	SELECT name, price, stock FROM products WHERE id = $1
	`

	insertTxQuery = `
	INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id
	`

	updateStockQuery = `
	UPDATE products SET stock = stock - $1 WHERE id = $2
	`

	getTodayTransaction = `
	SELECT
		t.id,
		t.total_amount,
		t.created_at,

		td.id,
		td.transaction_id,
		td.product_id,
		p.name,
		td.quantity,
		td.subtotal

	FROM transactions t
	JOIN transaction_details td
		ON td.transaction_id = t.id
	JOIN products p
		ON p.id = td.product_id

	WHERE DATE(t.created_at) = CURRENT_DATE

	ORDER BY t.created_at DESC
	`
)
