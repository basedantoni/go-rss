package main

import "net/http"

func errorHandler(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, http.StatusInternalServerError, response{Error: "Internal Server Error"})
}
