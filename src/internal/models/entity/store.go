package entity

import "time"

type Store struct {
	ID          string     `db:"id" json:"id"`
	UserID      string     `db:"user_id" json:"user_id"`
	StoreName   string     `db:"store_name" json:"store_name"`
	Slug        string     `db:"slug" json:"slug"`
	Description string     `db:"description" json:"description"`
	LogoURL     string     `db:"logo_url" json:"logo_url"`
	BannerURL   string     `db:"banner_url" json:"banner_url"`
	IsVerified  bool       `db:"is_verified" json:"is_verified"`
	RatingAvg   float64    `db:"rating_avg" json:"rating_avg"`
	TotalSales  int        `db:"total_sales" json:"total_sales"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}
