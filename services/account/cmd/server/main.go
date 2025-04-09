package server

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	config "github.com/ashish19912009/zrms/services/account/config/account_config"
	"github.com/ashish19912009/zrms/services/account/internal/auth"
	"github.com/ashish19912009/zrms/services/account/internal/handler"
	"github.com/ashish19912009/zrms/services/account/internal/middleware"
	"github.com/ashish19912009/zrms/services/account/internal/repository"
	"github.com/ashish19912009/zrms/services/account/internal/service"
	"github.com/ashish19912009/zrms/services/account/pb"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func loadEnv() {
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = "local" // default
	}

	// load env file accordingly
	envFile := fmt.Sprintf(".env.%s", appEnv)
	if err := godotenv.Load(envFile); err != nil {
		log.Printf("No %s file found, continuing without it", envFile)
	}
	log.Printf("🌍 Loaded %s environment", appEnv)
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

func mockValidateToken(token string) (map[string]interface{}, error) {
	if token == "valid-token" {
		return map[string]interface{}{"role": "admin"}, nil
	}
	return nil, errors.New("unauthorized")
}

func main() {
	loadEnv()
	config.LoadConfig()
	db, err := connectDB()
	if err != nil {
		log.Fatalf("DB error: %v", err)
	}
	defer db.Close()

	repo := repository.NewRepository(db)
	svc := service.NewAccountService(repo)
	jwtManager := auth.NewJWTManager("secret-key", time.Hour)
	grpcHandler := handler.NewGRPCHandler(svc, jwtManager)

	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort != "" {
		grpcPort = "50051"
	}

	listener, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", grpcPort, err)
	}

	authInterceptor := middleware.NewAuthInterceptor("my-secret-key", "admin", "superadmin")

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor.Unary()),
	)
	pb.RegisterAccountServiceServer(grpcServer, grpcHandler)

	log.Printf("gRPC server running on port %s in %s enviroment", grpcPort, os.Getenv("APP_ENV"))
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
