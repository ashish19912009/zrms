package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/ashish19912009/services/auth/internal/constants"
	"github.com/ashish19912009/services/auth/internal/logger"
	"github.com/ashish19912009/services/auth/internal/repository"
	"github.com/ashish19912009/services/auth/internal/token"
	"github.com/ashish19912009/services/auth/pb"
)

type AuthService interface {
	Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error)
	RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.LoginResponse, error)
	ValidateToken(ctx context.Context, req *pb.VerifyTokenRequest) (*pb.VerifyTokenResponse, error)
	Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error)
}

type authService struct {
	tokenManager token.TokenManager
	tokenRepo    repository.TokenRepository
	userRepo     repository.UserRepository
}

func NewAuthSerice(tokenManager token.TokenManager, tokenRepo repository.TokenRepository, userRepo repository.UserRepository) AuthService {
	return &authService{
		tokenManager: tokenManager,
		tokenRepo:    tokenRepo,
		userRepo:     userRepo,
	}
}

func getTokenTimer() (time.Duration, time.Duration) {
	var defaultAccessTokenTimer time.Duration = 24 * time.Hour      // 24Hours
	var defaultRefreshTokenTimer time.Duration = 24 * 7 * time.Hour // 24Hours
	accessTokenTimer, err := time.ParseDuration(os.Getenv(constants.EnvVariable.ACCESS_TOKEN_TIME))
	if err != nil {
		fmt.Printf("Can't convert env access token timer: %v\n", err)
		accessTokenTimer = defaultAccessTokenTimer * time.Hour
	}
	if accessTokenTimer == 0 {
		accessTokenTimer = defaultAccessTokenTimer * time.Hour
	}
	refreshTokenTimer, err := time.ParseDuration(os.Getenv(constants.EnvVariable.REFRESH_TOKEN_TIME))
	if err != nil {
		fmt.Printf("Can't convert env refresh token timer: %v\n", err)
		refreshTokenTimer = defaultRefreshTokenTimer * time.Hour
	}
	if refreshTokenTimer == 0 {
		refreshTokenTimer = defaultRefreshTokenTimer * time.Hour
	}
	return accessTokenTimer, refreshTokenTimer
}

func (s *authService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	if req.LoginId == "" || req.Password == "" || req.AccountType == "" {
		logger.Error(constants.ValidationMissingCredentials, nil, map[string]interface{}{
			"method": constants.Methods.Login,
		})
		return nil, errors.New(constants.ValidationMissingCredentials)
	}

	// Getting user information from database
	userDetails, err := s.userRepo.GetUser(ctx, req.LoginId, req.AccountType)
	if err != nil {
		return nil, errors.New(constants.WrongUsernamePassword)
	}

	if userDetails.Password != "" {
		isCorrectPassword := s.userRepo.VerifyPassword(userDetails.Password, req.Password)
		if isCorrectPassword {
			accessTokenTime, refreshTokenTime := getTokenTimer()
			accessToken, err := s.tokenManager.GenerateAccessToken(userDetails.EmployeeID, userDetails.AccountID, userDetails.MobileNo, userDetails.AccountType, userDetails.Name, userDetails.Permissions, accessTokenTime)
			if err != nil {
				logger.Error(constants.FailedToGenerateAct, err, map[string]interface{}{
					"method": constants.Methods.Login,
					"step":   constants.GenerateAccessToken,
				})
				return nil, fmt.Errorf(constants.FailedToGenerateAct, err)
			}

			refreshToken, err := s.tokenManager.GenerateRefreshToken(userDetails.AccountID, userDetails.AccountType, userDetails.Permissions, refreshTokenTime)
			if err != nil {
				logger.Error(constants.FailedToGenerateRsh, err, map[string]interface{}{
					"method": constants.Methods.Login,
					"step":   constants.GenerateRefreshToken,
				})
				return nil, fmt.Errorf(constants.FailedToGenerateRsh, err)
			}

			// Store refresh token and access token in "in memory DB"

			err = s.tokenRepo.StoreToken(ctx, constants.Access_token, userDetails.AccountID, accessToken, accessTokenTime)
			if err != nil {
				logger.Error(constants.FailedToStoreRshToken, err, map[string]interface{}{
					"method": constants.Methods.Login,
				})
				return nil, fmt.Errorf("%s: %w", constants.FailedToStoreRshToken, err)
			}
			err = s.tokenRepo.StoreToken(ctx, constants.Refresh_token, userDetails.AccountID, refreshToken, refreshTokenTime)
			if err != nil {
				logger.Error(constants.FailedToStoreRshToken, err, map[string]interface{}{
					"method": constants.Methods.Login,
				})
				return nil, fmt.Errorf("%s: %w", constants.FailedToStoreRshToken, err)
			}

			logger.Info(constants.SuccessfulLogin, map[string]interface{}{
				"method":  constants.Methods.Login,
				"user_id": userDetails.AccountID + "_" + userDetails.EmployeeID,
			})

			return &pb.LoginResponse{
				AccountId:    userDetails.AccountID,
				EmployeeId:   userDetails.EmployeeID,
				AccountType:  userDetails.AccountType,
				Name:         userDetails.Name,
				MobileNo:     userDetails.MobileNo,
				Permissions:  userDetails.Permissions,
				AccessToken:  accessToken,
				RefreshToken: refreshToken,
			}, nil
		}
		return nil, fmt.Errorf(constants.WrongUsernamePassword)
	}
	logger.Fatal(constants.PasswordMissingFromServer, nil, map[string]interface{}{
		"method":  constants.Methods.Login,
		"user_id": userDetails.AccountID + "_" + userDetails.EmployeeID,
	})
	return nil, fmt.Errorf(constants.WrongUsernamePassword)
}

