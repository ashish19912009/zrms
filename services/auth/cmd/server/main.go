package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/ashish19912009/services/auth/internal/logger"
	"github.com/ashish19912009/services/auth/internal/repository"
	"github.com/ashish19912009/services/auth/internal/store"
	"github.com/ashish19912009/services/auth/internal/token"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
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

func getAppConfig()

func main() {
	// load enviroment
	loadEnv()
	//Initialize logger
	logger.InitLogger("auth-service", zerolog.DebugLevel, "logs/auth_service.log")
	logger.Info("Starting authentication service", map[string]interface{}{
		"env": "local",
	})

	// get DB URL
	db, err := connectDB()
	if err != nil {
		log.Fatalf("DB error: %v", err)
	}
	defer db.Close()

	// Initialize in-memory store using config
	storeManager, err := store.NewStoreManager("../../config/config.yaml")
	if err != nil {
		errors.New("Failed to initialize store manager")
	}
	inMemoryStore := storeManager.Store()

	userRepo := repository.NewUserRepository(db)
	tokenRepo := repository.NewTokenRepository(inMemoryStore)
	tokenManger, err := token.NewjwtManager()

	//authService := service.NewAuthService(inMemoryStore)

}

/**
Example usage
logger.Info("Starting authentication service", map[string]interface{}{
	"env": "local",
})
logger.Debug("Debugging mode active", nil)
logger.Warn("High CPU usage detected", nil)
logger.Error("Database connection failed", nil, nil)
*/
