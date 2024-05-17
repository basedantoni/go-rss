package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/basedantoni/go-rss/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) createFeedHandler(w http.ResponseWriter, r *http.Request, user database.User) {
	type bodyParams struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}

	params := bodyParams{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error while decoding body params")
		return
	}

	feed, err := cfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:     uuid.New(),
		Name:   params.Name,
		Url:    params.Url,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Could not create new feed")
	}

	feedFollow, err := cfg.DB.CreateUserFeed(r.Context(), database.CreateUserFeedParams{
		ID:        uuid.New(),
		UserID:    user.ID,
		FeedID:    feed.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Could not create new feed")
	}

	type response struct {
		Feed       Feed               `json:"feed"`
		FeedFollow database.UsersFeed `json:"feed_follow"`
	}

	respondWithJSON(
		w, http.StatusOK,
		response{
			Feed:       databaseFeedToFeed(feed),
			FeedFollow: feedFollow,
		},
	)
}
