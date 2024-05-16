package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/basedantoni/go-rss/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) createUserFeedHandler(w http.ResponseWriter, r *http.Request, user database.User) {
	type bodyParams struct {
		FeedId string `json:"feed_id"`
	}

	params := bodyParams{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error while decoding body params")
		return
	}

	feedFollow, err := cfg.DB.CreateUserFeed(r.Context(), database.CreateUserFeedParams{
		ID:        uuid.New(),
		UserID:    user.ID,
		FeedID:    uuid.MustParse(params.FeedId),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Could not create new feed")
	}

	respondWithJSON(w, http.StatusOK, feedFollow)
}
