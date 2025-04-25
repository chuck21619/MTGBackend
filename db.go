// db.go
package main
import "os"
import "fmt"

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB

func initDB() {
	password := os.Getenv("DB_PASSWORD")
	connStr := fmt.Sprintf("postgresql://userdb_zvh9_user:%s@dpg-d05s8015pdvs73em8j5g-a.oregon-postgres.render.com/userdb_zvh9", password)
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to DB:", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal("Could not ping DB:", err)
	}
}
