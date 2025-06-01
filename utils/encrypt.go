package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a plain password.
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

// CheckPassword verifies a password against a bcrypt hash.
func CheckPassword(hash, password string) error {
	return  bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
