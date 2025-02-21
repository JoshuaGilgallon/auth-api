package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	StatusActive   = "active"
	StatusInactive = "inactive"
	StatusLocked   = "locked"
)

type User struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FirstName       string             `bson:"first_name" json:"name"`
	LastName        string             `bson:"last_name" json:"last_name"`
	Email           string             `bson:"email" json:"email"`
	EmailHash       string             `bson:"email_hash" json:"email_hash"`
	PhoneNumber     string             `bson:"phone_number" json:"phone_number"`
	PhoneNumberHash string             `bson:"phone_number_hash" json:"phone_number_hash"`
	Password        string             `bson:"password" json:"password"`
	Bio             string             `bson:"bio" json:"bio"`
	BirthDate       time.Time          `bson:"birth_date" json:"birth_date"`
	Language        string             `bson:"language" json:"language"`
	CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time          `bson:"updated_at" json:"updated_at"`
	MFAEnabled      bool               `bson:"mfa_enabled" json:"mfa_enabled"`
	Status          string             `bson:"status" json:"status"`
}
