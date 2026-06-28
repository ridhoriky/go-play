package entity

import "time"

type Category struct {
	ID          string     `db:"id" json:"id"`
	Name        string     `db:"name" json:"name"`
	Description *string    `db:"description" json:"description"`
	ParentID    *string    `db:"parent_id" json:"parent_id,omitempty"`
	ImageURL    *string    `db:"image_url" json:"image_url,omitempty"`
	SortOrder   int        `db:"sort_order" json:"sort_order"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}

type CategoryWithCount struct {
	Category
	ProductCount int `db:"product_count" json:"product_count"`
}
