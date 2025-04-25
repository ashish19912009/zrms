package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/yourorg/authz-service/configs"
	"github.com/yourorg/authz-service/internal/app/authz"
	"github.com/yourorg/authz-service/internal/infrastructure/grpc"
	"github.com/yourorg/authz-service/internal/infrastructure/opa"
	"github.com/yourorg/authz-service/internal/infrastructure/postgres"
)

func main() {
	// Load configuration
	cfg, err := configs.LoadConfig("configs/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize OPA client
	opaClient, err := opa.NewClient(cfg.OPAPolicyPath)
	if err != nil {
		log.Fatalf("Failed to initialize OPA client: %v", err)
	}

	// Initialize PostgreSQL repository
	repo, err := postgres.NewPolicyRepository(cfg.PostgresDSN)
	if err != nil {
		log.Fatalf("Failed to initialize repository: %v", err)
	}

	// Create authz service
	authzService := authz.NewService(opaClient, repo)

	// Start gRPC server
	go func() {
		log.Printf("Starting gRPC server on %s", cfg.GRPCServerAddress)
		if err := grpc.StartGRPCServer(cfg.GRPCServerAddress, authzService); err != nil {
			log.Fatalf("Failed to start gRPC server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
}
