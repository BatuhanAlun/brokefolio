package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func InitDB(DB_CONNECTION_STRING string) *sql.DB {

	db, err := sql.Open("postgres", DB_CONNECTION_STRING)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	log.Println("Connection Successfully Established")
	return db
}
