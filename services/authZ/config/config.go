package config

import (
	"os"
)

type Config struct {
	GRPCServerAddress string
	RegoPath          string
}

func Load() (*Config, error) {
	addr := os.Getenv("grpc_server_address")
	if addr == "" {
		addr = ":50052" // default port
	}
	regoPath := os.Getenv("rego_policy_path")
	if regoPath == "" {
		regoPath = "../../policy/authz.rego" // default port
	}

	return &Config{
		GRPCServerAddress: addr,
		RegoPath:          regoPath,
	}, nil
}
