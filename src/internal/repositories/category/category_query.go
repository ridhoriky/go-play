package category

const (
	hasProductsQuery = `
	SELECT EXISTS (
		SELECT 1 FROM products 
		WHERE category_id = $1 AND deleted_at IS NULL
	)`

	createCategoryQuery = `
	INSERT INTO categories (name, description, parent_id, image_url, sort_order)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id`

	getCategoryByIDQuery = `
	SELECT 
		id, 
		name, 
		description,
		parent_id,
		image_url,
		sort_order,
		created_at,
		updated_at,
		deleted_at
	FROM categories 
	WHERE id=$1
	`

	GetCategoryByNameQuery = `
	SELECT 
		id, 
		name, 
		description,
		parent_id,
		image_url,
		sort_order,
		created_at,
		updated_at,
		deleted_at
	FROM categories 
	WHERE name=$1
	`

	updateCategoryQuery = `
	UPDATE categories
	SET 
		name=$1, 
		description=$2,
		parent_id=$3,
		image_url=$4,
		sort_order=$5
	WHERE id=$6
	`

	deleteCategoryQuery = `
	UPDATE categories
	SET deleted_at = NOW()
	WHERE id = $1 AND deleted_at IS NULL
	`

	getCategoryTreeQuery = `
	SELECT 
		c.id, c.name, c.description, c.parent_id, c.image_url, c.sort_order,
		c.created_at, c.updated_at, c.deleted_at,
		(SELECT COUNT(*) FROM products WHERE category_id = c.id AND deleted_at IS NULL AND is_active = TRUE) AS product_count
	FROM categories c
	WHERE c.deleted_at IS NULL
	ORDER BY c.sort_order, c.name;
	`

	getByParentIDQuery = `
	SELECT id, name, description, parent_id, image_url, sort_order, created_at, updated_at, deleted_at
	FROM categories
	WHERE parent_id = $1 AND deleted_at IS NULL
	ORDER BY sort_order, name;
	`
)
