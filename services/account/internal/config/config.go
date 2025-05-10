package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type AuthServiceConfig struct {
	Host string `yaml:"host_authz"`
	Port int    `yaml:"port_authz"`
}

type AppConfig struct {
	Env         string            `yaml:"env"`
	Port        string            `yaml:"port"`
	AuthService AuthServiceConfig `yaml:"auth_service"`
}

// LoadConfig reads the YAML config file and unmarshals it into a Config struct
func LoadConfig(path string) (*AppConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var cfg AppConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("error parsing YAML config: %w", err)
	}

	return &cfg, nil
}
