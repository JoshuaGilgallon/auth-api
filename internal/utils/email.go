package utils

import (
	"auth-api/internal/models"
	"bytes"
	"html/template"
	"log"
	"os"
	"path/filepath"

	"github.com/resendlabs/resend-go"
)

func RenderVerificationEmailTemplate(email models.Email) (string, error) {
	// get the base url from environment variables
	baseURL := os.Getenv("EMAIL_REDIRECT_BASE")

	verificationLink := baseURL + email.VerificationCode

	templateData := struct {
		VerificationLink string
	}{
		VerificationLink: verificationLink,
	}

	templatePath, err := filepath.Abs("templates/email/verify.html")
	if err != nil {
		return "", err
	}

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, templateData)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

// send the email with resend
func SendEmail(email models.Email) error {
	// render the HTML email template
	htmlContent, err := RenderVerificationEmailTemplate(email)
	if err != nil {
		log.Printf("Error rendering email template: %v", err)
		return err
	}

	apiKey := os.Getenv("RESEND_API_KEY")
	client := resend.NewClient(apiKey)

	params := &resend.SendEmailRequest{
		From:    "onboarding@resend.dev",
		To:      []string{email.Recipient},
		Subject: email.Subject,
		Html:    htmlContent,
	}

	_, err = client.Emails.Send(params)
	log.Printf("Err sending email: %v", err)
	return err
}
