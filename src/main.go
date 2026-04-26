package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/apimgr/quotes/src/admin"
	"github.com/apimgr/quotes/src/anime"
	"github.com/apimgr/quotes/src/chucknorris"
	"github.com/apimgr/quotes/src/config"
	"github.com/apimgr/quotes/src/dadjokes"
	"github.com/apimgr/quotes/src/mode"
	"github.com/apimgr/quotes/src/paths"
	"github.com/apimgr/quotes/src/programming"
	"github.com/apimgr/quotes/src/quotes"
)

//go:embed data/quotes.json
var quotesData []byte

//go:embed data/anime.json
var animeData []byte

//go:embed data/chucknorris.json
var chuckNorrisData []byte

//go:embed data/dadjokes.json
var dadJokesData []byte

//go:embed data/programming.json
var programmingData []byte

var (
	// Version information (set by build flags)
	Version   = "dev"
	Commit    = "unknown"
	BuildDate = "unknown"
)

func main() {
	// Command-line flags
	port := flag.String("port", "", "Server port (overrides config)")
	address := flag.String("address", "", "Server address (overrides config)")
	configFile := flag.String("config", "", "Path to config file")
	dataPath := flag.String("data", "", "Path to data directory")
	logsPath := flag.String("logs", "", "Path to logs directory")
	showVersion := flag.Bool("version", false, "Show version information")
	showStatus := flag.Bool("status", false, "Show status (for health checks)")
	showHelp := flag.Bool("help", false, "Show help message")
	serviceCmd := flag.String("service", "", "Service command: install, uninstall, start, stop, restart, status")
	maintenanceCmd := flag.String("maintenance", "", "Maintenance mode: on, off")
	modeFlag := flag.String("mode", "", "Application mode (dev/development, prod/production)")
	updateCmd := flag.String("update", "", "Update command (stable, beta, nightly)")
	flag.Parse()

	// Handle update command
	if *updateCmd != "" {
		handleUpdateCommand(*updateCmd)
		os.Exit(0)
	}

	// Initialize mode from CLI flag first, then env var
	if *modeFlag != "" {
		if err := mode.Set(*modeFlag); err != nil {
			log.Printf("Warning: invalid mode: %v", err)
		}
	}
	mode.Initialize() // Check env var if mode wasn't set via flag

	// Unused vars to satisfy compiler
	_ = dataPath
	_ = logsPath

	// Show help
	if *showHelp {
		printHelp()
		os.Exit(0)
	}

	// Show version
	if *showVersion {
		fmt.Printf("Quotes API v%s\n", Version)
		fmt.Printf("Commit: %s\n", Commit)
		fmt.Printf("Build Date: %s\n", BuildDate)
		fmt.Printf("Go Version: %s\n", runtime.Version())
		fmt.Printf("OS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
		os.Exit(0)
	}

	// Status check (for health checks)
	if *showStatus {
		fmt.Println("OK")
		os.Exit(0)
	}

	// Get directories
	dirs := paths.GetDirectories()

	// Handle service commands
	if *serviceCmd != "" {
		handleServiceCommand(*serviceCmd, dirs)
		os.Exit(0)
	}

	// Handle maintenance mode
	if *maintenanceCmd != "" {
		handleMaintenanceCommand(*maintenanceCmd, dirs)
		os.Exit(0)
	}

	// Ensure directories exist
	if err := paths.EnsureDirectories(dirs); err != nil {
		log.Fatalf("Failed to create directories: %v", err)
	}

	// Determine config file path
	cfgPath := *configFile
	if cfgPath == "" {
		cfgPath = filepath.Join(dirs.Config, "server.yml")
	}

	// Load configuration
	cfg, err := config.Load(cfgPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Apply command-line overrides
	serverPort := cfg.Server.Port
	if *port != "" {
		serverPort = *port
	}
	if serverPort == "" {
		serverPort = "8080"
	}

	serverAddress := cfg.Server.Address
	if *address != "" {
		serverAddress = *address
	}
	if serverAddress == "" {
		serverAddress = "0.0.0.0"
	}

	log.Printf("Starting Quotes API v%s", Version)
	log.Printf("Config directory: %s", dirs.Config)
	log.Printf("Data directory: %s", dirs.Data)
	log.Printf("Logs directory: %s", dirs.Logs)

	// Load all quote collections
	loadCollections()

	// Create HTTP server
	mux := http.NewServeMux()

	// Initialize admin handler
	adminHandler := admin.NewHandler(
		cfg.Server.Admin.Username,
		cfg.Server.Admin.Password,
		cfg.Server.Admin.APIToken,
		cfg.Server.Session.Timeout,
		false, // SSL enabled
		Version,
		Commit,
		BuildDate,
	)
	adminHandler.RegisterRoutes(mux)

	setupRoutes(mux, cfg)

	addr := fmt.Sprintf("%s:%s", serverAddress, serverPort)
	server := &http.Server{
		Addr:           addr,
		Handler:        corsMiddleware(mux, cfg),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1MB
	}

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Println("Shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Server shutdown error: %v", err)
		}
	}()

	log.Printf("Server listening on %s", addr)
	log.Printf("API endpoint: http://%s/api/v1/random", addr)
	log.Printf("Web UI: http://%s/", addr)

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Server failed: %v", err)
	}

	log.Println("Server stopped")
}

