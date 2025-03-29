package services

import (
	"auth-api/internal/models"
	"auth-api/internal/repositories"
	"auth-api/internal/utils"
	"os"

	"go.mongodb.org/mongo-driver/bson/primitive"
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

func VerifyEmail(code string) (string, error) { // url, err
	// verify the email code
	response, err := repositories.VerifyEmail(code)
	if err != nil {
		return "", err
	}

	if !response {
		return "", nil
	}

	// form the redirect
	url := os.Getenv("EMAIL_REDIRECT_BASE") + code

	return url, nil
}

func GetIdFromCode(code string) (primitive.ObjectID, error) {
	id, err := repositories.GetIdFromCode(code)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return id, nil
}
