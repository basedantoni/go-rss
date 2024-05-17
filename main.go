package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/basedantoni/go-rss/internal/database"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func fetchRSSFeed(url string) (Rss, error) {
	resp, err := http.Get(url)
	if err != nil {
		return Rss{}, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return Rss{}, err
	}

	var rss Rss
	err = xml.Unmarshal(data, &rss)
	if err != nil {
		return Rss{}, err
	}

	return rss, nil
}

func (cfg *apiConfig) feedFetcher() {
	ticker := time.NewTicker(60 * time.Minute)
	defer ticker.Stop()

	for ; true; <-ticker.C {
		handleTick(cfg)
	}
}

func handleTick(cfg *apiConfig) {
	feeds, err := cfg.DB.GetNextFeedsToFetch(context.Background(), 10)
	if err != nil {
		log.Printf("Error fetching feeds: %v", err)
		return
	}

	var wg sync.WaitGroup
	wg.Add(len(feeds))

	for _, feed := range feeds {
		go func(feed database.Feed) {
			defer wg.Done()
			rss, err := fetchRSSFeed(feed.Url)
			if err != nil {
				log.Printf("Error fetching RSS feed from %s: %v", feed.Url, err)
				return
			}

			for _, item := range rss.Channel.Items {
				parseTime, err := time.Parse(time.RFC1123, item.PubDate)
				if err != nil {
					log.Printf("Error parsing time %s: %v", item.PubDate, err)
					return
				}

				cfg.DB.CreatePost(context.Background(), database.CreatePostParams{
					ID:          uuid.New(),
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
					Title:       item.Title,
					Url:         item.Link,
					Description: sql.NullString{String: item.Description, Valid: true},
					PublishedAt: sql.NullTime{Time: parseTime, Valid: true},
					FeedID:      feed.ID,
				})
			}
		}(feed)
	}
	wg.Wait()
}

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

	go apiCfg.feedFetcher()

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
	// Posts
	mux.HandleFunc("GET /v1/posts", apiCfg.middlewareAuth(apiCfg.indexUserPostHandler))

	server := &http.Server{Handler: mux, Addr: ":" + port}

	log.Fatal(server.ListenAndServe())
}