func loadCollections() {
	// Load quotes
	log.Println("Loading quotes...")
	if err := quotes.LoadQuotes(quotesData); err != nil {
		log.Fatalf("Failed to load quotes: %v", err)
	}
	log.Printf("Loaded %d quotes", quotes.GetTotalCount())

	// Load anime quotes
	log.Println("Loading anime quotes...")
	if err := anime.LoadQuotes(animeData); err != nil {
		log.Fatalf("Failed to load anime quotes: %v", err)
	}
	log.Printf("Loaded %d anime quotes", anime.GetTotalCount())

	// Load Chuck Norris jokes
	log.Println("Loading Chuck Norris jokes...")
	if err := chucknorris.LoadJokes(chuckNorrisData); err != nil {
		log.Fatalf("Failed to load Chuck Norris jokes: %v", err)
	}
	log.Printf("Loaded %d Chuck Norris jokes", chucknorris.GetTotalCount())

	// Load dad jokes
	log.Println("Loading dad jokes...")
	if err := dadjokes.LoadJokes(dadJokesData); err != nil {
		log.Fatalf("Failed to load dad jokes: %v", err)
	}
	log.Printf("Loaded %d dad jokes", dadjokes.GetTotalCount())

	// Load programming jokes
	log.Println("Loading programming jokes...")
	if err := programming.LoadJokes(programmingData); err != nil {
		log.Fatalf("Failed to load programming jokes: %v", err)
	}
	log.Printf("Loaded %d programming jokes", programming.GetTotalCount())
}

