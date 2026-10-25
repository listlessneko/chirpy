package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"codeberg.org/kedarui/chirpy/internal/auth"
	"codeberg.org/kedarui/chirpy/internal/database"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		data := []byte("You are not authorized")
		respondWithError(w, http.StatusUnauthorized, data)
		return
	}

	dbRefreshToken, err := cfg.database.GetRefreshToken(r.Context(), refreshToken)
	if err != nil || dbRefreshToken.RevokedAt.Valid {
		data := []byte("You are not authorized")
		respondWithError(w, http.StatusUnauthorized, data)
		return
	}

	expiresIn := 60 * time.Minute
	accessToken, err := auth.MakeJWT(dbRefreshToken.UserID, cfg.tokenSecret, expiresIn)
	if err != nil {
		data := []byte("Error creating token")
		respondWithError(w, http.StatusInternalServerError, data)
		return
	}

	user := User{
		Token: accessToken,
	}
	data, err := json.Marshal(&user)
	if err != nil {
		marshalError(w, err)
		return
	}
	respondWithOK(w, data)
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		data := []byte("You are not authorized")
		respondWithError(w, http.StatusUnauthorized, data)
	}

	revokeTokenParams := database.RevokeRefreshTokenParams{
		Token: refreshToken,
		RevokedAt: sql.NullTime{
			Time: time.Now(),
			Valid: true,
		},
	}

	_, err = cfg.database.RevokeRefreshToken(r.Context(), revokeTokenParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, nil)
		return
	}
	respondWithNoContent(w, nil)
}
