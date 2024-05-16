package main

import (
	"net/http"

	"github.com/basedantoni/go-rss/internal/database"
)

func (cfg *apiConfig) indexUserFeedHandler(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollows, err := cfg.DB.GetUserFeedFollows(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error accessing user feeds from database")
		return
	}

	respondWithJSON(w, http.StatusOK, feedFollows)
}
