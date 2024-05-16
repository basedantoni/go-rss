package main

import (
	"net/http"

	"github.com/basedantoni/go-rss/internal/auth"
)

func (cfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		handler(w, r, user)
	})
}
