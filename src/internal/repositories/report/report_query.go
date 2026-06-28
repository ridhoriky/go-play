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

	getSellerSalesSummaryQuery = `
	SELECT
		COUNT(o.id) AS total_orders,
		COALESCE(SUM(o.total_amount), 0) AS total_revenue,
		(
			SELECT COALESCE(SUM(oi.quantity), 0)
			FROM order_items oi
			JOIN orders o2 ON oi.order_id = o2.id
			WHERE o2.store_id = $1
			  AND o2.status IN ('paid', 'processing', 'shipped', 'delivered')
			  AND o2.created_at BETWEEN $2 AND $3
		) AS total_items_sold
	FROM orders o
	WHERE o.store_id = $1
	  AND o.status IN ('paid', 'processing', 'shipped', 'delivered')
	  AND o.created_at BETWEEN $2 AND $3;
	`

	getSellerTopProductsQuery = `
	SELECT
		p.id, p.name,
		COALESCE(SUM(oi.quantity), 0) AS total_sold,
		COALESCE(SUM(oi.subtotal), 0) AS revenue,
		p.rating_avg
	FROM order_items oi
	JOIN products p ON oi.product_id = p.id
	JOIN orders o ON oi.order_id = o.id
	WHERE o.store_id = $1
	  AND o.status IN ('paid', 'processing', 'shipped', 'delivered')
	  AND o.created_at BETWEEN $2 AND $3
	GROUP BY p.id, p.name, p.rating_avg
	ORDER BY revenue DESC
	LIMIT $4 OFFSET $5;
	`

	getSellerRecentOrdersQuery = `
	SELECT o.id, o.total_amount, o.status, o.created_at, u.name AS buyer_name
	FROM orders o
	JOIN users u ON o.buyer_id = u.id
	WHERE o.store_id = $1
	ORDER BY o.created_at DESC
	LIMIT $2;
	`

	getSellerPendingOrdersCountQuery = `
	SELECT COUNT(*) FROM orders
	WHERE store_id = $1 AND status = 'pending';
	`

	getSellerLowStockProductsQuery = `
	SELECT id, name, stock FROM products
	WHERE store_id = $1 AND stock <= $2 AND is_active = TRUE;
	`
)
