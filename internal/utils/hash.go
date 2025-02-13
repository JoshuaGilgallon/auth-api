package utils

import (
	"crypto/subtle"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

const (
	MinPasswordLength = 5
	bcryptCost = 12
)

func HashPassword(password string) (string, error) {
	if len(password) < MinPasswordLength {
		return "", errors.New("password too short")
	}
	
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func ValidatePassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	
	// Convert bool to int32 properly
	if err == nil {
		return subtle.ConstantTimeEq(1, 1) == 1
	}
	return subtle.ConstantTimeEq(0, 1) == 1
}