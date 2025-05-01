// db.go
package db

import (
	"database/sql"
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
	connStr := os.Getenv("DB_CONN_STR")
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
	_, err := d.DB.Exec("UPDATE users SET email = $1 WHERE username = $2", newEmail, userID)
	return err
}

func (d *Database) GetRefreshTokenHash(username string) (string, error) {
	var hash string
	err := d.DB.QueryRow("SELECT refresh_token_hash FROM users WHERE username = $1", username).Scan(&hash)
	if err != nil {
		return "", err
	}
	return hash, nil
}

func (d *Database) StoreRefreshToken(username string, hashedToken string) error {
	_, err := d.DB.Exec("UPDATE users SET refresh_token_hash = $1 WHERE username = $2", hashedToken, username)
	return err
}
