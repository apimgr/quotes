package server

import (
	"context"
	"embed"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
)

//go:embed static templates
var content embed.FS

var (
	// Version information (set by build flags)
	Version   = "dev"
	Commit    = "unknown"
	BuildDate = "unknown"
)

// Server represents the HTTP server with SPEC-compliant configuration
type Server struct {
	router         *chi.Mux
	port           string
	address        string
	settingsCache  map[string]interface{}
	settingsMutex  sync.RWMutex
	rateLimiters   map[string]*httprate.RateLimiter
	server         *http.Server
}

// NewServer creates a new server instance with Chi router
func NewServer(port, address string) *Server {
	s := &Server{
		router:        chi.NewRouter(),
		port:          port,
		address:       address,
		settingsCache: make(map[string]interface{}),
		rateLimiters:  make(map[string]*httprate.RateLimiter),
	}

	// Initialize default settings
	s.initDefaultSettings()

	// Setup middleware and routes
	s.setupMiddleware()
	s.setupRoutes()

	return s
}

// initDefaultSettings initializes default server settings
func (s *Server) initDefaultSettings() {
	s.settingsMutex.Lock()
	defer s.settingsMutex.Unlock()

	// CORS settings (default: allow all)
	s.settingsCache["server.cors_enabled"] = true
	s.settingsCache["server.cors_origins"] = []string{"*"}
	s.settingsCache["server.cors_methods"] = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	s.settingsCache["server.cors_headers"] = []string{"Content-Type", "Authorization"}
	s.settingsCache["server.cors_credentials"] = false

	// Rate limiting (default: enabled)
	s.settingsCache["rate.enabled"] = true
	s.settingsCache["rate.global_rps"] = 100
	s.settingsCache["rate.global_burst"] = 200
	s.settingsCache["rate.api_rps"] = 50
	s.settingsCache["rate.api_burst"] = 100
	s.settingsCache["rate.admin_rps"] = 10
	s.settingsCache["rate.admin_burst"] = 20

	// Initialize rate limiters
	s.rateLimiters["global"] = httprate.NewRateLimiter(100, time.Second)
	s.rateLimiters["api"] = httprate.NewRateLimiter(50, time.Second)
	s.rateLimiters["admin"] = httprate.NewRateLimiter(10, time.Second)
}

// setupMiddleware configures all middleware
func (s *Server) setupMiddleware() {
	// Recovery middleware (must be first)
	s.router.Use(middleware.Recoverer)

	// Request ID middleware
	s.router.Use(middleware.RequestID)

	// Real IP middleware
	s.router.Use(middleware.RealIP)

	// Logger middleware
	s.router.Use(middleware.Logger)

	// Timeout middleware (60 seconds)
	s.router.Use(middleware.Timeout(60 * time.Second))

	// Throttle concurrent requests (max 1000)
	s.router.Use(middleware.Throttle(1000))

	// Security headers middleware
	s.router.Use(s.securityHeadersMiddleware)

	// CORS middleware
	s.router.Use(s.corsMiddleware)

	// Global rate limiting
	s.router.Use(s.rateLimitMiddleware("global"))
}

// securityHeadersMiddleware adds security headers to all responses
func (s *Server) securityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Prevent clickjacking
		w.Header().Set("X-Frame-Options", "DENY")

		// Prevent MIME sniffing
		w.Header().Set("X-Content-Type-Options", "nosniff")

		// XSS Protection
		w.Header().Set("X-XSS-Protection", "1; mode=block")

		// Content Security Policy
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; "+
			"script-src 'self' 'unsafe-inline'; "+
			"style-src 'self' 'unsafe-inline'; "+
			"img-src 'self' data:; "+
			"connect-src 'self'")

		// Referrer Policy
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// Permissions Policy
		w.Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		// HSTS (if using HTTPS)
		if r.TLS != nil {
			w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}

		next.ServeHTTP(w, r)
	})
}

