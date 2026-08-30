package storage_test

import (
	"os"
	"path/filepath"
	"testing"

	"server/internal/storage"
)

func TestOpenCreatesDatabaseAtUserConfigDir(t *testing.T) {
	configRoot := t.TempDir()
	setUserConfigRoot(t, configRoot)

	databasePath, err := storage.DatabasePath()
	if err != nil {
		t.Fatalf("DatabasePath() returned an error: %v", err)
	}

	database, err := storage.Open()
	if err != nil {
		t.Fatalf("Open() returned an error: %v", err)
	}
	defer database.Close()

	if _, err := os.Stat(databasePath); err != nil {
		t.Fatalf("database file was not created at %q: %v", databasePath, err)
	}

	var tableName string
	if err := database.QueryRow(
		"SELECT name FROM sqlite_master WHERE type = 'table' AND name = 'inboxes'",
	).Scan(&tableName); err != nil {
		t.Fatalf("inboxes table was not initialized: %v", err)
	}

	if tableName != "inboxes" {
		t.Fatalf("initialized table = %q, want %q", tableName, "inboxes")
	}

	if err := database.QueryRow(
		"SELECT name FROM sqlite_master WHERE type = 'table' AND name = 'inbox_events'",
	).Scan(&tableName); err != nil {
		t.Fatalf("inbox_events table was not initialized: %v", err)
	}

	if tableName != "inbox_events" {
		t.Fatalf("initialized event table = %q, want %q", tableName, "inbox_events")
	}

	if filepath.Base(databasePath) != "hooklens.db" {
		t.Fatalf("database filename = %q, want %q", filepath.Base(databasePath), "hooklens.db")
	}
}
