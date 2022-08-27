package auth

import (
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
)

func HashAndSalt(password string) string {
	h, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}

	return hex.EncodeToString(h)
}

func Compare(hash, password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return false
	}
	return true
}
