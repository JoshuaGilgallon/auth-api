package repositories

import (
	"auth-api/internal/models"
	"auth-api/internal/utils"
	"context"
	"log"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		UserID:           user.ID,
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

func VerifyEmail(code string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var email models.Email
	filter := bson.M{"code": code}
	err := emailCollection.FindOne(ctx, filter).Decode(&email)
	if err != nil {
		log.Printf("Error finding verification code: %v", err)
		return false, errors.Wrap(err, "failed to find verification code")
	}

	// check if the code is expired
	if time.Since(email.CreatedAt) > 30*time.Minute {
		log.Printf("Verification code expired: %v", err)
		InvalidateEmailToken(code)
		return false, errors.Wrap(err, "verification code expired")
	}

	// update the user email verified status
	user, err := GetUserByID(email.UserID.Hex())
	if err != nil {
		log.Printf("Error finding user: %v", err)
		return false, errors.Wrap(err, "failed to find user")
	}

	user.EmailVerified = true
	_, err = SaveUser(user)
	if err != nil {
		log.Printf("Error updating user: %v", err)
		return false, errors.Wrap(err, "failed to update user")
	}

	return true, nil
}

func GetIdFromCode(code string) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var email models.Email
	filter := bson.M{"code": code}
	err := emailCollection.FindOne(ctx, filter).Decode(&email)
	if err != nil {
		log.Printf("Error finding verification code: %v", err)
		return primitive.NilObjectID, errors.Wrap(err, "failed to find verification code")
	}

	return email.UserID, nil
}

func InvalidateEmailToken(code string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"code": code}
	_, err := emailCollection.DeleteOne(ctx, filter)
	if err != nil {
		log.Printf("Error deleting verification code: %v", err)
		return errors.Wrap(err, "failed to delete verification code")
	}

	return nil
}
