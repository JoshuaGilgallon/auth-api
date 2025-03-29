package repositories

import (
	"auth-api/internal/models"
	"auth-api/internal/utils"
	"context"
	"log"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var emailCollection *mongo.Collection

// initialises the user collection
func SetEmailCollection(collection *mongo.Collection) {
	emailCollection = collection
}

func ValidateRefCode(refCode string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"code": refCode}
	err := emailCollection.FindOne(ctx, filter).Decode(&refCode)
	if err != nil {
		log.Printf("Error finding referral code: %v", err)
		return false, errors.Wrap(err, "failed to find referral code")
	}

	return true, nil
}

func CreateVerificationEmail(input models.VerifEmailInput) (models.Email, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// search for user id and get the email if it exists

	log.Printf("User id: %s", input.UserID)

	user, err := GetUserByID(input.UserID)
	if err != nil {
		log.Printf("Error finding user: %v", err)
		return models.Email{}, errors.Wrap(err, "failed to find user")
	}

	// if the users email is already verified, return error message
	if user.EmailVerified {
		return models.Email{}, errors.Wrap(err, "Email already verified")
	}

	// generate a verification code
	vCode, err := utils.GenToken()
	if err != nil {
		log.Printf("Error generating verification code: %v", err)
		return models.Email{}, errors.Wrap(err, "failed to generate verification code")
	}

	email := models.Email{
		Recipient:        user.Email,
		Subject:          "Email Verification",
		Body:             "Please verify your email by clicking the link below.",
		VerificationCode: vCode,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	_, err = emailCollection.InsertOne(ctx, email)
	if err != nil {
		log.Printf("Error inserting email: %v", err)
		return models.Email{}, errors.Wrap(err, "failed to insert email")
	}

	return email, nil
}
