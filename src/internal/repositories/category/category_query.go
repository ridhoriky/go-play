package category

const (
	getAllCategoriesQuery = `
	SELECT 
		c.id, 
		c.name, 
		c.description,
		c.created_at,
		c.updated_at,
		c.deleted_at,
        COUNT(*) OVER() AS total_count
	FROM categories c
	WHERE 1=1`

	// countCategoriesQuery = `
	// SELECT COUNT(*)
	// FROM categories c
	// WHERE 1=1`

	createCategoryQuery = `
	INSERT INTO categories (name, description)
	VALUES ($1, $2)
	RETURNING id`

	getCategoryByIDQuery = `
	SELECT 
		id, 
		name, 
		description 
	FROM categories 
	WHERE id=$1
	`

	GetCategoryByNameQuery = `
	SELECT 
		id, 
		name, 
		description 
	FROM categories 
	WHERE name=$1
	`

	updateCategoryQuery = `
	UPDATE categories
	SET 
		name=$1, 
		description=$2 
	WHERE id=$3
	`

	deleteCategoryQuery = `
	UPDATE categories
	SET deleted_at = NOW()
	WHERE id = $1 AND deleted_at IS NULL
	`
)
