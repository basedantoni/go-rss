package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/basedantoni/go-rss/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Could not load environment variables")
	}

	port := os.Getenv("PORT")
	postgresUri := os.Getenv("POSTGRES_URI")

	db, err := sql.Open("postgres", postgresUri)
	if err != nil {
		log.Fatal("Could not connect to database")
	}

	dbQueries := database.New(db)

	apiCfg := &apiConfig{
		DB: dbQueries,
	}

	mux := http.NewServeMux()
	mux.Handle("/app/*", http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	mux.HandleFunc("GET /v1/readiness", healthHandler)
	mux.HandleFunc("GET /v1/err", errorHandler)
	// Users
	mux.HandleFunc("GET /v1/users", apiCfg.middlewareAuth(apiCfg.showUserHandler))
	mux.HandleFunc("POST /v1/users", apiCfg.createUserHandler)
	// Feeds
	mux.HandleFunc("GET /v1/feeds", apiCfg.indexFeedHandler)
	mux.HandleFunc("POST /v1/feeds", apiCfg.middlewareAuth(apiCfg.createFeedHandler))
	// Feed follows
	mux.HandleFunc("GET /v1/feed_follows", apiCfg.middlewareAuth(apiCfg.indexUserFeedHandler))
	mux.HandleFunc("POST /v1/feed_follows", apiCfg.middlewareAuth(apiCfg.createUserFeedHandler))
	mux.HandleFunc("DELETE /v1/feed_follows/{feedFollowID}", apiCfg.middlewareAuth(apiCfg.removeUserFeedHandler))

	server := &http.Server{Handler: mux, Addr: ":" + port}

	log.Fatal(server.ListenAndServe())
}