func setupRoutes(mux *http.ServeMux, cfg *config.Config) {
	// Health check
	mux.HandleFunc("/healthz", handleHealth)
	mux.HandleFunc("/health", handleHealth)

	// API v1 - Quotes
	mux.HandleFunc("/api/v1/random", handleRandomQuote)
	mux.HandleFunc("/api/v1/random.txt", handleRandomQuoteTxt)
	mux.HandleFunc("/api/v1/quotes", handleAllQuotes)
	mux.HandleFunc("/api/v1/quotes.txt", handleAllQuotesTxt)
	mux.HandleFunc("/api/v1/quotes/", handleQuoteByID)
	mux.HandleFunc("/api/v1/stats", handleStats)
	mux.HandleFunc("/api/v1/stats.txt", handleStatsTxt)
	mux.HandleFunc("/api/v1/count", handleCount)
	mux.HandleFunc("/api/v1/count.txt", handleCountTxt)

	// API v1 - Anime
	mux.HandleFunc("/api/v1/anime/random", handleRandomAnime)
	mux.HandleFunc("/api/v1/anime/random.txt", handleRandomAnimeTxt)
	mux.HandleFunc("/api/v1/anime", handleAllAnime)
	mux.HandleFunc("/api/v1/anime.txt", handleAllAnimeTxt)

	// API v1 - Chuck Norris
	mux.HandleFunc("/api/v1/chucknorris/random", handleRandomChuckNorris)
	mux.HandleFunc("/api/v1/chucknorris/random.txt", handleRandomChuckNorrisTxt)
	mux.HandleFunc("/api/v1/chucknorris", handleAllChuckNorris)
	mux.HandleFunc("/api/v1/chucknorris.txt", handleAllChuckNorrisTxt)

	// API v1 - Dad Jokes
	mux.HandleFunc("/api/v1/dadjokes/random", handleRandomDadJoke)
	mux.HandleFunc("/api/v1/dadjokes/random.txt", handleRandomDadJokeTxt)
	mux.HandleFunc("/api/v1/dadjokes", handleAllDadJokes)
	mux.HandleFunc("/api/v1/dadjokes.txt", handleAllDadJokesTxt)

	// API v1 - Programming
	mux.HandleFunc("/api/v1/programming/random", handleRandomProgramming)
	mux.HandleFunc("/api/v1/programming/random.txt", handleRandomProgrammingTxt)
	mux.HandleFunc("/api/v1/programming", handleAllProgramming)
	mux.HandleFunc("/api/v1/programming.txt", handleAllProgrammingTxt)

	// Shorthand routes
	mux.HandleFunc("/random", handleRandomQuote)
	mux.HandleFunc("/random.txt", handleRandomQuoteTxt)

	// Web UI
	mux.HandleFunc("/", handleHome)
	mux.HandleFunc("/docs", handleDocs)

	// PWA support
	mux.HandleFunc("/manifest.json", handleManifest)
	mux.HandleFunc("/sw.js", handleServiceWorker)

	// robots.txt and security.txt
	mux.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		handleRobotsTxt(w, r, cfg)
	})
	mux.HandleFunc("/security.txt", func(w http.ResponseWriter, r *http.Request) {
		handleSecurityTxt(w, r, cfg)
	})
	mux.HandleFunc("/.well-known/security.txt", func(w http.ResponseWriter, r *http.Request) {
		handleSecurityTxt(w, r, cfg)
	})
}

// CORS middleware
func corsMiddleware(next http.Handler, cfg *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cors := cfg.WebSecurity.CORS
		if cors == "" {
			cors = "*"
		}
		w.Header().Set("Access-Control-Allow-Origin", cors)
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Health check handler
func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "ok",
		"version": Version,
	})
}

// Quote handlers
func handleRandomQuote(w http.ResponseWriter, r *http.Request) {
	quote, err := quotes.GetRandomQuote()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    quote,
	})
}

func handleRandomQuoteTxt(w http.ResponseWriter, r *http.Request) {
	quote, err := quotes.GetRandomQuote()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "\"%s\" - %s", quote.Quote, quote.Author)
}

func handleAllQuotes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    quotes.GetAllQuotes(),
		"count":   quotes.GetTotalCount(),
	})
}

func handleAllQuotesTxt(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	for _, q := range quotes.GetAllQuotes() {
		fmt.Fprintf(w, "\"%s\" - %s\n", q.Quote, q.Author)
	}
}

func handleQuoteByID(w http.ResponseWriter, r *http.Request) {
	// Extract ID from path /api/v1/quotes/{id}
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/quotes/")
	path = strings.TrimSuffix(path, ".txt")
	id, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, "Invalid quote ID", http.StatusBadRequest)
		return
	}

	quote, err := quotes.GetQuoteByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if strings.HasSuffix(r.URL.Path, ".txt") {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintf(w, "\"%s\" - %s", quote.Quote, quote.Author)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    quote,
	})
}

// Anime handlers
func handleRandomAnime(w http.ResponseWriter, r *http.Request) {
	quote, err := anime.GetRandomQuote()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    quote,
	})
}

func handleRandomAnimeTxt(w http.ResponseWriter, r *http.Request) {
	quote, err := anime.GetRandomQuote()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "\"%s\" - %s (%s)", quote.Quote, quote.Character, quote.Anime)
}

func handleAllAnime(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    anime.GetAllQuotes(),
		"count":   anime.GetTotalCount(),
	})
}

