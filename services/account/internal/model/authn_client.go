package model

import (
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Token struct {
	Token string `json:"access_token"`
}

type AuthCalims struct {
	EmployeeID       string            `json:"employee_id" protobuf:"bytes,1,opt,name=employee_id"`
	FranchiseID      string            `json:"franchise_id" protobuf:"bytes,2,opt,name=franchise_id"`
	AccountType      string            `json:"account_type" protobuf:"bytes,3,opt,name=account_type"`
	Name             string            `json:"name" protobuf:"bytes,4,opt,name=name"`
	MobileNo         string            `json:"mobile_no" protobuf:"bytes,5,opt,name=mobile_no"`
	RegisteredClaims *RegisteredClaims `json:"registered_claims" protobuf:"bytes,6,opt,name=registered_claims"`
}

type RegisteredClaims struct {
	ID        string                 `json:"jti,omitempty" protobuf:"bytes,1,opt,name=id"`         // JWT ID (UUID)
	Subject   string                 `json:"sub,omitempty" protobuf:"bytes,2,opt,name=subject"`    // User/Account ID
	Issuer    string                 `json:"iss,omitempty" protobuf:"bytes,3,opt,name=issuer"`     // Token issuer
	Audience  jwt.ClaimStrings       `json:"aud,omitempty" protobuf:"bytes,4,rep,name=audience"`   // Intended audience
	IssuedAt  *timestamppb.Timestamp `json:"iat,omitempty" protobuf:"bytes,5,opt,name=issued_at"`  // Issued at (Proto format)
	ExpiresAt *timestamppb.Timestamp `json:"exp,omitempty" protobuf:"bytes,6,opt,name=expires_at"` // Expiration time (Proto format)
}
