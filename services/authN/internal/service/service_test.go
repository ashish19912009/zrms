package service

// import (
// 	"context"
// 	"errors"
// 	"testing"
// 	"time"

// 	"github.com/ashish19912009/zrms/services/authN/internal/constants"
// 	"github.com/ashish19912009/zrms/services/authN/internal/models"
// 	"github.com/ashish19912009/zrms/services/authN/pb"
// 	mockRepo "github.com/ashish19912009/zrms/services/authN/test/mocks/repository"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// )

// // MockTokenManager is a mock implementation of the TokenManager interface.
// type MockTokenManager struct {
// 	mock.Mock
// }

// func (m *MockTokenManager) GenerateAccessToken(employeeID, accountID, mobileNo, accountType, name string, permissions models.PermissionsArray, duration time.Duration) (string, error) {
// 	args := m.Called(employeeID, accountID, mobileNo, accountType, name, permissions, duration)
// 	return args.String(0), args.Error(1)
// }

// func (m *MockTokenManager) GenerateRefreshToken(accountID, accountType string, permissions models.PermissionsArray, duration time.Duration) (string, error) {
// 	args := m.Called(accountID, accountType, permissions, duration)
// 	return args.String(0), args.Error(1)
// }

// func (m *MockTokenManager) VerifyToken(tokenString string) (*models.AuthClaims, error) {
// 	args := m.Called(tokenString)
// 	return args.Get(0).(*models.AuthClaims), args.Error(1)
// }

// // MockTokenRepository is a mock implementation of the TokenRepository interface.
// type MockTokenRepository struct {
// 	mock.Mock
// }

// func (m *MockTokenRepository) StoreToken(ctx context.Context, tokenType, accountID, token string, duration time.Duration) error {
// 	args := m.Called(ctx, tokenType, accountID, token, duration)
// 	return args.Error(0)
// }

// func (m *MockTokenRepository) CheckToken(ctx context.Context, tokenType, accountID, token string) (bool, error) {
// 	args := m.Called(ctx, tokenType, accountID, token)
// 	return args.Bool(0), args.Error(1)
// }

// func (m *MockTokenRepository) DeleteToken(ctx context.Context, tokenType, accountID string) error {
// 	args := m.Called(ctx, tokenType, accountID)
// 	return args.Error(0)
// }

// // MockUserRepository is a mock implementation of the UserRepository interface.
// type MockUserRepository struct {
// 	mock.Mock
// }

// func (m *MockUserRepository) GetUser(ctx context.Context, loginID, accountType string) (*models.User, error) {
// 	args := m.Called(ctx, loginID, accountType)
// 	return args.Get(0).(*models.User), args.Error(1)
// }

// func (m *MockUserRepository) VerifyPassword(hashedPassword, password string) bool {
// 	args := m.Called(hashedPassword, password)
// 	return args.Bool(0)
// }

// // func hashPassword(pwd string) string {
// // 	hashed, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
// // 	return string(hashed)
// // }

// func TestLogin(t *testing.T) {

// 	ctx := context.Background()
// 	loginReq := &pb.LoginRequest{
// 		LoginId:     "testuser",
// 		Password:    "password",
// 		AccountType: "standard",
// 	}

// 	accessTTL := 15 * time.Minute
// 	refreshTTL := 7 * 24 * time.Hour

// 	t.Run("Successful Login", func(t *testing.T) {
// 		mockUserRepo := new(mockRepo.UserRepository)
// 		mockTokenRepo := new(mockRepo.TokenRepository)
// 		mockTokenManager := new(mockRepo.TokenManager)
// 		authSvc := NewAuthServiceWithTTL(mockTokenManager, mockTokenRepo, mockUserRepo, accessTTL, refreshTTL)

// 		mockPermissions := models.PermissionsArray{
// 			{
// 				UserAccount: models.UserAccountPermission{
// 					Create: true,
// 					Read:   true,
// 					Update: true,
// 					Delete: false,
// 				},
// 			},
// 		}

