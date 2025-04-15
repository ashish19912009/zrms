package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthClaims struct {
	AccountID   string           `json:"account_id"`
	EmployeeID  string           `json:"employee_id"`
	AccountType string           `json:"account_type"`
	Name        string           `json:"name"`
	MobileNo    string           `json:"mobile_no"`
	Role        string           `json:"role"`
	Permissions PermissionsArray `json:"permissions"`
	jwt.RegisteredClaims
}

type TokenDetails struct {
	AccessToken   string    `json:"access_token"`
	RefreshToken  string    `json:"refresh_token"`
	AccessExpiry  time.Time `json:"access_exp"`
	RefreshExpiry time.Time `json:"refresh_exp"`
}

type RefreshTokenInput struct {
	AccountID    string `json:"account_id"`
	RefreshToken string `json:"refresh_token"`
}

type ValidateTokenInput struct {
	AccessToken string `json:"access_token"`
}
