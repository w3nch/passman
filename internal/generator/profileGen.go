package generator

import (
	"crypto/rand"
	"math/big"
)

type Profile struct {
	Username string
	Password string
}

func GenerateProfile() (string, string) {
	adjectives := []string{"quick", "dark", "silent", "wild", "cool", "swift", "bold", "neat"}
	nouns := []string{"fox", "wolf", "hawk", "bear", "lion", "eagle", "tiger", "shadow"}
	suffixes := []string{"123", "99", "77", "42", "x", "z"}

	randomWord := func(words []string) string {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(words))))
		return words[n.Int64()]
	}

	username := randomWord(adjectives) + randomWord(nouns) + randomWord(suffixes)
	password := Generate(16)

	return username, password
}