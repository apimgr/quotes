package httpapp

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"mime"
	"net"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/apimgr/quotes/src/catalog"
	"github.com/apimgr/quotes/src/config"
	"github.com/graphql-go/graphql"
)

const apiVersion = "v1"
const maxRequestBytes = 1 << 20
const defaultPageSize = 250
const maximumPageSize = 1000

type response struct {
	OK         bool        `json:"ok"`
	Data       interface{} `json:"data,omitempty"`
	Pagination *pagination `json:"pagination,omitempty"`
}

type pagination struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
	Pages int `json:"pages"`
}

type errorResponse struct {
	OK      bool                   `json:"ok"`
	Error   string                 `json:"error"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

type rateCounter struct {
	start time.Time
	count int
}

type rateLimiter struct {
	mu       sync.Mutex
	requests map[string]rateCounter
	limit    int
	window   time.Duration
	enabled  bool
}

func newRateLimiter(limit int, window time.Duration) *rateLimiter {
	return &rateLimiter{requests: make(map[string]rateCounter), limit: limit, window: window}
}

func (l *rateLimiter) allow(key string, now time.Time) (bool, time.Duration) {
	l.mu.Lock()
	defer l.mu.Unlock()
	entry := l.requests[key]
	if entry.start.IsZero() || now.Sub(entry.start) >= l.window {
		if entry.start.IsZero() && len(l.requests) >= 10000 {
			for ip, counter := range l.requests {
				if now.Sub(counter.start) >= l.window {
					delete(l.requests, ip)
				}
			}
			if len(l.requests) >= 10000 {
				return false, l.window
			}
		}
		entry = rateCounter{start: now}
	}
	if entry.count >= l.limit {
		l.requests[key] = entry
		return false, l.window - now.Sub(entry.start)
	}
	entry.count++
	l.requests[key] = entry
	if len(l.requests) > 10000 {
		for ip, counter := range l.requests {
			if now.Sub(counter.start) >= l.window {
				delete(l.requests, ip)
			}
		}
	}
	return true, 0
}

// NewHandler builds the public web, REST, OpenAPI, and GraphQL surfaces.
func NewHandler(store *catalog.Catalog, version string, limits config.RateLimitConfig) http.Handler {
	mux := http.NewServeMux()
	if limits.Requests <= 0 {
		limits.Requests = 120
	}
	if limits.WindowSeconds <= 0 {
		limits.WindowSeconds = 60
	}
	limiter := newRateLimiter(limits.Requests, time.Duration(limits.WindowSeconds)*time.Second)
	limiter.enabled = limits.Enabled
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			writeError(w, http.StatusNotFound, "NOT_FOUND", "Page not found", nil)
			return
		}
		home(w, r, store)
	})
	mux.HandleFunc("GET /random", randomQuotePage(store))
	mux.HandleFunc("GET /quotes/{id}", quotePage(store))
	mux.HandleFunc("GET /server/about", staticPage("About Quotes", "Quotes is a public, read-only collection API for inspirational quotes, anime, Chuck Norris facts, dad jokes, and programming quotes."))
	mux.HandleFunc("GET /server/help", staticPage("Help", "Browse categories, filter by author, open a quote ID for its detail page, or choose a random quote. The REST, GraphQL, and command-line interfaces are also public."))
	mux.HandleFunc("GET /server/privacy", staticPage("Privacy", "Quotes stores no visitor accounts or personal quote data. Requests are served from a dataset embedded in the application binary."))
	mux.HandleFunc("GET /server/terms", staticPage("Terms", "Quote data is provided for reference. Respect applicable copyrights and attribution requirements when reusing individual quotations."))
	serverPages := map[string]map[string]string{
		"about": {"title": "About Quotes", "content": "Quotes is a public, read-only collection API for inspirational quotes, anime, Chuck Norris facts, dad jokes, and programming quotes."},
		"help": {"title": "Help", "content": "Browse categories, filter by author, open a quote ID for its detail page, or choose a random quote. The REST, GraphQL, and command-line interfaces are also public."},
		"privacy": {"title": "Privacy", "content": "Quotes stores no visitor accounts or personal quote data. Requests are served from a dataset embedded in the application binary."},
		"terms": {"title": "Terms", "content": "Quote data is provided for reference. Respect applicable copyrights and attribution requirements when reusing individual quotations."},
	}
	for page, values := range serverPages {
		page, values := page, values
		mux.HandleFunc("GET /api/"+apiVersion+"/server/"+page, func(w http.ResponseWriter, r *http.Request) {
			writeSuccess(w, values)
		})
	}
	health := func(w http.ResponseWriter, r *http.Request) {
		status := healthStatus(store, version)
		accept := r.Header.Get("Accept")
		switch {
		case strings.Contains(accept, "text/html"):
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_ = healthTemplate.Execute(w, status)
		case strings.Contains(accept, "text/plain"):
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			_, _ = fmt.Fprintf(w, "Status: %s\nVersion: %s\nGo: %s\nQuotes: %d\n", status.Status, status.Version, status.GoVersion, status.Total)
		default:
			writeJSON(w, http.StatusOK, status)
		}
	}
	mux.HandleFunc("GET /server/healthz", health)
	healthJSON := func(w http.ResponseWriter, r *http.Request) { writeJSON(w, http.StatusOK, healthStatus(store, version)) }
	mux.HandleFunc("GET /api/healthz", healthJSON)
	mux.HandleFunc("GET /api/"+apiVersion+"/server/healthz", healthJSON)
	apiRoot := func(w http.ResponseWriter, r *http.Request) {
		writeSuccess(w, map[string]interface{}{"name": "Quotes API", "version": version, "categories": store.Categories()})
	}
	mux.HandleFunc("GET /api/"+apiVersion, apiRoot)
	mux.HandleFunc("GET /api/"+apiVersion+"/categories", func(w http.ResponseWriter, r *http.Request) {
		writeSuccess(w, store.Categories())
	})
	listQuotes := func(forcedCategory string) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			category := forcedCategory
			if category == "" {
				category = r.URL.Query().Get("category")
			}
			author := strings.TrimSpace(r.URL.Query().Get("author"))
			page, limit, ok := parsePagination(w, r)
			if !ok {
				return
			}
			quotes, total, err := store.List(category, author, page, limit)
			if err != nil {
				writeError(w, http.StatusBadRequest, "BAD_REQUEST", err.Error(), nil)
				return
			}
			pages := (total + limit - 1) / limit
			writeList(w, quotes, pagination{Page: page, Limit: limit, Total: total, Pages: pages})
		}
	}
	mux.Handle("GET /api/"+apiVersion+"/quotes", negotiateText(listQuotes("")))
	mux.Handle("GET /api/"+apiVersion+"/categories/{category}/quotes", negotiateText(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := store.Categories()[r.PathValue("category")]; !ok {
			writeError(w, http.StatusNotFound, "NOT_FOUND", "Category not found", nil)
			return
		}
		listQuotes(r.PathValue("category")).ServeHTTP(w, r)
	})))
	mux.Handle("GET /api/"+apiVersion+"/quotes/random", negotiateText(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		quote, err := store.Random(r.URL.Query().Get("category"))
		if err != nil {
			writeError(w, http.StatusBadRequest, "BAD_REQUEST", err.Error(), nil)
			return
		}
		writeSuccess(w, quote)
	})))
	mux.Handle("GET /api/"+apiVersion+"/quotes/{id}", negotiateText(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil || id < 1 {
			writeError(w, http.StatusBadRequest, "BAD_REQUEST", "Quote ID must be a positive integer", map[string]interface{}{"field": "id"})
			return
		}
		quote, err := store.Get("", id)
		if err != nil {
			writeError(w, http.StatusNotFound, "NOT_FOUND", "Quote not found", nil)
			return
		}
		writeSuccess(w, quote)
	})))
	mux.HandleFunc("GET /api/"+apiVersion+"/quotes.txt", func(w http.ResponseWriter, r *http.Request) {
		listText(w, r, store, "")
	})
	mux.HandleFunc("GET /api/"+apiVersion+"/quotes/random.txt", func(w http.ResponseWriter, r *http.Request) {
		quote, err := store.Random(r.URL.Query().Get("category"))
		if err != nil {
			writeError(w, http.StatusBadRequest, "BAD_REQUEST", err.Error(), nil)
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = fmt.Fprintf(w, "\"%s\" - %s\n", quote.Quote, quote.Author)
	})
	mux.HandleFunc("GET /api/"+apiVersion+"/server/swagger", openAPISpec)
	mux.HandleFunc("GET /api/swagger", openAPISpec)
	mux.HandleFunc("GET /server/docs/swagger", swaggerUI)
	graphqlHandler := newGraphQLHandler(store)
	mux.HandleFunc("POST /api/"+apiVersion+"/server/graphql", graphqlHandler)
	mux.HandleFunc("POST /api/graphql", graphqlHandler)
	mux.HandleFunc("GET /server/docs/graphql", graphQLUI)
	mux.HandleFunc("GET /manifest.webmanifest", manifest)
	mux.HandleFunc("GET /sw.js", serviceWorker)
	mux.HandleFunc("GET /.well-known/security.txt", securityTxt)
	mux.HandleFunc("GET /security.txt", securityTxt)

	var handler http.Handler = mux
	handler = methodGuard(handler)
	handler = cors(handler)
	handler = limitRequests(handler, limiter)
	handler = requestLimit(handler)
	handler = normalizeURL(handler)
	return handler
}


func normalizeURL(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" && strings.HasSuffix(r.URL.Path, "/") && !strings.Contains(r.URL.Path[strings.LastIndex(r.URL.Path, "/"):], ".") {
			canonical := *r.URL
			canonical.Path = strings.TrimSuffix(canonical.Path, "/")
			http.Redirect(w, r, canonical.String(), http.StatusMovedPermanently)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func requestLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(r.URL.Path)+len(r.URL.RawQuery) > 4096 {
			writeError(w, http.StatusRequestURITooLong, "REQUEST_URI_TOO_LONG", "Request path and query exceed the 4 KiB limit", nil)
			return
		}
		if r.ContentLength > maxRequestBytes {
			writeError(w, http.StatusRequestEntityTooLarge, "REQUEST_TOO_LARGE", "Request body exceeds the 1 MiB limit", nil)
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, maxRequestBytes)
		next.ServeHTTP(w, r)
	})
}

func limitRequests(next http.Handler, limiter *rateLimiter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.enabled {
			next.ServeHTTP(w, r)
			return
		}
		ip := clientIP(r)
		allowed, retryAfter := limiter.allow(ip, time.Now())
		if !allowed {
			seconds := int(retryAfter.Seconds()) + 1
			w.Header().Set("Retry-After", strconv.Itoa(seconds))
			writeError(w, http.StatusTooManyRequests, "RATE_LIMITED", "Too many requests", nil)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return host
	}
	return r.RemoteAddr
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, HEAD, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("X-Frame-Options", "DENY")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func parsePagination(w http.ResponseWriter, r *http.Request) (int, int, bool) {
	page, limit := 1, defaultPageSize
	if value := r.URL.Query().Get("page"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed < 1 {
			writeError(w, http.StatusBadRequest, "BAD_REQUEST", "page must be a positive integer", map[string]interface{}{"field": "page"})
			return 0, 0, false
		}
		page = parsed
	}
	if value := r.URL.Query().Get("limit"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed < 1 || parsed > maximumPageSize {
			writeError(w, http.StatusBadRequest, "BAD_REQUEST", "limit must be between 1 and 1000", map[string]interface{}{"field": "limit"})
			return 0, 0, false
		}
		limit = parsed
	}
	return page, limit, true
}

func writeSuccess(w http.ResponseWriter, value interface{}) {
	writeJSON(w, http.StatusOK, response{OK: true, Data: value})
}

func writeList(w http.ResponseWriter, value interface{}, page pagination) {
	writeJSON(w, http.StatusOK, struct {
		OK         bool        `json:"ok"`
		Data       interface{} `json:"data"`
		Pagination pagination  `json:"pagination"`
	}{OK: true, Data: value, Pagination: page})
}

func writeError(w http.ResponseWriter, status int, code, message string, details map[string]interface{}) {
	writeJSON(w, status, errorResponse{OK: false, Error: code, Message: message, Details: details})
}

func writeJSON(w http.ResponseWriter, status int, value interface{}) {
	body, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		http.Error(w, `{"ok":false,"error":"INTERNAL_ERROR","message":"Unable to encode response"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	_, _ = w.Write(append(body, '\n'))
}

func staticPage(title, content string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		renderPage(w, title, template.HTML(template.HTMLEscapeString(content)), nil, 1, 1)
	}
}

