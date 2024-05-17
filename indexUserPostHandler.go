package main

import (
	"net/http"

	"github.com/basedantoni/go-rss/internal/database"
)

func (cfg *apiConfig) indexUserPostHandler(w http.ResponseWriter, r *http.Request, user database.User) {
	feeds, err := cfg.DB.GetUserFeedFollows(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error fetching user's feeds")
	}

	posts := []database.Post{}
	for _, f := range feeds {
		post, err := cfg.DB.GetPost(r.Context(), f.FeedID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Error fetching user's posts")
		}
		posts = append(posts, post...)
	}

	respondWithJSON(w, http.StatusOK, posts)
}
