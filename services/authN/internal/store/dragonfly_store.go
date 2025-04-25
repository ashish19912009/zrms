package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9" // Using Redis client for Dragonfly compatibility
)

type DragonflyStore struct {
	client     *redis.Client
	defaultTTL time.Duration
}

func NewDragonflyStore(config *DragonflyConfig) (*DragonflyStore, error) {
	if config == nil {
		return nil, errors.New("dragonfly config cannot be nil")
	}

	opts := &redis.Options{
		Addr:         config.Address,
		Password:     config.Password,
		DB:           config.DB,
		PoolSize:     config.PoolSize,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
	}

	client := redis.NewClient(opts)

	// Verify connection
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("dragonfly connection failed: %w", err)
	}

	// Set default TTL
	ttl := config.TTL
	if ttl == 0 {
		ttl = 24 * time.Hour // Default 24 hour TTL
	}

	return &DragonflyStore{
		client:     client,
		defaultTTL: ttl,
	}, nil
}

// Implement all InMemoryStore interface methods...
func (d *DragonflyStore) Set(key string, value interface{}) error {
	return d.client.Set(context.Background(), key, value, d.defaultTTL).Err()
}

func (d *DragonflyStore) SetWithTTL(key string, value interface{}, ttl time.Duration) error {
	return d.client.Set(context.Background(), key, value, ttl).Err()
}

func (d *DragonflyStore) Get(key string) (interface{}, error) {
	val, err := d.client.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return nil, ErrKeyNotFound
	}
	return val, err
}

func (d *DragonflyStore) Delete(key string) error {
	return d.client.Del(context.Background(), key).Err()
}

func (d *DragonflyStore) Exists(key string) (bool, error) {
	exists, err := d.client.Exists(context.Background(), key).Result()
	return exists > 0, err
}

func (d *DragonflyStore) Keys(pattern string) ([]string, error) {
	return d.client.Keys(context.Background(), pattern).Result()
}

func (d *DragonflyStore) Close() error {
	return d.client.Close()
}

func (d *DragonflyStore) FlushAll() error {
	if err := d.client.FlushAll(context.Background()).Err(); err != nil {
		return fmt.Errorf("failed to flush dragonfly store: %w", err)
	}
	return nil
}

// Interface compliance check
var _ InMemoryStore = (*DragonflyStore)(nil)
