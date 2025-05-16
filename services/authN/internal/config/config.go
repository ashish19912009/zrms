package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	Env               string `yaml:"env"`
	Port              string `yaml:"port"`
	JWTPrivateKeyPath string `yaml:"jwtPrivateKeyPath"`
	JWTPublicKeyPath  string `yaml:"jwtPublicKeyPath"`
	KeyID             string `yaml:"keyID"`
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
