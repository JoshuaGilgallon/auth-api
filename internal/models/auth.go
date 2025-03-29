package models

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RefLoginInput struct {
	RefCode string `json:"ref_code"`
}

type CompleteSignupResponse struct {
	Success     bool   `json:"success"`
	AccessToken string `json:"access_token"`
}
