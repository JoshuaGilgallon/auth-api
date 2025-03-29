package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Email struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID           primitive.ObjectID `bson:"user_id" json:"user_id"`
	Recipient        string             `bson:"recipient" json:"recipient"`
	Subject          string             `bson:"subject" json:"subject"`
	Body             string             `bson:"body" json:"body"`
	VerificationCode string             `bson:"code" json:"verification_id"`
	CreatedAt        time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time          `bson:"updated_at" json:"updated_at"`
}

type VerifEmailInput struct {
	UserID string `bson:"user_id" json:"user_id"`
}

type EmailCodeInput struct {
	Code string `bson:"code" json:"code"`
}