func randomQuotePage(store *catalog.Catalog) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		category := r.URL.Query().Get("category")
		quote, err := store.Random(category)
		if err != nil {
			renderQuoteNotFound(w)
			return
		}
		renderQuotePage(w, quote, "Random quote", themeFor(r), category)
	}
}

func quotePage(store *catalog.Catalog) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil || id < 1 {
			renderQuoteNotFound(w)
			return
		}
		quote, err := store.Get("", id)
		if err != nil {
			renderQuoteNotFound(w)
			return
		}
		renderQuotePage(w, quote, fmt.Sprintf("Quote #%d", quote.ID), themeFor(r), "")
	}
}

func renderQuoteNotFound(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	if err := pageTemplate.Execute(w, map[string]interface{}{"title": "Quote not found", "content": "That quote is not available. Browse the collection to find another.", "quotes": nil, "page": 1, "pages": 1}); err != nil {
		http.Error(w, "Page rendering failed", http.StatusInternalServerError)
	}
}

func renderQuotePage(w http.ResponseWriter, quote catalog.Quote, title, theme, randomCategory string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := quoteTemplate.Execute(w, map[string]interface{}{"quote": quote, "title": title, "theme": theme, "randomCategory": randomCategory}); err != nil {
		http.Error(w, "Page rendering failed", http.StatusInternalServerError)
	}
}