func (s *authService) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.LoginResponse, error) {
	if req.RefreshToken == "" {
		logger.Error(constants.AuthRefreshRequired, nil, map[string]interface{}{
			"method": constants.Methods.RefreshToken,
		})
		return nil, fmt.Errorf(constants.AuthRefreshRequired)
	}

	claims, err := s.tokenManager.VerifyToken(req.RefreshToken)
	if err != nil {
		logger.Error(constants.AuthRshTokenInvalid, err, map[string]interface{}{
			"method": constants.Methods.RefreshToken,
		})
		return nil, fmt.Errorf(constants.AuthRshTokenInvalid, err)
	}

	// Check if refresh token exists in in_memory
	exists, err := s.tokenRepo.CheckToken(ctx, constants.Refresh_token, req.AccountId, req.RefreshToken)
	if err != nil || !exists {
		logger.Error(constants.AuthRshTokenInvalid, err, map[string]interface{}{
			"method": constants.Methods.RefreshToken,
			"check":  constants.RefreshTokenExistence,
		})
		return nil, fmt.Errorf("%s: %w", constants.AuthRshTokenInvalid, err)
	}
	if exists {
		userInfoFromToken, err := s.tokenManager.VerifyToken(req.RefreshToken)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", constants.AuthRefreshFailure, err)
		}
		// Getting user information from database
		userDetails, err := s.userRepo.GetUser(ctx, userInfoFromToken.AccountID, userInfoFromToken.AccountType)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", constants.WrongUsernamePassword, err)
		}
		accessTokenTime, refreshTokenTime := getTokenTimer()
		accessToken, err := s.tokenManager.GenerateAccessToken(userDetails.AccountID, userDetails.EmployeeID, userDetails.MobileNo, userDetails.AccountType, userDetails.Name, userDetails.Permissions, accessTokenTime)
		if err != nil {
			logger.Error(constants.FailedToGenerateAct, err, map[string]interface{}{
				"method": constants.Methods.RefreshToken,
				"step":   constants.GenerateAccessToken,
			})
			return nil, fmt.Errorf(constants.FailedToGenerateAct, err)
		}

		newRefreshToken, err := s.tokenManager.GenerateRefreshToken(userDetails.AccountID, userDetails.AccountType, userDetails.Permissions, refreshTokenTime)
		if err != nil {
			logger.Error(constants.FailedToGenerateRsh, err, map[string]interface{}{
				"method": constants.Methods.RefreshToken,
			})
			return nil, fmt.Errorf(constants.FailedToGenerateRsh, err)
		}

		// Update refresh token in Redis
		// make a abstraction layer to store token in any in memory DB i.e. Redis, memchached, Dragonfly
		err = s.tokenRepo.StoreToken(ctx, constants.Access_token, userDetails.AccountID, accessToken, accessTokenTime)
		if err != nil {
			logger.Error(constants.FailedToStoreRshToken, err, map[string]interface{}{
				"method": constants.Methods.RefreshToken,
			})
			return nil, err
		}
		err = s.tokenRepo.StoreToken(ctx, constants.Refresh_token, userDetails.AccountID, newRefreshToken, refreshTokenTime)
		if err != nil {
			logger.Error(constants.FailedToStoreRshToken, err, map[string]interface{}{
				"method": constants.Methods.RefreshToken,
			})
			return nil, err
		}

		logger.Info(constants.SuccessfulRefreshToken, map[string]interface{}{
			"method": constants.Methods.RefreshToken,
			"user":   claims.EmployeeID,
		})

		return &pb.LoginResponse{
			AccessToken:  accessToken,
			RefreshToken: newRefreshToken,
		}, nil
	}
	return nil, errors.New("couldn't refresh token")
}

