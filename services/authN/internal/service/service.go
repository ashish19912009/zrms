package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/ashish19912009/zrms/services/authN/internal/client"
	"github.com/ashish19912009/zrms/services/authN/internal/constants"
	"github.com/ashish19912009/zrms/services/authN/internal/logger"
	"github.com/ashish19912009/zrms/services/authN/internal/mapper"
	"github.com/ashish19912009/zrms/services/authN/internal/model"
	"github.com/ashish19912009/zrms/services/authN/internal/repository"
	"github.com/ashish19912009/zrms/services/authN/internal/token"
	"github.com/ashish19912009/zrms/services/authN/pb"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthService interface {
	Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error)
	RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.LoginResponse, error)
	VerifyAccessToken(ctx context.Context, req *pb.VerifyTokenRequest) (*pb.AuthClaims, error)
	Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error)
}

type authService struct {
	tokenManager token.TokenManager
	tokenRepo    repository.TokenRepository
	userRepo     repository.UserRepository
	accessTTL    time.Duration
	refreshTTL   time.Duration
	client       client.AuthZClient
}

func NewAuthService(tokenManager token.TokenManager, tokenRepo repository.TokenRepository, userRepo repository.UserRepository) AuthService {
	accessTokenTime, refreshTokenTime := getTokenTimer(constants.EnvVariable.ACCESS_TOKEN_TTL, constants.EnvVariable.REFRESH_TOKEN_TTL)
	return &authService{
		tokenManager: tokenManager,
		tokenRepo:    tokenRepo,
		userRepo:     userRepo,
		accessTTL:    accessTokenTime,
		refreshTTL:   refreshTokenTime,
	}
}

// Constructor for testing and flexibility
func NewAuthServiceWithTTL(
	tm token.TokenManager,
	tr repository.TokenRepository,
	ur repository.UserRepository,
	accessTTL, refreshTTL time.Duration,
) AuthService {
	return &authService{
		tokenManager: tm,
		tokenRepo:    tr,
		userRepo:     ur,
		accessTTL:    accessTTL,
		refreshTTL:   refreshTTL,
	}
}

func verifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		logger.Warn(constants.ErrInvalidPassword, map[string]interface{}{
			"method": constants.Methods.VerifyPassword,
		})
		return false
	}
	return true
}

