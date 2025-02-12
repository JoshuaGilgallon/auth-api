package models

type User struct {
	ID 			int    `json:"id"`
	FirstName 	string `json:"name"`
	LastName 	string `json:"last_name"`
	Email 		string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Password 	string `json:"password"`
}