func handleAllAnimeTxt(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	for _, q := range anime.GetAllQuotes() {
		fmt.Fprintf(w, "\"%s\" - %s (%s)\n", q.Quote, q.Character, q.Anime)
	}
}

// Chuck Norris handlers
func handleRandomChuckNorris(w http.ResponseWriter, r *http.Request) {
	joke, err := chucknorris.GetRandomJoke()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    joke,
	})
}

func handleRandomChuckNorrisTxt(w http.ResponseWriter, r *http.Request) {
	joke, err := chucknorris.GetRandomJoke()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, joke.Joke)
}

func handleAllChuckNorris(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    chucknorris.GetAllJokes(),
		"count":   chucknorris.GetTotalCount(),
	})
}

func handleAllChuckNorrisTxt(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	for _, j := range chucknorris.GetAllJokes() {
		fmt.Fprintln(w, j.Joke)
	}
}

// Dad Jokes handlers
func handleRandomDadJoke(w http.ResponseWriter, r *http.Request) {
	joke, err := dadjokes.GetRandomJoke()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    joke,
	})
}

func handleRandomDadJokeTxt(w http.ResponseWriter, r *http.Request) {
	joke, err := dadjokes.GetRandomJoke()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, joke.Joke)
}

func handleAllDadJokes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    dadjokes.GetAllJokes(),
		"count":   dadjokes.GetTotalCount(),
	})
}

func handleAllDadJokesTxt(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	for _, j := range dadjokes.GetAllJokes() {
		fmt.Fprintln(w, j.Joke)
	}
}

// Programming handlers
func handleRandomProgramming(w http.ResponseWriter, r *http.Request) {
	joke, err := programming.GetRandomJoke()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    joke,
	})
}

func handleRandomProgrammingTxt(w http.ResponseWriter, r *http.Request) {
	joke, err := programming.GetRandomJoke()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, joke.Joke)
}

func handleAllProgramming(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    programming.GetAllJokes(),
		"count":   programming.GetTotalCount(),
	})
}

func handleAllProgrammingTxt(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	for _, j := range programming.GetAllJokes() {
		fmt.Fprintln(w, j.Joke)
	}
}

// Stats handlers
func handleStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"quotes":       quotes.GetTotalCount(),
			"anime":        anime.GetTotalCount(),
			"chucknorris":  chucknorris.GetTotalCount(),
			"dadjokes":     dadjokes.GetTotalCount(),
			"programming":  programming.GetTotalCount(),
			"total":        quotes.GetTotalCount() + anime.GetTotalCount() + chucknorris.GetTotalCount() + dadjokes.GetTotalCount() + programming.GetTotalCount(),
			"version":      Version,
			"commit":       Commit,
			"build_date":   BuildDate,
		},
	})
}

func handleStatsTxt(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	total := quotes.GetTotalCount() + anime.GetTotalCount() + chucknorris.GetTotalCount() + dadjokes.GetTotalCount() + programming.GetTotalCount()
	fmt.Fprintf(w, "Quotes: %d\n", quotes.GetTotalCount())
	fmt.Fprintf(w, "Anime: %d\n", anime.GetTotalCount())
	fmt.Fprintf(w, "Chuck Norris: %d\n", chucknorris.GetTotalCount())
	fmt.Fprintf(w, "Dad Jokes: %d\n", dadjokes.GetTotalCount())
	fmt.Fprintf(w, "Programming: %d\n", programming.GetTotalCount())
	fmt.Fprintf(w, "Total: %d\n", total)
	fmt.Fprintf(w, "Version: %s\n", Version)
}

func handleCount(w http.ResponseWriter, r *http.Request) {
	total := quotes.GetTotalCount() + anime.GetTotalCount() + chucknorris.GetTotalCount() + dadjokes.GetTotalCount() + programming.GetTotalCount()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"count":   total,
	})
}

func handleCountTxt(w http.ResponseWriter, r *http.Request) {
	total := quotes.GetTotalCount() + anime.GetTotalCount() + chucknorris.GetTotalCount() + dadjokes.GetTotalCount() + programming.GetTotalCount()
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "%d", total)
}

