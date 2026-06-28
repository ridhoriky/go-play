package entity

import "time"

type ProductImage struct {
	ID        string    `db:"id" json:"id"`
	ProductID string    `db:"product_id" json:"product_id"`
	URL       string    `db:"url" json:"url"`
	AltText   string    `db:"alt_text" json:"alt_text"`
	SortOrder int       `db:"sort_order" json:"sort_order"`
	IsPrimary bool      `db:"is_primary" json:"is_primary"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