func home(w http.ResponseWriter, r *http.Request, store *catalog.Catalog) {
	category := r.URL.Query().Get("category")
	author := strings.TrimSpace(r.URL.Query().Get("author"))
	page, limit, ok := parsePagination(w, r)
	if !ok {
		return
	}
	quotes, total, err := store.List(category, author, page, limit)
	if err != nil {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", err.Error(), nil)
		return
	}
	data := map[string]interface{}{"quotes": quotes, "categories": store.Categories(), "category": category, "author": author, "total": total, "page": page, "pages": (total + limit - 1) / limit, "theme": themeFor(r)}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := homeTemplate.Execute(w, data); err != nil {
		http.Error(w, "Page rendering failed", http.StatusInternalServerError)
	}
}

func renderPage(w http.ResponseWriter, title string, content template.HTML, quotes []catalog.Quote, page, pages int) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := pageTemplate.Execute(w, map[string]interface{}{"title": title, "content": content, "quotes": quotes, "page": page, "pages": pages}); err != nil {
		http.Error(w, "Page rendering failed", http.StatusInternalServerError)
	}
}

func themeFor(r *http.Request) string {
	if theme := r.URL.Query().Get("theme"); theme == "light" || theme == "dark" || theme == "auto" {
		return theme
	}
	return "auto"
}

