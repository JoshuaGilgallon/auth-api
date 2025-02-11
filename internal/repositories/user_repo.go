package repositories

import (
	"auth-api/internal/models"
)

func GetAllUsers() ([]models.User, error) {
	// Simulate a database fetch
	return []models.User{
		{ID: 1, Name: "Alice", Email: "alice@example.com"},
	}, nil
}

func SaveUser(user models.User) (models.User, error) {
	// Simulate saving user
	user.ID = 2
	return user, nil
}
