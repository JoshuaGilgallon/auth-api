package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FirstName   string            `bson:"first_name" json:"name"`
	LastName    string            `bson:"last_name" json:"last_name"`
	Email       string            `bson:"email" json:"email"`
	PhoneNumber string            `bson:"phone_number" json:"phone_number"`
	Password    string            `bson:"password" json:"password"`
}