package ssl

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

func TestTLSConfigurationAndCertificates(t *testing.T) {
	if cfg, err := NewManager(Config{}).GetTLSConfig([]string{"example.test"}); err != nil || cfg != nil {
		t.Fatalf("disabled TLS = %v, %v", cfg, err)
	}
	manager := NewManager(Config{Enabled: true, CertPath: t.TempDir()})
	if _, err := manager.GetTLSConfig([]string{"example.test"}); err == nil {
		t.Fatal("missing certificates accepted")
	}
	if _, err := manager.GetTLSConfig([]string{"example.test"}); err == nil {
		t.Fatal("missing certificates accepted twice")
	}
	if got := manager.GetHTTPHandler(http.NotFoundHandler()); got == nil {
		t.Fatal("fallback handler missing")
	}
	if _, err := manager.getLetsEncryptTLSConfig([]string{"example.test"}); err != nil {
		t.Fatal(err)
	}
	if _, err := filepath.Abs(filepath.Join(manager.config.CertPath, "autocert")); err != nil {
		t.Fatal(err)
	}
}

func TestChallengeServer(t *testing.T) {
	server := NewChallengeServer()
	request := httptest.NewRequest(http.MethodGet, "/unrelated", nil)
	recorder := httptest.NewRecorder()
	if server.ServeHTTP(recorder, request) {
		t.Fatal("handled a non-challenge request")
	}
	request = httptest.NewRequest(http.MethodGet, "/.well-known/acme-challenge/missing", nil)
	if !server.ServeHTTP(recorder, request) || recorder.Code != http.StatusNotFound {
		t.Fatal("missing token was not rejected")
	}
	server.SetToken("token", "proof")
	request = httptest.NewRequest(http.MethodGet, "/.well-known/acme-challenge/token", nil)
	recorder = httptest.NewRecorder()
	if !server.ServeHTTP(recorder, request) || recorder.Body.String() != "proof" || recorder.Header().Get("Content-Type") != "text/plain" {
		t.Fatal("challenge response mismatch")
	}
	server.ClearToken("token")
	if _, ok := server.tokens["token"]; ok {
		t.Fatal("ClearToken retained token")
	}
}

func TestParseChallenge(t *testing.T) {
	for input, want := range map[string]string{"http": "http-01", "HTTP01": "http-01", " tls ": "tls-alpn-01", "tls-alpn-01": "tls-alpn-01", "dns": "dns-01", "dns01": "dns-01", "unknown": "http-01"} {
		if got := ParseChallenge(input); got != want {
			t.Errorf("ParseChallenge(%q) = %q, want %q", input, got, want)
		}
	}
}
