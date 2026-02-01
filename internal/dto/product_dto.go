package dto

import (
	"errors"

	"ne-project/internal/models"
)

type ProductDTO struct {
	Name       string `json:"name"`
	Price      int    `json:"price"`
	Stock      int    `json:"stock"`
	CategoryID *int   `json:"category_id,omitempty"` // Optional
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

func (p *ProductDTO) ToModel() models.Product {
	var categoryID int
	if p.CategoryID != nil {
		categoryID = *p.CategoryID
	}
	return models.Product{
		Name:       p.Name,
		Price:      p.Price,
		Stock:      p.Stock,
		CategoryID: categoryID,
	}
}
