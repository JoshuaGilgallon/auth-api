package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Session struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID       primitive.ObjectID `bson:"user_id" json:"user_id"`
	Token        string             `bson:"token" json:"token"`
	ExpiresAt    time.Time          `bson:"expires_at" json:"expires_at"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	IPAddress    string             `bson:"ip_address" json:"ip_address"`
	UserAgent    string             `bson:"user_agent" json:"user_agent"`
	LastActivity time.Time          `bson:"last_activity" json:"last_activity"`
	IssuedAt     time.Time          `bson:"issued_at" json:"issued_at"`
	RefreshToken string             `bson:"refresh_token" json:"refresh_token"`
}
