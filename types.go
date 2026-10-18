package main

import (
	"sync/atomic"

	"codeberg.org/kedarui/chirpy/internal/database"
)

type apiConfig struct {
	database       *database.Queries
	fileserverHits atomic.Int32
}

type chirp struct {
	Body string `json:"body"`
}

type cleaned_chirp struct {
	CleanedBody string `json:"cleaned_body"`
}

type errResp struct {
	Error string `json:"error"`
}

type validResp struct {
	Valid bool `json:"valid"`
}