func getTokenTimer(ACCESS_TOKEN_TTL string, REFRESH_TOKEN_TTL string) (time.Duration, time.Duration) {
	var defaultAccessTokenTimer time.Duration = 24 * time.Hour      // 24Hours
	var defaultRefreshTokenTimer time.Duration = 24 * 7 * time.Hour // 7 Days
	accessTokenTimer, err := time.ParseDuration(os.Getenv(ACCESS_TOKEN_TTL))
	if err != nil {
		fmt.Printf("Can't convert env access token timer: %v\n", err)
		accessTokenTimer = defaultAccessTokenTimer * time.Hour
	}
	if accessTokenTimer == 0 {
		accessTokenTimer = defaultAccessTokenTimer * time.Hour
	}
	refreshTokenTimer, err := time.ParseDuration(os.Getenv(REFRESH_TOKEN_TTL))
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
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	input := mapper.LoginRequest(req)
	if input.LoginID == "" || input.Password == "" || input.AccountType == "" {
		logger.Error(constants.ValidationMissingCredentials, nil, map[string]interface{}{
			"method": constants.Methods.Login,
		})
		return nil, errors.New(constants.ValidationMissingCredentials)
	}

	// Getting user information from database
	var userDetails *model.User
	userDetails, err := s.userRepo.GetUser(ctx, input.LoginID, input.AccountType)
	if err != nil {
		return nil, errors.New(constants.WrongUsernamePassword)
	}
	if userDetails == nil {
		return nil, status.Error(codes.Internal, constants.UserDataMissing)
	}
	if userDetails.Password != "" {
		isCorrectPassword := verifyPassword(userDetails.Password, input.Password)
		if isCorrectPassword {
			//allResources, err := s.userRepo.GetFranchiseRolePermissions(ctx, userDetails.FranchiseID)
			//getBatchPermission, err := s.client.BatchCheckAccess(ctx, userDetails.AccountID, userDetails.FranchiseID, allResources)
			accessToken, err := s.tokenManager.GenerateAccessToken(userDetails.EmployeeID, userDetails.FranchiseID, userDetails.AccountID, userDetails.MobileNo, userDetails.AccountType, userDetails.Name, s.accessTTL)
			if err != nil {
				logger.Error(constants.FailedToGenerateAct, err, map[string]interface{}{
					"method": constants.Methods.Login,
					"step":   constants.GenerateAccessToken,
				})
				return nil, fmt.Errorf(constants.FailedToGenerateAct, err)
			}

			refreshToken, err := s.tokenManager.GenerateRefreshToken(userDetails.AccountID, userDetails.AccountType, s.refreshTTL)
			if err != nil {
				logger.Error(constants.FailedToGenerateRsh, err, map[string]interface{}{
					"method": constants.Methods.Login,
					"step":   constants.GenerateRefreshToken,
				})
				return nil, fmt.Errorf(constants.FailedToGenerateRsh, err)
			}

			// Store refresh token and access token in "in memory DB"

			err = s.tokenRepo.StoreToken(ctx, constants.Access_token, userDetails.AccountID, accessToken, s.accessTTL)
			if err != nil {
				logger.Error(constants.FailedToStoreRshToken, err, map[string]interface{}{
					"method": constants.Methods.Login,
				})
				return nil, fmt.Errorf("%s: %w", constants.FailedToStoreRshToken, err)
			}
			err = s.tokenRepo.StoreToken(ctx, constants.Refresh_token, userDetails.AccountID, refreshToken, s.refreshTTL)
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

			return mapper.LoginResponse(userDetails, accessToken, refreshToken), nil
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
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	var input = mapper.RefreshTokenRequest(req)
	if input.RefreshToken == "" {
		logger.Error(constants.AuthRefreshRequired, nil, map[string]interface{}{
			"method": constants.Methods.RefreshToken,
		})
		return nil, fmt.Errorf(constants.AuthRefreshRequired)
	}

	claims, err := s.tokenManager.VerifyRefreshToken(input.RefreshToken)
	if err != nil {
		logger.Error(constants.AuthRshTokenInvalid, err, map[string]interface{}{
			"method": constants.Methods.RefreshToken,
		})
		return nil, fmt.Errorf(constants.AuthRshTokenInvalid, err)
	}

	// Check if refresh token exists in in_memory
	exists, err := s.tokenRepo.CheckToken(ctx, constants.Refresh_token, claims.RegisteredClaims.Subject, input.RefreshToken)
	if err != nil || !exists {
		logger.Error(constants.AuthRshTokenInvalid, err, map[string]interface{}{
			"method": constants.Methods.RefreshToken,
			"check":  constants.RefreshTokenExistence,
		})
		return nil, fmt.Errorf("%s: %w", constants.AuthRshTokenInvalid, err)
	}
	if exists {
		// Getting user information from database
		var userDetails *model.User
		userDetails, err = s.userRepo.GetUser(ctx, claims.RegisteredClaims.Subject, claims.AccountType)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", constants.WrongUsernamePassword, err)
		}
		if userDetails == nil {
			return nil, status.Error(codes.Internal, constants.UserDataMissing)
		}
		accessToken, err := s.tokenManager.GenerateAccessToken(userDetails.AccountID, userDetails.FranchiseID, userDetails.EmployeeID, userDetails.MobileNo, userDetails.AccountType, userDetails.Name, s.accessTTL)
		if err != nil {
			logger.Error(constants.FailedToGenerateAct, err, map[string]interface{}{
				"method": constants.Methods.RefreshToken,
				"step":   constants.GenerateAccessToken,
			})
			return nil, fmt.Errorf(constants.FailedToGenerateAct, err)
		}

		newRefreshToken, err := s.tokenManager.GenerateRefreshToken(userDetails.AccountID, userDetails.AccountType, s.refreshTTL)
		if err != nil {
			logger.Error(constants.FailedToGenerateRsh, err, map[string]interface{}{
				"method": constants.Methods.RefreshToken,
			})
			return nil, fmt.Errorf(constants.FailedToGenerateRsh, err)
		}

		// Update refresh token in Redis
		// make a abstraction layer to store token in any in memory DB i.e. Redis, memchached, Dragonfly
		err = s.tokenRepo.StoreToken(ctx, constants.Access_token, userDetails.AccountID, accessToken, s.accessTTL)
		if err != nil {
			logger.Error(constants.FailedToStoreRshToken, err, map[string]interface{}{
				"method": constants.Methods.RefreshToken,
			})
			return nil, err
		}
		err = s.tokenRepo.StoreToken(ctx, constants.Refresh_token, userDetails.AccountID, newRefreshToken, s.refreshTTL)
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

		return mapper.LoginResponse(userDetails, accessToken, newRefreshToken), nil
	}
	return nil, errors.New("couldn't refresh token")
}

func (s *authService) VerifyAccessToken(ctx context.Context, req *pb.VerifyTokenRequest) (*pb.AuthClaims, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	input := mapper.VerifyTokenRequest(req)
	if input.AccessToken == "" {
		logger.Error(constants.AuthAccessRequired, nil, map[string]interface{}{
			"method": constants.Methods.VerifyToken,
		})
		return nil, fmt.Errorf(constants.AuthAccessRequired)
	}

	claims, err := s.tokenManager.VerifyAccessToken(input.AccessToken)
	if err != nil {
		logger.Error(constants.AuthTokenVeriFailed, err, map[string]interface{}{
			"method": constants.Methods.VerifyToken,
		})
		return nil, fmt.Errorf(constants.AuthTokenVeriFailed)
	}

	logger.Info(constants.SuccessfulTokenValidation, map[string]interface{}{
		"method": constants.Methods.VerifyToken,
		"user":   claims.EmployeeID,
	})

	return mapper.VerifyTokenResponse(claims), nil
}

func (s *authService) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	input := mapper.LogoutRequestInput(req)
	if input.RefreshToken == "" {
		logger.Error(constants.ErrMissingRefreshToken, nil, map[string]interface{}{
			"method": constants.Methods.Logout,
		})
		return nil, fmt.Errorf(constants.ErrMissingRefreshToken)
	}

	claims, err := s.tokenManager.VerifyAccessToken(input.RefreshToken)
	if err != nil {
		logger.Error(constants.AuthRshTokenInvalid, err, map[string]interface{}{
			"method": constants.Methods.Logout,
		})
		return nil, fmt.Errorf(constants.AuthRshTokenInvalid, err)
	}

	// Remove refresh token from Redis
	err = s.tokenRepo.DeleteToken(ctx, constants.Access_token, claims.RegisteredClaims.Subject)
	if err != nil {
		logger.Error(constants.ErrTokenDeletionFailed, err, map[string]interface{}{
			"method": constants.Methods.Logout,
		})
		return nil, err
	}

	err = s.tokenRepo.DeleteToken(ctx, constants.Refresh_token, claims.RegisteredClaims.Subject)
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

	return mapper.LogoutResponse(true), nil
}
