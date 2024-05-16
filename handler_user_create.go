package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/basedantoni/go-rss/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	type bodyParams struct {
		Name string `json:"name"`
	}

	params := bodyParams{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error while decoding body params")
		return
	}

	user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      params.Name,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Could not create new user")
	}

	respondWithJSON(w, http.StatusOK, user)
}
