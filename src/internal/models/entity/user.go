package entity

import "time"

type User struct {
	ID        string    `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Email     string    `db:"email" json:"email"`
	Password  string    `db:"password_hash" json:"-"`
	Role      string    `db:"role" json:"role"`
	IsActive  bool      `db:"is_active" json:"is_active"`
	AvatarURL *string   `db:"avatar_url" json:"avatar_url"`
	Phone     *string   `db:"phone" json:"phone"`
	Address   *string   `db:"address" json:"address"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
