package main

import (
	"github.com/ashish19912009/services/auth/internal/logger"
	"github.com/rs/zerolog"
)

func main() {
	//Initialize logger
	logger.InitLogger("auth-service", zerolog.DebugLevel, "logs/auth_service.log")

	logger.Info("Starting authentication service", map[string]interface{}{
		"env": "local",
	})

	/**
	Example usage
	logger.Info("Starting authentication service", map[string]interface{}{
		"env": "local",
	})
	logger.Debug("Debugging mode active", nil)
	logger.Warn("High CPU usage detected", nil)
	logger.Error("Database connection failed", nil, nil)
	*/
}
