package main

import (
	"crypto/rand"
	_ "embed"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/apimgr/quotes/src/anime"
	"github.com/apimgr/quotes/src/chucknorris"
	"github.com/apimgr/quotes/src/dadjokes"
	"github.com/apimgr/quotes/src/database"
	"github.com/apimgr/quotes/src/paths"
	"github.com/apimgr/quotes/src/programming"
	"github.com/apimgr/quotes/src/quotes"
	"github.com/apimgr/quotes/src/server"
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
	port := flag.String("port", getEnv("PORT", "8080"), "Server port")
	address := flag.String("address", getEnv("ADDRESS", "0.0.0.0"), "Server address")
	showVersion := flag.Bool("version", false, "Show version information")
	showStatus := flag.Bool("status", false, "Show status (for health checks)")
	flag.Parse()

	// Show version
	if *showVersion {
		fmt.Printf("Quotes API v%s\n", Version)
		fmt.Printf("Commit: %s\n", Commit)
		fmt.Printf("Build Date: %s\n", BuildDate)
		os.Exit(0)
	}

	// Status check (for health checks)
	if *showStatus {
		fmt.Println("OK")
		os.Exit(0)
	}

	log.Printf("Starting Quotes API v%s", Version)

	// Get directories
	configDir := paths.GetConfigDir()
	dataDir := paths.GetDataDir()
	logsDir := paths.GetLogsDir()
	dbPath := paths.GetDBPath()

	log.Printf("Config directory: %s", configDir)
	log.Printf("Data directory: %s", dataDir)
	log.Printf("Logs directory: %s", logsDir)
	log.Printf("Database path: %s", dbPath)

	// Ensure directories exist
	for _, dir := range []string{configDir, dataDir, logsDir} {
		if err := paths.EnsureDir(dir); err != nil {
			log.Fatalf("Failed to create directory %s: %v", dir, err)
		}
	}

	// Initialize database
	if err := database.InitDB(dbPath); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// Check if admin exists, if not create one
	adminExists, err := database.AdminExists()
	if err != nil {
		log.Fatalf("Failed to check admin existence: %v", err)
	}

	if !adminExists {
		adminUser := getEnv("ADMIN_USER", "administrator")
		adminPassword := getEnv("ADMIN_PASSWORD", generateRandomPassword())
		adminToken := getEnv("ADMIN_TOKEN", generateRandomToken())

		if err := database.CreateAdmin(adminUser, adminPassword, adminToken); err != nil {
			log.Fatalf("Failed to create admin: %v", err)
		}

		log.Printf("✅ Admin user created: %s", adminUser)

		// Save credentials to file
		creds := &database.AdminCredentials{
			Username: adminUser,
			Token:    adminToken,
		}

		if err := database.SaveCredentialsToFile(creds, configDir, *port); err != nil {
			log.Printf("⚠️  Warning: Failed to save credentials file: %v", err)
		} else {
			credFile := fmt.Sprintf("%s/admin-credentials.txt", configDir)
			log.Printf("⚠️  Admin credentials saved to: %s", credFile)
		}
	}

	// Load quotes from embedded data
	log.Println("Loading quotes...")
	if err := quotes.LoadQuotes(quotesData); err != nil {
		log.Fatalf("Failed to load quotes: %v", err)
	}
	log.Printf("✅ Loaded %d quotes", quotes.GetTotalCount())

	// Load anime quotes from embedded data
	log.Println("Loading anime quotes...")
	if err := anime.LoadQuotes(animeData); err != nil {
		log.Fatalf("Failed to load anime quotes: %v", err)
	}
	log.Printf("✅ Loaded %d anime quotes", anime.GetTotalCount())

	// Load Chuck Norris jokes from embedded data
	log.Println("Loading Chuck Norris jokes...")
	if err := chucknorris.LoadJokes(chuckNorrisData); err != nil {
		log.Fatalf("Failed to load Chuck Norris jokes: %v", err)
	}
	log.Printf("✅ Loaded %d Chuck Norris jokes", chucknorris.GetTotalCount())

	// Load dad jokes from embedded data
	log.Println("Loading dad jokes...")
	if err := dadjokes.LoadJokes(dadJokesData); err != nil {
		log.Fatalf("Failed to load dad jokes: %v", err)
	}
	log.Printf("✅ Loaded %d dad jokes", dadjokes.GetTotalCount())

	// Load programming jokes from embedded data
	log.Println("Loading programming jokes...")
	if err := programming.LoadJokes(programmingData); err != nil {
		log.Fatalf("Failed to load programming jokes: %v", err)
	}
	log.Printf("✅ Loaded %d programming jokes", programming.GetTotalCount())

	// Set version information in server
	server.Version = Version
	server.Commit = Commit
	server.BuildDate = BuildDate

	// Start server
	srv := server.NewServer(*port, *address)
	if err := srv.Start(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// generateRandomPassword generates a random password
func generateRandomPassword() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "changeme123"
	}
	return hex.EncodeToString(bytes)
}

// generateRandomToken generates a random token
func generateRandomToken() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "default-token-please-change"
	}
	return hex.EncodeToString(bytes)
}
