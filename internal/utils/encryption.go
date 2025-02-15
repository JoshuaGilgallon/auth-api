package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// load the encryption key from its environment variable
func getKey() []byte {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	keyHex := os.Getenv("USER_AES_KEY")
	if keyHex == "" {
		log.Fatal("USER_AES_KEY must be set in the .env file")
	}

	// decode it into raw bytes
	keyBytes, err := hex.DecodeString(keyHex)
	if err != nil {
		log.Fatal("USER_AES_KEY is not a valid hex string")
	}

	// ensure the key is 32 bytes - AES-256 requirement
	if len(keyBytes) != 32 {
		log.Fatal("USER_AES_KEY must be exactly 32 bytes long")
	}

	return keyBytes
}

func Encrypt(plainText string) (string, error) {
	key := getKey()
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, 12)
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	cipherText := aead.Seal(nonce, nonce, []byte(plainText), nil)
	return hex.EncodeToString(cipherText), nil
}

func Decrypt(cipherText string) (string, error) {
	key := getKey()
	cipherTextBytes, err := hex.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	nonce := cipherTextBytes[:12]
	encryptedData := cipherTextBytes[12:]

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	plainText, err := aead.Open(nil, nonce, encryptedData, nil)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}
