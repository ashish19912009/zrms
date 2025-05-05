package store

import (
	"errors"
	"sync"
	"time"

	"github.com/ashish19912009/zrms/services/authZ/internal/constants"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	ErrKeyAlreadyExists = errors.New("key already exists")
	metricsHitCount     = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "lightningdb_hits_total",
		Help: "Total cache hits",
	})
	metricsMissCount = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "lightningdb_misses_total",
		Help: "Total cache misses",
	})
	metricsItemCount = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "lightningdb_items_current",
		Help: "Current number of cached items",
	})
)

// LightningDB implements InMemoryStore with a simple in-memory cache
type LightningDB struct {
	store  map[string]item
	mu     sync.RWMutex
	config *LightningConfig
}

type item struct {
	value      interface{}
	expiration time.Time
}

func init() {
	prometheus.MustRegister(metricsHitCount, metricsMissCount, metricsItemCount)
}

// NewLightningDB creates a new in-memory store
func NewLightningDB(config *LightningConfig) *LightningDB {
	if config == nil {
		config = &LightningConfig{
			InitialCapacity: 1000,
			MaxItems:        0, // 0 means unlimited
			CleanupInterval: 5 * time.Minute,
		}
	}
	db := &LightningDB{
		store:  make(map[string]item, config.InitialCapacity),
		config: config,
	}
	// Start background cleanup if interval is set
	if config.CleanupInterval > 0 {
		go db.startCleanup(config.CleanupInterval)
	}
	return db
}

// Set stores a value with optional TTL (0 means no expiration)
func (l *LightningDB) Set(key string, value interface{}) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.config.MaxItems > 0 && len(l.store) >= l.config.MaxItems {
		if l.config.FallbackStore != nil {
			return l.config.FallbackStore.Set(key, value)
		}
		return errors.New(constants.CapacityReached)
	}

	l.store[key] = item{value: value, expiration: time.Time{}}
	metricsItemCount.Inc()
	return nil
}

// SetWithTTL stores a value with time-to-live in seconds
func (l *LightningDB) SetWithTTL(key string, value interface{}, ttl time.Duration) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.config.MaxItems > 0 && len(l.store) >= l.config.MaxItems {
		if l.config.FallbackStore != nil {
			return l.config.FallbackStore.SetWithTTL(key, value, ttl)
		}
		return errors.New(constants.CapacityReached)
	}
	expiration := time.Now().Add(ttl)
	l.store[key] = item{
		value:      value,
		expiration: expiration,
	}
	metricsItemCount.Inc()
	return nil
}

// Get retrieves a value if it exists and isn't expired
func (l *LightningDB) Get(key string) (interface{}, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	item, found := l.store[key]
	if !found {
		metricsMissCount.Inc()
		return nil, ErrKeyNotFound
	}

	if !item.expiration.IsZero() && time.Now().After(item.expiration) {
		metricsMissCount.Inc()
		go l.Delete(key) // Async cleanup
		return nil, ErrKeyNotFound
	}

	metricsHitCount.Inc()
	return item.value, nil
}

// Delete removes a key
func (l *LightningDB) Delete(key string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if _, exists := l.store[key]; !exists {
		return ErrKeyNotFound
	}

	delete(l.store, key)
	return nil
}

// Exists checks if a key exists and isn't expired
func (l *LightningDB) Exists(key string) (bool, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	item, found := l.store[key]
	if !found {
		return false, nil
	}

	if !item.expiration.IsZero() && time.Now().After(item.expiration) {
		return false, nil
	}

	return true, nil
}

// Keys returns all non-expired keys (warning: not scalable for large datasets)
func (l *LightningDB) Keys(pattern string) ([]string, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	var keys []string
	now := time.Now()

	for k, item := range l.store {
		if !item.expiration.IsZero() && now.After(item.expiration) {
			continue // Skip expired items
		}
		keys = append(keys, k)
	}

	return keys, nil
}

// Close cleans up resources (no-op for in-memory store)
func (l *LightningDB) Close() error {
	return nil
}

// SetIfNotExists only sets the key if it doesn't already exist
func (l *LightningDB) SetIfNotExists(key string, value interface{}) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if _, exists := l.store[key]; exists {
		return ErrKeyAlreadyExists
	}

	l.store[key] = item{value: value}
	return nil
}

func (l *LightningDB) startCleanup(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		l.mu.Lock()
		now := time.Now()
		for k, item := range l.store {
			if !item.expiration.IsZero() && now.After(item.expiration) {
				delete(l.store, k)
			}
		}
		l.mu.Unlock()
	}
}

func (l *LightningDB) FlushAll() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.store = make(map[string]item, l.config.InitialCapacity)
	metricsItemCount.Set(0)
	return nil
}

// Interface compliance check
var _ InMemoryStore = (*LightningDB)(nil)
