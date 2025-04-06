package token

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/ashish19912009/services/auth/internal/constants"
	"github.com/ashish19912009/services/auth/internal/logger"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	EmployeeID  string   `json:"employee_id"`
	AccountID   string   `json:"account_id"`
	AccountType string   `json:"account_type"`
	Name        string   `json:"name"`
	MobileNo    string   `json:"mobile_no"`
	Permissions []string `json:"permissions"`
	jwt.RegisteredClaims
}

type TokenManager interface {
	GenerateAccessToken(accountID, employeeID, mobileNo, accountType, name string, permissions []string, duration time.Duration) (string, error)
	GenerateRefreshToken(accountID, accountType string, permissions []string, duration time.Duration) (string, error)
	VerifyToken(tokenString string) (*Claims, error)
}

type jwtManager struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	issuer     string
	audience   string
}

func NewjwtManager(privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey) *jwtManager {
	return &jwtManager{
		issuer:     os.Getenv("JWT_ISSUER"),
		audience:   os.Getenv("JWT_AUDIENCE"),
		privateKey: privateKey,
		publicKey:  publicKey,
	}
}

// GenerateToken creates a new access token
func (j *jwtManager) GenerateAccessToken(employeeID, accountID, mobileNo, accountType, name string, permissions []string, duration time.Duration) (string, error) {
	if employeeID == "" || accountID == "" || mobileNo == "" || accountType == "" || name == "" {
		logger.Error(constants.TokenParamMissing, nil, map[string]interface{}{
			"method": constants.Methods.GenerateAccToken,
		})
		return "", fmt.Errorf(constants.TokenParamMissing)
	}
	claims := Claims{
		EmployeeID:  employeeID,
		AccountID:   accountID,
		AccountType: accountType,
		Name:        name,
		MobileNo:    mobileNo,
		Permissions: permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.issuer,
			Audience:  jwt.ClaimStrings{j.audience},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}

	tokenChan := make(chan string)
	errChan := make(chan error)

	go func() {
		token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
		signedToken, err := token.SignedString(j.privateKey)
		if err != nil {
			errChan <- err
			return
		}
		tokenChan <- signedToken
	}()

	select {
	case signedToken := <-tokenChan:
		return signedToken, nil
	case err := <-errChan:
		return "", err
	}
}

// GenerateRefreshToken creates a new refresh token
func (j *jwtManager) GenerateRefreshToken(accountID, accountType string, permissions []string, duration time.Duration) (string, error) {
	if accountID == "" {
		logger.Error(constants.TokenParamMissing, nil, map[string]interface{}{
			"method": constants.Methods.GenerateRefreshToken,
		})
		return "", fmt.Errorf(constants.TokenParamMissing)
	}

	claims := Claims{
		AccountID:   accountID,
		AccountType: accountType,
		Permissions: permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.issuer,
			Audience:  jwt.ClaimStrings{j.audience},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}
	tokenChan := make(chan string)
	errChan := make(chan error)

	go func() {
		token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
		signedToken, err := token.SignedString(j.privateKey)
		if err != nil {
			errChan <- err
			return
		}
		tokenChan <- signedToken
	}()

	select {
	case signedToken := <-tokenChan:
		return signedToken, nil
	case err := <-errChan:
		return "", err
	}
}

// ValidateToken verifies the signature and expiration of an access token

func (j *jwtManager) VerifyToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New(constants.ErrUnexpectedSigningMethod)
		}
		return j.publicKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New(constants.ErrInvalidToken)
	}

	return claims, nil
}
