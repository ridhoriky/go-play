package database

import (
	"log"

	"ne-project/src/internal/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func InitDB(cfg *config.Config) (*sqlx.DB, error) {

	db, err := sqlx.Connect("postgres", cfg.Database.GetConnectionString())
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	log.Println("DB connection successfully")

	return db, nil
}
