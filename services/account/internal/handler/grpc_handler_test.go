// grpc_handler_test.go
package handler_test

/*
import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ashish19912009/zrms/services/account/internal/handler"
	"github.com/ashish19912009/zrms/services/account/internal/model"
	"github.com/ashish19912009/zrms/services/account/internal/service/mocks"
	pb "github.com/ashish19912009/zrms/services/account/pb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func timePtr(t time.Time) *time.Time {
	return &t
}

func TestCreateAccount_Success(t *testing.T) {
	mockSvc := new(mocks.AccountService)
	h := handler.NewGRPCHandler(mockSvc)

	acc := &model.Account{
		ID:         "acc-001",
		MobileNo:   "9876543210",
		Name:       "Test User",
		Role:       "admin",
		Status:     "active",
		EmployeeID: "EMP001",
		CreatedAt:  ptrTime(time.Now()),
	}

	mockSvc.On("CreateAccount", mock.Anything, mock.AnythingOfType("*model.Account")).Return(acc, nil)

	resp, err := h.CreateAccount(context.Background(), &pb.CreateAccountRequest{
		Id:       acc.ID,
		MobileNo: acc.MobileNo,
		Name:     acc.Name,
		Role:     acc.Role,
		Status:   acc.Status,
		EmpId:    acc.EmployeeID,
	})

	assert.NoError(t, err)
	assert.Equal(t, acc.ID, resp.Account.Id)
	mockSvc.AssertExpectations(t)
}

func TestCreateAccount_Error(t *testing.T) {
	mockSvc := new(mocks.AccountService)
	h := handler.NewGRPCHandler(mockSvc)

	mockSvc.On("CreateAccount", mock.Anything, mock.AnythingOfType("*model.Account")).
		Return(nil, errors.New("create error"))

	req := &pb.CreateAccountRequest{
		Id:       "acc-001",
		MobileNo: "9876543210",
		Name:     "John Doe",
		Role:     "admin",
		Status:   "active",
		EmpId:    "EMP001",
	}

	_, err := h.CreateAccount(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "create error")
	mockSvc.AssertExpectations(t)
}

func TestCreateAccount_MissingRequiredFields(t *testing.T) {
	mockSvc := new(mocks.AccountService)
	h := handler.NewGRPCHandler(mockSvc)

	_, err := h.CreateAccount(context.Background(), &pb.CreateAccountRequest{
		Id:       "acc-002",
		MobileNo: "",
		EmpId:    "",
		Name:     "",
		Role:     "",
		Status:   "",
	})
	assert.Error(t, err)
}

func TestCreateAccount_MissingNameAndRole(t *testing.T) {
	mockSvc := new(mocks.AccountService)
	h := handler.NewGRPCHandler(mockSvc)

	_, err := h.CreateAccount(context.Background(), &pb.CreateAccountRequest{
		Id:       "acc-003",
		MobileNo: "9876543210",
		Name:     "",
		Role:     "",
	})
	assert.Error(t, err)
}

func TestGetAccountByID_Success(t *testing.T) {

	mockSvc := new(mocks.AccountService)
	h := handler.NewGRPCHandler(mockSvc)

	acc := &model.Account{ID: "acc-001", CreatedAt: timePtr(time.Now())}
	mockSvc.On("GetAccountByID", mock.Anything, "acc-001").Return(acc, nil)

	resp, err := h.GetAccountByID(context.Background(), &pb.GetAccountByIDRequest{Id: "acc-001"})
	assert.NoError(t, err)
	assert.Equal(t, "acc-001", resp.Account.Id)
	mockSvc.AssertExpectations(t)
}

func TestGetAccountByID_Error(t *testing.T) {
	mockSvc := new(mocks.AccountService)
	h := handler.NewGRPCHandler(mockSvc)

	mockSvc.On("GetAccountByID", mock.Anything, "acc-001").Return(nil, errors.New("not found"))

	_, err := h.GetAccountByID(context.Background(), &pb.GetAccountByIDRequest{Id: "acc-001"})
	assert.Error(t, err)
	mockSvc.AssertExpectations(t)
}

func TestGetAccountByID_EmptyID(t *testing.T) {
	mockSvc := new(mocks.AccountService)
	h := handler.NewGRPCHandler(mockSvc)

	_, err := h.GetAccountByID(context.Background(), &pb.GetAccountByIDRequest{Id: ""})
	assert.Error(t, err)
}

func TestUpdateAccount_Success(t *testing.T) {
	mockSvc := new(mocks.AccountService)
	h := handler.NewGRPCHandler(mockSvc)

	acc := &model.Account{ID: "acc-001", Name: "Updated User"}
	mockSvc.On("UpdateAccount", mock.Anything, mock.AnythingOfType("*model.Account")).Return(acc, nil)

	resp, err := h.UpdateAccount(context.Background(), &pb.UpdateAccountRequest{
		Id:   acc.ID,
		Name: acc.Name,
	})
	assert.NoError(t, err)
	assert.Equal(t, "acc-001", resp.Account.Id)
	mockSvc.AssertExpectations(t)
}

func TestUpdateAccount_Error(t *testing.T) {
	mockSvc := new(mocks.AccountService)
	h := handler.NewGRPCHandler(mockSvc)

	mockSvc.On("UpdateAccount", mock.Anything, mock.Anything).Return(nil, errors.New("update failed"))

	_, err := h.UpdateAccount(context.Background(), &pb.UpdateAccountRequest{Id: "acc-001", MobileNo: "1234567890"})
	assert.Error(t, err)
	mockSvc.AssertExpectations(t)
}

func TestUpdateAccount_MissingID(t *testing.T) {
	mockSvc := new(mocks.AccountService)
	h := handler.NewGRPCHandler(mockSvc)

	_, err := h.UpdateAccount(context.Background(), &pb.UpdateAccountRequest{Id: ""})
	assert.Error(t, err)
}

func TestUpdateAccount_NoChanges(t *testing.T) {
	mockSvc := new(mocks.AccountService)
	h := handler.NewGRPCHandler(mockSvc)

	_, err := h.UpdateAccount(context.Background(), &pb.UpdateAccountRequest{
		Id: "acc-004",
	})
	assert.Error(t, err)
}

func TestListAccounts_Success(t *testing.T) {
	mockSvc := new(mocks.AccountService)
	h := handler.NewGRPCHandler(mockSvc)

	accounts := []*model.Account{{ID: "acc-001"}}
	mockSvc.On("ListAccounts", mock.Anything, uint64(0), uint64(10)).Return(accounts, nil)

	resp, err := h.GetAccounts(context.Background(), &pb.GetAccountsRequest{Skip: 0, Take: 10})
	assert.NoError(t, err)
	assert.Len(t, resp.Accounts, 1)
	assert.Equal(t, "acc-001", resp.Accounts[0].Id)
	mockSvc.AssertExpectations(t)
}

func TestListAccounts_Error(t *testing.T) {
	mockSvc := new(mocks.AccountService)
	h := handler.NewGRPCHandler(mockSvc)

	mockSvc.On("ListAccounts", mock.Anything, uint64(0), uint64(10)).Return(nil, errors.New("list error"))

	_, err := h.GetAccounts(context.Background(), &pb.GetAccountsRequest{Skip: 0, Take: 10})
	assert.Error(t, err)
	mockSvc.AssertExpectations(t)
}

func TestListAccounts_ExtremePagination(t *testing.T) {
	mockSvc := new(mocks.AccountService)
	h := handler.NewGRPCHandler(mockSvc)

	mockSvc.On("ListAccounts", mock.Anything, uint64(100000), uint64(1000)).Return([]*model.Account{}, nil)

	resp, err := h.GetAccounts(context.Background(), &pb.GetAccountsRequest{Skip: 100000, Take: 1000})
	assert.NoError(t, err)
	assert.Len(t, resp.Accounts, 0)
	mockSvc.AssertExpectations(t)
}

func TestListAccounts_ZeroTake(t *testing.T) {
	mockSvc := new(mocks.AccountService)
	h := handler.NewGRPCHandler(mockSvc)

	_, err := h.GetAccounts(context.Background(), &pb.GetAccountsRequest{Skip: 0, Take: 0})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "take must be greater than zero")

	mockSvc.AssertNotCalled(t, "GetAccounts")
}

func ptrTime(t time.Time) *time.Time {
	return &t
}

*/
