package server

import (
	"embed"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

//go:embed static templates
var content embed.FS

var (
	// Version information (set by build flags)
	Version   = "dev"
	Commit    = "unknown"
	BuildDate = "unknown"
)

// Server represents the HTTP server
type Server struct {
	router *mux.Router
	port   string
}

// NewServer creates a new server instance
func NewServer(port string) *Server {
	s := &Server{
		router: mux.NewRouter(),
		port:   port,
	}

	s.setupRoutes()
	return s
}

// setupRoutes configures all the routes
func (s *Server) setupRoutes() {
	// API routes
	api := s.router.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/random", handleRandomQuote).Methods("GET")
	api.HandleFunc("/quotes", handleAllQuotes).Methods("GET")
	api.HandleFunc("/quotes/{id:[0-9]+}", handleQuoteByID).Methods("GET")
	api.HandleFunc("/quotes/category/{category}", handleQuotesByCategory).Methods("GET")
	api.HandleFunc("/quotes/author/{author}", handleQuotesByAuthor).Methods("GET")
	api.HandleFunc("/status", handleStatus).Methods("GET")

	// Anime API routes
	api.HandleFunc("/anime", handleAllAnimeQuotes).Methods("GET")
	api.HandleFunc("/anime/random", handleRandomAnimeQuote).Methods("GET")
	api.HandleFunc("/anime/{id:[0-9]+}", handleAnimeQuoteByID).Methods("GET")
	api.HandleFunc("/anime/category/{category}", handleAnimeQuotesByCategory).Methods("GET")
	api.HandleFunc("/anime/show/{anime}", handleAnimeQuotesByAnime).Methods("GET")
	api.HandleFunc("/anime/character/{character}", handleAnimeQuotesByCharacter).Methods("GET")

	// Chuck Norris API routes
	api.HandleFunc("/chucknorris", handleAllChuckNorrisJokes).Methods("GET")
	api.HandleFunc("/chucknorris/random", handleRandomChuckNorrisJoke).Methods("GET")

	// Dad Jokes API routes
	api.HandleFunc("/dadjokes", handleAllDadJokes).Methods("GET")
	api.HandleFunc("/dadjokes/random", handleRandomDadJoke).Methods("GET")

	// Programming Jokes API routes
	api.HandleFunc("/programming", handleAllProgrammingJokes).Methods("GET")
	api.HandleFunc("/programming/random", handleRandomProgrammingJoke).Methods("GET")

	// JSON file endpoints
	api.HandleFunc("/{file:.*\\.json}", handleJSONFile).Methods("GET")

	// Shorthand routes (without /api/v1 prefix)
	s.router.HandleFunc("/anime", handleAllAnimeQuotes).Methods("GET")
	s.router.HandleFunc("/anime/random", handleRandomAnimeQuote).Methods("GET")
	s.router.HandleFunc("/chucknorris", handleAllChuckNorrisJokes).Methods("GET")
	s.router.HandleFunc("/chucknorris/random", handleRandomChuckNorrisJoke).Methods("GET")
	s.router.HandleFunc("/dadjokes", handleAllDadJokes).Methods("GET")
	s.router.HandleFunc("/dadjokes/random", handleRandomDadJoke).Methods("GET")
	s.router.HandleFunc("/programming", handleAllProgrammingJokes).Methods("GET")
	s.router.HandleFunc("/programming/random", handleRandomProgrammingJoke).Methods("GET")

	// Admin routes (protected)
	admin := s.router.PathPrefix("/api/v1/admin").Subrouter()
	admin.Use(authMiddleware)
	admin.HandleFunc("/settings", handleGetSettings).Methods("GET")
	admin.HandleFunc("/settings", handleSetSetting).Methods("POST")
	admin.HandleFunc("/settings/{key}", handleDeleteSetting).Methods("DELETE")

	// Static files
	s.router.PathPrefix("/static/").Handler(http.FileServer(http.FS(content)))

	// Web UI routes
	s.router.HandleFunc("/", handleHome).Methods("GET")
	s.router.HandleFunc("/admin", handleAdminPage).Methods("GET")

	// Health check
	s.router.HandleFunc("/health", handleHealth).Methods("GET")
}

// Start starts the HTTP server
func (s *Server) Start(address string) error {
	addr := fmt.Sprintf("%s:%s", address, s.port)
	log.Printf("Starting Quotes API Server v%s", Version)
	log.Printf("Server listening on %s", addr)
	log.Printf("API endpoint: http://%s/api/v1/random", addr)
	log.Printf("Web UI: http://%s/", addr)
	log.Printf("Admin panel: http://%s/admin", addr)

	return http.ListenAndServe(addr, s.router)
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
