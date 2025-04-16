package repository_test

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/ashish19912009/zrms/services/account/internal/model"
	"github.com/ashish19912009/zrms/services/account/internal/repository"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *repository.NewRepository) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	repo := repository.NewRepository(db)
	return db, mock, repo
}

func TestUpdateAccount(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	account := &model.FranchiseAccount{
		EmployeeID:  "EMP123",
		Name:        "John Doe",
		MobileNo:    "9999999999",
		AccountType: "staff",
		Status:      "active",
	}

	mock.ExpectExec("UPDATE accounts").
		WithArgs(account.Name, account.MobileNo, account.AccountType, account.Status, sqlmock.AnyArg(), account.EmployeeID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.UpdateAccount(account)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetFranchiseByID_NotFound(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	mock.ExpectQuery("SELECT (.+) FROM franchises").
		WithArgs("fr123").
		WillReturnRows(sqlmock.NewRows([]string{})) // No columns/rows returned

	fr, err := repo.GetFranchiseByID("fr123")
	assert.Error(t, err)
	assert.Nil(t, fr)
}

func TestGetFranchiseOwner_DBError(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	mock.ExpectQuery("SELECT (.+) FROM accounts").
		WithArgs("fr123", "owner").
		WillReturnError(errors.New("db error"))

	acc, err := repo.GetFranchiseOwner("fr123")
	assert.Error(t, err)
	assert.Nil(t, acc)
}

func TestGetFranchiseDocuments_Success(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "franchise_id", "doc_type", "doc_url", "created_at", "updated_at"}).
		AddRow("doc1", "fr123", "license", "http://example.com/license.pdf", time.Now(), time.Now())

	mock.ExpectQuery("SELECT (.+) FROM franchise_documents").
		WithArgs("fr123").
		WillReturnRows(rows)

	docs, err := repo.GetFranchiseDocuments("fr123")
	assert.NoError(t, err)
	assert.Len(t, docs, 1)
	assert.Equal(t, "doc1", docs[0].ID)
}

func TestGetFranchiseAccounts_EmptyResult(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	mock.ExpectQuery("SELECT (.+) FROM accounts").
		WithArgs("fr123").
		WillReturnRows(sqlmock.NewRows([]string{
			"employee_id", "full_name", "mobile_no", "account_type", "status", "created_at", "updated_at",
		}))

	accs, err := repo.GetFranchiseAccounts("fr123")
	assert.NoError(t, err)
	assert.Len(t, accs, 0)
}
