// db.go
package db

import (
	"database/sql"
	"fmt"
	"os"
	"log"
	_ "github.com/lib/pq"
)

// Database struct to hold the db connection
type Database struct {
	*sql.DB
}

// NewDatabase initializes the database connection and returns a Database instance
func NewDatabase() *Database {
	password := os.Getenv("DB_PASSWORD")
	connStr := fmt.Sprintf("postgresql://userdb_zvh9_user:%s@dpg-d05s8015pdvs73em8j5g-a.oregon-postgres.render.com/userdb_zvh9", password)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to DB:", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("Could not ping DB:", err)
	}

	return &Database{DB: db}
}

func (d *Database) UpdateUserEmail(userID string, newEmail string) error {
	_, err := d.DB.Exec("UPDATE users SET email = $1 WHERE id = $2", newEmail, userID)
	return err
}
