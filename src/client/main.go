package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	Version = "dev"
	CommitID = "unknown"
	BuildDate  = "N/A"
	BuildEpoch = "0"
	OfficialSite = ""
	serverURL string
	outputFormat string
)

type quote struct {
	ID int `json:"id"`
	Quote string `json:"quote"`
	Author string `json:"author"`
	Category string `json:"category"`
}

type apiResponse struct {
	OK bool `json:"ok"`
	Data json.RawMessage `json:"data"`
	Error string `json:"error"`
	Message string `json:"message"`
}

func main() {
	root := &cobra.Command{
		Use: filepath.Base(os.Args[0]),
		Short: "Retrieve quotes from the Quotes API",
		SilenceUsage: true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			configURL, err := loadServerURL()
			if err != nil {
				return err
			}
			serverChanged := cmd.Flags().Changed("server") || cmd.InheritedFlags().Changed("server")
			environmentURL := os.Getenv("QUOTES_SERVER_PRIMARY")
			if environmentURL == "" {
				environmentURL = os.Getenv("QUOTES_SERVER")
			}
			if serverChanged {
				if configURL == "" {
					if err := saveServerURL(serverURL); err != nil {
						return err
					}
				}
			} else if environmentURL != "" {
				serverURL = environmentURL
			} else if serverURL == "" {
				serverURL = configURL
				if serverURL == "" {
					serverURL = OfficialSite
				}
			}
			if serverURL == "" {
				return errors.New("no API server configured; set --server or QUOTES_SERVER")
			}
			parsed, err := url.Parse(serverURL)
			if err != nil || parsed.Host == "" || (parsed.Scheme != "https" && parsed.Scheme != "http") {
				return errors.New("--server must be an absolute HTTP or HTTPS URL")
			}
			return nil
		},
	}
	root.Version = fmt.Sprintf("%s (%s)", Version, CommitID)
	root.PersistentFlags().StringVar(&serverURL, "server", "", "Quotes API server URL")
	root.PersistentFlags().StringVar(&outputFormat, "output", "plain", "Output format: plain or json")
	root.AddCommand(randomCommand(), categoriesCommand(), quoteCommand())
	if err := root.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func randomCommand() *cobra.Command {
	var category string
	command := &cobra.Command{
		Use: "random",
		Short: "Get a random quote",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			query := url.Values{}
			if category != "" {
				query.Set("category", category)
			}
			var result quote
			if err := getJSON("/api/v1/quotes/random?"+query.Encode(), &result); err != nil {
				return err
			}
			return printQuote(cmd.OutOrStdout(), result, outputFormat)
		},
	}
	command.Flags().StringVar(&category, "category", "", "Filter by category")
	return command
}

func categoriesCommand() *cobra.Command {
	return &cobra.Command{
		Use: "categories",
		Short: "List quote categories and counts",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			var counts map[string]int
			if err := getJSON("/api/v1/categories", &counts); err != nil {
				return err
			}
			if outputFormat == "json" {
				encoder := json.NewEncoder(cmd.OutOrStdout())
				encoder.SetIndent("", "  ")
				return encoder.Encode(counts)
			}
			if outputFormat != "plain" {
				return fmt.Errorf("unsupported output format %q", outputFormat)
			}
			for _, name := range []string{"quotes", "anime", "chucknorris", "dadjokes", "programming"} {
				fmt.Fprintf(cmd.OutOrStdout(), "%s\t%d\n", name, counts[name])
			}
			return nil
		},
	}
}

func quoteCommand() *cobra.Command {
	return &cobra.Command{
		Use: "quote ID",
		Short: "Get a quote by its ID",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := url.PathUnescape(args[0])
			if err != nil || strings.Contains(id, "/") {
				return errors.New("quote ID must be a single URL path segment")
			}
			var result quote
			if err := getJSON("/api/v1/quotes/"+url.PathEscape(id), &result); err != nil {
				return err
			}
			return printQuote(cmd.OutOrStdout(), result, outputFormat)
		},
	}
}

func getJSON(path string, output interface{}) error {
	requestURL := strings.TrimRight(serverURL, "/") + path
	request, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	request.Header.Set("Accept", "application/json")
	request.Header.Set("User-Agent", "quotes-cli/"+Version)
	client := &http.Client{Timeout: 15 * time.Second}
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("connect to Quotes API: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		var apiErr apiResponse
		_ = json.NewDecoder(io.LimitReader(response.Body, 1<<20)).Decode(&apiErr)
		if apiErr.Message != "" {
			return fmt.Errorf("API returned %s: %s", response.Status, apiErr.Message)
		}
		return fmt.Errorf("API returned %s", response.Status)
	}
	var envelope apiResponse
	if err := json.NewDecoder(io.LimitReader(response.Body, 1<<20)).Decode(&envelope); err != nil {
		return fmt.Errorf("decode API response: %w", err)
	}
	if !envelope.OK {
		return fmt.Errorf("API response was unsuccessful: %s", envelope.Message)
	}
	if err := json.Unmarshal(envelope.Data, output); err != nil {
		return fmt.Errorf("decode quote data: %w", err)
	}
	return nil
}

func printQuote(writer io.Writer, item quote, format string) error {
	switch format {
	case "plain":
		if item.Author == "" {
			_, err := fmt.Fprintf(writer, "%s\n", item.Quote)
			return err
		}
		_, err := fmt.Fprintf(writer, "\"%s\" - %s\n", item.Quote, item.Author)
		return err
	case "json":
		encoder := json.NewEncoder(writer)
		encoder.SetIndent("", "  ")
		return encoder.Encode(item)
	default:
		return fmt.Errorf("unsupported output format %q", format)
	}
}

func configFilePath() (string, error) {
	base := os.Getenv("XDG_CONFIG_HOME")
	if base == "" {
		userConfig, err := os.UserConfigDir()
		if err != nil {
			return "", err
		}
		base = userConfig
	}
	return filepath.Join(base, "apimgr", "quotes", "cli.yml"), nil
}


type cliConfig struct {
	Server struct {
		Primary string `yaml:"primary"`
	} `yaml:"server"`
}

func loadServerURL() (string, error) {
	path, err := configFilePath()
	if err != nil {
		return "", fmt.Errorf("locate CLI config: %w", err)
	}
	content, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		if err := saveServerURL(""); err != nil {
			return "", err
		}
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("read CLI config: %w", err)
	}
	var config cliConfig
	if err := yaml.Unmarshal(content, &config); err != nil {
		return "", fmt.Errorf("parse CLI config: %w", err)
	}
	return strings.TrimSpace(config.Server.Primary), nil
}

func saveServerURL(server string) error {
	path, err := configFilePath()
	if err != nil {
		return fmt.Errorf("locate CLI config: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return fmt.Errorf("create CLI config directory: %w", err)
	}
	config := cliConfig{}
	config.Server.Primary = server
	content, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("encode CLI config: %w", err)
	}
	if err := os.WriteFile(path, content, 0600); err != nil {
		return fmt.Errorf("write CLI config: %w", err)
	}
	return nil
}
