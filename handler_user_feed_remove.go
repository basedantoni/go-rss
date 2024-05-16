package main

import (
	"net/http"

	"github.com/basedantoni/go-rss/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) removeUserFeedHandler(w http.ResponseWriter, r *http.Request, user database.User) {
	id := uuid.MustParse(r.PathValue("feedFollowID"))

	err := cfg.DB.DeleteUserFeed(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Could not create new feed")
	}

	respondWithJSON(w, http.StatusNoContent, "")
}
