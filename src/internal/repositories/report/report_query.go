package report

const (
	getSummaryQuery = `
	SELECT
		COALESCE(SUM(t.total_amount), 0)  AS total_revenue,
		COUNT(t.id)                       AS total_transactions,
		COALESCE(SUM(td.total_qty), 0)    AS total_items_sold,
		COALESCE(AVG(t.total_amount), 0)  AS average_transaction
	FROM transactions t
	LEFT JOIN (
		SELECT transaction_id, SUM(quantity) AS total_qty
		FROM transaction_details
		GROUP BY transaction_id
	) td ON td.transaction_id = t.id
	WHERE t.status = 'paid'
	  AND t.created_at BETWEEN $1 AND $2
	`

	getTopProductsQuery = `
	SELECT
		td.product_id,
		td.product_name,
		SUM(td.quantity) AS total_quantity,
		SUM(td.quantity * td.price) AS total_revenue
	FROM transaction_details td
	JOIN transactions t ON t.id = td.transaction_id
	WHERE t.status = 'paid'
	  AND t.created_at BETWEEN $1 AND $2
	GROUP BY td.product_id, td.product_name
	ORDER BY total_quantity DESC
	LIMIT $3
	`
)
