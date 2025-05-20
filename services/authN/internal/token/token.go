package token

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/ashish19912009/zrms/services/authN/internal/config"
	"github.com/ashish19912009/zrms/services/authN/internal/constants"
	"github.com/ashish19912009/zrms/services/authN/internal/jwtkeys"
	"github.com/ashish19912009/zrms/services/authN/internal/logger"
	"github.com/ashish19912009/zrms/services/authN/internal/model"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenManager interface {
	GenerateAccessToken(accountID, FranchiseID, employeeID, mobileNo, accountType, name string, duration time.Duration) (string, error)
	GenerateRefreshToken(accountID, accountType string, duration time.Duration) (string, error)
	VerifyAccessToken(tokenString string) (*model.AuthClaims, error)
	VerifyRefreshToken(tokenString string) (*model.AuthClaims, error)
}

type jwtManager struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	issuer     string
	audience   string
	alg        string
	typ        string
	kid        string
}

func NewjwtManager(privateKeyPath, publicKeyPath string, jwtHeader config.JWTHeaderConfig) (*jwtManager, error) {
	keys, err := jwtkeys.LoadKeys(privateKeyPath, publicKeyPath)
	if err != nil {
		return nil, err
	}
	_ = keys.PublicKey
	return &jwtManager{
		issuer:     os.Getenv("JWT_ISSUER"),
		audience:   os.Getenv("JWT_AUDIENCE"),
		privateKey: keys.PrivateKey,
		publicKey:  keys.PublicKey,
		alg:        jwtHeader.Alg,
		typ:        jwtHeader.Typ,
		kid:        jwtHeader.Kid,
	}, nil
}

func (jm *jwtManager) PublicKey() *rsa.PublicKey {
	return jm.publicKey
}

// GenerateToken creates a new access token
func (j *jwtManager) GenerateAccessToken(employeeID, FranchiseID, accountID, mobileNo, accountType, name string, duration time.Duration) (string, error) {
	if employeeID == "" || accountID == "" || mobileNo == "" || accountType == "" || name == "" {
		logger.Error(constants.TokenParamMissing, nil, map[string]interface{}{
			"method": constants.Methods.GenerateAccToken,
		})
		return "", fmt.Errorf(constants.TokenParamMissing)
	}
	claims := model.AuthClaims{
		EmployeeID:  employeeID,
		FranchiseID: FranchiseID,
		AccountType: accountType,
		Name:        name,
		MobileNo:    mobileNo,
		TokenType:   "access",
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
		token.Header["alg"] = j.alg
		token.Header["typ"] = j.typ
		token.Header["kid"] = j.kid
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
		AccountType: accountType,
		TokenType:   "refresh",
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

// ValidateToken verifies the signature and expiration of an access token

func (j *jwtManager) VerifyAccessToken(tokenString string) (*model.AuthClaims, error) {
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

// VerifyRefreshToken verifies the signature and expiration of a refresh token.
func (j *jwtManager) VerifyRefreshToken(tokenString string) (*model.AuthClaims, error) {
	// Parse the refresh token with the expected signing method
	token, err := jwt.ParseWithClaims(tokenString, &model.AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Check if the token's signing method is ECDSA (for ES256)
		if token.Method != jwt.SigningMethodRS256 {
			return nil, errors.New(constants.ErrUnexpectedSigningMethod)
		}
		return j.publicKey, nil
	})

	if err != nil {
		return nil, err
	}

	// Extract the claims and check if the token is valid
	claims, ok := token.Claims.(*model.AuthClaims)
	if !ok || !token.Valid {
		return nil, errors.New(constants.ErrInvalidToken)
	}

	// Optionally: Add further refresh token checks, e.g., make sure it's not expired or revoked
	if claims.RegisteredClaims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New(constants.ErrTokenExpired)
	}

	return claims, nil
}
