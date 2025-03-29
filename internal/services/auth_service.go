package services

import (
	"auth-api/internal/errors"
	"auth-api/internal/models"
	"auth-api/internal/repositories"
	"auth-api/internal/utils"
)

func Login(input models.LoginInput) (models.Session, error) {
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

func LoginWithEmailRef(input models.RefLoginInput) (models.Session, error) {
	var user models.User

	valid, err := repositories.ValidateRefCode(input.RefCode)
	if err != nil || !valid {
		return models.Session{}, errors.NewInvalidCredentialsError("Invalid Referral Code. Please resend a verification email.", nil)
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

func CompleteSignup(input models.SetupUserInput) (models.Session, error) {
	ID, err := GetIdFromCode(input.Token)
	if err != nil {
		return models.Session{}, errors.NewInvalidCredentialsError("Invalid Token. Please resend a verification email.", nil)
	}

	// Create session
	session, err := CreateSession(ID)
	if err != nil {
		return models.Session{}, errors.NewFailedToCreateError("Failed to create session. Please wait a moment and try again, or contact support for assistance.", nil)
	}

	return session, nil
}
