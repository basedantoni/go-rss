package main

import (
	"net/http"
)

func (cfg *apiConfig) indexFeedHandler(w http.ResponseWriter, r *http.Request) {
	databaseFeeds, err := cfg.DB.GetAllFeeds(r.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Could index feeds")
	}

	feeds := []Feed{}
	for _, f := range databaseFeeds {
		feeds = append(feeds, databaseFeedToFeed(f))
	}

	respondWithJSON(w, http.StatusOK, feeds)
}
