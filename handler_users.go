package main

import (
	"database/sql"
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
		ID:          createdUser.ID,
		CreatedAt:   createdUser.CreatedAt,
		UpdatedAt:   createdUser.UpdatedAt,
		Email:       createdUser.Email,
		IsChirpyRed: false,
	}

	data, err := json.Marshal(&user)
	if err != nil {
		marshalError(w, err)
		return
	}
	respondWithCreated(w, data)
}

func (cfg *apiConfig) handlerLoginUser(w http.ResponseWriter, r *http.Request) {
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

	accessTokenExpiresIn := 3600 * time.Second
	tokenString, err := auth.MakeJWT(dbUser.ID, cfg.tokenSecret, accessTokenExpiresIn)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, nil)
		return
	}

	refreshTokenString := auth.MakeRefreshToken()
	refreshTokenExpiresIn := 1440 * time.Hour
	refreshTokenExpiresAt := time.Now().Add(refreshTokenExpiresIn)
	createRefreshTokenParams := database.CreateRefreshTokenParams{
		Token:     refreshTokenString,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    dbUser.ID,
		ExpiresAt: refreshTokenExpiresAt,
		RevokedAt: sql.NullTime{},
	}

	refreshToken, err := cfg.database.CreateRefreshToken(r.Context(), createRefreshTokenParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, nil)
		return
	}

	user := User{
		ID:           dbUser.ID,
		CreatedAt:    dbUser.CreatedAt,
		UpdatedAt:    dbUser.UpdatedAt,
		Email:        dbUser.Email,
		Token:        tokenString,
		RefreshToken: refreshToken.Token,
		IsChirpyRed:  dbUser.IsChirpyRed,
	}

	data, err := json.Marshal(&user)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, nil)
		return
	}

	respondWithOK(w, data)
}

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	userId, err := validateToken(w, r.Header, cfg.tokenSecret)
	if err != nil {
		return
	}

	var ac accountCreds
	err = json.NewDecoder(r.Body).Decode(&ac)
	if err != nil {
		marshalError(w, err)
		return
	}

	newHashedPassword, err := auth.HashPassword(ac.Password)

	updateUserPrams := database.UpdateUserParams{
		ID:        userId,
		UpdatedAt: time.Now(),
		Email: sql.NullString{
			String: ac.Email,
			Valid:  true,
		},
		HashedPassword: sql.NullString{
			String: newHashedPassword,
			Valid:  true,
		},
	}

	updatedDBUser, err := cfg.database.UpdateUser(r.Context(), updateUserPrams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, nil)
		return
	}

	updatedUser := User{
		ID:          updatedDBUser.ID,
		CreatedAt:   updatedDBUser.CreatedAt,
		UpdatedAt:   updatedDBUser.UpdatedAt,
		Email:       updatedDBUser.Email,
		IsChirpyRed: updatedDBUser.IsChirpyRed,
	}

	data, err := json.Marshal(&updatedUser)
	if err != nil {
		marshalError(w, err)
		return
	}
	respondWithOK(w, data)
}
