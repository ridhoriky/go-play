package database

import (
	"database/sql"
	"log"
	"ne-project/internal/config"

	_ "github.com/lib/pq"
)

func InitDB(cfg *config.Config) (*sql.DB, error) {
	//open DB
	db , err := sql.Open("postgres", cfg.Database.GetConnectionString())
	
	if err != nil {
		return nil , err
	}
	
	//Test connection 
	err = db.Ping()
	if err != nil {
		return nil , err
	}
	
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	log.Println("DB connection successfully")

	return db, nil
}