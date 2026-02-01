package repositories

import (
	"database/sql"
	"ne-project/internal/models"
)

type CategoryRepository struct {
	 db *sql.DB
}

func NewCategoryRepository (db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db:db }
}

func (repo *CategoryRepository) GetAll() ([]models.Category, error){
	query := "SELECT id, name, description FROM categories"
	rows, err := repo.db.Query(query)

	if err != nil {
		return nil , err 
	}
	defer rows.Close()

	categories := make([]models.Category, 0)
	for rows.Next(){
		var c models.Category
		err := rows.Scan(&c.ID, &c.Name, &c.Description)
		if err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return categories, nil
}

func (repo *CategoryRepository) Create(category *models.Category) error{
	query := "INSERT INTO categories (name, description) VALUES ($1, $2) RETURNING id"
	err := repo.db.QueryRow(query, category.Name, category.Description).Scan(&category.ID)
	return err
}

func (repo *CategoryRepository) GetByID(id int) (*models.Category, error){
	query := "SELECT id, name, description FROM categories WHERE id=$1"
	var c models.Category
	err := repo.db.QueryRow(query, id).Scan(&c.ID, &c.Name, &c.Description)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (repo *CategoryRepository) Update(category *models.Category) error{
	query := "UPDATE categories SET name=$1, description=$2 WHERE id=$3"
	_, err := repo.db.Exec(query, category.Name, category.Description, category.ID)
	return err
}

func (repo *CategoryRepository) Delete(id int) error{
	query := "DELETE FROM categories WHERE id=$1"
	_, err := repo.db.Exec(query, id)
	return err
}