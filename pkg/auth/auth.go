package auth

import (
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
)

// HashAndSalt hashes the password with bcrypt and returns a string
func HashAndSalt(password string) string {
	h, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}

	return hex.EncodeToString(h)
}

// Compare compares a hash and password.
// Returns true or false.
func Compare(hash, password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return false
	}
	return true
}
