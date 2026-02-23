package dto

import (
	"errors"

	"ne-project/internal/models/entity"
)

type CategoryDTO struct {
	ID          int    `json:"id,omitempty"` // Omit in create requests
	Name        string `json:"name"`
	Description string `json:"description,omitempty"` // Optional
}

func (c *CategoryDTO) Validate() error {
	if c.Name == "" {
		return errors.New("category name is required")
	}
	if len(c.Name) > 255 {
		return errors.New("category name too long (max 255 chars)")
	}
	if len(c.Description) > 1000 {
		return errors.New("category description too long (max 1000 chars)")
	}
	return nil
}

func (c *CategoryDTO) ToModel() entity.Category {
	return entity.Category{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
	}
}

func (d *CategoryDTO) ToModelPtr() *entity.Category {
	c := d.ToModel()
	return &c
}
