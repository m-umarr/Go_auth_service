package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

const (
	dsn = "postgres://me:root@localhost/authdb?sslmode=disable"
)

func Connect() (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Println("Connected to the database successfully")
	return db, nil
}