func manifest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/manifest+json")
	_, _ = w.Write([]byte("{\"name\":\"Quotes API\",\"short_name\":\"Quotes\",\"start_url\":\"/\",\"display\":\"standalone\",\"background_color\":\"#101720\",\"theme_color\":\"#d78245\"}\n"))
}

func serviceWorker(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	_, _ = w.Write([]byte("self.addEventListener('install', event => event.waitUntil(caches.open('quotes-v1').then(cache => cache.addAll(['/']))));\nself.addEventListener('fetch', event => { if (event.request.method === 'GET' && new URL(event.request.url).origin === location.origin) event.respondWith(caches.match(event.request).then(cached => cached || fetch(event.request))); });\n"))
}

func securityTxt(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = fmt.Fprintln(w, "Contact: https://github.com/apimgr/quotes/security/advisories/new")
	_, _ = fmt.Fprintln(w, "Preferred-Languages: en")
}

func openAPISpec(w http.ResponseWriter, r *http.Request) {
	spec := map[string]interface{}{
		"openapi": "3.1.0",
		"info": map[string]string{"title": "Quotes API", "version": apiVersion, "description": "Public API for five embedded quote and joke collections."},
		"servers": []map[string]string{{"url": "/api/" + apiVersion}},
		"paths": map[string]interface{}{
			"/quotes": map[string]interface{}{"get": map[string]interface{}{"summary": "List and filter quotes", "parameters": []map[string]interface{}{{"name": "category", "in": "query", "schema": map[string]string{"type": "string"}}, {"name": "author", "in": "query", "schema": map[string]string{"type": "string"}}, {"name": "page", "in": "query", "schema": map[string]string{"type": "integer", "minimum": "1"}}, {"name": "limit", "in": "query", "schema": map[string]string{"type": "integer", "minimum": "1", "maximum": "1000"}}}, "responses": map[string]interface{}{"200": map[string]string{"description": "Paginated quote list"}}}},
			"/categories/{category}/quotes": map[string]interface{}{"get": map[string]interface{}{"summary": "List quotes in a category", "parameters": []map[string]interface{}{{"name": "category", "in": "path", "required": true, "schema": map[string]string{"type": "string"}}}, "responses": map[string]interface{}{"200": map[string]string{"description": "Paginated category quotes"}}}},
		"/quotes/random": map[string]interface{}{"get": map[string]interface{}{"summary": "Get a random quote", "parameters": []map[string]interface{}{{"name": "category", "in": "query", "schema": map[string]string{"type": "string"}}}, "responses": map[string]interface{}{"200": map[string]string{"description": "Random quote"}}}},
			"/quotes/{id}": map[string]interface{}{"get": map[string]interface{}{"summary": "Get quote by ID", "parameters": []map[string]interface{}{{"name": "id", "in": "path", "required": true, "schema": map[string]string{"type": "integer"}}}, "responses": map[string]interface{}{"200": map[string]string{"description": "Quote"}, "404": map[string]string{"description": "Quote not found"}}}},
			"/categories": map[string]interface{}{"get": map[string]interface{}{"summary": "List quote categories and counts", "responses": map[string]interface{}{"200": map[string]string{"description": "Category counts"}}}},
		},
	}
	writeJSON(w, http.StatusOK, spec)
}

