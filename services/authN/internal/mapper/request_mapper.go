package mapper

import (
	"github.com/ashish19912009/services/auth/internal/models"
	"github.com/ashish19912009/services/auth/pb"
)

func LoginRequest(req *pb.LoginRequest) *models.LoginInput {
	return &models.LoginInput{
		LoginID:     req.LoginId,
		AccountType: req.AccountType,
		Password:    req.Password,
	}
}

func RefreshTokenRequest(req *pb.RefreshTokenRequest) *models.RefreshTokenInput {
	return &models.RefreshTokenInput{
		AccountID:    req.AccountId,
		RefreshToken: req.RefreshToken,
	}
}

func LogoutRequestInput(req *pb.LogoutRequest) *models.LogoutInput {
	return &models.LogoutInput{
		RefreshToken: req.RefreshToken,
	}
}

func VerifyTokenRequest(req *pb.VerifyTokenRequest) *models.ValidateTokenInput {
	return &models.ValidateTokenInput{
		AccessToken: req.AccessToken,
	}
}
