package server

import (
	"context"
	"errors"
	"time"

	"github.com/ashish19912009/zrms/services/authZ/internal/constants"
	"github.com/ashish19912009/zrms/services/authZ/internal/model"
	"github.com/ashish19912009/zrms/services/authZ/internal/service"
	"github.com/ashish19912009/zrms/services/authZ/internal/validations"
	"github.com/ashish19912009/zrms/services/authZ/pb"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthZServer struct {
	pb.UnimplementedAuthZServiceServer
	service service.AuthZService
}

func NewAuthZServer(svc service.AuthZService) *AuthZServer {
	return &AuthZServer{
		service: svc,
	}
}

func (s *AuthZServer) Register(grpcServer *grpc.Server) {
	pb.RegisterAuthZServiceServer(grpcServer, s)
}

// CheckAccess implements the CheckAccess RPC
func (s *AuthZServer) CheckAccess(ctx context.Context, req *pb.CheckAccessRequest) (*pb.CheckAccessResponse, error) {
	token := ctx.Value("user").(jwt.Token)
	role, _ := token.Get("account_type")
	if role == "superAdmin" || role == "SuperAdmin" || role == "super_admin" || role == "superadmin" {
		decison := &pb.Decision{
			Allowed:       true,
			Reason:        "Super Admin",
			IssuedAt:      time.Now().Unix(),
			ExpiresAt:     24 * time.Now().Unix(),
			PolicyVersion: "",
		}
		return &pb.CheckAccessResponse{
			Decision: decison,
		}, nil
	}
	aM, err := model.CheckAccessFromPbToModel(req)
	if err != nil {
		return nil, err
	}
	err = validations.ValidateCheckAccess(aM)
	if err != nil {
		return nil, err
	}
	allowed, reason, issued_at, expires_at, policy_version, err := s.service.IsAuthorized(ctx, aM.FranchiseID, aM.AccountID, aM.Resource, aM.Action, aM.Context)
	if err != nil {
		return nil, err
	}
	accessR := &model.CheckAccessResponse{
		Allowed:       allowed,
		Reason:        reason,
		IssuedAt:      issued_at,
		ExpiresAt:     expires_at,
		PolicyVersion: policy_version,
	}
	accessPb := model.CheckAccessFromModelToPb(accessR)
	return accessPb, nil
}

func (s *AuthZServer) BatchCheckAccess(ctx context.Context, req *pb.BatchCheckAccessRequest) (*pb.BatchCheckAccessResponse, error) {
	aM, err := model.BatchCheckAccessFromPbToModel(req)
	if err != nil {
		return nil, err
	}
	if len(aM.Resources) == 0 {
		return nil, errors.New(constants.NoResources)
	}
	for _, r := range aM.Resources {
		access := &model.CheckAccess{
			AccountID:   aM.AccountID,
			FranchiseID: aM.FranchiseID,
			Resource:    r.Resource,
			Action:      r.Action,
		}
		err = validations.ValidateCheckAccess(access)
		if err != nil {
			return nil, err
		}
	}

	batchRespose, err := s.service.IsAuthorizedBatch(ctx, aM.FranchiseID, aM.AccountID, aM.Resources)
	if err != nil {
		return nil, err
	}

	// Ensure response length matches input
	if len(batchRespose) != len(aM.Resources) {
		return nil, status.Error(codes.Internal, "response length mismatch")
	}

	accessPb := model.BatchCheckAccessFromModelToPb(batchRespose)
	return accessPb, nil
}
