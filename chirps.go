package main

import (
	"encoding/json"
	"log"
	"net/http"

	"codeberg.org/kedarui/chirpy/internal/database"
)

func (cfg *apiConfig) handlerPostChirps(w http.ResponseWriter, r *http.Request) {
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
		log.Printf("Error: %s", err)
		respondWithError(w, http.StatusInternalServerError, data)
		return
	}
	cleaned_body := redactProfanity(c.Body)

	createChirpParams := database.CreateChirpParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Body:   cleaned_body,
		UserID: c.UserID,
	}

	createdChirp, err := cfg.database.CreateChirp(r.Context(), createChirpParams)
	if err != nil {
		log.Printf("Error: %s", err)
		respondWithError(w, http.StatusBadRequest, nil)
		return
	}

	newChirp := chirp{
		Body:   createdChirp.Body,
		UserID: createdChirp.UserID,
	}
	data, err := json.Marshal(&newChirp)
	if err != nil {
		marshalError(w, err)
		return
	}
	respondWithCreated(w, data)
}
