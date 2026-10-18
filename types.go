package main

import "sync/atomic"

type apiConfig struct {
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
