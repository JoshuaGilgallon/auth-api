package utils

import (
	"auth-api/internal/models"
	"os"

	"github.com/resendlabs/resend-go"
)

// send the email with resend
func SendEmail(email models.Email) error {
	apiKey := os.Getenv("RESEND_API_KEY")
	client := resend.NewClient(apiKey)

	params := &resend.SendEmailRequest{
		From:    "onboarding@resend.dev",
		To:      []string{email.Recipient},
		Subject: email.Subject,
		Html:    email.Body,
	}

	_, err := client.Emails.Send(params)
	return err
}
