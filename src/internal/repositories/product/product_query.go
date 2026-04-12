package product

const (
	getAllProductsQuery = `
	SELECT 
		p.id, 
		p.name, 
		p.price, 
		p.stock,
		p.category_id,
		coalesce(c.name, '') as category_name,
		COUNT(*) over() as total_count
	FROM
		products p
	LEFT JOIN categories c ON
		c.id = p.category_id
	WHERE 
		p.deleted_at is null
	`

	// countProductsQuery = `
	// SELECT COUNT(*)
	// FROM products p
	// WHERE 1=1`

	insertProductQuery = `
	INSERT INTO products (name, price, stock, category_id) 
	VALUES ($1, $2, $3, $4) 
	RETURNING id`

	updateProductQuery = `
	UPDATE products 
	SET name=$1, price=$2, stock=$3, category_id=$4 
	WHERE id=$5`

	deleteProductQuery = `
	UPDATE products
	SET deleted_at = NOW()
	WHERE id = $1 AND deleted_at IS NULL
	`

	getProductByIDQuery = `
	SELECT 
		p.id, 
		p.name, 
		p.price, 
		p.stock, 
		p.category_id,
		COALESCE(c.name, '') 
	FROM products p
	LEFT JOIN categories c ON p.category_id = c.id
	WHERE p.id = $1`

	insertBulkProductQuery = `
	INSERT INTO products (%s) 
	VALUES %s 
	RETURNING %s`
)
