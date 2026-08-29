package storage

import (
	"os"
	"path/filepath"
)

const (
	applicationDirectory = "HookLens"
	databaseFilename     = "hooklens.db"
)

// DatabasePath returns the path for HookLens's user-local database.
//
// UserConfigDir supplies the platform-specific parent directory:
//   - macOS: ~/Library/Application Support
//   - Linux: $XDG_CONFIG_HOME or ~/.config
//   - Windows: %AppData%
func DatabasePath() (string, error) {
	baseDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	return databasePathFrom(baseDir)
}

func databasePathFrom(baseDir string) (string, error) {
	appDir := filepath.Join(baseDir, applicationDirectory)

	if err := os.MkdirAll(appDir, 0o755); err != nil {
		return "", err
	}

	return filepath.Join(appDir, databaseFilename), nil
}
