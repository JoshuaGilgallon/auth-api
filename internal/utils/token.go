package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"time"
)

const (
	tokenLength = 32
)

func GenToken() (string, error) {
	b := make([]byte, tokenLength)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	hash := sha256.New()
	hash.Write(b)
	hash.Write([]byte(time.Now().String()))

	return hex.EncodeToString(hash.Sum(nil)), nil
}
