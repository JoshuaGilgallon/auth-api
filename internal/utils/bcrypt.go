package utils

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

const (
	MinPasswordLength = 5
	bcryptCost        = 12
)

func HashBcrypt(plainText string) (string, error) {
	if len(plainText) < MinPasswordLength {
		return "", errors.New("plainText too short for hasing")
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(plainText), bcryptCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func ValidateBcrypt(plainText, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plainText))
	return err == nil
}
