package mapper

import (
	"github.com/ashish19912009/zrms/services/authN/internal/model"
	"github.com/ashish19912009/zrms/services/authN/pb"
)

func LoginResponse(user *model.User, accessToken, refreshToken string) *pb.LoginResponse {
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

func VerifyTokenResponse(user *model.User, is_valid bool) *pb.VerifyTokenResponse {
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
