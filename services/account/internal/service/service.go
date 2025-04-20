package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/ashish19912009/zrms/services/account/internal/model"
	"github.com/ashish19912009/zrms/services/account/internal/repository"
	"github.com/ashish19912009/zrms/services/account/internal/validations"
)

/*
Responsibilities of AccountService:
Validate input data (e.g., check if mobileNo is valid before inserting).

Interact with Repository to fetch/update accounts.

Use goroutines & WaitGroups where needed for parallel execution.

Handle errors properly.

Return structured responses.
*/
// mockery --name=AccountService --dir=services/account/internal/service --output=services/account/internal/service/mocks --outpkg=mocks

// AccountService defines business logic for accounts
type AccountService interface {
	UpdateAccount(ctx context.Context, id string, account *model.FranchiseAccount) (*model.FranchiseAccountResponse, error)
	GetFranchiseByID(ctx context.Context, id string) (*model.FranchiseResponse, error)
	GetFranchiseOwner(ctx context.Context, id string) (*model.FranchiseOwnerResponse, error)
	GetAccountByID(ctx context.Context, id string) (*model.FranchiseAccountResponse, error)
	GetFranchiseDocuments(ctx context.Context, id string) ([]model.FranchiseDocumentResponse, error)
	GetFranchiseAccounts(ctx context.Context, id string) ([]model.FranchiseAccountResponse, error)
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
	if s.repo == nil {
		log.Fatal("repo is nil")
	}
	if err := validations.ValidateAccount(account); err != nil {
		return nil, err
	}

	// Call repository to create account
	return s.repo.CreateAccount(ctx, account)
}

// UpdateAccount modify an existing account
func (s *accountService) UpdateAccount(ctx context.Context, account *model.Account) (*model.Account, error) {
	if err := validations.ValidateAccountUpdate(account); err != nil {
		return nil, err
	}

	existing, err := s.repo.GetAccountByID(ctx, account.ID)
	if err != nil {
		log.Fatalf("After req.MobileNo in service: '%s'\n", existing.MobileNo)
		return nil, err
	}

	if account.Name == "" {
		account.Name = existing.Name
	}
	if account.MobileNo == "" {
		account.MobileNo = existing.MobileNo
	}
	if account.Role == "" {
		account.Role = existing.Role
	}
	if account.Status == "" {
		account.Status = existing.Status
	}
	return s.repo.UpdateAccount(ctx, account)
}

// GetAccountByID fetchs account details
func (s *accountService) GetAccountByID(ctx context.Context, id string) (*model.Account, error) {

	//check if id exist or not
	if id == "" {
		return nil, errors.New(validations.ErrAccountIDRequired.Error())
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
