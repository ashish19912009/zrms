package main

import (
	"log"
	"net"

	"github.com/ashish19912009/zrms/services/authZ/internal/server"
	"google.golang.org/grpc"
	"honnef.co/go/tools/config"
)

func main() {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Start gRPC server
	lis, err := net.Listen("tcp", cfg.GRPCServerAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	authzServer := server.NewAuthZServer()

	authzServer.Register(s)

	log.Printf("AuthZ gRPC server running on %s", cfg.GRPCServerAddress)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
