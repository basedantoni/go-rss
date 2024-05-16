package main

import (
	"net/http"

	"github.com/basedantoni/go-rss/internal/database"
)

func (cfg *apiConfig) showUserHandler(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWithJSON(w, http.StatusOK, user)
}
