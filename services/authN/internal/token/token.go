package token

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/ashish19912009/zrms/services/authN/internal/constants"
	"github.com/ashish19912009/zrms/services/authN/internal/logger"
	"github.com/ashish19912009/zrms/services/authN/internal/model"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenManager interface {
	GenerateAccessToken(accountID, employeeID, mobileNo, accountType, name string, permissions, duration time.Duration) (string, error)
	GenerateRefreshToken(accountID, accountType string, permissions, duration time.Duration) (string, error)
	VerifyToken(tokenString string) (*model.AuthClaims, error)
}

type jwtManager struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	issuer     string
	audience   string
}

func NewjwtManager(privateKeyPath, publicKeyPath string) (*jwtManager, error) {
	privateKey, err := loadPrivateKeyFromFile(privateKeyPath)
	if err != nil {
		return nil, err
	}

	publicKey, err := loadPublicKeyFromFile(publicKeyPath)
	if err != nil {
		return nil, err
	}
	return &jwtManager{
		issuer:     os.Getenv("JWT_ISSUER"),
		audience:   os.Getenv("JWT_AUDIENCE"),
		privateKey: privateKey,
		publicKey:  publicKey,
	}, nil
}

func loadPrivateKeyFromFile(path string) (*rsa.PrivateKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key file: %w", err)
	}
	block, _ := pem.Decode(data)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("invalid or missing PEM block for RSA private key")
	}
	privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSA private key: %w", err)
	}
	return privKey, nil
}

func loadPublicKeyFromFile(path string) (*rsa.PublicKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(data)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.New("invalid public key data")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not an RSA public key")
	}
	return rsaPub, nil
}

// GenerateToken creates a new access token
func (j *jwtManager) GenerateAccessToken(employeeID, accountID, mobileNo, accountType, name string, duration time.Duration) (string, error) {
	if employeeID == "" || accountID == "" || mobileNo == "" || accountType == "" || name == "" {
		logger.Error(constants.TokenParamMissing, nil, map[string]interface{}{
			"method": constants.Methods.GenerateAccToken,
		})
		return "", fmt.Errorf(constants.TokenParamMissing)
	}
	claims := model.AuthClaims{
		EmployeeID:  employeeID,
		AccountType: accountType,
		Name:        name,
		MobileNo:    mobileNo,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			Subject:   accountID,
			Issuer:    j.issuer,
			Audience:  jwt.ClaimStrings{j.audience},
			IssuedAt:  jwt.NewNumericDate(time.Now()),
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
func (j *jwtManager) GenerateRefreshToken(accountID, accountType string, duration time.Duration) (string, error) {
	if accountID == "" {
		logger.Error(constants.TokenParamMissing, nil, map[string]interface{}{
			"method": constants.Methods.GenerateRefreshToken,
		})
		return "", fmt.Errorf(constants.TokenParamMissing)
	}

	claims := model.AuthClaims{
		AccountID:   accountID,
		AccountType: accountType,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.issuer,
			Audience:  jwt.ClaimStrings{j.audience},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
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

func (j *jwtManager) VerifyToken(tokenString string) (*model.AuthClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &model.AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New(constants.ErrUnexpectedSigningMethod)
		}
		return j.publicKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*model.AuthClaims)
	if !ok || !token.Valid {
		return nil, errors.New(constants.ErrInvalidToken)
	}

	return claims, nil
}
