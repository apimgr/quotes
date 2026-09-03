package chucknorris

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

// Joke represents a Chuck Norris joke
type Joke struct {
	ID       int    `json:"id"`
	Joke     string `json:"joke"`
	Category string `json:"category"`
}

var (
	jokes []Joke
	rng   *rand.Rand
)

func init() {
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// LoadJokes loads Chuck Norris jokes from embedded JSON data
func LoadJokes(jsonData []byte) error {
	err := json.Unmarshal(jsonData, &jokes)
	if err != nil {
		return fmt.Errorf("failed to parse chucknorris.json: %w", err)
	}

	if len(jokes) == 0 {
		return fmt.Errorf("no Chuck Norris jokes found")
	}

	return nil
}

// GetRandomJoke returns a random Chuck Norris joke
func GetRandomJoke() (*Joke, error) {
	if len(jokes) == 0 {
		return nil, fmt.Errorf("no jokes available")
	}

	index := rng.Intn(len(jokes))
	return &jokes[index], nil
}

// GetAllJokes returns all Chuck Norris jokes
func GetAllJokes() []Joke {
	return jokes
}

// GetJokeByID returns a joke by its ID
func GetJokeByID(id int) (*Joke, error) {
	for _, joke := range jokes {
		if joke.ID == id {
			return &joke, nil
		}
	}
	return nil, fmt.Errorf("joke with ID %d not found", id)
}

// GetTotalCount returns the total number of jokes
func GetTotalCount() int {
	return len(jokes)
}