// 		mockUser := &models.User{
// 			EmployeeID:  "E123",
// 			AccountID:   "A123",
// 			AccountType: "standard",
// 			Name:        "Test User",
// 			MobileNo:    "1234567890",
// 			Password:    "hashedpassword",
// 			Permissions: mockPermissions,
// 			Status:      "active",
// 		}

// 		mockUserRepo.On("GetUser", ctx, "testuser", "standard").Return(mockUser, nil)
// 		mockUserRepo.On("VerifyPassword", "hashedpassword", "password").Return(true)

// 		accessToken := "access-token"
// 		refreshToken := "refresh-token"

// 		mockTokenManager.On("GenerateAccessToken", mockUser.EmployeeID, mockUser.AccountID, mockUser.MobileNo, mockUser.AccountType, mockUser.Name, mockPermissions, accessTTL).Return(accessToken, nil)
// 		mockTokenManager.On("GenerateRefreshToken", mockUser.AccountID, mockUser.AccountType, mockPermissions, refreshTTL).Return(refreshToken, nil)

// 		mockTokenRepo.On("StoreToken", ctx, constants.Access_token, mockUser.AccountID, accessToken, accessTTL).Return(nil)
// 		mockTokenRepo.On("StoreToken", ctx, constants.Refresh_token, mockUser.AccountID, refreshToken, refreshTTL).Return(nil)

// 		resp, err := authSvc.Login(ctx, loginReq)
// 		assert.NoError(t, err)
// 		assert.NotNil(t, resp)
// 		assert.Equal(t, mockUser.AccountID, resp.AccountId)
// 		assert.Equal(t, mockUser.EmployeeID, resp.EmployeeId)
// 		assert.Equal(t, mockUser.AccountType, resp.AccountType)
// 		assert.Equal(t, mockUser.Name, resp.Name)
// 		assert.Equal(t, mockUser.MobileNo, resp.MobileNo)

// 		// ✅ Now properly check permissions
// 		assert.Equal(t, mockPermissions, resp.Permissions)

// 		assert.Equal(t, accessToken, resp.AccessToken)
// 		assert.Equal(t, refreshToken, resp.RefreshToken)
// 	})

// 	// t.Run("Inactive Account", func(t *testing.T) {
// 	// 	mockRepo := new(MockUserRepository)
// 	// 	mockTokenRepo := new(MockTokenRepository)
// 	// 	mockTokenManager := new(MockTokenManager)
// 	// 	authSvc := NewAuthServiceWithTTL(mockTokenManager, mockTokenRepo, mockRepo, 15*time.Minute, 24*time.Hour)

// 	// 	user := &models.User{
// 	// 		AccountID:   "A123",
// 	// 		EmployeeID:  "E123",
// 	// 		AccountType: "standard",
// 	// 		Name:        "Test User",
// 	// 		MobileNo:    "1234567890",
// 	// 		Password:    hashPassword("secure-password"),
// 	// 		Role:        "user",
// 	// 		Permissions: []string{"read", "write"},
// 	// 		Status:      "inactive", // Important!
// 	// 	}

// 	// 	mockRepo.On("GetUser", mock.Anything, "testuser", "standard").Return(user, nil)

// 	// 	req := &pb.LoginRequest{
// 	// 		LoginId:     "testuser",
// 	// 		AccountType: "standard",
// 	// 		Password:    "secure-password",
// 	// 	}

// 	// 	_, err := authSvc.Login(context.Background(), req)
// 	// 	require.Error(t, err)
// 	// 	require.Equal(t, codes.PermissionDenied, status.Code(err))
// 	// })

// 	t.Run("User Not Found", func(t *testing.T) {
// 		mockTokenManager := new(MockTokenManager)
// 		mockTokenRepo := new(MockTokenRepository)
// 		mockUserRepo := new(MockUserRepository)
// 		authSvc := NewAuthServiceWithTTL(mockTokenManager, mockTokenRepo, mockUserRepo, accessTTL, refreshTTL)

