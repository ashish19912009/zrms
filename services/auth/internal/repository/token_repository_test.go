package repository_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/ashish19912009/services/auth/internal/repository"
	"github.com/ashish19912009/services/auth/internal/store"
)

func setupTestRepo() repository.TokenRepository {
	memStore := store.NewLightningDB(nil) // in-memory store
	return repository.NewTokenRepository(memStore)
}

func TestStoreToken_And_CheckToken(t *testing.T) {
	repo := setupTestRepo()
	ctx := context.Background()

	keyName := "refresh"
	accountID := "emp123"
	token := "sample_token"

	// Store the token
	err := repo.StoreToken(ctx, keyName, accountID, token, 0)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// Check token should return true
	ok, err := repo.CheckToken(ctx, keyName, accountID, token)
	if err != nil {
		t.Fatalf("expected no error checking token, got: %v", err)
	}
	if !ok {
		t.Fatalf("expected token to match")
	}
}

func TestCheckToken_NotFound(t *testing.T) {
	repo := setupTestRepo()
	ctx := context.Background()

	ok, err := repo.CheckToken(ctx, "refresh", "nonexistent", "whatever")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if ok {
		t.Fatal("expected token check to return false for nonexistent key")
	}
}

func TestCheckToken_InvalidType(t *testing.T) {
	memStore := store.NewLightningDB(nil)
	_ = memStore.Set("emp123:refresh", 12345) // not a string

	repo := repository.NewTokenRepository(memStore)
	ctx := context.Background()

	ok, err := repo.CheckToken(ctx, "refresh", "emp123", "12345")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if ok {
		t.Fatal("expected token check to return false due to invalid type")
	}
}

func TestStoreToken_WithTTL(t *testing.T) {
	repo := setupTestRepo()
	ctx := context.Background()

	keyName := "refresh"
	accountID := "emp123"
	token := "sample_token"

	err := repo.StoreToken(ctx, keyName, accountID, token, 1*time.Second)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	time.Sleep(2 * time.Second)

	ok, err := repo.CheckToken(ctx, keyName, accountID, token)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if ok {
		t.Fatal("expected token to be expired")
	}
}

func TestDeleteToken(t *testing.T) {
	repo := setupTestRepo()
	ctx := context.Background()

	keyName := "refresh"
	accountID := "emp123"
	token := "to_delete"

	err := repo.StoreToken(ctx, keyName, accountID, token, 0)
	if err != nil {
		t.Fatalf("store failed: %v", err)
	}

	err = repo.DeleteToken(ctx, keyName, accountID)
	if err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	ok, err := repo.CheckToken(ctx, keyName, accountID, token)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if ok {
		t.Fatal("expected token to be deleted")
	}
}

func TestCheckToken_TokenMismatch(t *testing.T) {
	repo := setupTestRepo()
	ctx := context.Background()

	_ = repo.StoreToken(ctx, "refresh", "emp123", "token1", 0)

	ok, err := repo.CheckToken(ctx, "refresh", "emp123", "wrong_token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatal("expected token mismatch to return false")
	}
}

func TestConcurrentAccess(t *testing.T) {
	repo := setupTestRepo()
	ctx := context.Background()
	keyName, accountID := "refresh", "empRace"
	token := "race_token"

	// Store the token
	err := repo.StoreToken(ctx, keyName, accountID, token, 0)
	if err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(2)

		go func() {
			defer wg.Done()
			repo.CheckToken(ctx, keyName, accountID, token)
		}()

		go func() {
			defer wg.Done()
			repo.DeleteToken(ctx, keyName, accountID)
		}()
	}

	wg.Wait()
}
