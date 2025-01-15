package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func openDB() error {
	connect := "host=localhost port=5432 user=nurymalibekov dbname=hotel sslmode=disable"

	var err error
	DB, err = sql.Open("postgres", connect)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	fmt.Println("Connected to the database successfully!")
	return nil
}

func closeDB() error {
	if err := DB.Close(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Closing database!")
	return nil
}