// 		mockUserRepo.On("GetUser", mock.Anything, "testuser", "standard").Return((*models.User)(nil), errors.New(constants.ErrUserNotFound))

// 		resp, err := authSvc.Login(ctx, loginReq)
// 		assert.Error(t, err)
// 		assert.Nil(t, resp)
// 		assert.Equal(t, constants.WrongUsernamePassword, err.Error())
// 	})

// 	t.Run("Incorrect Password", func(t *testing.T) {
// 		mockTokenManager := new(MockTokenManager)
// 		mockTokenRepo := new(MockTokenRepository)
// 		mockUserRepo := new(MockUserRepository)
// 		authSvc := NewAuthServiceWithTTL(mockTokenManager, mockTokenRepo, mockUserRepo, accessTTL, refreshTTL)

// 		mockUser := &models.User{
// 			EmployeeID:  "E123",
// 			AccountID:   "A123",
// 			AccountType: "standard",
// 			Password:    "hashedpassword",
// 		}

// 		mockUserRepo.On("GetUser", ctx, "testuser", "standard").Return(mockUser, nil)
// 		mockUserRepo.On("VerifyPassword", "hashedpassword", "password").Return(false)

// 		resp, err := authSvc.Login(ctx, loginReq)
// 		assert.Error(t, err)
// 		assert.Nil(t, resp)
// 		assert.Equal(t, constants.WrongUsernamePassword, err.Error())
// 	})

// 	t.Run("Token Generation Failure", func(t *testing.T) {
// 		mockTokenManager := new(MockTokenManager)
// 		mockTokenRepo := new(MockTokenRepository)
// 		mockUserRepo := new(MockUserRepository)
// 		authSvc := NewAuthServiceWithTTL(mockTokenManager, mockTokenRepo, mockUserRepo, accessTTL, refreshTTL)

// 		mockPermissions := models.PermissionsArray{
// 			{
// 				UserAccount: models.UserAccountPermission{
// 					Create: true,
// 					Read:   true,
// 					Update: true,
// 					Delete: false,
// 				},
// 			},
// 		}

// 		mockUser := &models.User{
// 			EmployeeID:  "E123",
// 			AccountID:   "A123",
// 			AccountType: "standard",
// 			Password:    "hashedpassword",
// 			Permissions: mockPermissions,
// 			MobileNo:    "9876543210",
// 			Name:        "TokenFail",
// 		}

// 		mockUserRepo.On("GetUser", ctx, "testuser", "standard").Return(mockUser, nil)
// 		mockUserRepo.On("VerifyPassword", "hashedpassword", "password").Return(true)

// 		mockTokenManager.On("GenerateAccessToken", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return("", errors.New("access token error"))

// 		resp, err := authSvc.Login(ctx, loginReq)
// 		assert.Error(t, err)
// 		assert.Nil(t, resp)
// 		assert.Equal(t, "failed to generate access token: access token error", err.Error())
// 	})

// 	t.Run("Token Storage Failure", func(t *testing.T) {
// 		mockTokenManager := new(MockTokenManager)
// 		mockTokenRepo := new(MockTokenRepository)
// 		mockUserRepo := new(MockUserRepository)
// 		authSvc := NewAuthServiceWithTTL(mockTokenManager, mockTokenRepo, mockUserRepo, accessTTL, refreshTTL)
// 		mockPermissions := models.PermissionsArray{
// 			{
// 				UserAccount: models.UserAccountPermission{
// 					Create: true,
// 					Read:   true,
// 					Update: true,
// 					Delete: false,
// 				},
// 			},
// 		}
// 		mockUser := &models.User{
// 			EmployeeID:  "E123",
// 			AccountID:   "A123",
// 			AccountType: "standard",
// 			Password:    "hashedpassword",
// 			Permissions: mockPermissions,
// 			MobileNo:    "9876543210",
// 			Name:        "StorageFail",
// 		}