func swaggerUI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(`<!doctype html>
<html lang="en"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width, initial-scale=1"><title>OpenAPI · Quotes</title><style>body{font:1rem/1.6 system-ui,sans-serif;max-width:60rem;margin:3rem auto;padding:0 1rem}code{background:#eee;padding:.2rem}</style></head><body><main><p><a href="/">Quotes</a></p><h1>OpenAPI</h1><p>OpenAPI 3.1 specification for the public Quotes API.</p><p><a href="/api/v1/server/swagger">View the machine-readable JSON specification</a></p><h2>GET /api/v1/quotes</h2><p>List quotes with optional <code>category</code>, <code>author</code>, <code>page</code>, and <code>limit</code> filters.</p><h2>GET /api/v1/quotes/random</h2><p>Retrieve a random quote; optionally filter by category.</p><h2>GET /api/v1/quotes/{id}</h2><p>Look up a quote by its globally unique integer ID.</p><h2>GET /api/v1/categories</h2><p>List category names and collection counts.</p></main></body></html>
`))
}

func graphQLUI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(`<!doctype html>
<html lang="en"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width, initial-scale=1"><title>GraphQL · Quotes</title><style>body{font:1rem/1.5 system-ui,sans-serif;max-width:60rem;margin:2rem auto;padding:0 1rem}textarea{width:100%;min-height:10rem;font:1rem monospace}pre{white-space:pre-wrap;background:#eee;padding:1rem}button{padding:.7rem}</style></head><body><main><p><a href="/">Quotes</a></p><h1>GraphQL explorer</h1><form id="query"><label for="source">Query</label><textarea id="source">{ randomQuote { id quote author category } }</textarea><button>Run query</button></form><pre id="result" aria-live="polite"></pre><script>document.getElementById('query').addEventListener('submit', async event => { event.preventDefault(); const response = await fetch('/api/v1/server/graphql', {method:'POST', headers:{'Content-Type':'application/json'}, body:JSON.stringify({query:document.getElementById('source').value})}); document.getElementById('result').textContent = JSON.stringify(await response.json(), null, 2); });</script></main></body></html>
`))
}

