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

type UserAdvancedSearchCriteria struct {
	FirstName       string     `bson:"first_name" json:"first_name"`
	LastName        string     `bson:"last_name" json:"last_name"`
	Email           string     `bson:"email" json:"email"`
	PhoneNumber     string     `bson:"phone_number" json:"phone_number"`
	StartTime       *time.Time `bson:"start_time" json:"start_time"`
	EndTime         *time.Time `bson:"end_time" json:"end_time"`
	UpdateStartTime *time.Time `bson:"update_start_time" json:"update_start_time"`
	UpdateEndTime   *time.Time `bson:"update_end_time" json:"update_end_time"`
	PageNumber      string     `bson:"page_number" json:"page_number"`
	PageSize        string     `bson:"page_size" json:"page_size"`
}

type UserSearchCriteria struct {
	SearchTerm string `form:"search_term" json:"search_term"`
	PageNumber int64  `form:"page_number" json:"page_number"`
	PageSize   int64  `form:"page_size" json:"page_size"`
}
