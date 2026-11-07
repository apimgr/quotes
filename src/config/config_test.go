package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadSaveAndDefaults(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "server.yml")
	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Server.Address != "0.0.0.0" || cfg.WebUI.Theme != "dark" {
		t.Fatalf("unexpected defaults: %#v", cfg)
	}
	if err := Save(); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatal(err)
	}
	if got := GetTheme(); got != "dark" {
		t.Fatalf("theme = %q", got)
	}
	if got := GetCORS(); got != "*" {
		t.Fatalf("CORS = %q", got)
	}

	if err := os.WriteFile(path, []byte("web-ui:\n  theme: light\nweb-security:\n  cors: https://example.test\n"), 0600); err != nil {
		t.Fatal(err)
	}
	cfg, err = Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.WebUI.Theme != "light" || GetCORS() != "https://example.test" {
		t.Fatalf("loaded values: %#v", cfg)
	}
	if err := Save(); err != nil {
		t.Fatal(err)
	}
}

func TestLoadMigratesYAMLAndRejectsMalformedConfig(t *testing.T) {
	dir := t.TempDir()
	legacy := filepath.Join(dir, "server.yaml")
	current := filepath.Join(dir, "server.yml")
	if err := os.WriteFile(legacy, []byte("web-ui:\n  theme: light\n"), 0600); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(current)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.WebUI.Theme != "light" {
		t.Fatalf("theme = %q", cfg.WebUI.Theme)
	}
	if _, err := os.Stat(legacy); !os.IsNotExist(err) {
		t.Fatalf("legacy file still exists: %v", err)
	}
	if err := os.WriteFile(current, []byte("web-ui: ["), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(current); err == nil {
		t.Fatal("Load accepted malformed YAML")
	}
}

func TestFormatStringSlice(t *testing.T) {
	if got := formatStringSlice(nil); got != "[]" {
		t.Fatalf("empty slice = %q", got)
	}
	if got := formatStringSlice([]string{"/", "/api"}); got != `["/", "/api"]` {
		t.Fatalf("slice = %q", got)
	}
	if !strings.Contains(generateConfigYAML(DefaultConfig()), "web-security:") {
		t.Fatal("generated config omits web security")
	}
}
