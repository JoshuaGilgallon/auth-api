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
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FirstName   string             `bson:"first_name" json:"name"`
	LastName    string             `bson:"last_name" json:"last_name"`
	Email       string             `bson:"email" json:"email"`
	PhoneNumber string             `bson:"phone_number" json:"phone_number"`
	Password    string             `bson:"password" json:"password"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
	MFAEnabled  bool               `bson:"mfa_enabled" json:"mfa_enabled"`
	Status      string             `bson:"status" json:"status"`
}
