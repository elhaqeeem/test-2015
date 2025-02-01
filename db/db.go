package db

import (
	"database/sql"
	"fmt"
	"golang-echo-postgresql/config"
	"log"

	_ "github.com/lib/pq"
)

// InitDB initializes and returns a database connection
func InitDB() *sql.DB {
	// Load config values from environment
	cfg := config.LoadConfig()

	// Build the connection string using loaded values
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)

	// Open a connection to the PostgreSQL database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}

	// Ping the database to check if it's reachable
	if err = db.Ping(); err != nil {
		log.Fatal("Database is not reachable: ", err)
	}

	fmt.Println("Connected to the database successfully")
	return db
}
