package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCheckPasswordHash(t *testing.T) {
	password1 := "1password!"
	password2 := "2passwords?"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name          string
		password      string
		hash          string
		wantErr       bool
		matchPassword bool
	}{
		{
			name:          "Correct password",
			password:      password1,
			hash:          hash1,
			wantErr:       false,
			matchPassword: true,
		},
		{
			name:          "Incorrect password",
			password:      "wrongPassword",
			hash:          hash1,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Password doesn't match different hash",
			password:      password1,
			hash:          hash2,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Empty password",
			password:      "",
			hash:          hash1,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Invalid hash",
			password:      password1,
			hash:          "invalidhash",
			wantErr:       true,
			matchPassword: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, err := CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && match != tt.matchPassword {
				t.Errorf("CheckPasswordHash() expects %v, got %v", tt.matchPassword, match)
			}
		})
	}
}

func TestMakeJWT(t *testing.T) {
	userId1 := uuid.New()
	tokenSecret1 := "ThisIsAVerySecretToken"
	tokenSecret2 := "ThisTooIsAVerySecretToken"
	posOneMinute := time.Minute
	negOneMinute := -time.Minute

	tokenString1, err := MakeJWT(userId1, tokenSecret1, posOneMinute)
	if err != nil {
		t.Fatalf("MakeJWT() failed: %v", err)
	}
	tokenString2, err := MakeJWT(userId1, tokenSecret2, posOneMinute)
	if err != nil {
		t.Fatalf("MakeJWT() failed: %v", err)
	}
	expiredTokenString, err := MakeJWT(userId1, tokenSecret1, negOneMinute)
	if err != nil {
		t.Fatalf("MakeJWT() failed: %v", err)
	}

	tests := []struct {
		name        string
		tokenString string
		tokenSecret string
		wantErr     bool
		matchUserId uuid.UUID
	}{
		{
			name:        "Correct token",
			tokenString: tokenString1,
			tokenSecret: tokenSecret1,
			wantErr:     false,
			matchUserId: userId1,
		},
		{
			name:        "Incorrect token",
			tokenString: "This is definitely a token",
			tokenSecret: tokenSecret1,
			wantErr:     true,
			matchUserId: uuid.Nil,
		},
		{
			name:        "Token doesn't match token with a different secret",
			tokenString: tokenString2,
			tokenSecret: tokenSecret1,
			wantErr:     true,
			matchUserId: uuid.Nil,
		},
		{
			name:        "Expired token",
			tokenString: expiredTokenString,
			tokenSecret: tokenSecret1,
			wantErr:     true,
			matchUserId: uuid.Nil,
		},
		{
			name:        "Empty token string",
			tokenString: "",
			tokenSecret: tokenSecret1,
			wantErr:     true,
			matchUserId: uuid.Nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, err := ValidateJWT(tt.tokenString, tt.tokenSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && match != tt.matchUserId {
				t.Errorf("ValidateJWT() expects %v, got %v", tt.matchUserId, match)
			}
		})
	}
}
