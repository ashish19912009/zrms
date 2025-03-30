package repository_test

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ashish19912009/zrms/services/account/internal/model"
	"github.com/ashish19912009/zrms/services/account/internal/repository"
	"github.com/stretchr/testify/assert"
)

// --- Test Case ---
func TestCreateAccount_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := repository.NewRepository(db)
	time := time.Now()
	acc := &model.Account{
		ID:         "6e040204-3b92-4499-a4f1-0fdb74f593f4",
		MobileNo:   "9876543210",
		Name:       "Test User",
		Role:       "admin",
		Status:     "active",
		EmployeeID: "EMP001",
		CreatedAt:  &time,
	}

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO users.accounts (id, mobile_no, employee_id, name, role, status, created_at)`)).
		WithArgs(acc.ID, acc.MobileNo, acc.EmployeeID, acc.Name, acc.Role, acc.Status, acc.CreatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	result, err := repo.CreateAccount(context.Background(), acc)

	assert.NoError(t, err)
	assert.Equal(t, acc.ID, result.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateAccount_EmptyFields(t *testing.T) {
	db, _, _ := sqlmock.New()
	defer db.Close()

	repo := repository.NewRepository(db)
	time := time.Now()
	acc := &model.Account{
		ID:        "acc_001",
		MobileNo:  "",
		Name:      "",
		Role:      "admin",
		Status:    "active",
		CreatedAt: &time,
	}

	_, err := repo.CreateAccount(context.Background(), acc)
	assert.Error(t, err)
}

func TestUpdateAccount_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewRepository(db)

	acc := &model.Account{
		ID:         "acc-001",
		MobileNo:   "9999999999",
		Name:       "Admin User",
		Role:       "admin",
		Status:     "active",
		EmployeeID: "EMP002",
	}

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE users.accounts SET mobile_no = $1, name = $2, role = $3, status = $4, employee_id = $5 WHERE id = $6`)).
		WithArgs(acc.MobileNo, acc.Name, acc.Role, acc.Status, acc.EmployeeID, acc.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	updated, err := repo.UpdateAccount(context.Background(), acc)
	assert.NoError(t, err)
	assert.Equal(t, acc.ID, updated.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateAccount_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewRepository(db)
	acc := &model.Account{ID: "acc-001"}

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE users.accounts SET mobile_no = $1, name = $2, role = $3, status = $4, employee_id = $5 WHERE id = $6`)).
		WithArgs(acc.MobileNo, acc.Name, acc.Role, acc.Status, acc.EmployeeID, acc.ID).
		WillReturnError(errors.New("update failed"))

	_, err = repo.UpdateAccount(context.Background(), acc)
	assert.Error(t, err)
	assert.EqualError(t, err, "failed to update account: update failed")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateAccount_InvalidInput(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewRepository(db)

	acc := &model.Account{} // Missing required fields
	_, err = repo.UpdateAccount(context.Background(), acc)
	assert.Error(t, err)
}

func TestGetAccountByID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewRepository(db)

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "mobile_no", "name", "role", "status", "employee_id", "created_at"}).
		AddRow("acc-001", "9999999999", "Admin User", "admin", "active", "EMP002", now)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, mobile_no, name, role, status, employee_id, created_at FROM users.accounts WHERE id = $1`)).
		WithArgs("acc-001").
		WillReturnRows(rows)

	acc, err := repo.GetAccountByID(context.Background(), "acc-001")
	assert.NoError(t, err)
	assert.NotNil(t, acc)
	assert.Equal(t, "acc-001", acc.ID)
	assert.Equal(t, "admin", acc.Role)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAccountByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, mobile_no, name, role, status, employee_id, created_at FROM users.accounts WHERE id = $1`)).
		WithArgs("invalid-id").
		WillReturnError(sql.ErrNoRows)

	acc, err := repo.GetAccountByID(context.Background(), "invalid-id")
	assert.Error(t, err)
	assert.Nil(t, acc)
	assert.Equal(t, sql.ErrNoRows, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAccountByID_InvalidInput(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewRepository(db)
	_, err = repo.GetAccountByID(context.Background(), "")
	assert.Error(t, err)
}

func TestListAccounts_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewRepository(db)
	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "mobile_no", "name", "role", "status", "employee_id", "created_at"}).
		AddRow("acc-001", "9999999999", "Admin User", "admin", "active", "EMP002", now)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, mobile_no, name, role, status, employee_id, created_at FROM users.accounts 
			ORDER BY created_at DESC OFFSET $1 LIMIT $2`)).
		WithArgs(uint64(0), uint64(100)).
		WillReturnRows(rows)

	accounts, err := repo.ListAccounts(context.Background(), 0, 0)
	assert.NoError(t, err)
	assert.Len(t, accounts, 1)
	assert.Equal(t, "acc-001", accounts[0].ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestListAccounts_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, mobile_no, name, role, status, employee_id, created_at FROM users.accounts 
			ORDER BY created_at DESC OFFSET $1 LIMIT $2`)).
		WithArgs(uint64(0), uint64(100)).
		WillReturnError(errors.New("query failed"))

	accounts, err := repo.ListAccounts(context.Background(), 0, 0)
	assert.Error(t, err)
	assert.Nil(t, accounts)
	assert.EqualError(t, err, "failed to list accounts: query failed")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestListAccounts_InvalidPagination(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewRepository(db)
	_, err = repo.ListAccounts(context.Background(), ^uint64(0), ^uint64(0)) // pass max uint64 values
	assert.Error(t, err)
}
