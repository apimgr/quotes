package httpapp

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/apimgr/quotes/src/catalog"
	"github.com/apimgr/quotes/src/config"
)

func testHandler(t *testing.T) http.Handler {
	t.Helper()
	store, err := catalog.LoadEmbedded()
	if err != nil {
		t.Fatal(err)
	}
	return NewHandler(store, "test", config.DefaultConfig().Server.RateLimit)
}

func TestPublicRoutes(t *testing.T) {
	handler := testHandler(t)
	cases := []struct {
		path string
		wantStatus int
		wantBody string
	}{
		{"/", http.StatusOK, "Quotes"},
		{"/?category=programming", http.StatusOK, "value=\"programming\""},
		{"/random?category=programming", http.StatusOK, "Another programming quote"},
		{"/random?category=unknown", http.StatusNotFound, "Quote not found"},
		{"/quotes/1", http.StatusOK, "Quote #1"},
		{"/quotes/999999", http.StatusNotFound, "Quote not found"},
		{"/server/about", http.StatusOK, "About Quotes"},
		{"/server/help", http.StatusOK, "Help"},
		{"/server/privacy", http.StatusOK, "Privacy"},
		{"/server/terms", http.StatusOK, "Terms"},
		{"/server/healthz", http.StatusOK, `"status": "healthy"`},
		{"/api/healthz", http.StatusOK, `"status": "healthy"`},
		{"/api/v1/server/healthz", http.StatusOK, `"version": "test"`},
		{"/api/v1/categories", http.StatusOK, `"programming"`},
		{"/api/v1/quotes", http.StatusOK, `"pagination"`},
		{"/api/v1/quotes?category=anime&limit=2", http.StatusOK, `"limit": 2`},
		{"/api/v1/quotes?author=Character%201", http.StatusOK, `"total": 1`},
		{"/api/v1/categories/anime/quotes", http.StatusOK, `"category": "anime"`},
		{"/api/v1/quotes/random?category=programming", http.StatusOK, `"category": "programming"`},
		{"/api/v1/quotes/1", http.StatusOK, `"id": 1`},
		{"/api/v1/quotes/invalid", http.StatusBadRequest, `"error": "BAD_REQUEST"`},
		{"/api/v1/quotes/999999", http.StatusNotFound, `"error": "NOT_FOUND"`},
		{"/api/v1/quotes?limit=1001", http.StatusBadRequest, `"error": "BAD_REQUEST"`},
		{"/api/v1/categories/unknown/quotes", http.StatusNotFound, `"error": "NOT_FOUND"`},
		{"/api/v1/server/swagger", http.StatusOK, `"openapi": "3.1.0"`},
		{"/api/swagger", http.StatusOK, `"openapi": "3.1.0"`},
		{"/server/docs/swagger", http.StatusOK, "machine-readable"},
		{"/server/docs/graphql", http.StatusOK, "GraphQL explorer"},
		{"/manifest.webmanifest", http.StatusOK, `"name":"Quotes API"`},
		{"/sw.js", http.StatusOK, "addEventListener"},
		{"/.well-known/security.txt", http.StatusOK, "security/advisories"},
		{"/missing", http.StatusNotFound, `"error": "NOT_FOUND"`},
	}
	for _, tc := range cases {
		t.Run(tc.path, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, tc.path, nil)
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, request)
			if response.Code != tc.wantStatus {
				t.Fatalf("status = %d, body = %s, want %d", response.Code, response.Body.String(), tc.wantStatus)
			}
			if !strings.Contains(response.Body.String(), tc.wantBody) {
				t.Fatalf("body does not contain %q: %s", tc.wantBody, response.Body.String())
			}
			if response.Header().Get("Access-Control-Allow-Origin") != "*" {
				t.Fatal("public CORS policy is missing")
			}
		})
	}
}

func TestGraphQL(t *testing.T) {
	handler := testHandler(t)
	for _, tc := range []struct {
		path string
		query string
		want string
	}{
		{"/api/graphql", `{ randomQuote(category: "anime") { id category } }`, `"category": "anime"`},
		{"/api/v1/server/graphql", `{ quote(id: 1) { id quote author category } }`, `"id": 1`},
		{"/api/graphql", `{ categories { name count } }`, `"name": "quotes"`},
		{"/api/graphql", `{ quotes(category: "programming", page: 1, limit: 1) { id quote category } }`, `"category": "programming"`},
	} {
		t.Run(tc.path+tc.query, func(t *testing.T) {
			body, err := json.Marshal(map[string]string{"query": tc.query})
			if err != nil {
				t.Fatal(err)
			}
			request := httptest.NewRequest(http.MethodPost, tc.path, bytes.NewReader(body))
			request.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, request)
			if response.Code != http.StatusOK || !strings.Contains(response.Body.String(), tc.want) {
				t.Fatalf("status=%d body=%s; want %q", response.Code, response.Body.String(), tc.want)
			}
		})
	}
}

func TestErrorAndTextResponses(t *testing.T) {
	handler := testHandler(t)
	request := httptest.NewRequest(http.MethodGet, "/api/v1/quotes/random?category=anime", nil)
	request.Header.Set("Accept", "text/plain")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusOK || !strings.HasSuffix(response.Body.String(), "\n") || !strings.HasPrefix(response.Header().Get("Content-Type"), "text/plain") {
		t.Fatalf("plain response: status=%d headers=%v body=%q", response.Code, response.Header(), response.Body.String())
	}
	request = httptest.NewRequest(http.MethodGet, "/api/v1/quotes/random?category=invalid", nil)
	response = httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	var problem errorResponse
	if err := json.Unmarshal(response.Body.Bytes(), &problem); err != nil {
		t.Fatal(err)
	}
	if response.Code != http.StatusBadRequest || problem.OK || problem.Error != "BAD_REQUEST" || problem.Message == "" {
		t.Fatalf("error response = status %d, %#v", response.Code, problem)
	}
}

func TestSecurityMiddleware(t *testing.T) {
	handler := testHandler(t)
	request := httptest.NewRequest(http.MethodPut, "/api/v1/quotes", nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusMethodNotAllowed || response.Header().Get("Allow") == "" || !strings.Contains(response.Body.String(), "METHOD_NOT_ALLOWED") {
		t.Fatalf("method guard response: %d %s", response.Code, response.Body.String())
	}
	large := strings.NewReader(strings.Repeat("x", maxRequestBytes+1))
	request = httptest.NewRequest(http.MethodPost, "/api/graphql", large)
	request.Header.Set("Content-Type", "application/json")
	response = httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusRequestEntityTooLarge || !strings.Contains(response.Body.String(), "REQUEST_TOO_LARGE") {
		t.Fatalf("request limit response: %d %s", response.Code, response.Body.String())
	}
	request = httptest.NewRequest(http.MethodOptions, "/api/v1/quotes", nil)
	response = httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusNoContent {
		t.Fatalf("preflight response status = %d", response.Code)
	}
}

func TestRateLimit(t *testing.T) {
	handler := testHandler(t)
	var response *httptest.ResponseRecorder
	for i := 0; i < 121; i++ {
		request := httptest.NewRequest(http.MethodGet, "/api/v1/categories", nil)
		request.RemoteAddr = "192.0.2.1:12345"
		response = httptest.NewRecorder()
		handler.ServeHTTP(response, request)
	}
	if response.Code != http.StatusTooManyRequests || response.Header().Get("Retry-After") == "" || !strings.Contains(response.Body.String(), "RATE_LIMITED") {
		t.Fatalf("rate limit response: %d %s", response.Code, response.Body.String())
	}
}
