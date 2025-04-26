package service

import (
	"context"

	pb "github.com/ashish19912009/zrms/services/authZ/api"
)

type AuthZService interface {
	IsAuthorized(ctx context.Context, req *pb.AuthorizationRequest) (*pb.AuthorizationResponse, error)
}

type authZService struct{}

func NewAuthZService() AuthZService {
	return &authZService{}
}

func (s *authZService) IsAuthorized(ctx context.Context, req *pb.AuthorizationRequest) (*pb.AuthorizationResponse, error) {
	// 🔥 Dummy response for now
	return &pb.AuthorizationResponse{
		Allowed: true,
		Reason:  "Always allowed (stub)",
	}, nil
}
