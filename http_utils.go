package main

import (
	"log"
	"net/http"

	"codeberg.org/kedarui/chirpy/internal/auth"
	"github.com/google/uuid"
)

func marshalError(w http.ResponseWriter, err error) {
	log.Printf("Error marshalling json: %s", err)
	w.WriteHeader(http.StatusInternalServerError)
}

func respondWithError(w http.ResponseWriter, statusCode int, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if !(len(data) > 0) {
		data = []byte("Something went wrong\n")
	}
	w.Write(data)
}

func validateToken(w http.ResponseWriter, headers http.Header, tokenSecret string) (uuid.UUID, error) {
	var userId uuid.UUID
	bearerToken, err := auth.GetBearerToken(headers)
	if err == nil {
		userId, err = auth.ValidateJWT(bearerToken, tokenSecret)
	}
	if err != nil {
		data := []byte("You are not authorized")
		respondWithError(w, http.StatusUnauthorized, data)
	}
	return userId, nil
}

func respondWithOK(w http.ResponseWriter, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	data = append(data, byte('\n'))
	w.Write(data)
}

func respondWithCreated(w http.ResponseWriter, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	data = append(data, byte('\n'))
	w.Write(data)
}

func respondWithNoContent(w http.ResponseWriter, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	data = append(data, byte('\n'))
	w.Write(data)
}
