package store

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type BadgerStoreTestSuite struct {
	suite.Suite
	store   InMemoryStore
	tempDir string
}

func TestBadgerStore(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "badger-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	store, err := NewBadgerStore(&BadgerConfig{
		Dir:        tempDir,
		SyncWrites: false,
	})
	if err != nil {
		t.Fatalf("Failed to create Badger store: %v", err)
	}

	suite.Run(t, &BadgerStoreTestSuite{
		store:   store,
		tempDir: tempDir,
	})
}

func (s *BadgerStoreTestSuite) TearDownSuite() {
	if err := s.store.Close(); err != nil {
		s.T().Errorf("Failed to close store: %v", err)
	}
	os.RemoveAll(s.tempDir)
}

func (s *BadgerStoreTestSuite) TestPersistence() {
	err := s.store.Set("persistent", "data")
	assert.NoError(s.T(), err)

	// Recreate store to test persistence
	newStore, err := NewBadgerStore(&BadgerConfig{
		Dir: s.tempDir,
	})
	if err != nil {
		s.T().Fatalf("Failed to recreate Badger store: %v", err)
	}

	val, err := newStore.Get("persistent")
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "data", val)

	if err := newStore.Close(); err != nil {
		s.T().Errorf("Failed to close new store: %v", err)
	}
}
