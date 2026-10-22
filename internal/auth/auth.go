package auth

import (
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
			Issuer: "chirpy-access",
			IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
			Subject: userId.String(),
		},
		)
	s, err := t.SignedString(tokenSecret)
	if err != nil {
		return "", err
	}
	return s, nil
}
