package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/ashish19912009/zrms/services/account/internal/client"
	config "github.com/ashish19912009/zrms/services/account/internal/config"
	"github.com/ashish19912009/zrms/services/account/internal/handler"
	"github.com/ashish19912009/zrms/services/account/internal/logger"
	"github.com/ashish19912009/zrms/services/account/internal/middleware"
	"github.com/ashish19912009/zrms/services/account/internal/repository"
	"github.com/ashish19912009/zrms/services/account/internal/service"
	"github.com/ashish19912009/zrms/services/account/pb"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var appEnv string

func loadEnv() {
	appEnv = os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = "local" // default
	}

	// load env file accordingly
	envFile := fmt.Sprintf("../../env/.env.%s", appEnv)
	if err := godotenv.Load(envFile); err != nil {
		log.Printf("No %s file found, continuing without it", envFile)
	}
	log.Printf("🌍 Loaded %s environment - %s", appEnv, envFile)
}

func connectDB() (*sql.DB, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is not set")
	}
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open DB: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping DB: %w", err)
	}
	return db, nil
}

func main() {
	loadEnv()
	config_yaml_path := fmt.Sprintf("../../config/config.%s.yaml", appEnv)
	var cfg *config.AppConfig
	cfg, err := config.LoadConfig(config_yaml_path)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	db, err := connectDB()
	if err != nil {
		log.Fatalf("DB error: %v", err)
	}
	defer db.Close()

	// connecting with authZ (authorization) micro service
	authzClient, err := client.NewAuthZServiceClient(cfg.AuthZService.Host, cfg.AuthZService.Port)
	if err != nil {
		logger.Fatal("failed to connect to authz service: %v", err, nil)
	}
	defer authzClient.Close()

	// connecting with authN (authentication) micro service
	authnClient, err := client.NewAuthNServiceClient(cfg.AuthNService.Host, cfg.AuthNService.Port)
	if err != nil {
		logger.Fatal("failed to connect to authn service: %v", err, nil)
	}
	defer authnClient.Close()

	repo := repository.NewRepository(db)
	adminRepo := repository.NewAdminRepository(db)
	svc := service.NewAccountService(repo, authzClient)
	svcAdmin, err := service.NewAdminService(adminRepo, repo, authzClient)
	if err != nil {
		logger.Fatal("failed to connect to admin auth service: %v", err, nil)
	}
	grpcHandler := handler.NewGRPCHandler(svc, svcAdmin)

	// Create listner
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", cfg.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	authInterceptor := middleware.NewAuthInterceptor(authzClient, authnClient)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor.Unary()),
	)
	pb.RegisterAccountServiceServer(grpcServer, grpcHandler)
	if cfg.Env != "production" {
		reflection.Register(grpcServer)
	}
	log.Printf("✅ Account gRPC server running on %s in %s enviroment", cfg.Port, appEnv)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
