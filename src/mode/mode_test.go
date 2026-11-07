package mode

import (
	"errors"
	"os"
	"testing"
)

func TestModeBehavior(t *testing.T) {
	if err := Set("invalid"); err == nil {
		t.Fatal("Set accepted invalid mode")
	}
	if got, err := ParseMode(" DEV "); err != nil || got != Development {
		t.Fatalf("ParseMode(dev) = %q, %v", got, err)
	}
	if got, err := ParseMode("PRODUCTION"); err != nil || got != Production {
		t.Fatalf("ParseMode(production) = %q, %v", got, err)
	}
	if _, err := ParseMode("staging"); err == nil {
		t.Fatal("ParseMode accepted staging")
	}
	if err := Set("development"); err != nil {
		t.Fatal(err)
	}
	if !IsDevelopment() || IsProduction() || !ShouldShowDebugEndpoints() || !ShouldEnableAutoReload() || !ShouldEnableProfiling() || ShouldCacheTemplates() || ShouldCacheStaticFiles() {
		t.Fatal("development mode flags are inconsistent")
	}
	if GetLogLevel() != "debug" || GetErrorDetail(errors.New("specific")) != "specific" {
		t.Fatal("development details were hidden")
	}
	if err := Set("production"); err != nil {
		t.Fatal(err)
	}
	if IsDevelopment() || !IsProduction() || ShouldShowDebugEndpoints() || ShouldEnableAutoReload() || ShouldEnableProfiling() || !ShouldCacheTemplates() || !ShouldCacheStaticFiles() {
		t.Fatal("production mode flags are inconsistent")
	}
	if GetLogLevel() != "info" || GetErrorDetail(errors.New("specific")) == "specific" || GetErrorDetail(nil) != "" {
		t.Fatal("production error details are incorrect")
	}
}

func TestInitializeAndModeHelpers(t *testing.T) {
	previous := os.Getenv("MODE")
	t.Cleanup(func() { _ = os.Setenv("MODE", previous) })
	if err := os.Setenv("MODE", "dev"); err != nil {
		t.Fatal(err)
	}
	Initialize()
	if !IsDevelopment() {
		t.Fatal("Initialize ignored MODE")
	}
	if err := os.Setenv("MODE", "invalid"); err != nil {
		t.Fatal(err)
	}
	Initialize()
	if !IsDevelopment() {
		t.Fatal("invalid MODE changed current mode")
	}
	if err := os.Unsetenv("MODE"); err != nil {
		t.Fatal(err)
	}
	if got := GetCacheHeaders(); got["Cache-Control"] != "no-cache, no-store, must-revalidate" || got["Pragma"] != "no-cache" || got["Expires"] != "0" {
		t.Fatalf("development cache headers = %#v", got)
	}
	if Production.String() != "production" || Development.String() != "development" {
		t.Fatal("mode string conversion failed")
	}
	if err := Set("production"); err != nil {
		t.Fatal(err)
	}
	if got := GetCacheHeaders(); got["Cache-Control"] != "public, max-age=31536000, immutable" {
		t.Fatalf("production cache headers = %#v", got)
	}
}
