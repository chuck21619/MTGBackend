package main

import (
	"log"
	"os"

	"database/sql"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func migrate() {
	// Ensure that the DB connection string is set in the environment variables
	connStr := os.Getenv("DB_CONN_STR")
	if connStr == "" {
		log.Fatal("DB_CONN_STR not set in environment variables")
	}

	// Open the database connection using pq driver
	db, err := sql.Open("postgres", connStr)
	// db, err := goose.OpenDBWithDriver("postgres", connStr)
	if err != nil {
		log.Fatal("Error opening database:", err)
	}
	defer db.Close()

	if err := os.Chdir("/app"); err != nil {
		log.Fatal("Error changing directory to /app:", err)
	}
	
	// Run the migrations
	if err := goose.Up(db, "./migrations"); err != nil {
		log.Fatal("Error running migrations:", err)
	}

	log.Println("Migrations applied successfully!")
}

func main() {
	migrate()
}
