package storage_test

import (
	"os"
	"path/filepath"
	"testing"

	"server/internal/storage"
)

func TestDatabasePathUsesUserConfigDir(t *testing.T) {
	configRoot := t.TempDir()
	setUserConfigRoot(t, configRoot)

	databasePath, err := storage.DatabasePath()
	if err != nil {
		t.Fatalf("DatabasePath() returned an error: %v", err)
	}

	baseDir, err := os.UserConfigDir()
	if err != nil {
		t.Fatalf("os.UserConfigDir() returned an error: %v", err)
	}

	wantPath := filepath.Join(baseDir, "HookLens", "hooklens.db")
	if databasePath != wantPath {
		t.Fatalf("DatabasePath() = %q, want %q", databasePath, wantPath)
	}

	if _, err := os.Stat(filepath.Dir(databasePath)); err != nil {
		t.Fatalf("database directory was not created: %v", err)
	}

	if _, err := os.Stat(databasePath); !os.IsNotExist(err) {
		t.Fatalf("DatabasePath() should not create the database file; stat error = %v", err)
	}
}

func setUserConfigRoot(t *testing.T, configRoot string) {
	t.Helper()

	// UserConfigDir reads a different variable on each supported OS. Setting
	// all three keeps this test isolated without encoding a platform path.
	t.Setenv("AppData", configRoot)
	t.Setenv("HOME", configRoot)
	t.Setenv("XDG_CONFIG_HOME", configRoot)
}
