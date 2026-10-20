package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"codeberg.org/kedarui/chirpy/internal/auth"
	"codeberg.org/kedarui/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	var ac accountCreds
	err := json.NewDecoder(r.Body).Decode(&ac)
	if err != nil {
		marshalError(w, err)
		return
	}

	hashedPassword, err := auth.HashPassword(ac.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, nil)
	}

	createUserParams := database.CreateUserParams{
		ID:             uuid.New(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		Email:          ac.Email,
		HashedPassword: hashedPassword,
	}

	createdUser, err := cfg.database.CreateUser(r.Context(), createUserParams)
	if err != nil {
		log.Printf("Error creating user: %s", err)
		respondWithError(w, http.StatusInternalServerError, nil)
		return
	}

	user := User{
		ID:        createdUser.ID,
		CreatedAt: createdUser.CreatedAt,
		UpdatedAt: createdUser.UpdatedAt,
		Email:     createdUser.Email,
	}

	data, err := json.Marshal(&user)
	if err != nil {
		marshalError(w, err)
		return
	}
	respondWithCreated(w, data)
}
