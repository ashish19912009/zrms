package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/ashish19912009/zrms/services/account/internal/model"
	"github.com/ashish19912009/zrms/services/account/internal/repository"
)

/*
Responsibilities of AccountService:
Validate input data (e.g., check if mobileNo is valid before inserting).

Interact with Repository to fetch/update accounts.

Use goroutines & WaitGroups where needed for parallel execution.

Handle errors properly.

Return structured responses.
*/

// AccountService defines business logic for accounts
type AccountService interface {
	CreateAccount(ctx context.Context, account *model.Account) (*model.Account, error)
	UpdateAccount(ctx context.Context, account *model.Account) (*model.Account, error)
	GetAccountByID(ctx context.Context, id string) (*model.Account, error)
	ListAccounts(ctx context.Context, skip, take uint64) ([]*model.Account, error)
}

// accountService implements AccountService
type accountService struct {
	repo repository.Repository
}

// NewAccountService initializes service with a repository
func NewAccountService(repo repository.Repository) AccountService {
	return &accountService{repo: repo}
}

// CreateAccount adds a new account
func (s *accountService) CreateAccount(ctx context.Context, account *model.Account) (*model.Account, error) {
	if account.MobileNo == "" || account.Name == "" {
		return nil, fmt.Errorf("mobile no and name are required")
	}

	// Call repository to create account
	return s.repo.CreateAccount(ctx, account)
}

// UpdateAccount modify an existing account
func (s *accountService) UpdateAccount(ctx context.Context, account *model.Account) (*model.Account, error) {

	if account.ID == "" {
		return nil, fmt.Errorf("account ID is required")
	}
	return s.repo.UpdateAccount(ctx, account)
}

// GetAccountByID fetchs account details
func (s *accountService) GetAccountByID(ctx context.Context, id string) (*model.Account, error) {

	//check if id exist or not
	if id == "" {
		return nil, fmt.Errorf("id can't be empty")
	}
	return s.repo.GetAccountByID(ctx, id)
}

// ListAccounts retrieves multiple accounts
func (s *accountService) ListAccounts(ctx context.Context, skip, take uint64) ([]*model.Account, error) {
	if take > 100 || (skip == 0 && take == 0) {
		take = 100
	}

	// Parallelize the call if needed in future
	var wg sync.WaitGroup
	var mu sync.Mutex

	wg.Add(1)

	var accounts []*model.Account
	var err error

	go func() {
		defer wg.Done()
		result, queryErr := s.repo.ListAccounts(ctx, skip, take)

		mu.Lock() // Lock before modifying shared variables
		accounts = result
		err = queryErr
		mu.Unlock() // // Unlock after modification
	}()

	wg.Wait()

	if err != nil {
		return nil, fmt.Errorf("failed to list accounts:%w", err)
	}
	return accounts, nil
}
