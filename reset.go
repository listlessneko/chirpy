package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, _ *http.Request) {
	data := []byte("You are forbidden from using this endpoint")
	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, data)
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	cfg.fileserverHits.Store(0)
	reset := fmt.Sprintf("Reset to: %d\n", cfg.fileserverHits.Load())
	w.Write([]byte(reset))
}
