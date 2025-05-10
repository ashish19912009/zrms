package mapper

import (
	"github.com/ashish19912009/zrms/services/account/internal/model"
	authn_pb "github.com/ashish19912009/zrms/services/authN/pb"
)

func VerifyTokenFromModelToPb(token model.Token) *authn_pb.VerifyTokenRequest {
	return &authn_pb.VerifyTokenRequest{
		AccessToken: token.Token,
	}
}

func VerifyTokenFromPbToModel(authClaim *authn_pb.VerifyTokenResponse) *model.AuthCalims {
	return &model.AuthCalims{
		EmployeeID: authClaim,
	}
}
