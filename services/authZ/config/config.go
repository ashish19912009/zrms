package config

import (
	"os"
)

type Config struct {
	GRPCServerAddress string
}

func Load() (*Config, error) {
	addr := os.Getenv("GRPC_SERVER_ADDRESS")
	if addr == "" {
		addr = ":50051" // default port
	}

	return &Config{
		GRPCServerAddress: addr,
	}, nil
}
