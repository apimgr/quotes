package server

import (
	"encoding/json"
	"html/template"
	"net/http"
	"os"
	"strconv"

	"github.com/apimgr/quotes/src/anime"
	"github.com/apimgr/quotes/src/chucknorris"
	"github.com/apimgr/quotes/src/dadjokes"
	"github.com/apimgr/quotes/src/programming"
	"github.com/apimgr/quotes/src/quotes"
	"github.com/gorilla/mux"
)

// APIResponse represents a standard API response
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// handleRandomQuote returns a random quote
func handleRandomQuote(w http.ResponseWriter, r *http.Request) {
	quote, err := quotes.GetRandomQuote()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    quote,
	})
}

// handleAllQuotes returns all quotes
func handleAllQuotes(w http.ResponseWriter, r *http.Request) {
	allQuotes := quotes.GetAllQuotes()

	respondWithJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    allQuotes,
	})
}

// handleQuoteByID returns a quote by ID
func handleQuoteByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid quote ID")
		return
	}

	quote, err := quotes.GetQuoteByID(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    quote,
	})
}

// handleQuotesByCategory returns quotes by category
func handleQuotesByCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	category := vars["category"]

	categoryQuotes := quotes.GetQuotesByCategory(category)
	if len(categoryQuotes) == 0 {
		respondWithError(w, http.StatusNotFound, "No quotes found for this category")
		return
	}

	respondWithJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    categoryQuotes,
	})
}

// handleQuotesByAuthor returns quotes by author
func handleQuotesByAuthor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	author := vars["author"]

	authorQuotes := quotes.GetQuotesByAuthor(author)
	if len(authorQuotes) == 0 {
		respondWithError(w, http.StatusNotFound, "No quotes found for this author")
		return
	}

	respondWithJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    authorQuotes,
	})
}

// handleHome renders the home page
func handleHome(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFS(content, "templates/base.html", "templates/home.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Title":   "Quotes API",
		"Version": Version,
		"Count":   quotes.GetTotalCount(),
	}

	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

// handleAdminPage renders the admin page
func handleAdminPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFS(content, "templates/base.html", "templates/admin.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Title":   "Admin Panel - Quotes API",
		"Version": Version,
	}

	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

// handleJSONFile serves raw JSON files
func handleJSONFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filename := vars["file"]

	// Read the JSON file
	data, err := os.ReadFile("./src/data/" + filename)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "JSON file not found")
		return
	}

	// Send raw JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// respondWithJSON sends a JSON response
func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}

// respondWithError sends an error response
func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	respondWithJSON(w, statusCode, APIResponse{
		Success: false,
		Error:   message,
	})
}

// Anime quote handlers

// handleRandomAnimeQuote returns a random anime quote
func handleRandomAnimeQuote(w http.ResponseWriter, r *http.Request) {
	quote, err := anime.GetRandomQuote()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    quote,
	})
}

// handleAllAnimeQuotes returns all anime quotes
func handleAllAnimeQuotes(w http.ResponseWriter, r *http.Request) {
	allQuotes := anime.GetAllQuotes()

	respondWithJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    allQuotes,
	})
}

// handleAnimeQuoteByID returns an anime quote by ID
func handleAnimeQuoteByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid quote ID")
		return
	}

	quote, err := anime.GetQuoteByID(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    quote,
	})
}

// handleAnimeQuotesByCategory returns anime quotes by category
func handleAnimeQuotesByCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	category := vars["category"]

	categoryQuotes := anime.GetQuotesByCategory(category)
	if len(categoryQuotes) == 0 {
		respondWithError(w, http.StatusNotFound, "No anime quotes found for this category")
		return
	}

	respondWithJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    categoryQuotes,
	})
}

// handleAnimeQuotesByAnime returns quotes from a specific anime
func handleAnimeQuotesByAnime(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	animeName := vars["anime"]

	animeQuotes := anime.GetQuotesByAnime(animeName)
	if len(animeQuotes) == 0 {
		respondWithError(w, http.StatusNotFound, "No quotes found for this anime")
		return
	}

	respondWithJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    animeQuotes,
	})
}

// handleAnimeQuotesByCharacter returns quotes by a specific character
func handleAnimeQuotesByCharacter(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	character := vars["character"]

	characterQuotes := anime.GetQuotesByCharacter(character)
	if len(characterQuotes) == 0 {
		respondWithError(w, http.StatusNotFound, "No quotes found for this character")
		return
	}

	respondWithJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    characterQuotes,
	})
}

// Chuck Norris joke handlers

// handleRandomChuckNorrisJoke returns a random Chuck Norris joke
func handleRandomChuckNorrisJoke(w http.ResponseWriter, r *http.Request) {
	joke, err := chucknorris.GetRandomJoke()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    joke,
	})
}

// handleAllChuckNorrisJokes returns all Chuck Norris jokes
func handleAllChuckNorrisJokes(w http.ResponseWriter, r *http.Request) {
	allJokes := chucknorris.GetAllJokes()

	respondWithJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    allJokes,
	})
}

// Dad joke handlers

// handleRandomDadJoke returns a random dad joke
func handleRandomDadJoke(w http.ResponseWriter, r *http.Request) {
	joke, err := dadjokes.GetRandomJoke()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    joke,
	})
}

// handleAllDadJokes returns all dad jokes
func handleAllDadJokes(w http.ResponseWriter, r *http.Request) {
	allJokes := dadjokes.GetAllJokes()

	respondWithJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    allJokes,
	})
}

// Programming joke handlers

// handleRandomProgrammingJoke returns a random programming joke
func handleRandomProgrammingJoke(w http.ResponseWriter, r *http.Request) {
	joke, err := programming.GetRandomJoke()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    joke,
	})
}

// handleAllProgrammingJokes returns all programming jokes
func handleAllProgrammingJokes(w http.ResponseWriter, r *http.Request) {
	allJokes := programming.GetAllJokes()

	respondWithJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    allJokes,
	})
}
