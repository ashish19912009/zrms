package store

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type MemcachedStoreTestSuite struct {
	StoreTestSuite
}

func TestMemcachedStore(t *testing.T) {
	if os.Getenv("INTEGRATION") != "true" {
		t.Skip("Set INTEGRATION=true to run Memcached tests")
	}

	store, err := NewRedisStore(&RedisConfig{
		Address: "localhost:11211",
		DB:      1, // Use separate DB for tests
	})
	if err != nil {
		t.Fatalf("Failed to create memcached store: %v", err)
	}

	suite.Run(t, &MemcachedStoreTestSuite{
		StoreTestSuite: StoreTestSuite{
			store: store,
		},
	})
}
