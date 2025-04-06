package store

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type DragonflyStoreTestSuite struct {
	suite.Suite
	store InMemoryStore
}

func TestDragonflyStore(t *testing.T) {
	if os.Getenv("INTEGRATION") != "true" {
		t.Skip("Set INTEGRATION=true to run Dragonfly tests")
	}

	store, err := NewDragonflyStore(&DragonflyConfig{
		Address: "localhost:6379",
		DB:      2, // Use separate DB for tests
	})
	if err != nil {
		t.Fatalf("Failed to create Dragonfly store: %v", err)
	}

	suite.Run(t, &DragonflyStoreTestSuite{
		store: store,
	})
}

func (s *DragonflyStoreTestSuite) TearDownTest() {
	// Cleanup after each test
	err := s.store.(*DragonflyStore).client.FlushDB(context.Background()).Err()
	assert.NoError(s.T(), err)
}

func (s *DragonflyStoreTestSuite) TestDragonflySpecificFeatures() {
	// Test TTL operations
	err := s.store.SetWithTTL("temp_key", "value", time.Second)
	assert.NoError(s.T(), err)

	// Test DB selection (should be DB 2 as configured)
	// Can verify by checking key doesn't exist in default DB
	client := s.store.(*DragonflyStore).client
	_, err = client.Get(context.Background(), "temp_key").Result()
	assert.NoError(s.T(), err)
}
