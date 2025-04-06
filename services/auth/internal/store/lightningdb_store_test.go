package store

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type LightningDBTestSuite struct {
	StoreTestSuite
}

func TestLightningDB(t *testing.T) {
	suite.Run(t, &LightningDBTestSuite{
		StoreTestSuite: StoreTestSuite{
			store: NewLightningDB(&LightningConfig{
				InitialCapacity: 100,
				CleanupInterval: time.Minute,
			}),
		},
	})
}

func (s *LightningDBTestSuite) TestMemoryLimits() {
	store := s.store.(*LightningDB)
	store.config.MaxItems = 2 // Set small limit

	err := store.Set("key1", "value1")
	assert.NoError(s.T(), err)
	err = store.Set("key2", "value2")
	assert.NoError(s.T(), err)

	err = store.Set("key3", "value3")
	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "cache capacity reached")
}

func (s *LightningDBTestSuite) TestFallback() {
	store := s.store.(*LightningDB)
	store.config.MaxItems = 1
	store.config.FallbackStore = NewLightningDB(nil) // Simple fallback

	err := store.Set("key1", "value1")
	assert.NoError(s.T(), err)
	err = store.Set("key2", "value2") // Should go to fallback
	assert.NoError(s.T(), err)

	_, err = store.Get("key2")
	assert.NoError(s.T(), err)
}
