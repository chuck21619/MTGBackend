// db.go
package db

import (
	"database/sql"
	"os"
	"log"
	_ "github.com/lib/pq"
	"github.com/chuck21619/MTGBackend/models"
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

func (d *Database) UpdateGoogleSheet(userID string, new_google_sheet string) error {
	_, err := d.DB.Exec("UPDATE users SET google_sheet = $1 WHERE username = $2", new_google_sheet, userID)
	return err
}

func (d *Database) GetGoogleSheet(username string) (string, error) {
	var result string
	err := d.DB.QueryRow("SELECT google_sheet FROM users WHERE username = $1", username).Scan(&result)
	if err != nil {
		return "", err
	}
	return result, nil
}

func (d *Database) GetRefreshTokenHash(username string) (string, error) {
	var hash string
	err := d.DB.QueryRow("SELECT refresh_token_hash FROM users WHERE username = $1", username).Scan(&hash)
	if err != nil {
		return "", err
	}
	return hash, nil
}

func (d *Database) GetProfileInfo(username string) (*models.User, error) {
    var user models.User
    query := `SELECT email, google_sheet FROM users WHERE username = $1`
    err := d.DB.QueryRow(query, username).Scan(&user.Email, &user.GoogleSheet)
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (d *Database) StoreRefreshToken(username string, hashedToken string) error {
	_, err := d.DB.Exec("UPDATE users SET refresh_token_hash = $1 WHERE username = $2", hashedToken, username)
	return err
}

func (db *Database) DeleteRefreshToken(username string) error {
	_, err := db.Exec("UPDATE users SET refresh_token_hash = NULL WHERE username = $1", username)
	return err
}
