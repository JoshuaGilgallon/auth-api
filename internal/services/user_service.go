package services

import (
	"auth-api/internal/models"
	"auth-api/internal/repositories"
)

type UserInput struct {
	FirstName 	string `json:"name"`
	LastName 	string `json:"last_name"`
	Email 		string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Password 	string `json:"password"`
}

func CreateUser(input UserInput) (models.User, error) {
	user := models.User{
		FirstName:   input.FirstName,
		LastName:    input.LastName,
		Email:       input.Email,
		PhoneNumber: input.PhoneNumber,
		Password:    input.Password,
	}
	return repositories.SaveUser(user)
}

func GetUser(id string) (models.User, error) {
	return repositories.GetUserByID(id)
}
