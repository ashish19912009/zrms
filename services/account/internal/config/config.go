package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type AuthZServiceConfig struct {
	Host string `yaml:"host_authz"`
	Port string `yaml:"port_authz"`
}

type AuthNServiceConfig struct {
	Host string `yaml:"host_authn"`
	Port string `yaml:"port_authn"`
}

type AppConfig struct {
	Env          string             `yaml:"env"`
	Port         string             `yaml:"port"`
	AuthZService AuthZServiceConfig `yaml:"authz_service"`
	AuthNService AuthNServiceConfig `yaml:"authn_service"`
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
