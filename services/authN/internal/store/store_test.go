package store

import (
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type StoreTestSuite struct {
	suite.Suite
	store InMemoryStore
}

// Common tests that run for all stores
func (s *StoreTestSuite) TestBasicCRUD() {
	// Test Set and Get
	err := s.store.Set("key1", "value1")
	assert.NoError(s.T(), err)

	val, err := s.store.Get("key1")
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "value1", val)

	// Test Exists
	exists, err := s.store.Exists("key1")
	assert.NoError(s.T(), err)
	assert.True(s.T(), exists)

	// Test Delete
	err = s.store.Delete("key1")
	assert.NoError(s.T(), err)

	// Verify deletion
	_, err = s.store.Get("key1")
	assert.ErrorIs(s.T(), err, ErrKeyNotFound)
}

func (s *StoreTestSuite) TestTTL() {
	ttlStore, ok := s.store.(interface {
		SetWithTTL(key string, value interface{}, ttl time.Duration) error
	})
	if !ok {
		s.T().Skip("Store doesn't support TTL")
		return
	}

	err := ttlStore.SetWithTTL("temp", "data", time.Second)
	assert.NoError(s.T(), err)

	val, err := s.store.Get("temp")
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "data", val)

	time.Sleep(1100 * time.Millisecond)

	_, err = s.store.Get("temp")
	assert.ErrorIs(s.T(), err, ErrKeyNotFound)
}

func (s *StoreTestSuite) TestKeys() {
	err := s.store.Set("prefix_key1", "value1")
	assert.NoError(s.T(), err)
	err = s.store.Set("prefix_key2", "value2")
	assert.NoError(s.T(), err)

	keys, err := s.store.Keys("prefix_*")
	assert.NoError(s.T(), err)
	assert.GreaterOrEqual(s.T(), len(keys), 2)
}

func (s *StoreTestSuite) TestClose() {
	err := s.store.Close()
	assert.NoError(s.T(), err)
}
