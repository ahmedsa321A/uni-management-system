// Package store contains all the database interaction logic.
package store

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3" // Driver import
)

// Connect opens a connection to the SQLite database file.
func Connect(dbPath string) (*sql.DB, error) {
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("database file does not exist at path: %s", dbPath)
	}

	// THE FIX: Add "?_loc=auto" to the connection string.
	// This tells the SQLite driver to automatically parse time-related
	// columns (DATE, TIME, DATETIME) from strings into time.Time objects.
	dsn := fmt.Sprintf("%s?_loc=auto", dbPath)

	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Ping the database to verify the connection is alive.
	if err = db.Ping(); err != nil {
		db.Close() // Close the connection if ping fails.
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}
