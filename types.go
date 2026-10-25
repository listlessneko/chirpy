package main

import (
	"sync/atomic"
	"time"

	"codeberg.org/kedarui/chirpy/internal/database"
	"github.com/google/uuid"
)

type apiConfig struct {
	database       *database.Queries
	platform       string
	tokenSecret    string
	fileserverHits atomic.Int32
}

type accountCreds struct {
	Email            string `json:"email"`
	Password         string `json:"password"`
	ExpiresInSeconds int    `json:"expires_in_seconds"`
}

type User struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
}

type chirp struct {
	ID        uuid.UUID `json:"id"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
}

type cleaned_chirp struct {
	CleanedBody string `json:"cleaned_body"`
}

type errResp struct {
	Error string `json:"error"`
}

type validResp struct {
	Valid bool `json:"valid"`
}
