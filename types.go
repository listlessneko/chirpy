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
	polkaKey       string
	fileserverHits atomic.Int32
}

type accountCreds struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type polkaRequest struct {
	Event string `json:"event"`
	Data  struct {
		UserID string `json:"user_id"`
	} `json:"data"`
}

type User struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
}

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
}

type ErrResp struct {
	Error string `json:"error"`
}
