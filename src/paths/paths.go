package paths

import (
	"os"
	"path/filepath"
	"runtime"
)

// GetConfigDir returns the configuration directory for the application
func GetConfigDir() string {
	// Check environment variable first
	if dir := os.Getenv("CONFIG_DIR"); dir != "" {
		return dir
	}

	switch runtime.GOOS {
	case "windows":
		if appData := os.Getenv("APPDATA"); appData != "" {
			return filepath.Join(appData, "quotes")
		}
		return filepath.Join(os.Getenv("USERPROFILE"), ".config", "quotes")
	case "darwin":
		if home := os.Getenv("HOME"); home != "" {
			return filepath.Join(home, "Library", "Application Support", "quotes")
		}
		return "/usr/local/etc/quotes"
	default: // Linux and other Unix-like systems
		if xdgConfig := os.Getenv("XDG_CONFIG_HOME"); xdgConfig != "" {
			return filepath.Join(xdgConfig, "quotes")
		}
		if home := os.Getenv("HOME"); home != "" {
			return filepath.Join(home, ".config", "quotes")
		}
		return "/etc/quotes"
	}
}

// GetDataDir returns the data directory for the application
func GetDataDir() string {
	// Check environment variable first
	if dir := os.Getenv("DATA_DIR"); dir != "" {
		return dir
	}

	switch runtime.GOOS {
	case "windows":
		if appData := os.Getenv("LOCALAPPDATA"); appData != "" {
			return filepath.Join(appData, "quotes")
		}
		return filepath.Join(os.Getenv("USERPROFILE"), ".local", "share", "quotes")
	case "darwin":
		if home := os.Getenv("HOME"); home != "" {
			return filepath.Join(home, "Library", "Application Support", "quotes")
		}
		return "/usr/local/var/quotes"
	default: // Linux and other Unix-like systems
		if xdgData := os.Getenv("XDG_DATA_HOME"); xdgData != "" {
			return filepath.Join(xdgData, "quotes")
		}
		if home := os.Getenv("HOME"); home != "" {
			return filepath.Join(home, ".local", "share", "quotes")
		}
		return "/var/lib/quotes"
	}
}

// GetLogsDir returns the logs directory for the application
func GetLogsDir() string {
	// Check environment variable first
	if dir := os.Getenv("LOGS_DIR"); dir != "" {
		return dir
	}

	switch runtime.GOOS {
	case "windows":
		return filepath.Join(GetDataDir(), "logs")
	case "darwin":
		if home := os.Getenv("HOME"); home != "" {
			return filepath.Join(home, "Library", "Logs", "quotes")
		}
		return "/usr/local/var/log/quotes"
	default: // Linux and other Unix-like systems
		return "/var/log/quotes"
	}
}

// GetDBPath returns the database file path
func GetDBPath() string {
	// Check environment variable first
	if dbPath := os.Getenv("DB_PATH"); dbPath != "" {
		return dbPath
	}

	dataDir := GetDataDir()
	return filepath.Join(dataDir, "db", "quotes.db")
}

// EnsureDir creates a directory if it doesn't exist
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}
