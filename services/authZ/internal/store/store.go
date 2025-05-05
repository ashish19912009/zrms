package store

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/ashish19912009/zrms/services/authZ/internal/constants"
	"github.com/ashish19912009/zrms/services/authZ/internal/logger"
	"github.com/klauspost/compress/zstd"
	"gopkg.in/yaml.v3"
)

var (
	ErrUnsupportedDatabase = errors.New(constants.ErrUnsupportedDatabase)
	ErrKeyNotFound         = errors.New(constants.ErrKeyNotFound)
	ErrInvalidConfig       = errors.New(constants.ErrInvalidConfig)
)

// InMemoryStore defines the interface for all in-memory databases
type InMemoryStore interface {
	Set(key string, value interface{}) error
	SetWithTTL(key string, value interface{}, ttl time.Duration) error
	Get(key string) (interface{}, error)
	Delete(key string) error
	Exists(key string) (bool, error)
	Keys(pattern string) ([]string, error)
	FlushAll() error
	Close() error
}

// Config holds the YAML configuration for selecting the database
type Config struct {
	Type      string           `yaml:"type"`
	Redis     *RedisConfig     `yaml:"redis,omitempty"`
	Memcached *MemcachedConfig `yaml:"memcached,omitempty"`
	Dragonfly *DragonflyConfig `yaml:"dragonfly,omitempty"`
	Badger    *BadgerConfig    `yaml:"badger,omitempty"`
	Lightning *LightningConfig `yaml:"lightning,omitempty"`
}

type LightningConfig struct {
	InitialCapacity int           `yaml:"initial_capacity"`
	MaxItems        int           `yaml:"max_items"` // 0 = unlimited
	CleanupInterval time.Duration `yaml:"cleanup_interval"`
	FallbackStore   InMemoryStore `yaml:"-"` // For runtime fallback
}

type RedisConfig struct {
	Address      string        `yaml:"address"`
	Password     string        `yaml:"password"`
	DB           int           `yaml:"db"`
	PoolSize     int           `yaml:"pool_size"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	TTL          time.Duration `yaml:"ttl"`
}

type MemcachedConfig struct {
	Addresses    []string      `yaml:"addresses"`
	Timeout      time.Duration `yaml:"timeout"`
	MaxIdleConns int           `yaml:"max_idle_conns"`
}

type DragonflyConfig struct {
	Address      string        `yaml:"address"`
	Password     string        `yaml:"password"`
	DB           int           `yaml:"db"`
	PoolSize     int           `yaml:"pool_size"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	TTL          time.Duration `yaml:"ttl"`
}

type BadgerConfig struct {
	Dir        string `yaml:"dir"`
	SyncWrites bool   `yaml:"sync_writes"`
	Logger     bool   `yaml:"logger"`
}

// LoadConfig reads the YAML file and returns the config
func LoadConfig(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		logger.Error(constants.FailedToParse, err, map[string]interface{}{"file_path": path})
		return nil, fmt.Errorf(constants.FailedToRead, err)
	}

	var config Config
	if err := yaml.Unmarshal(file, &config); err != nil {
		logger.Error(constants.FailedToParse, err, map[string]interface{}{"file_path": path})
		return nil, fmt.Errorf(constants.FailedToUnmarshal, err)
	}

	if config.Type == "" {
		return nil, fmt.Errorf(constants.YAMLFileIssue)
	}

	// Override with environment variable if set
	envStoreType := os.Getenv(constants.EnvVariable.IN_MEMORY_STORE_TYPE)
	if envStoreType != "" {
		logger.Info(constants.ConfigOverride, map[string]interface{}{"store_type": envStoreType})
		config.Type = envStoreType
	}

	if err := config.Validate(); err != nil {
		logger.Error(constants.TypeSpecify, err, map[string]interface{}{"config_type": config.Type})
		return nil, err
	}

	return &config, nil
}

func (c *Config) Validate() error {
	if c.Type == "" {
		logger.Error(constants.TypeSpecify, nil, map[string]interface{}{constants.TypeKey: "missing"})
		return ErrInvalidConfig
	}

	switch c.Type {
	case constants.LightningType:
		if c.Lightning == nil {
			c.Lightning = &LightningConfig{InitialCapacity: 1000, MaxItems: 0}
		}
	case constants.RedisType:
		if c.Redis == nil {
			logger.Error(constants.InvalidRedisConfig, nil, map[string]interface{}{constants.TypeKey: c.Type})
			return ErrInvalidConfig
		}
	case constants.MemcachedType:
		if c.Memcached == nil {
			logger.Error(constants.InvalidMemcachedConfig, nil, map[string]interface{}{constants.TypeKey: c.Type})
			return ErrInvalidConfig
		}
	case constants.DragonflyType:
		if c.Dragonfly == nil {
			logger.Error(constants.InvalidDragonflyConfig, nil, map[string]interface{}{constants.TypeKey: c.Type})
			return ErrInvalidConfig
		}
	case constants.BadgerType:
		if c.Badger == nil {
			logger.Error(constants.InvalidBadgerConfig, nil, map[string]interface{}{constants.TypeKey: c.Type})
			return ErrInvalidConfig
		}
	}
	return nil
}

// StoreManager manages the selected in-memory store
type StoreManager struct {
	store   InMemoryStore
	encoder *zstd.Encoder
	decoder *zstd.Decoder
}

func (sm *StoreManager) Store() InMemoryStore {
	return sm.store
}

// NewStoreManager initializes the store based on the config
func NewStoreManager(configPath string) (*StoreManager, *Config, error) {
	config, err := LoadConfig(configPath)
	if err != nil {
		logger.Info(constants.FallbackLightning, map[string]interface{}{"error": err.Error()})
		return &StoreManager{store: NewLightningDB(nil)}, nil, err
	}

	store, err := NewStoreFromConfig(config)
	if err != nil {
		logger.Info(constants.FallbackLightningDueToFailure, map[string]interface{}{"error": err.Error()})
		store = NewLightningDB(nil)
	}

	return &StoreManager{store: store}, config, nil
}

// NewStoreFromConfig creates a store based on config
func NewStoreFromConfig(config *Config) (InMemoryStore, error) {
	if config == nil {
		config = &Config{}
	}
	storeType := config.Type
	if storeType == "" {
		storeType = constants.LightningType
	}
	switch storeType {
	case constants.LightningType, "":
		return NewLightningDB(config.Lightning), nil
	default:
		logger.Error(constants.UnsupportedDatabaseType, nil, map[string]interface{}{constants.TypeKey: config.Type})
		return nil, ErrUnsupportedDatabase
	}
}
