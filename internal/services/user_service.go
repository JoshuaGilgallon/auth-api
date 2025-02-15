package services

import (
	"auth-api/internal/models"
	"auth-api/internal/repositories"
	"auth-api/internal/utils"
	"time"
)

type UserInput struct {
	FirstName   string `json:"name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
	MFAEnabled  bool   `json:"mfa_enabled"`
}

func CreateUser(input UserInput) (models.User, error) {
	// hash password before saving
	hashedPassword, err := utils.HashBcrypt(input.Password)

	if err != nil {
		return models.User{}, err
	}

	// encrypt the email and phone number before saving
	encryptedPhoneNumber, err := utils.Encrypt(input.PhoneNumber)
	if err != nil {
		return models.User{}, err
	}

	hashedPhoneNumber := utils.HashSHA(input.PhoneNumber)

	encryptedEmail, err := utils.Encrypt(input.Email)
	if err != nil {
		return models.User{}, err
	}

	hashedEmail := utils.HashSHA(input.Email)

	now := time.Now()
	user := models.User{
		FirstName:       input.FirstName,
		LastName:        input.LastName,
		Email:           encryptedEmail,
		EmailHash:       hashedEmail,
		PhoneNumber:     encryptedPhoneNumber,
		PhoneNumberHash: hashedPhoneNumber,
		Password:        hashedPassword,
		CreatedAt:       now,
		UpdatedAt:       now,
		MFAEnabled:      input.MFAEnabled,
		Status:          models.StatusActive, // set the status to active by default
	}
	return repositories.SaveUser(user)
}

func GetUser(id string) (models.User, error) {
	user, err := repositories.GetUserByID(id)
	if err != nil {
		return models.User{}, err
	}

	// decrypt email
	decryptedEmail, err := utils.Decrypt(user.Email)
	if err != nil {
		return models.User{}, err
	}
	user.Email = decryptedEmail

	// decrypt phone number
	decryptedPhone, err := utils.Decrypt(user.PhoneNumber)
	if err != nil {
		return models.User{}, err
	}
	user.PhoneNumber = decryptedPhone

	return user, nil
}