func (s *authService) ValidateToken(ctx context.Context, req *pb.VerifyTokenRequest) (*pb.VerifyTokenResponse, error) {
	if req.AccessToken == "" {
		logger.Error(constants.AuthAccessRequired, nil, map[string]interface{}{
			"method": constants.Methods.ValidateToken,
		})
		return nil, fmt.Errorf(constants.AuthAccessRequired)
	}

	claims, err := s.tokenManager.VerifyToken(req.AccessToken)
	if err != nil {
		logger.Error(constants.AuthTokenVeriFailed, err, map[string]interface{}{
			"method": constants.Methods.ValidateToken,
		})
		return nil, fmt.Errorf(constants.AuthTokenVeriFailed)
	}

	logger.Info(constants.SuccessfulTokenValidation, map[string]interface{}{
		"method": constants.Methods.ValidateToken,
		"user":   claims.EmployeeID,
	})

	return &pb.VerifyTokenResponse{
		EmployeeId:  claims.EmployeeID,
		MobileNo:    claims.MobileNo,
		Permissions: claims.Permissions,
		IsValid:     true,
	}, nil
}

func (s *authService) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	if req.RefreshToken == "" {
		logger.Error(constants.ErrMissingRefreshToken, nil, map[string]interface{}{
			"method": constants.Methods.Logout,
		})
		return nil, fmt.Errorf(constants.ErrMissingRefreshToken)
	}

	claims, err := s.tokenManager.VerifyToken(req.RefreshToken)
	if err != nil {
		logger.Error(constants.AuthRshTokenInvalid, err, map[string]interface{}{
			"method": constants.Methods.Logout,
		})
		return nil, fmt.Errorf(constants.AuthRshTokenInvalid, err)
	}

	// Remove refresh token from Redis
	err = s.tokenRepo.DeleteToken(ctx, constants.Access_token, claims.AccountID)
	if err != nil {
		logger.Error(constants.ErrTokenDeletionFailed, err, map[string]interface{}{
			"method": constants.Methods.Logout,
		})
		return nil, err
	}

	err = s.tokenRepo.DeleteToken(ctx, constants.Refresh_token, claims.AccountID)
	if err != nil {
		logger.Error(constants.ErrTokenDeletionFailed, err, map[string]interface{}{
			"method": constants.Methods.Logout,
		})
		return nil, err
	}

	logger.Info(constants.MsgLogoutSuccess, map[string]interface{}{
		"method": constants.Methods.Logout,
		"user":   claims.EmployeeID,
	})

	return &pb.LogoutResponse{Success: true}, nil
}
