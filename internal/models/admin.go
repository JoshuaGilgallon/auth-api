package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AdminUser struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username       string             `bson:"username" json:"username"`
	Password       string             `bson:"password" json:"password"`
	ClearanceLevel int                `bson:"clearance_level" json:"clearance_level"`
}
