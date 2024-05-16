package main

import (
	"net/http"
)

func (cfg *apiConfig) indexFeedHandler(w http.ResponseWriter, r *http.Request) {
	feeds, err := cfg.DB.GetAllFeeds(r.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Could index feeds")
	}

	respondWithJSON(w, http.StatusOK, feeds)
}
