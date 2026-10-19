package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"codeberg.org/kedarui/chirpy/internal/database"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func (cfg *apiConfig) handlerPostChirp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var c chirp
	err := decoder.Decode(&c)
	if len(c.Body) > 140 {
		tooLongResp := errResp{Error: "Chirp is too long"}
		data, err := json.Marshal(&tooLongResp)
		if err != nil {
			marshalError(w, err)
			return
		}
		log.Printf("Error: %v", err)
		respondWithError(w, http.StatusBadRequest, data)
		return
	}
	if err != nil {
		errResp := errResp{Error: "Something went wrong"}
		data, err := json.Marshal(&errResp)
		if err != nil {
			marshalError(w, err)
			return
		}
		log.Printf("Error: %v", err)
		respondWithError(w, http.StatusInternalServerError, data)
		return
	}
	cleaned_body := redactProfanity(c.Body)

	createChirpParams := database.CreateChirpParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Body:      cleaned_body,
		UserID:    c.UserID,
	}

	createdChirp, err := cfg.database.CreateChirp(r.Context(), createChirpParams)
	if err != nil {
		log.Printf("Error: %v", err)
		respondWithError(w, http.StatusBadRequest, nil)
		return
	}

	newChirp := chirp{
		ID:        createdChirp.ID,
		Body:      createdChirp.Body,
		CreatedAt: createdChirp.CreatedAt,
		UpdatedAt: createdChirp.UpdatedAt,
		UserID:    createdChirp.UserID,
	}
	data, err := json.Marshal(&newChirp)
	if err != nil {
		marshalError(w, err)
		return
	}
	respondWithCreated(w, data)
}

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.database.GetAllChirps(r.Context())
	if err != nil {
		log.Printf("Error retrieving all chirps: %s\n", err)
		respondWithError(w, http.StatusInternalServerError, nil)
	}

	var allChirps []chirp

	for _, c := range dbChirps {
		chirp := chirp{
			ID:        c.ID,
			Body:      c.Body,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			UserID:    c.UserID,
		}
		allChirps = append(allChirps, chirp)
	}

	data, err := json.Marshal(&allChirps)
	if err != nil {
		marshalError(w, err)
		return
	}
	respondWithOK(w, data)
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	chirpId, err := uuid.Parse(r.PathValue("chirp_id"))
	if err != nil {
		log.Printf("Invalid uuid: %v\n", err)
		data := []byte(err.Error() + "\n")
		respondWithError(w, http.StatusBadRequest, data)
		return
	}

	dbChirp, err := cfg.database.GetChirpById(r.Context(), chirpId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			data := []byte("Chirp not found")
			respondWithError(w, http.StatusNotFound, data)
			return
		}
		log.Printf("Error retrieving this chirp: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, nil)
		return
	}

	respChirp := chirp{
		ID:        dbChirp.ID,
		Body:      dbChirp.Body,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		UserID:    dbChirp.UserID,
	}

	data, err := json.Marshal(&respChirp)
	if err != nil {
		marshalError(w, err)
		return
	}
	respondWithOK(w, data)
}
