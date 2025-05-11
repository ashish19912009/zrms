package mapper

import (
	"github.com/ashish19912009/zrms/services/authN/internal/model"
	"github.com/ashish19912009/zrms/services/authN/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func LoginResponse(user *model.User, accessToken, refreshToken string) *pb.LoginResponse {
	return &pb.LoginResponse{
		AccountId:    user.AccountID,
		FranchiseId:  user.FranchiseID,
		EmployeeId:   user.EmployeeID,
		AccountType:  user.AccountType,
		Name:         user.Name,
		MobileNo:     user.MobileNo,
		Email:        user.Email,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}

func VerifyTokenResponse(usrClaims *model.AuthClaims) *pb.AuthClaims {
	var rClaims = &pb.RegisteredClaims{
		Id:        usrClaims.RegisteredClaims.ID,
		Subject:   usrClaims.RegisteredClaims.Subject,
		Issuer:    usrClaims.RegisteredClaims.Issuer,
		Audience:  usrClaims.RegisteredClaims.Audience,
		IssuedAt:  timestamppb.New(usrClaims.RegisteredClaims.IssuedAt.Time),
		ExpiresAt: timestamppb.New(usrClaims.RegisteredClaims.ExpiresAt.Time),
	}
	return &pb.AuthClaims{
		EmployeeId:       usrClaims.EmployeeID,
		FranchiseId:      usrClaims.FranchiseID,
		AccountType:      usrClaims.AccountType,
		Name:             usrClaims.Name,
		MobileNo:         usrClaims.MobileNo,
		RegisteredClaims: rClaims,
	}
}

func LogoutResponse(success bool) *pb.LogoutResponse {
	return &pb.LogoutResponse{
		Success: success,
	}
}
