package main

import (
	"fmt"
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	data := []byte("You are forbidden from using this endpoint")
	if cfg.platform != "dev" {
		log.Println("User is forbidden")
		respondWithError(w, http.StatusForbidden, data)
	}

	cfg.fileserverHits.Store(0)

	cfg.database.DeleteAllUsers(r.Context())

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	reset := fmt.Sprintf("Reset to: %d\n", cfg.fileserverHits.Load())
	w.Write([]byte(reset))
}
