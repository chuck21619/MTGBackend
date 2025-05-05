package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	dsn := os.Getenv("DB_CONN_STR")
	if dsn == "" {
		log.Fatal("DB_CONN_STR is not set")
	}

	for i := 0; i < 30; i++ {
		db, err := sql.Open("postgres", dsn)
		if err == nil {
			err = db.Ping()
			if err == nil {
				fmt.Println("Database is ready!")
				return
			}
		}

		fmt.Println("Waiting for database...")
		time.Sleep(2 * time.Second)
	}

	log.Fatal("Timed out waiting for the database")
}
