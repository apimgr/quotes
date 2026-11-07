package catalog

import (
	"crypto/rand"
	"math/big"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"sort"
	"strings"

	"github.com/apimgr/quotes/src/data"
)

const (
	CategoryQuotes      = "quotes"
	CategoryAnime       = "anime"
	CategoryChuckNorris = "chucknorris"
	CategoryDadJokes    = "dadjokes"
	CategoryProgramming = "programming"
)

var categoryNames = [...]string{CategoryQuotes, CategoryAnime, CategoryChuckNorris, CategoryDadJokes, CategoryProgramming}

// Names returns the supported category names in stable order.
func Names() []string {
	return append([]string(nil), categoryNames[:]...)
}

// Quote is the stable, public representation shared by every collection.
type Quote struct {
	ID       int    `json:"id"`
	Quote    string `json:"quote"`
	Author   string `json:"author"`
	Category string `json:"category"`
}

// Catalog contains immutable quote records after construction.
type Catalog struct {
	byCategory map[string][]Quote
	byID       map[int]Quote
	all        []Quote
}

type sourceRecord struct {
	ID        int    `json:"id"`
	Quote     string `json:"quote"`
	Joke      string `json:"joke"`
	Author    string `json:"author"`
	Character string `json:"character"`
}

// LoadEmbedded reads all quote collections embedded in the application binary.
func LoadEmbedded() (*Catalog, error) {
	collections := make(map[string][]byte, len(categoryNames))
	for _, category := range categoryNames {
		content, err := fs.ReadFile(data.EmbeddedData, categoryFile(category))
		if err != nil {
			return nil, fmt.Errorf("read embedded %s collection: %w", category, err)
		}
		collections[category] = content
	}
	return New(collections)
}

// New parses and validates every category before publishing any records.
func New(collections map[string][]byte) (*Catalog, error) {
	c := &Catalog{
		byCategory: make(map[string][]Quote, len(categoryNames)),
		byID:       make(map[int]Quote),
	}
	for _, category := range categoryNames {
		content, ok := collections[category]
		if !ok {
			return nil, fmt.Errorf("missing %s collection", category)
		}
		var records []sourceRecord
		if err := json.Unmarshal(content, &records); err != nil {
			return nil, fmt.Errorf("parse %s collection: %w", category, err)
		}
		if len(records) == 0 {
			return nil, fmt.Errorf("%s collection is empty", category)
		}
		for _, record := range records {
			text := record.Quote
			if text == "" {
				text = record.Joke
			}
			author := record.Author
			if author == "" {
				author = record.Character
			}
			if record.ID < 1 || strings.TrimSpace(text) == "" {
				return nil, fmt.Errorf("%s collection contains an invalid record", category)
			}
			quote := Quote{ID: len(c.all) + 1, Quote: text, Author: author, Category: category}
			c.byID[quote.ID] = quote
			c.byCategory[category] = append(c.byCategory[category], quote)
			c.all = append(c.all, quote)
		}
		sort.Slice(c.byCategory[category], func(i, j int) bool {
			return c.byCategory[category][i].ID < c.byCategory[category][j].ID
		})
	}
	return c, nil
}

// Categories returns category names and record counts in stable order.
func (c *Catalog) Categories() map[string]int {
	counts := make(map[string]int, len(c.byCategory))
	for _, category := range categoryNames {
		counts[category] = len(c.byCategory[category])
	}
	return counts
}

// List returns a filtered, paginated copy of the catalog.
func (c *Catalog) List(category, author string, page, limit int) ([]Quote, int, error) {
	if category != "" {
		if _, ok := c.byCategory[category]; !ok {
			return nil, 0, fmt.Errorf("unknown category %q", category)
		}
	}
	matches := make([]Quote, 0)
	for _, quote := range c.all {
		if category != "" && quote.Category != category {
			continue
		}
		if author != "" && !strings.EqualFold(quote.Author, author) {
			continue
		}
		matches = append(matches, quote)
	}
	total := len(matches)
	if page < 1 || limit < 1 {
		return nil, 0, errors.New("page and limit must be positive")
	}
	if page-1 > total/limit {
		return []Quote{}, total, nil
	}
	start := (page - 1) * limit
	if start >= total {
		return []Quote{}, total, nil
	}
	end := start + limit
	if end > total {
		end = total
	}
	return append([]Quote(nil), matches[start:end]...), total, nil
}

// Get returns the record identified by category and ID.
func (c *Catalog) Get(category string, id int) (Quote, error) {
	quote, ok := c.byID[id]
	if !ok || (category != "" && quote.Category != category) {
		return Quote{}, errors.New("quote not found")
	}
	return quote, nil
}

// Random returns an unbiased cryptographically selected quote from a category or all records.
func (c *Catalog) Random(category string) (Quote, error) {
	quotes := c.all
	if category != "" {
		var ok bool
		quotes, ok = c.byCategory[category]
		if !ok {
			return Quote{}, fmt.Errorf("unknown category %q", category)
		}
	}
	if len(quotes) == 0 {
		return Quote{}, errors.New("quote catalog is empty")
	}
	index, err := rand.Int(rand.Reader, big.NewInt(int64(len(quotes))))
	if err != nil {
		return Quote{}, fmt.Errorf("select random quote: %w", err)
	}
	return quotes[index.Int64()], nil
}

func categoryFile(category string) string {
	if category == CategoryQuotes {
		return "quotes.json"
	}
	return category + ".json"
}
