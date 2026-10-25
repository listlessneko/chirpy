package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func HashPassword(password string) (string, error) {
	var hash string
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return hash, err
	}
	return hash, nil
}

func CheckPasswordHash(password string, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(password, hash)
}

func MakeJWT(userId uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	t := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			Issuer:    "chirpy-access",
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
			Subject:   userId.String(),
		},
	)
	s, err := t.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	return s, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	var userId uuid.UUID
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return userId, err
	}

	subject, err := token.Claims.GetSubject()
	if err != nil {
		return userId, err
	}

	userId, err = uuid.Parse(subject)
	if err != nil {
		return userId, err
	}

	return userId, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	var bearer string
	bearer = headers.Get("Authorization")
	if bearer == "" {
		return bearer, errors.New("No bearer found")
	}
	prefix := "Bearer "
	return bearer[len(prefix):], nil
}

func MakeRefreshToken() string {
	key := make([]byte, 32)
	rand.Read(key)
	return hex.EncodeToString(key)
}

func GetAPIKey(headers http.Header) (string, error) {
	var key string
	key = headers.Get("Authorization")
	if key == "" {
		return key, errors.New("No key found")
	}
	prefix := "apiKey "
	return key[len(prefix):], nil
}
