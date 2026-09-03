package mode

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

const appName = "quotes"

// Mode represents the application's operating mode
type Mode string

const (
	// Production mode - optimized for production with minimal logging and caching enabled
	Production Mode = "production"
	// Development mode - verbose logging, debug endpoints, and no caching
	Development Mode = "development"
)

var (
	currentMode Mode = Production // Default mode
	modeMutex   sync.RWMutex
)

// Get returns the current application mode
func Get() Mode {
	modeMutex.RLock()
	defer modeMutex.RUnlock()
	return currentMode
}

// Set sets the application mode
func Set(mode string) error {
	parsed, err := ParseMode(mode)
	if err != nil {
		return err
	}

	modeMutex.Lock()
	defer modeMutex.Unlock()
	currentMode = parsed
	return nil
}

// ParseMode parses a string into a Mode constant
// Accepts: "dev", "development", "prod", "production" (case-insensitive)
func ParseMode(s string) (Mode, error) {
	normalized := strings.ToLower(strings.TrimSpace(s))

	switch normalized {
	case "dev", "development":
		return Development, nil
	case "prod", "production":
		return Production, nil
	default:
		return Production, fmt.Errorf("invalid mode: %s (valid options: dev, development, prod, production)", s)
	}
}

// IsDevelopment returns true if the current mode is Development
func IsDevelopment() bool {
	return Get() == Development
}

// IsProduction returns true if the current mode is Production
func IsProduction() bool {
	return Get() == Production
}

// Initialize sets up the mode based on priority:
// 1. CLI flag (if provided via Set())
// 2. MODE environment variable
// 3. Default to Production
func Initialize() {
	// Check if MODE environment variable is set
	if envMode := os.Getenv("MODE"); envMode != "" {
		if err := Set(envMode); err != nil {
			// If invalid mode in env var, log warning and use default
			fmt.Fprintf(os.Stderr, "Warning: Invalid MODE environment variable: %s. Using default (production)\n", envMode)
		}
	}
	// If Set() was called with --mode flag before Initialize(), it will already be set
	// and this won't override it
}

// GetErrorDetail returns error details based on the current mode
// In development mode: returns full error details
// In production mode: returns generic error message
func GetErrorDetail(err error) string {
	if err == nil {
		return ""
	}

	if IsDevelopment() {
		return err.Error()
	}

	return "An internal error occurred. Please try again later."
}

// ShouldShowDebugEndpoints returns true if debug endpoints should be enabled
// Debug endpoints include /debug/pprof/* and /debug/vars
func ShouldShowDebugEndpoints() bool {
	return IsDevelopment()
}

// GetCacheHeaders returns appropriate cache control headers based on mode
// Development: no-cache to ensure fresh content on every request
// Production: cache headers for optimal performance
func GetCacheHeaders() map[string]string {
	if IsDevelopment() {
		return map[string]string{
			"Cache-Control": "no-cache, no-store, must-revalidate",
			"Pragma":        "no-cache",
			"Expires":       "0",
		}
	}

	// Production: cache static files for 1 year (immutable content)
	// For dynamic content, this should be overridden appropriately
	return map[string]string{
		"Cache-Control": "public, max-age=31536000, immutable",
	}
}

// GetLogLevel returns the appropriate log level for the current mode
func GetLogLevel() string {
	if IsDevelopment() {
		return "debug"
	}
	return "info"
}

// ShouldCacheTemplates returns true if templates should be cached
func ShouldCacheTemplates() bool {
	return IsProduction()
}

// ShouldCacheStaticFiles returns true if static files should be cached
func ShouldCacheStaticFiles() bool {
	return IsProduction()
}

// ShouldEnableAutoReload returns true if auto-reload should be enabled
func ShouldEnableAutoReload() bool {
	return IsDevelopment()
}

// ShouldEnableProfiling returns true if profiling should be enabled
func ShouldEnableProfiling() bool {
	return IsDevelopment()
}

// String returns the string representation of the Mode
func (m Mode) String() string {
	return string(m)
}
