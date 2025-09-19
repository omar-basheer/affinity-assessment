package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var DB *sql.DB

func InitDB(path string) error {
	var err error
	DB, err = sql.Open("sqlite3", path)
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	// ensure vouchers table exists
	schema := `
	CREATE TABLE IF NOT EXISTS vouchers (
		id TEXT PRIMARY KEY,
		customer_id INTEGER NOT NULL,
		customer_name TEXT NOT NULL,
		order_value REAL NOT NULL,
		amount REAL NOT NULL,
		created_at DATETIME NOT NULL,
		expires_at DATETIME NOT NULL
	);`

	_, err = DB.Exec(schema)
	if err != nil {
		log.Fatal("failed to run migration:", err)
	}
	return nil
}
