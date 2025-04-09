package mapper

import (
	"github.com/ashish19912009/services/auth/internal/models"
	"github.com/ashish19912009/services/auth/pb"
)

func LoginResponse(user *models.User, accessToken, refreshToken string) *pb.LoginResponse {
	return &pb.LoginResponse{
		AccountId:    user.AccountID,
		EmployeeId:   user.EmployeeID,
		AccountType:  user.AccountType,
		Name:         user.Name,
		MobileNo:     user.MobileNo,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}

func VerifyTokenResponse(user *models.User, is_valid bool) *pb.VerifyTokenResponse {
	return &pb.VerifyTokenResponse{
		AccountId:   user.AccountID,
		AccountType: user.AccountType,
		Permissions: user.Permissions,
		Role:        user.Role,
		IsValid:     is_valid,
	}
}

func LogoutResponse(success bool) *pb.LogoutResponse {
	return &pb.LogoutResponse{
		Success: success,
	}
}
