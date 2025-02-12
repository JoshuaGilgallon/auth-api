package services

import "auth-api/internal/models"

type UserInput struct {
	FirstName 	string `json:"name"`
	LastName 	string `json:"last_name"`
	Email 		string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Password 	string `json:"password"`
}

func CreateUser(input UserInput) models.User {
	// Simulate saving to DB
	return models.User{
		ID: 2, 
		FirstName: input.FirstName, 
		LastName: input.LastName, 
		Email: input.Email, 
		PhoneNumber: input.PhoneNumber,
		Password: input.Password,
	}
}
