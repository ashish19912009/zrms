// config/test_config/e2e_config.go
package config

import (
	"os"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func LoadTestDBConfig() DBConfig {
	return DBConfig{
		Host:     os.Getenv("TEST_DB_HOST"),
		Port:     os.Getenv("TEST_DB_PORT"),
		User:     os.Getenv("TEST_DB_USER"),
		Password: os.Getenv("TEST_DB_PASSWORD"),
		DBName:   os.Getenv("TEST_DB_NAME"),
		SSLMode:  "disable",
	}
}
