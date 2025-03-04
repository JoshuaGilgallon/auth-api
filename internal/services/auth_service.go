package services

import (
	"auth-api/internal/errors"
	"auth-api/internal/models"
	"auth-api/internal/repositories"
	"auth-api/internal/utils"
)

type LoginInput struct {
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

func Login(input LoginInput) (models.Session, error) {
	var user models.User

	// Try to get the user by Email first
	user, err := repositories.GetUserByEmail(input.Email)
	if err != nil || user.ID.IsZero() {
		// If not found, try by phone number
		user, err = repositories.GetUserByPhoneNumber(input.PhoneNumber)
		if err != nil || user.ID.IsZero() {
			return models.Session{}, errors.NewInvalidCredentialsError("Incorrect Login Credentials", nil)
		}
	} else {
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