// Web UI handlers
func handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	total := quotes.GetTotalCount() + anime.GetTotalCount() + chucknorris.GetTotalCount() + dadjokes.GetTotalCount() + programming.GetTotalCount()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Quotes API</title>
    <link rel="manifest" href="/manifest.json">
    <style>
        :root { --bg: #1a1a2e; --fg: #eee; --accent: #e94560; --card: #16213e; }
        * { box-sizing: border-box; margin: 0; padding: 0; }
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background: var(--bg); color: var(--fg); min-height: 100vh; display: flex; flex-direction: column; align-items: center; padding: 2rem; }
        h1 { color: var(--accent); margin-bottom: 1rem; }
        .stats { background: var(--card); padding: 1.5rem 2rem; border-radius: 12px; margin: 1rem 0; }
        .stats p { margin: 0.5rem 0; }
        .endpoints { background: var(--card); padding: 1.5rem 2rem; border-radius: 12px; margin: 1rem 0; max-width: 600px; width: 100%%; }
        .endpoints h2 { color: var(--accent); margin-bottom: 1rem; }
        .endpoints a { color: #6dd5ed; text-decoration: none; display: block; padding: 0.3rem 0; }
        .endpoints a:hover { text-decoration: underline; }
        footer { margin-top: auto; padding-top: 2rem; opacity: 0.7; }
    </style>
</head>
<body>
    <h1>Quotes API</h1>
    <p>Version %s</p>
    <div class="stats">
        <p><strong>Total Entries:</strong> %d</p>
        <p>Quotes: %d | Anime: %d | Chuck Norris: %d</p>
        <p>Dad Jokes: %d | Programming: %d</p>
    </div>
    <div class="endpoints">
        <h2>API Endpoints</h2>
        <a href="/api/v1/random">/api/v1/random</a>
        <a href="/api/v1/random.txt">/api/v1/random.txt</a>
        <a href="/api/v1/anime/random">/api/v1/anime/random</a>
        <a href="/api/v1/chucknorris/random">/api/v1/chucknorris/random</a>
        <a href="/api/v1/dadjokes/random">/api/v1/dadjokes/random</a>
        <a href="/api/v1/programming/random">/api/v1/programming/random</a>
        <a href="/api/v1/stats">/api/v1/stats</a>
        <a href="/docs">/docs</a>
    </div>
    <footer>Quotes API - apimgr</footer>
</body>
</html>`, Version, total, quotes.GetTotalCount(), anime.GetTotalCount(), chucknorris.GetTotalCount(), dadjokes.GetTotalCount(), programming.GetTotalCount())
}

func handleDocs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Quotes API Documentation</title>
    <style>
        :root { --bg: #1a1a2e; --fg: #eee; --accent: #e94560; --card: #16213e; }
        * { box-sizing: border-box; margin: 0; padding: 0; }
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background: var(--bg); color: var(--fg); padding: 2rem; line-height: 1.6; }
        h1, h2 { color: var(--accent); margin: 1rem 0; }
        pre { background: var(--card); padding: 1rem; border-radius: 8px; overflow-x: auto; margin: 1rem 0; }
        code { color: #6dd5ed; }
        .endpoint { background: var(--card); padding: 1rem; border-radius: 8px; margin: 1rem 0; }
        .method { color: #4ade80; font-weight: bold; }
    </style>
</head>
<body>
    <h1>Quotes API Documentation</h1>
    <h2>Collections</h2>
    <p>This API provides access to 5 collections: quotes, anime, chucknorris, dadjokes, and programming.</p>

    <h2>Endpoints</h2>
    <div class="endpoint">
        <p><span class="method">GET</span> <code>/api/v1/random</code> - Random inspirational quote</p>
        <p><span class="method">GET</span> <code>/api/v1/random.txt</code> - Plain text format</p>
    </div>
    <div class="endpoint">
        <p><span class="method">GET</span> <code>/api/v1/anime/random</code> - Random anime quote</p>
        <p><span class="method">GET</span> <code>/api/v1/anime/random.txt</code> - Plain text format</p>
    </div>
    <div class="endpoint">
        <p><span class="method">GET</span> <code>/api/v1/chucknorris/random</code> - Random Chuck Norris joke</p>
        <p><span class="method">GET</span> <code>/api/v1/chucknorris/random.txt</code> - Plain text format</p>
    </div>
    <div class="endpoint">
        <p><span class="method">GET</span> <code>/api/v1/dadjokes/random</code> - Random dad joke</p>
        <p><span class="method">GET</span> <code>/api/v1/dadjokes/random.txt</code> - Plain text format</p>
    </div>
    <div class="endpoint">
        <p><span class="method">GET</span> <code>/api/v1/programming/random</code> - Random programming joke</p>
        <p><span class="method">GET</span> <code>/api/v1/programming/random.txt</code> - Plain text format</p>
    </div>
    <div class="endpoint">
        <p><span class="method">GET</span> <code>/api/v1/stats</code> - Collection statistics</p>
        <p><span class="method">GET</span> <code>/healthz</code> - Health check</p>
    </div>

    <h2>Response Format</h2>
    <pre><code>{
  "success": true,
  "data": {
    "id": 1,
    "quote": "The only way to do great work is to love what you do.",
    "author": "Steve Jobs",
    "category": "inspiration"
  }
}</code></pre>
</body>
</html>`)
}

