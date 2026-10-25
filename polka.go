package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"codeberg.org/kedarui/chirpy/internal/auth"
	"codeberg.org/kedarui/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerPolkaWebhook(w http.ResponseWriter, r *http.Request) {
	key, err := auth.GetAPIKey(r.Header)
	if err != nil || key != cfg.polkaKey {
		data := []byte("You are not authorized")
		respondWithError(w, http.StatusUnauthorized, data)
		return
	}

	var polka polkaRequest
	err = json.NewDecoder(r.Body).Decode(&polka)
	if err != nil {
		marshalError(w, err)
		return
	}

	if polka.Event != "user.upgraded" {
		respondWithNoContent(w, nil)
		return
	}

	userId, err := uuid.Parse(polka.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, nil)
		return
	}

	updateChirpyRedStatusParams := database.UpdateUserChirpyRedStatusParams{
		ID: userId,
		UpdatedAt: time.Now(),
		IsChirpyRed: true,
		
	}
	dbUser, err := cfg.database.UpdateUserChirpyRedStatus(r.Context(), updateChirpyRedStatusParams)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, nil)
			return
		}
		respondWithError(w, http.StatusInternalServerError, nil)
		return
	}
	log.Printf("dbUser: %v", dbUser)
	respondWithNoContent(w, nil)
}
