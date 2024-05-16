package main

import (
	"net/http"

	"github.com/basedantoni/go-rss/internal/database"
)

type apiConfig struct {
	DB *database.Queries
}

type authedHandler func(http.ResponseWriter, *http.Request, database.User)
