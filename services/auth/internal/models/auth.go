package models

import "time"

type LoginInput struct {
	LoginID     string
	AccountType string
	Password    string
}

type LoginResponse struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    time.Duration
}

type LogoutInput struct {
	RefreshToken string
}

type LogoutResponse struct {
	Success bool
}
