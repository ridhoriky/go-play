package admin

const (
	queryListUsersBase = `
		SELECT
			u.id, u.name, u.email, u.role, u.is_active, u.created_at,
			CASE WHEN s.id IS NOT NULL THEN TRUE ELSE FALSE END AS has_store,
			s.store_name
		FROM users u
		LEFT JOIN stores s ON u.id = s.user_id AND s.deleted_at IS NULL
		WHERE u.deleted_at IS NULL
	`

	queryCountUsersBase = `
		SELECT COUNT(*)
		FROM users u
		LEFT JOIN stores s ON u.id = s.user_id AND s.deleted_at IS NULL
		WHERE u.deleted_at IS NULL
	`

	queryGetUserByID = `
		SELECT
			u.id, u.name, u.email, u.role, u.is_active, u.created_at,
			CASE WHEN s.id IS NOT NULL THEN TRUE ELSE FALSE END AS has_store,
			s.store_name
		FROM users u
		LEFT JOIN stores s ON u.id = s.user_id AND s.deleted_at IS NULL
		WHERE u.id = $1 AND u.deleted_at IS NULL
	`

	queryUpdateUser = `
		UPDATE users
		SET role = COALESCE($2, role),
			is_active = COALESCE($3, is_active),
			updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	queryListSellersBase = `
		SELECT
			s.id AS store_id, s.store_name, s.slug, s.is_verified,
			s.rating_avg, s.total_sales,
			u.name AS owner_name, u.email AS owner_email,
			s.created_at,
			(SELECT COUNT(*) FROM products WHERE store_id = s.id AND deleted_at IS NULL) AS total_products
		FROM stores s
		JOIN users u ON s.user_id = u.id
		WHERE s.deleted_at IS NULL
	`

	queryCountSellersBase = `
		SELECT COUNT(*)
		FROM stores s
		JOIN users u ON s.user_id = u.id
		WHERE s.deleted_at IS NULL
	`

	queryUpdateStoreVerification = `
		UPDATE stores SET is_verified = $2, updated_at = NOW() WHERE id = $1 AND deleted_at IS NULL
	`

	queryGetPlatformSummary = `
		SELECT
			(SELECT COUNT(*) FROM users WHERE deleted_at IS NULL) AS total_users,
			(SELECT COUNT(*) FROM users WHERE role = 'buyer' AND deleted_at IS NULL) AS total_buyers,
			(SELECT COUNT(*) FROM users WHERE role = 'seller' AND deleted_at IS NULL) AS total_sellers,
			(SELECT COUNT(*) FROM users WHERE role = 'admin' AND deleted_at IS NULL) AS total_admins,
			(SELECT COUNT(*) FROM stores WHERE deleted_at IS NULL) AS total_stores,
			(SELECT COUNT(*) FROM stores WHERE is_verified = TRUE AND deleted_at IS NULL) AS verified_stores,
			(SELECT COUNT(*) FROM stores WHERE is_verified = FALSE AND deleted_at IS NULL) AS unverified_stores,
			(SELECT COUNT(*) FROM products WHERE deleted_at IS NULL) AS total_products,
			(SELECT COUNT(*) FROM products WHERE deleted_at IS NULL AND is_active = TRUE) AS active_products,
			(SELECT COUNT(*) FROM products WHERE deleted_at IS NULL AND is_active = FALSE) AS inactive_products,
			(SELECT COUNT(*) FROM orders) AS total_orders,
			(SELECT COUNT(*) FROM orders WHERE status = 'pending') AS pending_orders,
			(SELECT COUNT(*) FROM orders WHERE status = 'paid') AS paid_orders,
			(SELECT COUNT(*) FROM orders WHERE status = 'processing') AS processing_orders,
			(SELECT COUNT(*) FROM orders WHERE status = 'shipped') AS shipped_orders,
			(SELECT COUNT(*) FROM orders WHERE status = 'delivered') AS delivered_orders,
			(SELECT COUNT(*) FROM orders WHERE status = 'completed') AS completed_orders,
			(SELECT COUNT(*) FROM orders WHERE status = 'canceled') AS canceled_orders,
			(SELECT COALESCE(SUM(total_amount), 0) FROM orders WHERE status IN ('paid','processing','shipped','delivered','completed')) AS total_revenue,
			(SELECT COUNT(*) FROM users WHERE deleted_at IS NULL AND created_at >= NOW() - INTERVAL '7 days') AS new_users_this_week,
			(SELECT COUNT(*) FROM users WHERE deleted_at IS NULL AND created_at >= NOW() - INTERVAL '1 month') AS new_users_this_month,
			(SELECT COUNT(*) FROM orders WHERE created_at >= NOW() - INTERVAL '7 days') AS new_orders_this_week,
			(SELECT COUNT(*) FROM orders WHERE created_at >= NOW() - INTERVAL '1 month') AS new_orders_this_month
	`

	queryGetTopStoresBase = `
		SELECT s.id AS store_id, s.store_name, s.slug, s.is_verified,
			s.rating_avg, s.total_sales,
			u.name AS owner_name, u.email AS owner_email,
			s.created_at,
			(SELECT COUNT(*) FROM products WHERE store_id = s.id AND deleted_at IS NULL) AS total_products,
			(SELECT COALESCE(SUM(o.total_amount), 0) FROM orders o 
			 WHERE o.store_id = s.id AND o.status IN ('paid','processing','shipped','delivered')) AS revenue
		FROM stores s
		JOIN users u ON s.user_id = u.id
		WHERE s.deleted_at IS NULL
	`

	queryGetTopProductsBase = `
		SELECT 
			p.id AS product_id, p.name AS product_name,
			s.store_name, COALESCE(c.name, '') AS category_name, p.price,
			p.rating_avg,
			(SELECT COALESCE(SUM(oi.quantity), 0) 
			 FROM order_items oi 
			 JOIN orders o ON oi.order_id = o.id 
			 WHERE oi.product_id = p.id AND o.status IN ('paid','processing','shipped','delivered')) AS quantity_sold,
			(SELECT COALESCE(SUM(oi.subtotal), 0) 
			 FROM order_items oi 
			 JOIN orders o ON oi.order_id = o.id 
			 WHERE oi.product_id = p.id AND o.status IN ('paid','processing','shipped','delivered')) AS revenue
		FROM products p
		JOIN stores s ON p.store_id = s.id
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE p.deleted_at IS NULL AND p.is_active = TRUE
	`
)
