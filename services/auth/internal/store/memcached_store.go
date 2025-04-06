package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

var (
	ErrMemcachedConnection = errors.New("failed to connect to Memcached")
	ErrMemcachedOperation  = errors.New("memcached operation failed")
)

type MemcachedStore struct {
	client  *memcache.Client
	timeout time.Duration
}

func NewMemcachedStore(config *MemcachedConfig) (*MemcachedStore, error) {
	if config == nil {
		return nil, errors.New("memcached config cannot be nil")
	}

	client := memcache.New(config.Addresses...)
	client.Timeout = config.Timeout

	// Test connection
	_, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	err := client.Ping()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrMemcachedConnection, err)
	}

	return &MemcachedStore{
		client:  client,
		timeout: config.Timeout,
	}, nil
}

func (m *MemcachedStore) Set(key string, value interface{}) error {
	_, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	// Convert value to []byte
	var valueBytes []byte
	switch v := value.(type) {
	case []byte:
		valueBytes = v
	case string:
		valueBytes = []byte(v)
	default:
		return fmt.Errorf("%w: unsupported value type", ErrMemcachedOperation)
	}

	item := &memcache.Item{
		Key:        key,
		Value:      valueBytes,
		Expiration: int32(time.Hour.Seconds()), // Default 1 hour TTL
	}

	if err := m.client.Set(item); err != nil {
		return fmt.Errorf("%w: %v", ErrMemcachedOperation, err)
	}
	return nil
}

func (m *MemcachedStore) SetWithTTL(key string, value interface{}, ttl time.Duration) error {
	var valueBytes []byte

	switch v := value.(type) {
	case []byte:
		valueBytes = v
	case string:
		valueBytes = []byte(v)
	default:
		return fmt.Errorf("%w: unsupported value type", ErrMemcachedOperation)
	}

	// Memcached expects TTL in seconds (as int32)
	expiration := int32(ttl.Seconds())
	if expiration <= 0 {
		expiration = int32(m.timeout)
	}

	item := &memcache.Item{
		Key:        key,
		Value:      valueBytes,
		Expiration: expiration,
	}

	return m.client.Set(item)
}

func (m *MemcachedStore) Get(key string) (interface{}, error) {
	_, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	item, err := m.client.Get(key)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			return nil, ErrKeyNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrMemcachedOperation, err)
	}
	return item.Value, nil
}

func (m *MemcachedStore) Delete(key string) error {
	_, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	if err := m.client.Delete(key); err != nil && err != memcache.ErrCacheMiss {
		return fmt.Errorf("%w: %v", ErrMemcachedOperation, err)
	}
	return nil
}

func (m *MemcachedStore) Exists(key string) (bool, error) {
	_, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	_, err := m.client.Get(key)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			return false, nil
		}
		return false, fmt.Errorf("%w: %v", ErrMemcachedOperation, err)
	}
	return true, nil
}

func (m *MemcachedStore) Keys(pattern string) ([]string, error) {
	// Memcached doesn't natively support key pattern matching
	// This is a limitation compared to Redis
	return nil, fmt.Errorf("%w: key pattern matching not supported", ErrMemcachedOperation)
}

func (m *MemcachedStore) Close() error {
	// Memcached client doesn't have a close method
	// Connection pooling is handled automatically
	return nil
}

func (m *MemcachedStore) FlushAll() error {
	_, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	if err := m.client.FlushAll(); err != nil {
		return fmt.Errorf("%w: failed to flush all keys: %v", ErrMemcachedOperation, err)
	}
	return nil
}

// Interface compliance check
var _ InMemoryStore = (*MemcachedStore)(nil)