func newGraphQLHandler(store *catalog.Catalog) http.HandlerFunc {
	schema, schemaErr := buildGraphQLSchema(store)
	return func(w http.ResponseWriter, r *http.Request) {
		if schemaErr != nil {
			writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "GraphQL service is unavailable", nil)
			return
		}
		mediaType, _, mediaTypeErr := mime.ParseMediaType(r.Header.Get("Content-Type"))
		if mediaTypeErr != nil || mediaType != "application/json" {
			writeError(w, http.StatusUnsupportedMediaType, "UNSUPPORTED_MEDIA_TYPE", "GraphQL requests must use application/json", nil)
			return
		}
		var input struct {
			Query         string                 `json:"query"`
			Variables     map[string]interface{} `json:"variables"`
			OperationName string                 `json:"operationName"`
		}
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&input); err != nil {
			var maxBytesError *http.MaxBytesError
			if errors.As(err, &maxBytesError) {
				writeError(w, http.StatusRequestEntityTooLarge, "REQUEST_TOO_LARGE", "Request body exceeds the 1 MiB limit", nil)
				return
			}
			writeError(w, http.StatusBadRequest, "BAD_REQUEST", "A valid GraphQL request body is required", nil)
			return
		}
		if strings.TrimSpace(input.Query) == "" {
			writeError(w, http.StatusBadRequest, "BAD_REQUEST", "A GraphQL query is required", nil)
			return
		}
		result := graphql.Do(graphql.Params{Schema: schema, RequestString: input.Query, VariableValues: input.Variables, OperationName: input.OperationName})
		if len(result.Errors) > 0 {
			errors := make([]map[string]string, 0, len(result.Errors))
			for _, queryError := range result.Errors {
				errors = append(errors, map[string]string{"message": queryError.Message})
			}
			writeJSON(w, http.StatusOK, map[string]interface{}{"errors": errors, "data": result.Data})
			return
		}
		writeJSON(w, http.StatusOK, map[string]interface{}{"data": result.Data})
	}
}


func buildGraphQLSchema(store *catalog.Catalog) (graphql.Schema, error) {
	quoteType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Quote",
		Fields: graphql.Fields{
			"id":       &graphql.Field{Type: graphql.NewNonNull(graphql.Int)},
			"quote":    &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
			"author":   &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
			"category": &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		},
	})
	categoryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Category",
		Fields: graphql.Fields{
			"name":  &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
			"count": &graphql.Field{Type: graphql.NewNonNull(graphql.Int)},
		},
	})
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"quote": &graphql.Field{
				Type: quoteType,
				Args: graphql.FieldConfigArgument{"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)}},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id, ok := p.Args["id"].(int)
					if !ok {
						return nil, fmt.Errorf("id must be an integer")
					}
					quote, err := store.Get("", id)
					if err != nil {
						return nil, nil
					}
					return quote, nil
				},
			},
			"randomQuote": &graphql.Field{
				Type: quoteType,
				Args: graphql.FieldConfigArgument{"category": &graphql.ArgumentConfig{Type: graphql.String}},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					category, _ := p.Args["category"].(string)
					return store.Random(category)
				},
			},
			"quotes": &graphql.Field{
				Type: graphql.NewList(graphql.NewNonNull(quoteType)),
				Args: graphql.FieldConfigArgument{
					"category": &graphql.ArgumentConfig{Type: graphql.String},
					"author":   &graphql.ArgumentConfig{Type: graphql.String},
					"page":     &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 1},
					"limit":    &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: defaultPageSize},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					category, _ := p.Args["category"].(string)
					author, _ := p.Args["author"].(string)
					page, _ := p.Args["page"].(int)
					limit, _ := p.Args["limit"].(int)
					if page < 1 || limit < 1 || limit > maximumPageSize {
						return nil, fmt.Errorf("page must be positive and limit must be between 1 and %d", maximumPageSize)
					}
					quotes, _, err := store.List(category, author, page, limit)
					return quotes, err
				},
			},
			"categories": &graphql.Field{
				Type: graphql.NewList(graphql.NewNonNull(categoryType)),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					counts := store.Categories()
					categories := make([]map[string]interface{}, 0, len(catalog.Names()))
					for _, name := range catalog.Names() {
						categories = append(categories, map[string]interface{}{"name": name, "count": counts[name]})
					}
					return categories, nil
				},
			},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{Query: queryType})
	if err != nil {
		return graphql.Schema{}, fmt.Errorf("initialize GraphQL schema: %w", err)
	}
	return schema, nil
}

