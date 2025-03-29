package services

import (
	"auth-api/internal/models"
	"auth-api/internal/repositories"
	"auth-api/internal/utils"
)

func CreateVerifEmail(input models.VerifEmailInput) (string, error) {
	email, err := repositories.CreateVerificationEmail(input)
	if err != nil {
		return "", err
	}

	// Send the email using the email service
	err = utils.SendEmail(email)
	if err != nil {
		return "", err
	}

	return "Verification email sent successfully", nil
}
