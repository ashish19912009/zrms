package store

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisStore implements the InMemoryStore interface using Redis
type RedisStore struct {
	client *redis.Client
	ctx    context.Context
	ttl    time.Duration
}

// NewRedisStore initializes a new Redis store
func NewRedisStore(config *RedisConfig) (*RedisStore, error) {
	if config == nil {
		return nil, errors.New("redif config cannot be nil")
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
	ctx := context.Background()

	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("redis connection failed: %w", err)
	}

	ttl := config.TTL
	if ttl == 0 {
		ttl = time.Hour * 48 // Default TTL
	}

	return &RedisStore{
		client: client,
		ctx:    ctx,
		ttl:    ttl,
	}, nil
}

// Set stores a key-value pair in Redis
func (r *RedisStore) Set(key string, value interface{}) error {
	// Using a goroutine to optimize performance
	if err := r.client.Set(r.ctx, key, value, time.Hour).Err(); err != nil {
		log.Printf("Failed to set key %s: %v", key, err)
	}
	return nil
}

func (r *RedisStore) SetWithTTL(key string, value interface{}, ttl time.Duration) error {
	ctx, cancel := context.WithTimeout(r.ctx, r.ttl)
	defer cancel()
	return r.client.Set(ctx, key, value, ttl).Err()
}

// Get retrieves a value from Redis
func (r *RedisStore) Get(key string) (interface{}, error) {

	value, err := r.client.Get(r.ctx, key).Result()
	if err == redis.Nil {
		return nil, errors.New("key not found")
	}
	if err != nil {
		return nil, errors.New("error in redis operation")
	}
	return value, nil
}

// Delete removes a key from Redis
func (r *RedisStore) Delete(key string) error {
	if err := r.client.Del(r.ctx, key).Err(); err != nil {
		return errors.New("error in redis operation")
	}
	return nil
}

func (r *RedisStore) Exists(key string) (bool, error) {
	exists, err := r.client.Exists(r.ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("redis exists check failed: %w", err)
	}
	return exists == 1, nil
}

func (r *RedisStore) Keys(pattern string) ([]string, error) {
	keys, err := r.client.Keys(r.ctx, pattern).Result()
	if err != nil {
		return nil, fmt.Errorf("redis keys operation failed: %w", err)
	}
	return keys, nil
}

func (r *RedisStore) Close() error {
	return r.client.Close()
}

func (r *RedisStore) FlushAll() error {
	if err := r.client.FlushAll(r.ctx).Err(); err != nil {
		return fmt.Errorf("failed to flush redis store: %w", err)
	}
	return nil
}

// Interface compliance check
var _ InMemoryStore = (*RedisStore)(nil)
