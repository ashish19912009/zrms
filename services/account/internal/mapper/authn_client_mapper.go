package mapper

import (
	"errors"

	"github.com/ashish19912009/zrms/services/account/internal/model"
	authn_pb "github.com/ashish19912009/zrms/services/authN/pb"
)

func VerifyTokenFromModelToPb(token model.Token) *authn_pb.VerifyTokenRequest {
	return &authn_pb.VerifyTokenRequest{
		AccessToken: token.Token,
	}
}

func VerifyTokenFromPbToModel(authClaim *authn_pb.AuthClaims) (*model.AuthClaims, error) {

	if authClaim != nil {
		var regClaims *model.RegisteredClaims
		if authClaim.RegisteredClaims != nil {
			regClaims = &model.RegisteredClaims{
				ID:        authClaim.RegisteredClaims.Id,
				Subject:   authClaim.RegisteredClaims.Subject,
				Issuer:    authClaim.RegisteredClaims.Issuer,
				Audience:  authClaim.RegisteredClaims.Audience,
				IssuedAt:  authClaim.RegisteredClaims.IssuedAt,
				ExpiresAt: authClaim.RegisteredClaims.ExpiresAt,
			}
		}
		return &model.AuthClaims{
			EmployeeID:       authClaim.EmployeeId,
			FranchiseID:      authClaim.FranchiseId,
			AccountType:      authClaim.AccountType,
			Name:             authClaim.Name,
			MobileNo:         authClaim.MobileNo,
			RegisteredClaims: regClaims,
		}, nil
	}
	return nil, errors.New("Empty Auth Claims")
}
