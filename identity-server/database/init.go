package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func EnsureDatabase() error {
	// Connect to postgres default database
	connStr := "host=localhost port=5432 user=postgres password=password dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("error connecting to postgres: %v", err)
	}
	defer db.Close()

	// Check if database exists
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE datname = 'identity_server')").Scan(&exists)
	if err != nil {
		return fmt.Errorf("error checking database existence: %v", err)
	}

	// Create database if it doesn't exist
	if !exists {
		_, err = db.Exec("CREATE DATABASE identity_server")
		if err != nil {
			return fmt.Errorf("error creating database: %v", err)
		}
		fmt.Println("Identity server database created successfully")
	}

	return nil
}