// 		mockUserRepo.On("GetUser", ctx, "testuser", "standard").Return(mockUser, nil)
// 		mockUserRepo.On("VerifyPassword", "hashedpassword", "password").Return(true)
// 		mockTokenManager.On("GenerateAccessToken", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return("access-token", nil)
// 		mockTokenManager.On("GenerateRefreshToken", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return("refresh-token", nil)
// 		mockTokenRepo.On("StoreToken", ctx, constants.Access_token, mockUser.AccountID, "access-token", mock.Anything).Return(errors.New("storage error"))

// 		resp, err := authSvc.Login(ctx, loginReq)
// 		assert.Error(t, err)
// 		assert.Nil(t, resp)
// 		assert.Equal(t, "failed to store token in in_memory_DB: storage error", err.Error())
// 	})

// 	// t.Run("Access Token Generation Fails", func(t *testing.T) {
// 	// 	mockRepo := new(MockUserRepository)
// 	// 	mockTokenRepo := new(MockTokenRepository)
// 	// 	mockTokenManager := new(MockTokenManager)
// 	// 	authSvc := NewAuthServiceWithTTL(mockTokenManager, mockTokenRepo, mockRepo, 15*time.Minute, 24*time.Hour)

// 	// 	user := &models.User{
// 	// 		AccountID:   "A123",
// 	// 		EmployeeID:  "E123",
// 	// 		AccountType: "standard",
// 	// 		Name:        "Test User",
// 	// 		MobileNo:    "1234567890",
// 	// 		Password:    hashPassword("secure-password"),
// 	// 		Role:        "user",
// 	// 		Permissions: []string{"read", "write"},
// 	// 		Status:      "active",
// 	// 	}

// 	// 	mockRepo.On("GetUser", mock.Anything, "testuser", "standard").Return(user, nil)
// 	// 	mockTokenManager.
// 	// 		On("GenerateAccessToken", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
// 	// 		Return("", errors.New("token error"))

// 	// 	req := &pb.LoginRequest{
// 	// 		LoginId:     "testuser",
// 	// 		AccountType: "standard",
// 	// 		Password:    "secure-password",
// 	// 	}

// 	// 	_, err := authSvc.Login(context.Background(), req)
// 	// 	require.Error(t, err)
// 	// 	require.Equal(t, codes.Internal, status.Code(err))
// 	// })

// }

// func TestLogin_InvalidPassword(t *testing.T) {
// 	mockUserRepo := new(mockRepo.UserRepository)
// 	mockTokenRepo := new(mockRepo.TokenRepository)
// 	mockTokenManager := new(mockRepo.TokenManager)

// 	authSvc := NewAuthServiceWithTTL(
// 		mockTokenManager,
// 		mockTokenRepo,
// 		mockUserRepo,
// 		time.Minute*15,
// 		time.Hour*24,
// 	)

// 	req := &pb.LoginRequest{
// 		LoginId:     "testuser",
// 		Password:    "wrongpassword",
// 		AccountType: "admin",
// 	}
// 	mockPermissions := models.PermissionsArray{
// 		{
// 			UserAccount: models.UserAccountPermission{
// 				Create: true,
// 				Read:   true,
// 				Update: true,
// 				Delete: false,
// 			},
// 		},
// 	}
// 	user := &models.User{
// 		AccountID:   "acc123",
// 		EmployeeID:  "emp123",
// 		MobileNo:    "9999999999",
// 		AccountType: "admin",
// 		Name:        "Test User",
// 		Password:    "$2a$10$HashedPassword", // hashed
// 		Permissions: mockPermissions,
// 	}

// 	mockUserRepo.On("GetUser", mock.Anything, "testuser", "admin").Return(user, nil)
// 	mockUserRepo.On("VerifyPassword", user.Password, "wrongpassword").Return(false)

// 	resp, err := authSvc.Login(context.Background(), req)

// 	assert.Error(t, err)
// 	assert.Nil(t, resp)
// }