// PWA handlers
func handleManifest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/manifest+json")
	fmt.Fprint(w, `{
  "name": "Quotes API",
  "short_name": "Quotes",
  "description": "Random quotes and jokes API",
  "start_url": "/",
  "display": "standalone",
  "background_color": "#1a1a2e",
  "theme_color": "#e94560",
  "icons": [
    {
      "src": "data:image/svg+xml,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 100 100'><text y='.9em' font-size='90'>💬</text></svg>",
      "sizes": "any",
      "type": "image/svg+xml"
    }
  ]
}`)
}

func handleServiceWorker(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript")
	fmt.Fprint(w, `self.addEventListener('install', e => e.waitUntil(caches.open('quotes-v1').then(c => c.addAll(['/']))));
self.addEventListener('fetch', e => e.respondWith(caches.match(e.request).then(r => r || fetch(e.request))));`)
}

// robots.txt handler
func handleRobotsTxt(w http.ResponseWriter, r *http.Request, cfg *config.Config) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintln(w, "User-agent: *")
	for _, path := range cfg.WebRobots.Allow {
		fmt.Fprintf(w, "Allow: %s\n", path)
	}
	for _, path := range cfg.WebRobots.Deny {
		fmt.Fprintf(w, "Disallow: %s\n", path)
	}
}

// security.txt handler
func handleSecurityTxt(w http.ResponseWriter, r *http.Request, cfg *config.Config) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	admin := cfg.WebSecurity.Admin
	if admin == "" {
		admin = "admin@example.com"
	}
	fmt.Fprintf(w, "Contact: mailto:%s\n", admin)
	fmt.Fprintln(w, "Preferred-Languages: en")
	fmt.Fprintf(w, "Canonical: https://%s/.well-known/security.txt\n", cfg.Server.FQDN)
}

// Service command handler
func handleServiceCommand(cmd string, dirs paths.Directories) {
	switch cmd {
	case "install":
		installService(dirs)
	case "uninstall":
		uninstallService()
	case "start":
		startService()
	case "stop":
		stopService()
	case "restart":
		restartService()
	case "status":
		serviceStatus()
	default:
		fmt.Printf("Unknown service command: %s\n", cmd)
		fmt.Println("Available commands: install, uninstall, start, stop, restart, status")
		os.Exit(1)
	}
}

func installService(dirs paths.Directories) {
	if runtime.GOOS != "linux" {
		fmt.Println("Service installation is only supported on Linux")
		os.Exit(1)
	}

	execPath, _ := os.Executable()
	serviceContent := fmt.Sprintf(`[Unit]
Description=Quotes API Server
After=network.target

[Service]
Type=simple
ExecStart=%s
Restart=always
RestartSec=5
WorkingDirectory=%s

[Install]
WantedBy=multi-user.target
`, execPath, dirs.Config)

	servicePath := "/etc/systemd/system/quotes.service"
	if err := os.WriteFile(servicePath, []byte(serviceContent), 0644); err != nil {
		fmt.Printf("Failed to create service file: %v\n", err)
		os.Exit(1)
	}

	exec.Command("systemctl", "daemon-reload").Run()
	fmt.Println("Service installed successfully")
	fmt.Println("Run 'quotes --service start' to start the service")
}

