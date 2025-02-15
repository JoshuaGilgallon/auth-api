package services

import (
	"auth-api/internal/errors"
	"auth-api/internal/models"
	"auth-api/internal/repositories"
	"auth-api/internal/utils"
	"log"
)

type LoginInput struct {
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

func Login(input LoginInput) (models.Session, error) {
	var user models.User

	log.Printf("Login attempt - Email: %s, Phone: %s", input.Email, input.PhoneNumber)

	// Try to get the user by Email first
	user, err := repositories.GetUserByEmail(input.Email)
	log.Printf("GetUserByEmail returned - User: %+v, Error: %v", user, err)
	if err != nil || user.ID.IsZero() {
		log.Printf("User not found with email: %s, error: %v", input.Email, err)
		// If not found, try by phone number
		user, err = repositories.GetUserByPhoneNumber(input.PhoneNumber)
		if err != nil || user.ID.IsZero() {
			log.Printf("User not found with phone number: %s, error: %v", input.PhoneNumber, err)
			return models.Session{}, errors.NewInvalidCredentialsError("Incorrect Login Credentials", nil)
		}
		log.Printf("User found with phone number: %s", input.PhoneNumber)
	} else {
		log.Printf("User found with email: %s", input.Email)
	}

	// Validate password
	if !utils.ValidateBcrypt(input.Password, user.Password) {
		return models.Session{}, errors.NewInvalidCredentialsError("The Password is Incorrect", nil)
	}

	// Create session
	session, err := CreateSession(user.ID)
	if err != nil {
		return models.Session{}, errors.NewFailedToCreateError("Failed to create session. Please wait a moment and try again, or contact support for assistance.", nil)
	}

	return session, nil
}

func Logout(accessToken string) error {
	return InvalidateSessionByToken(accessToken)
}