package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB(filepath string) *sql.DB {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		log.Fatal(err)
	}

	createTables(db)
	return db
}

func createTables(db *sql.DB) {
	query := `
    CREATE TABLE IF NOT EXISTS fields (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT,
        area_ha REAL,
        region TEXT
    );
    
    CREATE TABLE IF NOT EXISTS sowings (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        field_id INTEGER,
        crop TEXT,
        sowed_at TEXT
    );
    
    CREATE TABLE IF NOT EXISTS harvests (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        field_id INTEGER,
        crop TEXT,
        yield_t_per_ha REAL
    );`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}
