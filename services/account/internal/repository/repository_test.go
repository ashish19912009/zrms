package repository_test

// import (
// 	"context"
// 	"database/sql"
// 	"errors"
// 	"testing"
// 	"time"

// 	"github.com/ashish19912009/zrms/services/account/internal/dbutils"
// 	"github.com/ashish19912009/zrms/services/account/internal/model"
// 	"github.com/ashish19912009/zrms/services/account/internal/repository"
// 	"github.com/stretchr/testify/assert"
// )

// type mockDB struct{}

// func (m *mockDB) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
// 	return nil, nil
// }

// func (m *mockDB) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
// 	return nil, nil
// }

// // Mocks for dbutils
// var (
// 	mockBuildUpdateQuery      func(method, schema, table string, columns []string, conditions map[string]any, opts *dbutils.QueryBuilderOptions) (string, []any, error)
// 	mockExecuteAndScanRow     func(ctx context.Context, method string, db *sql.DB, query string, args []any, dest ...any) error
// 	mockCheckDBConn           func(db *sql.DB, method string) error
// 	originalBuildUpdateQuery  = dbutils.BuildUpdateQuery
// 	originalExecuteAndScanRow = dbutils.ExecuteAndScanRow
// 	originalCheckDBConn       = dbutils.CheckDBConn
// )

// func restoreDBUtils() {
// 	dbutils.BuildUpdateQuery = originalBuildUpdateQuery
// 	dbutils.ExecuteAndScanRow = originalExecuteAndScanRow
// 	dbutils.CheckDBConn = originalCheckDBConn
// }

// func TestUpdateAccount_Success(t *testing.T) {
// 	defer restoreDBUtils()

// 	db := &sql.DB{}
// 	repo := repository.NewRepository(db)

// 	// Setup mocks
// 	dbutils.CheckDBConn = func(db *sql.DB, method string) error {
// 		return nil
// 	}

// 	dbutils.BuildUpdateQuery = func(method, schema, table string, columns []string, conditions map[string]any, opts *dbutils.QueryBuilderOptions) (string, []any, error) {
// 		return "UPDATE query", []any{"val1", "val2"}, nil
// 	}

// 	dbutils.ExecuteAndScanRow = func(ctx context.Context, method string, db *sql.DB, query string, args []any, dest ...any) error {
// 		// Simulate scan into FranchiseAccountResponse
// 		dest[0] = "acc-id"
// 		dest[1] = "franchise-id"
// 		dest[2] = "emp-id"
// 		dest[3] = "manager"
// 		dest[4] = "John Doe"
// 		dest[5] = "9876543210"
// 		dest[6] = "john@example.com"
// 		dest[7] = "admin"
// 		dest[8] = "active"
// 		dest[9] = time.Now()
// 		dest[10] = time.Now()
// 		return nil
// 	}

// 	ctx := context.Background()
// 	account := &model.FranchiseAccount{
// 		Name:      "John Doe",
// 		Email:     "john@example.com",
// 		MobileNo:  "9876543210",
// 		RoleID:    "admin",
// 		Status:    "active",
// 		UpdatedAt: time.Now(),
// 	}

// 	resp, err := repo.UpdateAccount(ctx, "acc-id", account)
// 	assert.NoError(t, err)
// 	assert.NotNil(t, resp)
// 	assert.Equal(t, "John Doe", resp.Name)
// 	assert.Equal(t, "john@example.com", resp.Email)
// 	assert.Equal(t, "admin", resp.RoleName)
// }

// func TestUpdateAccount_DBConnError(t *testing.T) {
// 	defer restoreDBUtils()

// 	db := &sql.DB{}
// 	repo := repository.NewRepository(db)

// 	dbutils.CheckDBConn = func(db *sql.DB, method string) error {
// 		return errors.New("db connection error")
// 	}

// 	ctx := context.Background()
// 	account := &model.FranchiseAccount{}

// 	resp, err := repo.UpdateAccount(ctx, "id", account)
// 	assert.Error(t, err)
// 	assert.Nil(t, resp)
// 	assert.Equal(t, "db connection error", err.Error())
// }
