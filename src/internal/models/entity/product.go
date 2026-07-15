package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

type Product struct {
	ID          string          `db:"id" json:"id"`
	StoreID     string          `db:"store_id" json:"store_id"`
	CategoryID  string          `db:"category_id" json:"category_id"`
	Name        string          `db:"name" json:"name"`
	Slug        string          `db:"slug" json:"slug"`
	Description string          `db:"description" json:"description"`
	Price       decimal.Decimal `db:"price" json:"price"`
	Stock       int             `db:"stock" json:"stock"`
	RatingAvg   float64         `db:"rating_avg" json:"rating_avg"`
	TotalSold   int             `db:"total_sold" json:"total_sold"`
	IsActive    bool            `db:"is_active" json:"is_active"`
	CreatedAt   time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time       `db:"updated_at" json:"updated_at"`
	DeletedAt   *time.Time      `db:"deleted_at" json:"deleted_at,omitempty"`
}

type ProductWithCategory struct {
	Product
	CategoryName    string `db:"category_name" json:"category_name"`
	StoreName       string `db:"store_name" json:"store_name"`
	StoreSlug       string `db:"store_slug" json:"store_slug"`
	StoreIsVerified bool   `db:"store_is_verified" json:"store_is_verified"`
	PrimaryImage    string `db:"primary_image" json:"primary_image"`
}

type ProductDetail struct {
	Product
	CategoryName    string  `db:"category_name" json:"category_name"`
	StoreName       string  `db:"store_name" json:"store_name"`
	StoreSlug       string  `db:"store_slug" json:"store_slug"`
	StoreIsVerified bool    `db:"store_is_verified" json:"store_is_verified"`
	StoreLogoURL    string  `db:"store_logo_url" json:"store_logo_url"`
	StoreRatingAvg  float64 `db:"store_rating_avg" json:"store_rating_avg"`
	TotalReviews    int     `db:"total_reviews" json:"total_reviews"`
	PrimaryImage    string  `db:"primary_image" json:"primary_image"`
}