func uninstallService() {
	if runtime.GOOS != "linux" {
		fmt.Println("Service uninstallation is only supported on Linux")
		os.Exit(1)
	}

	exec.Command("systemctl", "stop", "quotes").Run()
	exec.Command("systemctl", "disable", "quotes").Run()
	os.Remove("/etc/systemd/system/quotes.service")
	exec.Command("systemctl", "daemon-reload").Run()
	fmt.Println("Service uninstalled successfully")
}

func startService() {
	if runtime.GOOS != "linux" {
		fmt.Println("Service management is only supported on Linux")
		os.Exit(1)
	}
	exec.Command("systemctl", "start", "quotes").Run()
	fmt.Println("Service started")
}

func stopService() {
	if runtime.GOOS != "linux" {
		fmt.Println("Service management is only supported on Linux")
		os.Exit(1)
	}
	exec.Command("systemctl", "stop", "quotes").Run()
	fmt.Println("Service stopped")
}

func restartService() {
	if runtime.GOOS != "linux" {
		fmt.Println("Service management is only supported on Linux")
		os.Exit(1)
	}
	exec.Command("systemctl", "restart", "quotes").Run()
	fmt.Println("Service restarted")
}

func serviceStatus() {
	if runtime.GOOS != "linux" {
		fmt.Println("Service management is only supported on Linux")
		os.Exit(1)
	}
	cmd := exec.Command("systemctl", "status", "quotes")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

// Maintenance command handler
func handleMaintenanceCommand(cmd string, dirs paths.Directories) {
	maintenanceFile := filepath.Join(dirs.Data, "maintenance")
	switch cmd {
	case "on":
		if err := os.WriteFile(maintenanceFile, []byte("maintenance"), 0644); err != nil {
			fmt.Printf("Failed to enable maintenance mode: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Maintenance mode enabled")
	case "off":
		os.Remove(maintenanceFile)
		fmt.Println("Maintenance mode disabled")
	default:
		fmt.Printf("Unknown maintenance command: %s\n", cmd)
		fmt.Println("Available commands: on, off")
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Printf(`Quotes API v%s

Usage: quotes [options]

Options:
  --port PORT        Server port (default: 8080)
  --address ADDR     Server address (default: 0.0.0.0)
  --config PATH      Path to config file
  --data PATH        Path to data directory
  --logs PATH        Path to logs directory
  --mode MODE        Application mode (dev, prod)
  --update BRANCH    Update from branch (stable, beta, nightly)
  --version          Show version information
  --status           Show status (for health checks)
  --help             Show this help message

Service Management (Linux):
  --service install    Install as systemd service
  --service uninstall  Uninstall systemd service
  --service start      Start the service
  --service stop       Stop the service
  --service restart    Restart the service
  --service status     Show service status

Maintenance:
  --maintenance on     Enable maintenance mode
  --maintenance off    Disable maintenance mode

Examples:
  quotes --port 3000
  quotes --mode dev --port 8080
  quotes --config /etc/apimgr/quotes/server.yml
  quotes --service install
  quotes --update stable
`, Version)
}

func handleUpdateCommand(branch string) {
	validBranches := map[string]bool{
		"stable":  true,
		"beta":    true,
		"nightly": true,
	}

	if !validBranches[branch] {
		fmt.Printf("Error: invalid update branch %q (valid: stable, beta, nightly)\n", branch)
		os.Exit(1)
	}

	fmt.Printf("Updating Quotes API from %s branch...\n", branch)

	// Check if git is available
	if _, err := exec.LookPath("git"); err != nil {
		fmt.Println("Error: git is not installed")
		os.Exit(1)
	}

	// Perform update
	cmd := exec.Command("git", "pull", "origin", branch)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Update failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Update complete. Please rebuild the application.")
}
