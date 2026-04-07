package report

const getSummaryQuery = `
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
