package store

const (
	createStoreQuery = `
	INSERT INTO stores (id, user_id, store_name, slug, description, logo_url, banner_url, is_verified, rating_avg, total_sales, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	getStoreByIDQuery = `
	SELECT id, user_id, store_name, slug, COALESCE(description, ''), COALESCE(logo_url, ''), COALESCE(banner_url, ''), is_verified, rating_avg, total_sales, created_at, updated_at, deleted_at
	FROM stores
	WHERE id = $1 AND deleted_at IS NULL
	`
	getStoreByUserIDQuery = `
	SELECT id, user_id, store_name, slug, COALESCE(description, ''), COALESCE(logo_url, ''), COALESCE(banner_url, ''), is_verified, rating_avg, total_sales, created_at, updated_at, deleted_at
	FROM stores
	WHERE user_id = $1 AND deleted_at IS NULL
	`
	getStoreBySlugQuery = `
	SELECT id, user_id, store_name, slug, COALESCE(description, ''), COALESCE(logo_url, ''), COALESCE(banner_url, ''), is_verified, rating_avg, total_sales, created_at, updated_at, deleted_at
	FROM stores
	WHERE slug = $1 AND deleted_at IS NULL
	`
	updateStoreQuery = `
	UPDATE stores
	SET store_name = $1, description = $2, logo_url = $3, banner_url = $4, updated_at = NOW()
	WHERE id = $5 AND deleted_at IS NULL
	`
	softDeleteStoreQuery = `
	UPDATE stores
	SET deleted_at = NOW(), updated_at = NOW()
	WHERE id = $1 AND deleted_at IS NULL
	`
	checkSlugExistsQuery = `
	SELECT EXISTS(SELECT 1 FROM stores WHERE slug = $1 AND deleted_at IS NULL)
	`
	listStoresBaseQuery = `
	SELECT id, user_id, store_name, slug, COALESCE(description, ''), COALESCE(logo_url, ''), COALESCE(banner_url, ''), is_verified, rating_avg, total_sales, created_at, updated_at, deleted_at, COUNT(*) OVER() AS total_count
	FROM stores
	WHERE deleted_at IS NULL
	`
	getStoreStatsQuery = `
	SELECT
		(SELECT COUNT(*) FROM products WHERE store_id = $1 AND deleted_at IS NULL) AS total_products,
		(SELECT COUNT(*) FROM products WHERE store_id = $1 AND is_active = TRUE AND deleted_at IS NULL) AS active_products,
		(SELECT COUNT(*) FROM orders WHERE store_id = $1) AS total_orders,
		(SELECT COALESCE(SUM(total_amount), 0) FROM orders WHERE store_id = $1 AND status IN ('paid','processing','shipped','delivered')) AS total_revenue,
		s.rating_avg,
		(SELECT COUNT(*) FROM reviews r JOIN products p ON r.product_id = p.id WHERE p.store_id = $1) AS total_reviews
	FROM stores s
	WHERE s.id = $1 AND s.deleted_at IS NULL;
	`
)
