package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	StatusPending  = "pending"
	StatusActive   = "active"
	StatusInactive = "inactive"
	StatusLocked   = "locked"
)

type User struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FirstName     string             `bson:"first_name" json:"name"`
	LastName      string             `bson:"last_name" json:"last_name"`
	Email         string             `bson:"email" json:"email"`
	EmailVerified bool               `bson:"email_verified" json:"email_verified"`
	PhoneNumber   string             `bson:"phone_number" json:"phone_number"`
	Password      string             `bson:"password" json:"password"`
	Bio           string             `bson:"bio" json:"bio"`
	BirthDate     time.Time          `bson:"birth_date" json:"birth_date"`
	Language      string             `bson:"language" json:"language"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at"`
	MFAEnabled    bool               `bson:"mfa_enabled" json:"mfa_enabled"`
	Status        string             `bson:"status" json:"status"`
}

type UserInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SetupUserInput struct {
	FirstName   string `json:"name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	BirthDate   string `json:"birth_date"`
}

type FullUserInput struct {
	FirstName   string `json:"name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
	BirthDate   string `json:"birth_date"`
	Language    string `json:"language"`
	MFAEnabled  bool   `json:"mfa_enabled"`
}

type UserCreateReturn struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type UserAdvancedSearchCriteria struct {
	FirstName       string     `json:"first_name"`
	LastName        string     `json:"last_name"`
	Email           string     `json:"email"`
	PhoneNumber     string     `json:"phone_number"`
	StartTime       *time.Time `json:"start_time"`
	EndTime         *time.Time `json:"end_time"`
	UpdateStartTime *time.Time `json:"update_start_time"`
	UpdateEndTime   *time.Time `json:"update_end_time"`
	PageNumber      int64      `form:"page_number" json:"page_number"`
	PageSize        int64      `form:"page_size" json:"page_size"`
}

type UserSearchCriteria struct {
	SearchTerm string `form:"search_term" json:"search_term"`
	PageNumber int64  `form:"page_number" json:"page_number"`
	PageSize   int64  `form:"page_size" json:"page_size"`
}
