package services

import (
	"auth-api/internal/errors"
	"auth-api/internal/models"
	"auth-api/internal/repositories"
	"auth-api/internal/utils"
	"time"
)

func Login(input models.LoginInput) (models.Session, error) {
	var user models.User

	user, err := repositories.GetUserByEmail(input.Email)
	if err != nil || user.ID.IsZero() {
		return models.Session{}, errors.NewInvalidCredentialsError("Incorrect Login Credentials", nil)
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

	// Get the user by ID
	user, err := repositories.GetUserByID(ID.Hex())
	if err != nil {
		return models.Session{}, errors.NewFailedToCreateError("Failed to retrieve user information.", nil)
	}

	// Update user with provided information
	user.FirstName = input.FirstName
	user.LastName = input.LastName

	// Parse and set birth date
	birthDate, err := time.Parse("2006-01-02", input.BirthDate)
	if err != nil {
		return models.Session{}, errors.NewValidationError("Invalid birth date format.", nil)
	}
	user.BirthDate = birthDate

	// Set email as verified and status as active
	user.EmailVerified = true
	user.Status = models.StatusActive
	user.UpdatedAt = time.Now()

	// Save the updated user
	_, err = repositories.SaveUser(user)
	if err != nil {
		return models.Session{}, errors.NewFailedToCreateError("Failed to update user information.", nil)
	}

	// Create session
	session, err := CreateSession(ID)
	if err != nil {
		return models.Session{}, errors.NewFailedToCreateError("Failed to create session. Please wait a moment and try again, or contact support for assistance.", nil)
	}

	return session, nil
}
