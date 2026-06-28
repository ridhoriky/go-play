package order

const (
	insertOrderQuery = `
		INSERT INTO orders (id, buyer_id, store_id, order_number, total_amount, status, shipping_address, shipping_cost, payment_method, payment_ref, notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`

	insertOrderItemBulkQuery = `
		INSERT INTO order_items (%s) VALUES %s
	`

	getProductStockForUpdateQuery = `
		SELECT stock, is_active FROM products WHERE id = $1 FOR UPDATE
	`

	deductProductStockQuery = `
		UPDATE products SET stock = stock - $1, total_sold = total_sold + $1 WHERE id = $2
	`

	restoreProductStockQuery = `
		UPDATE products SET stock = stock + $1, total_sold = total_sold - $1 WHERE id = $2
	`

	getOrderByIDQuery = `
		SELECT id, buyer_id, store_id, order_number, total_amount, status, shipping_address, shipping_cost, payment_method, payment_ref, notes, created_at, updated_at
		FROM orders WHERE id = $1
	`

	getOrderByOrderNumberQuery = `
		SELECT id, buyer_id, store_id, order_number, total_amount, status, shipping_address, shipping_cost, payment_method, payment_ref, notes, created_at, updated_at
		FROM orders WHERE order_number = $1
	`

	getOrderItemsByOrderIDQuery = `
		SELECT id, order_id, product_id, product_name, product_image, quantity, price, subtotal, created_at
		FROM order_items WHERE order_id = $1
	`

	getOrdersByBuyerIDBaseQuery = `
		SELECT id, buyer_id, store_id, order_number, total_amount, status, shipping_address, shipping_cost, payment_method, payment_ref, notes, created_at, updated_at, COUNT(*) OVER() as total_count
		FROM orders WHERE buyer_id = $1
	`

	getOrdersByStoreIDBaseQuery = `
		SELECT id, buyer_id, store_id, order_number, total_amount, status, shipping_address, shipping_cost, payment_method, payment_ref, notes, created_at, updated_at, COUNT(*) OVER() as total_count
		FROM orders WHERE store_id = $1
	`

	updateOrderStatusQuery = `
		UPDATE orders SET status = $1, updated_at = NOW() WHERE id = $2
	`

	updateOrderPaymentQuery = `
		UPDATE orders SET payment_method = $1, payment_ref = $2, updated_at = NOW() WHERE id = $3
	`

	hasActiveOrdersForProductQuery = `
		SELECT EXISTS (
			SELECT 1 FROM order_items oi
			JOIN orders o ON oi.order_id = o.id
			WHERE oi.product_id = $1 AND o.status IN ('pending', 'paid', 'processing', 'shipped')
		)
	`

	getRecentOrdersByProductIDQuery = `
		SELECT o.id, u.name as buyer_name, o.total_amount, o.status, o.created_at
		FROM orders o
		JOIN order_items oi ON o.id = oi.order_id
		JOIN users u ON o.buyer_id = u.id
		WHERE oi.product_id = $1
		ORDER BY o.created_at DESC
		LIMIT $2
	`
)