func negotiateText(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept"), "text/plain") {
			next.ServeHTTP(w, r)
			return
		}
		buffer := newResponseBuffer()
		next.ServeHTTP(buffer, r)
		if buffer.status < 200 || buffer.status >= 300 || !strings.HasPrefix(buffer.header.Get("Content-Type"), "application/json") {
			copyResponse(w, buffer)
			return
		}
		var body struct {
			Data json.RawMessage `json:"data"`
		}
		if err := json.Unmarshal(buffer.body, &body); err != nil {
			copyResponse(w, buffer)
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		if len(body.Data) == 0 || string(body.Data) == "null" {
			_, _ = w.Write([]byte("\n"))
			return
		}
		var quotes []catalog.Quote
		if err := json.Unmarshal(body.Data, &quotes); err == nil {
			for _, quote := range quotes {
				writePlainQuote(w, quote)
			}
			return
		}
		var quote catalog.Quote
		if err := json.Unmarshal(body.Data, &quote); err != nil {
			copyResponse(w, buffer)
			return
		}
		writePlainQuote(w, quote)
	})
}

type responseBuffer struct {
	header http.Header
	status int
	body   []byte
}

func newResponseBuffer() *responseBuffer {
	return &responseBuffer{header: make(http.Header)}
}

func (b *responseBuffer) Header() http.Header {
	return b.header
}

func (b *responseBuffer) WriteHeader(status int) {
	if b.status == 0 {
		b.status = status
	}
}

func (b *responseBuffer) Write(content []byte) (int, error) {
	if b.status == 0 {
		b.status = http.StatusOK
	}
	b.body = append(b.body, content...)
	return len(content), nil
}

func copyResponse(w http.ResponseWriter, source *responseBuffer) {
	for key, values := range source.header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	status := source.status
	if status == 0 {
		status = http.StatusOK
	}
	w.WriteHeader(status)
	_, _ = w.Write(source.body)
}

func writePlainQuote(w http.ResponseWriter, quote catalog.Quote) {
	if quote.Author == "" {
		_, _ = fmt.Fprintf(w, "%s\n", quote.Quote)
		return
	}
	_, _ = fmt.Fprintf(w, "%q - %s\n", quote.Quote, quote.Author)
}

func listText(w http.ResponseWriter, r *http.Request, store *catalog.Catalog, category string) {
	page, limit, ok := parsePagination(w, r)
	if !ok {
		return
	}
	quotes, _, err := store.List(category, strings.TrimSpace(r.URL.Query().Get("author")), page, limit)
	if err != nil {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", err.Error(), nil)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	for _, quote := range quotes {
		writePlainQuote(w, quote)
	}
}

func methodGuard(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead && r.Method != http.MethodPost && r.Method != http.MethodOptions {
			w.Header().Set("Allow", "GET, HEAD, POST, OPTIONS")
			writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "The requested method is not supported", nil)
			return
		}
		if r.Method == http.MethodPost && r.URL.Path != "/api/"+apiVersion+"/server/graphql" && r.URL.Path != "/api/graphql" {
			w.Header().Set("Allow", "GET, HEAD, OPTIONS")
			writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "The requested method is not supported", nil)
			return
		}
		next.ServeHTTP(w, r)
	})
}

type healthResponse struct {
	Project   map[string]string `json:"project"`
	Status    string            `json:"status"`
	Version   string            `json:"version"`
	GoVersion string            `json:"go_version"`
	Uptime    string            `json:"uptime"`
	Mode      string            `json:"mode"`
	Timestamp string            `json:"timestamp"`
	Features  map[string]int    `json:"features"`
	Total     int               `json:"total_quotes"`
}

var processStarted = time.Now()

func healthStatus(store *catalog.Catalog, version string) healthResponse {
	counts := store.Categories()
	total := 0
	for _, count := range counts {
		total += count
	}
	return healthResponse{
		Project: map[string]string{"name": "quotes", "description": "A public, read-only quote and joke collection."},
		Status: "healthy",
		Version: version,
		GoVersion: runtime.Version(),
		Uptime: time.Since(processStarted).Round(time.Second).String(),
		Mode: "production",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Features: counts,
		Total: total,
	}
}
