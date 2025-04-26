package server

import (
	pb "github.com/ashish19912009/zrms/services/authZ/api" // path to your generated protobuf code
	"github.com/ashish19912009/zrms/services/authZ/internal/service"
	"google.golang.org/grpc"
)

type AuthZServer struct {
	pb.UnimplementedAuthZServiceServer
	service service.AuthZService
}

func NewAuthZServer() *AuthZServer {
	return &AuthZServer{
		service: service.NewAuthZService(),
	}
}

func (s *AuthZServer) Register(grpcServer *grpc.Server) {
	pb.RegisterAuthZServiceServer(grpcServer, s)
}
