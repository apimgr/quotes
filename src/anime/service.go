package anime

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

// AnimeQuote represents a single anime quote
type AnimeQuote struct {
	ID        int    `json:"id"`
	Quote     string `json:"quote"`
	Character string `json:"character"`
	Anime     string `json:"anime"`
	Category  string `json:"category"`
}

var (
	quotes []AnimeQuote
	rng    *rand.Rand
)

func init() {
	// Initialize random number generator with seed
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// LoadQuotes loads anime quotes from embedded JSON data
func LoadQuotes(jsonData []byte) error {
	err := json.Unmarshal(jsonData, &quotes)
	if err != nil {
		return fmt.Errorf("failed to parse anime.json: %w", err)
	}

	if len(quotes) == 0 {
		return fmt.Errorf("no anime quotes found in anime.json")
	}

	return nil
}

// GetRandomQuote returns a random anime quote from the loaded quotes
func GetRandomQuote() (*AnimeQuote, error) {
	if len(quotes) == 0 {
		return nil, fmt.Errorf("no anime quotes available, please load quotes first")
	}

	index := rng.Intn(len(quotes))
	return &quotes[index], nil
}

// GetAllQuotes returns all loaded anime quotes
func GetAllQuotes() []AnimeQuote {
	return quotes
}

// GetQuoteByID returns an anime quote by its ID
func GetQuoteByID(id int) (*AnimeQuote, error) {
	for _, quote := range quotes {
		if quote.ID == id {
			return &quote, nil
		}
	}
	return nil, fmt.Errorf("anime quote with ID %d not found", id)
}

// GetQuotesByCategory returns all anime quotes in a specific category
func GetQuotesByCategory(category string) []AnimeQuote {
	var result []AnimeQuote
	for _, quote := range quotes {
		if quote.Category == category {
			result = append(result, quote)
		}
	}
	return result
}

// GetQuotesByAnime returns all quotes from a specific anime
func GetQuotesByAnime(animeName string) []AnimeQuote {
	var result []AnimeQuote
	for _, quote := range quotes {
		if quote.Anime == animeName {
			result = append(result, quote)
		}
	}
	return result
}

// GetQuotesByCharacter returns all quotes by a specific character
func GetQuotesByCharacter(character string) []AnimeQuote {
	var result []AnimeQuote
	for _, quote := range quotes {
		if quote.Character == character {
			result = append(result, quote)
		}
	}
	return result
}

// GetTotalCount returns the total number of loaded anime quotes
func GetTotalCount() int {
	return len(quotes)
}
