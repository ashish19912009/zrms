package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/ashish19912009/zrms/services/authN/internal/config"
	"github.com/ashish19912009/zrms/services/authN/internal/handler"
	"github.com/ashish19912009/zrms/services/authN/internal/jwk"
	"github.com/ashish19912009/zrms/services/authN/internal/logger"
	"github.com/ashish19912009/zrms/services/authN/internal/repository"
	"github.com/ashish19912009/zrms/services/authN/internal/service"
	"github.com/ashish19912009/zrms/services/authN/internal/store"
	"github.com/ashish19912009/zrms/services/authN/internal/token"
	pb "github.com/ashish19912009/zrms/services/authN/pb"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
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
	// load enviroment
	loadEnv()

	// load config file
	configFilePath := fmt.Sprintf("../../config/config.%s.yaml", appEnv)
	var cfg *config.AppConfig
	cfg, err := config.LoadConfig(configFilePath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	//Initialize logger
	logger.InitLogger("auth-service", zerolog.DebugLevel, "../../log_report/authN_service.log")
	logger.Info("Starting authentication service", map[string]interface{}{
		"env": appEnv,
	})

	// DB Connection
	db, err := connectDB()
	logger.Info("Store db", map[string]interface{}{
		"db": db,
	})
	if err != nil {
		log.Fatalf("DB error: %v", err)
	}
	defer db.Close()

	// Create listner
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", cfg.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Initialize in-memory store using config
	storeManager, storeCfg, err := store.NewStoreManager(configFilePath)
	// logger.Info("Store Manager", map[string]interface{}{
	// 	"storeManager": storeManager,
	// })
	// logger.Info("Store Config", map[string]interface{}{
	// 	"storeCfg": storeCfg,
	// })
	if err != nil {
		log.Fatalf("Failed to initialize store manager: %v", err)
	}
	if storeCfg == nil {
		log.Fatal("Store config is nil after loading, cannot proceed")
	}
	inMemoryStore := storeManager.Store()

	userRepo := repository.NewUserRepository(db)
	tokenRepo := repository.NewTokenRepository(inMemoryStore)
	tokenManger, err := token.NewjwtManager(cfg.JWTPrivateKeyPath, cfg.JWTPublicKeyPath, cfg.JWTHeader)
	if err != nil {
		log.Fatalf("failed to create JWT manager: %v", err)
	}

	// Initialize the JWK set using the manager's public key
	if err := jwk.InitializeJWK(tokenManger.PublicKey(), cfg.JWTHeader.Kid, cfg.JWTHeader.Alg, cfg.JWTHeader.Use); err != nil {
		log.Fatalf("Failed to initialize JWK: %v", err)
	}
	// register grpc server
	grpcServer := grpc.NewServer()

	var authService service.AuthService
	accessTTL := os.Getenv("ACCESS_TOKEN_TTL")
	refreshTTL := os.Getenv("REFRESH_TOKEN_TTL")
	authService = service.NewAuthService(tokenManger, tokenRepo, userRepo)
	if accessTTL != "" && refreshTTL != "" {
		attl, err1 := time.ParseDuration(accessTTL)
		rttl, err2 := time.ParseDuration(refreshTTL)
		if err1 != nil && err2 != nil {
			log.Fatalf("Invalid access duration: %v or refresh duration: %v", err1, err2)
		} else {
			authService = service.NewAuthServiceWithTTL(tokenManger, tokenRepo, userRepo, attl, rttl)
		}
	}

	// grpc handler
	grpcHandler, err := handler.NewGRPCHandler(authService)
	if err != nil {
		log.Fatalf("Failed to initialize gRPC handler: %v", err)
	}

	pb.RegisterAuthServiceServer(grpcServer, grpcHandler)

	if cfg.Env != "production" {
		reflection.Register(grpcServer)
	}
	httpErrChan := make(chan error)
	go func() {
		http.HandleFunc("/.well-known/jwks.json", jwk.Handler)

		httpPort := os.Getenv("JWK_HTTP_PORT")
		if httpPort == "" {
			httpPort = "8080" // default fallback
		}

		logger.Info("Starting JWK HTTP server", map[string]interface{}{
			"port": httpPort,
		})

		if err := http.ListenAndServe(":"+httpPort, nil); err != nil {
			httpErrChan <- err
		}
	}()

	// Start gRPC server in main goroutine
	log.Printf("✅ AuthN gRPC server is running on port %s", cfg.Port)
	grpcErrChan := make(chan error)
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			grpcErrChan <- err
		}
	}()

	// Wait for either server to fail
	select {
	case err := <-httpErrChan:
		logger.Error("HTTP server failed", err, nil)
		os.Exit(1)
	case err := <-grpcErrChan:
		logger.Error("gRPC server failed", err, nil)
		os.Exit(1)
	}
}
