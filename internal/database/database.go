package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func InitDB(connectionString string) (*sql.DB, error) {
	//open DB
	db , err := sql.Open("postgres", connectionString)
	
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