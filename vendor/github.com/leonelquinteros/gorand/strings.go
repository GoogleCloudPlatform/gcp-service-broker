package gorand

import (
	"crypto/rand"
	"math/big"
)

const (
	lowercase = "abcdefghijklmnopqrstuvwxyz"
	uppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numbers   = "0123456789"
)

// GetRandomChars returns a fixed length string (n int) of chars contained on a collection string provided (c string).
func GetRandomChars(c string, n int) (string, error) {
	var r string

	// Get random chars one by one
	for i := 0; i < n; i++ {
		// Read random position
		p, err := rand.Int(rand.Reader, big.NewInt(int64(len(c))))
		if err != nil {
			return "", err
		}

		r += string(c[p.Int64()])
	}

	return r, nil
}

// GetAlphaNumString returns a fixed length (n int) string of random letters and numbers [a-z][A-Z][0-9]
func GetAlphaNumString(n int) (string, error) {
	return GetRandomChars(lowercase+uppercase+numbers, n)
}

// GetAlphaString returns a fixed length (n int) string of random letters [a-z][A-Z]
func GetAlphaString(n int) (string, error) {
	return GetRandomChars(lowercase+uppercase, n)
}

// GetNumString returns a fixed length (n int) string of random numbers [0-9]
func GetNumString(n int) (string, error) {
	return GetRandomChars(numbers, n)
}
