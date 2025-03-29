package models

type LoginInput struct {
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type RefLoginInput struct {
	RefCode string `json:"ref_code"`
}
