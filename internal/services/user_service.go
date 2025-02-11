package services

import "auth-api/internal/models"

type UserInput struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func FetchUsers() []models.User {
	// Simulated database call
	return []models.User{
		{ID: 1, Name: "John Doe", Email: "john@example.com"},
	}
}

func CreateUser(input UserInput) models.User {
	// Simulate saving to DB
	return models.User{ID: 2, Name: input.Name, Email: input.Email}
}
