package server

import (
	"encoding/json"
	"net/http"

	"github.com/apimgr/quotes/src/database"
	"github.com/gorilla/mux"
)

// handleGetSettings returns all settings
func handleGetSettings(w http.ResponseWriter, r *http.Request) {
	settings, err := database.GetAllSettings()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve settings")
		return
	}

	respondWithJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    settings,
	})
}

// handleSetSetting sets or updates a setting
func handleSetSetting(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Key == "" {
		respondWithError(w, http.StatusBadRequest, "Key is required")
		return
	}

	if err := database.SetSetting(req.Key, req.Value); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to set setting")
		return
	}

	respondWithJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    map[string]string{"message": "Setting updated successfully"},
	})
}

// handleDeleteSetting deletes a setting
func handleDeleteSetting(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	if err := database.DeleteSetting(key); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete setting")
		return
	}

	respondWithJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    map[string]string{"message": "Setting deleted successfully"},
	})
}
