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

func (cfg * apiConfig) handlerLoginUser(w http.ResponseWriter, r *http.Request) {
	var ac accountCreds
	err := json.NewDecoder(r.Body).Decode(&ac)
	if err != nil {
		marshalError(w, err)
		return
	}

	dbUser, err := cfg.database.GetUserByEmail(r.Context(), ac.Email)
	if err != nil {
		data := []byte("Incorrect email or password")
		respondWithError(w, http.StatusUnauthorized, data)
		return
	}

	ok, err := auth.CheckPasswordHash(ac.Password, dbUser.HashedPassword)
	if !ok || err != nil {
		data := []byte("Incorrect email or password")
		respondWithError(w, http.StatusUnauthorized, data)
		return
	}

	expiresIn := 3600 * time.Second

	if ac.ExpiresInSeconds != 0 && ac.ExpiresInSeconds < int(expiresIn) {
		expiresIn = time.Duration(ac.ExpiresInSeconds)
	}

	tokenString, err := auth.MakeJWT(dbUser.ID, cfg.tokenSecret, expiresIn)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, nil)
		return
	}

	user := User{
		ID: dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email: dbUser.Email,
		Token: tokenString,
	}

	data, err := json.Marshal(&user)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, nil)
		return
	}
	
	respondWithOK(w, data)
}
