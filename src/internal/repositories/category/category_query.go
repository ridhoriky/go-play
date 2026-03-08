package category

const (
	getAllCategoriesQuery = `SELECT id, name, description FROM categories`
	createCategoryQuery   = `INSERT INTO categories (name, description) VALUES ($1, $2) RETURNING id`
	getCategoryByIDQuery  = `SELECT id, name, description FROM categories WHERE id=$1`
	updateCategoryQuery   = `UPDATE categories SET name=$1, description=$2 WHERE id=$3`
	deleteCategoryQuery   = `DELETE FROM categories WHERE id=$1`
)