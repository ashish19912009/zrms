package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/ashish19912009/zrms/services/authZ/internal/config"
	"github.com/ashish19912009/zrms/services/authZ/internal/constants"
	"github.com/ashish19912009/zrms/services/authZ/internal/logger"
	"github.com/ashish19912009/zrms/services/authZ/internal/repository"
	"github.com/ashish19912009/zrms/services/authZ/internal/server"
	"github.com/ashish19912009/zrms/services/authZ/internal/service"
	"github.com/ashish19912009/zrms/services/authZ/internal/store"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
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
	// Load Enviroment
	loadEnv()
	logger.InitLogger("auth-service", zerolog.DebugLevel, "../../log_report/authz_service.log")
	logger.Info("Starting authorization services", map[string]interface{}{
		"env": appEnv,
	})

	// load config file
	configFilePath := fmt.Sprintf("../../config/config.%s.yaml", appEnv)
	var cfg *config.AppConfig
	cfg, err := config.LoadConfig(configFilePath)
	if err != nil {
		log.Fatalf(constants.FailedToLoadConfig, err)
	}

	// Connect Database
	db, err := connectDB()
	if err != nil {
		log.Fatalf("DB error: %v", err)
	}
	defer db.Close()

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", cfg.Port))
	if err != nil {
		logger.Fatal(constants.FailedToListen, err, nil)
	}

	// Select Cache DB
	storeManager, storeCfg, err := store.NewStoreManager(configFilePath)
	if err != nil {
		log.Fatalf(constants.FailedIniStrManager, err)
	}
	if storeCfg == nil {
		log.Fatal(constants.StoreConfigNil)
	}
	inMemoryStore := storeManager.Store()

	// Initialize repository
	repo := repository.NewAuthZRepository(db)
	cacheRepo, err := repository.NewCacheRepository(inMemoryStore)
	if err != nil {
		logger.Fatal(constants.FailedToStartCache, err, nil)
	}

	// Initialize service with precompiled Rego policy
	authzService, err := service.NewAuthZService(repo, cfg.RepoPolicyPath, cacheRepo)
	if err != nil {
		log.Fatalf(constants.FailedToStartService, err)
	}

	// Initialize gRPC server
	grpcServer := grpc.NewServer()

	// Register the gRPC server with the AuthZ service
	authzServer := server.NewAuthZServer(authzService)
	authzServer.Register(grpcServer)
	log.Printf(constants.GRPCServerRunning, cfg.Port, os.Getenv("APP_ENV"))

	// Start serving
	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatal(constants.FailedToStartServer, err, nil)
		log.Fatalf(constants.FailedToStartServer, err)
	}
}
