package generator

import (
	"crypto/rand"
	"math/big"
)

func Generate(passLength int) string {

	if passLength <= 0 {
		passLength = 32
	}

	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"

	password := ""

	for i := 0; i < passLength; i++ {

		randomIndex, err := rand.Int(
			rand.Reader,
			big.NewInt(int64(len(charset))),
		)

		if err != nil {
			panic(err)
		}

		password += string(charset[randomIndex.Int64()])
	}

	return password
}