// corsMiddleware handles CORS based on settings
func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.settingsMutex.RLock()
		enabled := s.settingsCache["server.cors_enabled"].(bool)
		origins := s.settingsCache["server.cors_origins"].([]string)
		methods := s.settingsCache["server.cors_methods"].([]string)
		headers := s.settingsCache["server.cors_headers"].([]string)
		s.settingsMutex.RUnlock()

		if !enabled {
			next.ServeHTTP(w, r)
			return
		}

		origin := r.Header.Get("Origin")

		// Check if origin is allowed
		allowed := false
		for _, o := range origins {
			if o == "*" || o == origin {
				allowed = true
				break
			}
		}

		if allowed {
			if origins[0] == "*" {
				w.Header().Set("Access-Control-Allow-Origin", "*")
			} else {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			}

			// Join methods and headers
			methodsStr := ""
			for i, m := range methods {
				if i > 0 {
					methodsStr += ", "
				}
				methodsStr += m
			}

			headersStr := ""
			for i, h := range headers {
				if i > 0 {
					headersStr += ", "
				}
				headersStr += h
			}

			w.Header().Set("Access-Control-Allow-Methods", methodsStr)
			w.Header().Set("Access-Control-Allow-Headers", headersStr)
		}

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// rateLimitMiddleware applies rate limiting
func (s *Server) rateLimitMiddleware(limiterName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			s.settingsMutex.RLock()
			enabled := s.settingsCache["rate.enabled"].(bool)
			limiter, exists := s.rateLimiters[limiterName]
			s.settingsMutex.RUnlock()

			if !enabled || !exists {
				next.ServeHTTP(w, r)
				return
			}

			// Apply rate limiting
			limiter.Handler(next).ServeHTTP(w, r)
		})
	}
}

// setupRoutes configures all the routes
func (s *Server) setupRoutes() {
	// API routes with stricter rate limiting
	s.router.Route("/api/v1", func(r chi.Router) {
		r.Use(s.rateLimitMiddleware("api"))

		// Quotes endpoints
		r.Get("/random", handleRandomQuote)
		r.Get("/quotes", handleAllQuotes)
		r.Get("/quotes/{id:[0-9]+}", handleQuoteByID)
		r.Get("/quotes/category/{category}", handleQuotesByCategory)
		r.Get("/quotes/author/{author}", handleQuotesByAuthor)
		r.Get("/status", handleStatus)

		// Anime endpoints
		r.Get("/anime", handleAllAnimeQuotes)
		r.Get("/anime/random", handleRandomAnimeQuote)
		r.Get("/anime/{id:[0-9]+}", handleAnimeQuoteByID)
		r.Get("/anime/category/{category}", handleAnimeQuotesByCategory)
		r.Get("/anime/show/{anime}", handleAnimeQuotesByAnime)
		r.Get("/anime/character/{character}", handleAnimeQuotesByCharacter)

		// Chuck Norris endpoints
		r.Get("/chucknorris", handleAllChuckNorrisJokes)
		r.Get("/chucknorris/random", handleRandomChuckNorrisJoke)

		// Dad Jokes endpoints
		r.Get("/dadjokes", handleAllDadJokes)
		r.Get("/dadjokes/random", handleRandomDadJoke)

		// Programming Jokes endpoints
		r.Get("/programming", handleAllProgrammingJokes)
		r.Get("/programming/random", handleRandomProgrammingJoke)

		// JSON file endpoints
		r.Get("/{file:.*\\.json}", handleJSONFile)

		// Admin routes (most restrictive rate limiting)
		r.Route("/admin", func(r chi.Router) {
			r.Use(s.rateLimitMiddleware("admin"))
			r.Use(authMiddleware)
			r.Get("/settings", handleGetSettings)
			r.Post("/settings", handleSetSetting)
			r.Delete("/settings/{key}", handleDeleteSetting)
		})
	})

	// Shorthand routes (without /api/v1 prefix)
	s.router.Get("/anime", handleAllAnimeQuotes)
	s.router.Get("/anime/random", handleRandomAnimeQuote)
	s.router.Get("/chucknorris", handleAllChuckNorrisJokes)
	s.router.Get("/chucknorris/random", handleRandomChuckNorrisJoke)
	s.router.Get("/dadjokes", handleAllDadJokes)
	s.router.Get("/dadjokes/random", handleRandomDadJoke)
	s.router.Get("/programming", handleAllProgrammingJokes)
	s.router.Get("/programming/random", handleRandomProgrammingJoke)

	// Static files
	fileServer := http.FileServer(http.FS(content))
	s.router.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	// Web UI routes
	s.router.Get("/", handleHome)
	s.router.Get("/admin", handleAdminPage)

	// Health check
	s.router.Get("/health", handleHealth)
	s.router.Get("/healthz", handleHealth)
}

// Start starts the HTTP server with graceful shutdown support
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%s", s.address, s.port)

	s.server = &http.Server{
		Addr:           addr,
		Handler:        s.router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1MB
	}

	log.Printf("Starting Quotes API Server v%s", Version)
	log.Printf("Server listening on %s", addr)
	log.Printf("API endpoint: http://%s/api/v1/random", addr)
	log.Printf("Web UI: http://%s/", addr)
	log.Printf("Admin panel: http://%s/admin", addr)

	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	if s.server != nil {
		return s.server.Shutdown(ctx)
	}
	return nil
}

// handleHealth handles health check requests
func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":"ok","version":"%s"}`, Version)
}

// handleStatus handles status check requests
func handleStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":"ok","version":"%s","commit":"%s","build_date":"%s"}`, Version, Commit, BuildDate)
}
