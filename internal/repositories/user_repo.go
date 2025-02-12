package repositories

import (
	"auth-api/internal/models"
)

func SaveUser(user models.User) (models.User, error) {
	// Simulate saving user
	user.ID = 2
	return user, nil
}
