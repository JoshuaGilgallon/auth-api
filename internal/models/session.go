package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Session struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID           primitive.ObjectID `bson:"user_id" json:"user_id"`
	AccessToken      string             `bson:"access_token" json:"access_token"`
	RefreshToken     string             `bson:"refresh_token" json:"refresh_token"`
	AccessExpiresAt  time.Time          `bson:"access_expires_at" json:"access_expires_at"`
	RefreshExpiresAt time.Time          `bson:"refresh_expires_at" json:"refresh_expires_at"`
	CreatedAt        time.Time          `bson:"created_at" json:"created_at"`
	LastActivity     time.Time          `bson:"last_activity" json:"last_activity"`
}
