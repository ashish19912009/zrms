package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTIssuer struct {
	secretKey     string
	tokenDuration time.Duration
}

type CustomClaims struct {
	AccountID string `json:"account_id"`
	Role      string `json:"role"`
	jwt.RegisteredClaims
}

func NewJWTIssuer(secretKey string, tokenDuration time.Duration) *JWTIssuer {
	return &JWTIssuer{
		secretKey:     secretKey,
		tokenDuration: tokenDuration,
	}
}

func (j *JWTIssuer) GenerateToken(accountID, role string) (string, error) {
	if accountID == "" || role == "" {
		return "", errors.New("accountID and role are required for token generation")
	}

	now := time.Now()
	claims := CustomClaims{
		AccountID: accountID,
		Role:      role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(j.tokenDuration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}
