package quotes

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

// Quote represents a single inspirational quote
type Quote struct {
	ID       int    `json:"id"`
	Quote    string `json:"quote"`
	Author   string `json:"author"`
	Category string `json:"category"`
}

var (
	quotes []Quote
	rng    *rand.Rand
)

func init() {
	// Initialize random number generator with seed
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// LoadQuotes loads quotes from embedded JSON data
func LoadQuotes(jsonData []byte) error {
	err := json.Unmarshal(jsonData, &quotes)
	if err != nil {
		return fmt.Errorf("failed to parse quotes.json: %w", err)
	}

	if len(quotes) == 0 {
		return fmt.Errorf("no quotes found in quotes.json")
	}

	return nil
}

// GetRandomQuote returns a random quote from the loaded quotes
func GetRandomQuote() (*Quote, error) {
	if len(quotes) == 0 {
		return nil, fmt.Errorf("no quotes available, please load quotes first")
	}

	index := rng.Intn(len(quotes))
	return &quotes[index], nil
}

// GetAllQuotes returns all loaded quotes
func GetAllQuotes() []Quote {
	return quotes
}

// GetQuoteByID returns a quote by its ID
func GetQuoteByID(id int) (*Quote, error) {
	for _, quote := range quotes {
		if quote.ID == id {
			return &quote, nil
		}
	}
	return nil, fmt.Errorf("quote with ID %d not found", id)
}

// GetQuotesByCategory returns all quotes in a specific category
func GetQuotesByCategory(category string) []Quote {
	var result []Quote
	for _, quote := range quotes {
		if quote.Category == category {
			result = append(result, quote)
		}
	}
	return result
}

// GetQuotesByAuthor returns all quotes by a specific author
func GetQuotesByAuthor(author string) []Quote {
	var result []Quote
	for _, quote := range quotes {
		if quote.Author == author {
			result = append(result, quote)
		}
	}
	return result
}

// GetTotalCount returns the total number of loaded quotes
func GetTotalCount() int {
	return len(quotes)
}

// GetAllJSON returns the raw JSON data
func GetAllJSON() ([]byte, error) {
	return json.Marshal(quotes)
}
