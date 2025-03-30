package service_test

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/ashish19912009/zrms/services/account/internal/model"
	"github.com/ashish19912009/zrms/services/account/internal/repository/mocks"
	"github.com/ashish19912009/zrms/services/account/internal/service"
	"github.com/ashish19912009/zrms/services/account/internal/validations"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMain(m *testing.M) {
	validations.SetAllowedRoles([]string{"admin", "manager", "delivery"})
	validations.SetAllowedStatuses([]string{"active", "inactive", "suspended"})
	os.Exit(m.Run())
}

func TestCreateAccount_Success(t *testing.T) {

	mockRepo := new(mocks.Repository)
	svc := service.NewAccountService(mockRepo)

	now := time.Now()
	input := &model.Account{
		ID:         "uuid-123",
		MobileNo:   "9876543210",
		Name:       "Test User",
		Role:       "admin",  // must be one of allowed roles
		Status:     "active", // must be one of allowed statuses
		EmployeeID: "EMP001",
		CreatedAt:  &now,
	}

	mockRepo.On("CreateAccount", mock.Anything, input).Return(input, nil)

	result, err := svc.CreateAccount(context.Background(), input)
	assert.NoError(t, err)
	assert.Equal(t, input.ID, result.ID)
	mockRepo.AssertExpectations(t)
}

func TestCreateAccount_RepoError(t *testing.T) {
	mockRepo := new(mocks.Repository)
	svc := service.NewAccountService(mockRepo)
	acc := &model.Account{ID: "uuid-123", MobileNo: "9876543210", Name: "Test User", Role: "admin", Status: "active"}

	mockRepo.On("CreateAccount", mock.Anything, acc).Return(nil, errors.New("repo error"))

	_, err := svc.CreateAccount(context.Background(), acc)
	assert.Error(t, err)
	assert.EqualError(t, err, "repo error")
	mockRepo.AssertExpectations(t)
}

func TestCreateAccount_Invalid_Other_Required_Fields(t *testing.T) {
	mockRepo := new(mocks.Repository)
	svc := service.NewAccountService(mockRepo)

	acc := &model.Account{ID: "acc-001", MobileNo: "", Name: ""}

	_, err := svc.CreateAccount(context.Background(), acc)
	assert.Error(t, err)
	assert.EqualError(t, err, validations.ErrMobileRequired.Error())

	mockRepo.AssertNotCalled(t, "CreateAccount", mock.Anything, mock.Anything)
}

func TestUpdateAccount_Success(t *testing.T) {
	mockRepo := new(mocks.Repository)
	svc := service.NewAccountService(mockRepo)
	acc := &model.Account{
		ID:       "uuid-123",
		MobileNo: "9876543210",
		Name:     "Updated Name",
		Role:     "admin",
		Status:   "active",
	}

	mockRepo.On("UpdateAccount", mock.Anything, acc).Return(acc, nil)

	updated, err := svc.UpdateAccount(context.Background(), acc)
	assert.NoError(t, err)
	assert.Equal(t, acc.ID, updated.ID)
	mockRepo.AssertExpectations(t)
}

func TestUpdateAccount_RepoError(t *testing.T) {
	mockRepo := new(mocks.Repository)
	svc := service.NewAccountService(mockRepo)
	acc := &model.Account{
		ID:       "uuid-123",
		MobileNo: "9876543210",
		Name:     "Updated Name",
		Role:     "admin",
		Status:   "active",
	}

	mockRepo.On("UpdateAccount", mock.Anything, acc).Return(nil, errors.New("update error"))

	_, err := svc.UpdateAccount(context.Background(), acc)
	assert.Error(t, err)
	assert.EqualError(t, err, "update error")
	mockRepo.AssertExpectations(t)
}

func TestUpdateAccount_InvalidInput(t *testing.T) {
	mockRepo := new(mocks.Repository)
	svc := service.NewAccountService(mockRepo)
	acc := &model.Account{ID: "", Name: ""}

	_, err := svc.UpdateAccount(context.Background(), acc)
	assert.Error(t, err)
	assert.EqualError(t, err, validations.ErrAccountIDRequired.Error())

	mockRepo.AssertNotCalled(t, "UpdateAccount", mock.Anything, mock.Anything)
}

func TestGetAccountByID_Success(t *testing.T) {
	mockRepo := new(mocks.Repository)
	svc := service.NewAccountService(mockRepo)
	acc := &model.Account{ID: "uuid-123"}

	mockRepo.On("GetAccountByID", mock.Anything, "uuid-123").Return(acc, nil)

	result, err := svc.GetAccountByID(context.Background(), "uuid-123")
	assert.NoError(t, err)
	assert.Equal(t, acc.ID, result.ID)
	mockRepo.AssertExpectations(t)
}

func TestGetAccountByID_RepoError(t *testing.T) {
	mockRepo := new(mocks.Repository)
	svc := service.NewAccountService(mockRepo)

	mockRepo.On("GetAccountByID", mock.Anything, "uuid-123").Return(nil, errors.New("not found"))

	_, err := svc.GetAccountByID(context.Background(), "uuid-123")
	assert.Error(t, err)
	assert.EqualError(t, err, "not found")
	mockRepo.AssertExpectations(t)
}

func TestGetAccountByID_InvalidInput(t *testing.T) {
	mockRepo := new(mocks.Repository)
	svc := service.NewAccountService(mockRepo)

	_, err := svc.GetAccountByID(context.Background(), "")
	assert.Error(t, err)
	assert.EqualError(t, err, validations.ErrAccountIDRequired.Error())

	mockRepo.AssertNotCalled(t, "GetAccountByID", mock.Anything, "")
}

func TestListAccounts_Success(t *testing.T) {
	mockRepo := new(mocks.Repository)
	svc := service.NewAccountService(mockRepo)
	mockAccounts := []*model.Account{{ID: "acc-001"}}

	mockRepo.On("ListAccounts", mock.Anything, uint64(0), uint64(100)).Return(mockAccounts, nil)

	res, err := svc.ListAccounts(context.Background(), 0, 0)
	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, "acc-001", res[0].ID)
	mockRepo.AssertExpectations(t)
}

func TestListAccounts_RepoError(t *testing.T) {
	mockRepo := new(mocks.Repository)
	svc := service.NewAccountService(mockRepo)

	mockRepo.On("ListAccounts", mock.Anything, uint64(0), uint64(100)).Return(nil, errors.New("db error"))

	res, err := svc.ListAccounts(context.Background(), 0, 0)
	assert.Error(t, err)
	assert.Nil(t, res)
	assert.EqualError(t, err, "failed to list accounts:db error")
	mockRepo.AssertExpectations(t)
}

func TestListAccounts_InvalidPagination(t *testing.T) {
	mockRepo := new(mocks.Repository)
	svc := service.NewAccountService(mockRepo)

	mockRepo.On("ListAccounts", mock.Anything, uint64(9999999999), uint64(100)).Return(nil, errors.New("invalid pagination"))

	res, err := svc.ListAccounts(context.Background(), 9999999999, 9999999999)
	assert.Error(t, err)
	assert.Nil(t, res)
	assert.EqualError(t, err, "failed to list accounts:invalid pagination")
	mockRepo.AssertExpectations(t)
}
