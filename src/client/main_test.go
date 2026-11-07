package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetJSON(t *testing.T) {
	original := serverURL
	t.Cleanup(func() { serverURL = original })
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/quote" || r.Header.Get("Accept") != "application/json" || !strings.HasPrefix(r.Header.Get("User-Agent"), "quotes-cli/") {
			t.Errorf("unexpected request: %s %#v", r.URL, r.Header)
		}
		_, _ = w.Write([]byte(`{"ok":true,"data":{"id":7,"quote":"hello","author":"world"}}`))
	}))
	defer server.Close()
	serverURL = server.URL
	var got quote
	if err := getJSON("/quote", &got); err != nil {
		t.Fatal(err)
	}
	if got.ID != 7 || got.Quote != "hello" || got.Author != "world" {
		t.Fatalf("decoded quote = %#v", got)
	}
}

func TestGetJSONRejectsErrors(t *testing.T) {
	original := serverURL
	t.Cleanup(func() { serverURL = original })
	for _, tc := range []struct {
		name string
		status int
		body string
	}{
		{"API message", http.StatusBadRequest, `{"message":"bad query"}`},
		{"empty server error", http.StatusServiceUnavailable, ""},
		{"invalid JSON", http.StatusOK, "not-json"},
		{"unsuccessful envelope", http.StatusOK, `{"ok":false,"message":"rejected"}`},
		{"invalid quote", http.StatusOK, `{"ok":true,"data":"wrong-shape"}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.status)
				_, _ = w.Write([]byte(tc.body))
			}))
			defer server.Close()
			serverURL = server.URL
			var got quote
			if err := getJSON("/quote", &got); err == nil {
				t.Fatalf("getJSON accepted %s", tc.name)
				}
		})
	}
	serverURL = "http://127.0.0.1:1"
	var got quote
	if err := getJSON("/quote", &got); err == nil {
		t.Fatal("getJSON accepted an unavailable server")
	}
}

func TestPrintQuote(t *testing.T) {
	for _, tc := range []struct { format, author, want string; wantErr bool }{
		{"plain", "", "a line\n", false},
		{"plain", "Ada", `"a line" - Ada` + "\n", false},
		{"json", "Ada", `"quote": "a line"`, false},
		{"xml", "Ada", "", true},
	} {
		var output bytes.Buffer
		err := printQuote(&output, quote{Quote: "a line", Author: tc.author}, tc.format)
		if (err != nil) != tc.wantErr {
			t.Fatalf("format %q err = %v", tc.format, err)
		}
		if !tc.wantErr && !strings.Contains(output.String(), tc.want) {
			t.Errorf("format %q output = %q, want %q", tc.format, output.String(), tc.want)
		}
	}
}

func TestCLIConfigRoundTripAndInvalidYAML(t *testing.T) {
	root := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", root)
	if err := saveServerURL("https://quotes.example"); err != nil {
		t.Fatal(err)
	}
	got, err := loadServerURL()
	if err != nil || got != "https://quotes.example" {
		t.Fatalf("loadServerURL() = %q, %v", got, err)
	}
	path, err := configFilePath()
	if err != nil { t.Fatal(err) }
	info, err := os.Stat(path)
	if err != nil || info.Mode().Perm() != 0600 {
		t.Fatalf("CLI config permissions = %v, error = %v", info, err)
	}
	if err := os.WriteFile(path, []byte("server: ["), 0600); err != nil { t.Fatal(err) }
	if _, err := loadServerURL(); err == nil {
		t.Fatal("invalid YAML configuration was accepted")
	}
	if filepath.Dir(path) == root {
		t.Fatal("configuration path omitted the application subdirectory")
	}
}
