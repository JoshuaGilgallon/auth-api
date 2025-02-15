package utils

import (
	"crypto/sha256"
	"encoding/base64"
)

// THIS IS UNSALTED AND SHOULD NOT BE USED FOR PASSWORDS OR SENSITIVE INFORMATION, IT IS LOW SECURITY

func HashSHA(plainText string) string {
	hash := sha256.New()
	hash.Write([]byte(plainText))
	hashedBytes := hash.Sum(nil)
	return base64.StdEncoding.EncodeToString(hashedBytes)
}

func ValidateSHA(plainText, storedHash string) bool {
	return HashSHA(plainText) == storedHash
}
