// db.go
package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB

func initDB() {
	connStr := "postgres://userdb_zvh9_user:qXoma9zbzVduyXtXdPiVPbb8RtqlS4c6@dpg-d05s8015pdvs73em8j5g-a.oregon-postgres.render.com/userdb_zvh9"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to DB:", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal("Could not ping DB:", err)
	}
}
