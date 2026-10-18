package main

import (
	"slices"
	"strings"
)

var theProfane = []string{"kerfuffle", "sharbert", "fornax"}

const (
	redacted = "****"
)

func redactProfanity(text string) string {
	words := strings.Fields(text)
	for i, w := range words {
		if slices.Contains(theProfane, strings.ToLower(w)) {
			words[i] = redacted
		}
	}
	return strings.Join(words, " ")
}
