package service_test

import (
	"context"
	"testing"
	"time"

	"zrms/services/account/internal/model"
	"zrms/services/account/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mock Repository ---
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) CreateAccount(ctx context.Context, acc *model.Account) (*model.Account, error) {
	args := m.Called(ctx, acc)
	return args.Get(0).(*model.Account), args.Error(1)
}

func (m *MockRepository) UpdateAccount(ctx context.Context, acc *model.Account) (*model.Account, error) {
	return nil, nil
}
func (m *MockRepository) GetAccountByID(ctx context.Context, id string) (*model.Account, error) {
	return nil, nil
}
func (m *MockRepository) ListAccounts(ctx context.Context, skip, take uint64) ([]*model.Account, error) {
	return nil, nil
}

// --- Test Case ---
func TestCreateAccount(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := service.NewAccountService(mockRepo)

	input := &model.Account{
		ID:         "uuid-123",
		MobileNo:   "9876543210",
		Name:       "Test User",
		Role:       "admin",
		Status:     "active",
		EmployeeID: "EMP001",
		CreatedAt:  time.Now(),
	}

	mockRepo.
		On("CreateAccount", mock.Anything, input).
		Return(input, nil)

	result, err := svc.CreateAccount(context.Background(), input)

	assert.NoError(t, err)
	assert.Equal(t, input.ID, result.ID)
	assert.Equal(t, "Test User", result.Name)
	mockRepo.AssertExpectations(t)
}
