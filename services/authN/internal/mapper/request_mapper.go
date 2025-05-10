package mapper

import (
	"github.com/ashish19912009/zrms/services/authN/internal/model"
	"github.com/ashish19912009/zrms/services/authN/pb"
)

func LoginRequest(req *pb.LoginRequest) *model.LoginInput {
	return &model.LoginInput{
		LoginID:     req.LoginId,
		AccountType: req.AccountType,
		Password:    req.Password,
	}
}

func RefreshTokenRequest(req *pb.RefreshTokenRequest) *model.RefreshTokenInput {
	return &model.RefreshTokenInput{
		AccountID:    req.AccountId,
		RefreshToken: req.RefreshToken,
	}
}

func LogoutRequestInput(req *pb.LogoutRequest) *model.LogoutInput {
	return &model.LogoutInput{
		RefreshToken: req.RefreshToken,
	}
}

func VerifyTokenRequest(req *pb.VerifyTokenRequest) *model.ValidateTokenInput {
	return &model.ValidateTokenInput{
		AccessToken: req.AccessToken,
	}
}
