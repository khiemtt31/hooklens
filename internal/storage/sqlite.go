package storage

import (
	"database/sql"
	"fmt"

	_ "embed"
	_ "modernc.org/sqlite"
)

//go:embed schema.sql
var schema string

// Open opens HookLens's user-local SQLite database and applies the schema.
func Open() (*sql.DB, error) {
	databasePath, err := DatabasePath()
	if err != nil {
		return nil, fmt.Errorf("resolve database path: %w", err)
	}

	return OpenAt(databasePath)
}

// OpenAt opens a SQLite database at an explicit path. It is useful for
// callers that need a separate database and for isolated tests.
func OpenAt(databasePath string) (*sql.DB, error) {
	database, err := sql.Open("sqlite", databasePath)
	if err != nil {
		return nil, fmt.Errorf("open sqlite database: %w", err)
	}

	database.SetMaxOpenConns(1)
	database.SetMaxIdleConns(1)

	if err := database.Ping(); err != nil {
		database.Close()

		return nil, fmt.Errorf("ping sqlite database: %w", err)
	}

	if _, err := database.Exec("PRAGMA foreign_keys = ON"); err != nil {
		database.Close()

		return nil, fmt.Errorf("configure sqlite database: %w", err)
	}

	if _, err := database.Exec(schema); err != nil {
		database.Close()

		return nil, fmt.Errorf("initialize sqlite schema: %w", err)
	}

	return database, nil
}
