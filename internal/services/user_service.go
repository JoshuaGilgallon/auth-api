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
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return models.User{}, err
	}

	now := time.Now()
	user := models.User{
		FirstName:   input.FirstName,
		LastName:    input.LastName,
		Email:       input.Email,
		PhoneNumber: input.PhoneNumber,
		Password:    hashedPassword,
		CreatedAt:   now,
		UpdatedAt:   now,
		MFAEnabled:  input.MFAEnabled,
		Status:      models.StatusActive, // Set default status to active
	}
	return repositories.SaveUser(user)
}

func GetUser(id string) (models.User, error) {
	return repositories.GetUserByID(id)
}
