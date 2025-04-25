package handler

import (
	"context"

	"github.com/ashish19912009/services/auth/internal/constants"
	"github.com/ashish19912009/services/auth/internal/logger"
	"github.com/ashish19912009/services/auth/internal/service"
	"github.com/ashish19912009/services/auth/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCHandler struct {
	pb.UnimplementedAuthServiceServer
	authService service.AuthService
}

func NewGRPCHandler(auth service.AuthService) (*GRPCHandler, error) {
	return &GRPCHandler{
		authService: auth,
	}, nil
}

func (h *GRPCHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	logger.Info(constants.LoginAttempt, map[string]interface{}{
		"loginID": req.LoginId,
	})
	if req.LoginId == "" || req.Password == "" || req.AccountType == "" {
		logger.Error(constants.ValidationMissingCredentials, nil, map[string]interface{}{
			"method":  constants.Methods.Login,
			"loginID": req.LoginId,
		})
		return nil, status.Error(codes.InvalidArgument, constants.ValidationMissingCredentials)
	}
	request := &pb.LoginRequest{
		LoginId:     req.LoginId,
		Password:    req.Password,
		AccountType: req.AccountType,
	}
	return h.authService.Login(ctx, request)
}

func (h *GRPCHandler) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.LoginResponse, error) {
	logger.Info(constants.AttemptRefreshToken, map[string]interface{}{
		"refresh_token": req.RefreshToken,
	})
	if req.RefreshToken == "" {
		logger.Error(constants.AuthRefreshRequired, nil, map[string]interface{}{
			"method":        constants.Methods.RefreshToken,
			"refresh_token": req.RefreshToken,
		})
		return nil, status.Error(codes.InvalidArgument, constants.AuthRefreshRequired)
	}
	return h.authService.RefreshToken(ctx, req)
}

func (h *GRPCHandler) ValidateToken(ctx context.Context, req *pb.VerifyTokenRequest) (*pb.VerifyTokenResponse, error) {
	if req.AccessToken == "" {
		logger.Error(constants.AuthAccessRequired, nil, map[string]interface{}{
			"method":        constants.Methods.AccessToken,
			"refresh_token": req.AccessToken,
		})
		return nil, status.Error(codes.InvalidArgument, constants.AuthAccessRequired)
	}
	return h.authService.VerifyToken(ctx, req)
}

func (h *GRPCHandler) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	if req.RefreshToken == "" {
		logger.Error(constants.ErrMissingRefreshToken, nil, map[string]interface{}{
			"method": constants.Methods.Logout,
		})
	}

	_, err := h.authService.Logout(ctx, req)
	if err != nil {
		logger.Error(constants.ErrTokenDeletionFailed, err, map[string]interface{}{
			"method":        constants.Methods.Logout,
			"refresh_token": req.RefreshToken,
		})
		return nil, err
	}

	logger.Info(constants.MsgLogoutSuccess, map[string]interface{}{
		"method":        constants.Methods.Logout,
		"refresh_token": req.RefreshToken,
	})

	return &pb.LogoutResponse{
		Success: true,
	}, nil
}
