package main

import (
	"net/http"

	"github.com/basedantoni/go-rss/internal/auth"
)

func (cfg *apiConfig) showUserHandler(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetApiKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusForbidden, "Invalid API Key")
		return
	}

	user, err := cfg.DB.GetUser(r.Context(), apiKey)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}
