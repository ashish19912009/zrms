package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type JWTHeaderConfig struct {
	Alg string `yaml:"alg"`
	Typ string `yaml:"typ"`
	Kid string `yaml:"keyID"`
	Use string `yaml:"use"`
}

type AppConfig struct {
	Env               string          `yaml:"env"`
	Port              string          `yaml:"port"`
	JWTPrivateKeyPath string          `yaml:"jwtPrivateKeyPath"`
	JWTPublicKeyPath  string          `yaml:"jwtPublicKeyPath"`
	JWTHeader         JWTHeaderConfig `yaml:"jwtHeader"`
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
