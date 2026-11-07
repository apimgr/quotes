package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/apimgr/quotes/src/catalog"
	"github.com/apimgr/quotes/src/config"
	"github.com/apimgr/quotes/src/httpapp"
	"github.com/apimgr/quotes/src/paths"
)

var (
	Version      = "dev"
	CommitID     = "unknown"
	BuildDate    = "N/A"
	BuildEpoch   = "0"
	OfficialSite = ""
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		log.Printf("Error: %v", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	flags := flag.NewFlagSet(filepath.Base(os.Args[0]), flag.ContinueOnError)
	flags.SetOutput(os.Stderr)
	port := flags.String("port", "", "HTTP port (overrides server.yml)")
	address := flags.String("address", "", "Bind address (overrides server.yml)")
	configPath := flags.String("config", "", "Path to server.yml")
	showVersion := flags.Bool("version", false, "Show version and build details")
	flags.BoolVar(showVersion, "v", false, "Show version and build details")
	debug := flags.Bool("debug", false, "Enable debug logging")
	color := flags.String("color", "auto", "Color output: auto, yes, no")
	showStatus := flags.Bool("status", false, "Return success when the service binary is available")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("unexpected arguments: %s", strings.Join(flags.Args(), " "))
	}
	switch strings.ToLower(*color) {
	case "auto":
	case "yes":
		os.Unsetenv("NO_COLOR")
	case "no":
		os.Setenv("NO_COLOR", "1")
	default:
		return fmt.Errorf("invalid --color value %q; choose auto, yes, or no", *color)
	}
	if *debug {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		os.Setenv("DEBUG", "1")
	}
	if *showVersion {
		fmt.Printf("quotes %s (%s) built %s\nGo: %s\nOS/Arch: %s/%s\nOfficial site: %s\n", Version, CommitID, BuildDate, runtime.Version(), runtime.GOOS, runtime.GOARCH, OfficialSite)
		return nil
	}
	if *showStatus {
		fmt.Println("OK")
		return nil
	}

	dirs := paths.GetDirectories()
	if err := paths.EnsureDirectories(dirs); err != nil {
		return fmt.Errorf("create application directories: %w", err)
	}
	cfgPath := *configPath
	if cfgPath == "" {
		cfgPath = filepath.Join(dirs.Config, "server.yml")
	}
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("load configuration: %w", err)
	}
	listenAddress := cfg.Server.Address
	if listenAddress == "" {
		listenAddress = "0.0.0.0"
	}
	listenPort := cfg.Server.Port
	if listenPort == "" {
		listenPort = "8080"
	}
	if *address != "" {
		listenAddress = *address
	}
	if *port != "" {
		listenPort = *port
	}
	if strings.Contains(listenPort, ":") {
		return fmt.Errorf("port must contain only a port number")
	}
	if err := validateListenAddress(listenAddress, listenPort); err != nil {
		return err
	}
	store, err := catalog.LoadEmbedded()
	if err != nil {
		return fmt.Errorf("load quote catalog: %w", err)
	}

	server := &http.Server{
		Addr:              net.JoinHostPort(listenAddress, listenPort),
		Handler:           httpapp.NewHandler(store, Version, cfg.Server.RateLimit),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       90 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}
	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("Quotes API v%s listening on %s", Version, server.Addr)
		serverErrors <- server.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(shutdown)
	select {
	case err := <-serverErrors:
		if err != nil && err != http.ErrServerClosed {
			return fmt.Errorf("serve HTTP: %w", err)
		}
		return nil
	case <-shutdown:
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("shut down HTTP server: %w", err)
	}
	if err := <-serverErrors; err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("serve HTTP: %w", err)
	}
	return nil
}

func validateListenAddress(address, port string) error {
	if address == "" {
		return fmt.Errorf("bind address cannot be empty")
	}
	portNumber, err := strconv.Atoi(port)
	if err != nil || portNumber < 1 || portNumber > 65535 {
		return fmt.Errorf("port must be an integer from 1 to 65535")
	}
	return nil
}
