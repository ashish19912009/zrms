package store

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type RedisStoreTestSuite struct {
	StoreTestSuite
}

func TestRedisStore(t *testing.T) {
	if os.Getenv("INTEGRATION") != "true" {
		t.Skip("Set INTEGRATION=true to run Redis tests")
	}

	store, err := NewRedisStore(&RedisConfig{
		Address: "localhost:6379",
		DB:      1, // Use separate DB for tests
	})
	if err != nil {
		t.Fatalf("Failed to create Redis store: %v", err)
	}

	suite.Run(t, &RedisStoreTestSuite{
		StoreTestSuite: StoreTestSuite{
			store: store,
		},
	})
}

func (s *RedisStoreTestSuite) TearDownTest() {
	// Cleanup after each test
	s.store.(*RedisStore).client.FlushDB(s.store.(*RedisStore).ctx)
}
