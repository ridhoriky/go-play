package dto

import (
	"errors"

	"ne-project/internal/models/entity"
)

type ProductDTO struct {
	ID         int    `db:"id" json:"id,omitempty"` // Omit in create requests
	Name       string `db:"name" json:"name"`
	Price      int    `db:"price" json:"price"`
	Stock      int    `db:"stock" json:"stock"`
	CategoryID *int   `db:"category_id" json:"category_id,omitempty"` // Optional
}

type ProductResponse struct {
	ID    int    `db:"id" json:"id"`
	Name  string `db:"name" json:"name"`
	Price int    `db:"price" json:"price"`
	Stock int    `db:"stock" json:"stock"`

	CategoryID          *int   `db:"category_id" json:"category_id"`
	CategoryName        string `db:"category_name" json:"category_name"`
	CategoryDescription string `db:"category_description" json:"category_description"`
}

type ProductFilterRequest struct {
	Name     string
	Category string
	Page     int
	Limit    int
}

func (p *ProductDTO) Validate() error {
	if p.Name == "" {
		return errors.New("product name is required")
	}
	if len(p.Name) > 255 {
		return errors.New("product name too long (max 255 chars)")
	}
	if p.Price <= 0 {
		return errors.New("product price must be greater than 0")
	}
	if p.Stock < 0 {
		return errors.New("product stock cannot be negative")
	}
	if p.CategoryID != nil && *p.CategoryID <= 0 {
		return errors.New("category ID must be positive if provided")
	}
	return nil
}

func (p *ProductDTO) ToModel() entity.Product {
	var categoryID int
	if p.CategoryID != nil {
		categoryID = *p.CategoryID
	}
	return entity.Product{
		ID:         p.ID,
		Name:       p.Name,
		Price:      p.Price,
		Stock:      p.Stock,
		CategoryID: categoryID,
	}
}

func (d *ProductDTO) ToModelPtr() *entity.Product {
	p := d.ToModel()
	return &p
}
