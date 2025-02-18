package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AdminUser struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username       string             `bson:"username" json:"username"`
	Password       string             `bson:"password" json:"password"`
	ClearanceLevel int                `bson:"clearance_level" json:"clearance_level"`
}

type AdminSession struct {
	AdminID         primitive.ObjectID `bson:"admin_id" json:"admin_id"`
	AccessToken     string             `bson:"access_token" json:"access_token"` // unliked the normal session structure, admin sessions do not refresh their access token, it simply expires after 30 minutes and they are logged out (for security)
	AccessExpiresAt time.Time          `bson:"access_expires_at" json:"access_expires_at"`
	CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
	LastActivity    time.Time          `bson:"last_activity" json:"last_activity"`
}