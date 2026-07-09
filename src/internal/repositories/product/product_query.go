package product

const (
	getAllProductsQuery = `
	SELECT 
		p.id, 
		p.store_id,
		p.category_id,
		p.name,
		p.slug,
		p.description,
		p.price, 
		p.stock,
		p.rating_avg,
		p.total_sold,
		p.is_active,
		coalesce(c.name, '') as category_name,
		coalesce(s.store_name, '') as store_name,
		coalesce(s.slug, '') as store_slug,
		coalesce(s.is_verified, false) as store_is_verified,
		coalesce((SELECT url FROM product_images WHERE product_id = p.id AND is_primary = true LIMIT 1), '') as primary_image,
		COUNT(*) over() as total_count
	FROM
		products p
	LEFT JOIN categories c ON
		c.id = p.category_id
	LEFT JOIN stores s ON
		s.id = p.store_id
	WHERE 
		p.deleted_at is null
	`

	insertProductQuery = `
	INSERT INTO products (store_id, category_id, name, slug, description, price, stock, is_active) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
	RETURNING id`

	updateProductQuery = `
	UPDATE products 
	SET category_id=$1, name=$2, description=$3, price=$4, stock=$5, is_active=$6
	WHERE id=$7 AND store_id=$8 AND deleted_at IS NULL`

	deleteProductQuery = `
	UPDATE products
	SET deleted_at = NOW()
	WHERE id = $1 AND deleted_at IS NULL
	`

	getProductByIDQuery = `
	SELECT 
		p.id, 
		p.store_id,
		p.category_id,
		p.name, 
		p.slug,
		p.description,
		p.price, 
		p.stock, 
		p.rating_avg,
		p.total_sold,
		p.is_active,
		p.created_at,
		p.updated_at,
		COALESCE(c.name, '') as category_name,
		s.store_name as store_name,
		s.slug as store_slug,
		COALESCE(s.is_verified, false) as store_is_verified,
		COALESCE(s.logo_url, '') as store_logo_url,
		s.rating_avg as store_rating_avg,
		(SELECT COUNT(*) FROM reviews r WHERE r.product_id = p.id) as total_reviews
	FROM products p
	LEFT JOIN categories c ON p.category_id = c.id
	JOIN stores s ON p.store_id = s.id
	WHERE p.id = $1 AND p.deleted_at IS NULL`

	getProductBySlugQuery = `
	SELECT 
		p.id, 
		p.store_id,
		p.category_id,
		p.name, 
		p.slug,
		p.description,
		p.price, 
		p.stock, 
		p.rating_avg,
		p.total_sold,
		p.is_active
	FROM products p
	WHERE p.slug = $1 AND p.deleted_at IS NULL`

	getProductDetailBySlugQuery = `
	SELECT 
		p.id, 
		p.store_id,
		p.category_id,
		p.name, 
		p.slug,
		p.description,
		p.price, 
		p.stock, 
		p.rating_avg,
		p.total_sold,
		p.is_active,
		p.created_at,
		p.updated_at,
		COALESCE(c.name, '') as category_name,
		s.store_name as store_name,
		s.slug as store_slug,
		COALESCE(s.is_verified, false) as store_is_verified,
		COALESCE(s.logo_url, '') as store_logo_url,
		s.rating_avg as store_rating_avg,
		(SELECT COUNT(*) FROM reviews r WHERE r.product_id = p.id) as total_reviews
	FROM products p
	LEFT JOIN categories c ON p.category_id = c.id
	JOIN stores s ON p.store_id = s.id
	WHERE p.slug = $1 AND p.deleted_at IS NULL`

	insertBulkProductQuery = `
	INSERT INTO products (%s) 
	VALUES %s 
	RETURNING %s`
)
