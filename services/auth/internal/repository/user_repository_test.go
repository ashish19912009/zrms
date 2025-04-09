package repository_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ashish19912009/services/auth/internal/constants"
	"github.com/ashish19912009/services/auth/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

// func insertUser(t *testing.T, mock sqlmock.Sqlmock, user models.User) {
// 	t.Helper()

// 	mock.ExpectExec("INSERT INTO accounts").
// 		WithArgs(
// 			user.AccountID,
// 			user.EmployeeID,
// 			user.AccountType,
// 			user.Name,
// 			user.MobileNo,
// 			user.Password,
// 			user.Role,
// 			user.Permissions, // nil or []string
// 			user.Status,
// 			nil, // deleted_at
// 		).
// 		WillReturnResult(sqlmock.NewResult(1, 1))
// }

func TestGetUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewUserRepository(db)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"account_id", "employee_id", "account_type", "name", "mobile_no", "password", "role", "permissions", "status",
		}).AddRow(
			"acc-123", "emp-456", "admin", "Ashish", "9876543210", "hashedpass", "admin", []byte(`["read","write"]`), "active",
		)

		mock.ExpectQuery("SELECT (.+) FROM accounts").
			WithArgs("acc-123", "admin").
			WillReturnRows(rows)

		user, err := repo.GetUser(ctx, "acc-123", "admin")
		assert.NoError(t, err)
		assert.Equal(t, "acc-123", user.AccountID)
		assert.Equal(t, []string{"read", "write"}, []string(user.Permissions))
	})

	t.Run("Missing Credentials", func(t *testing.T) {
		user, err := repo.GetUser(ctx, "", "")
		assert.Nil(t, user)
		assert.ErrorContains(t, err, constants.CredentialMissing)
	})

	t.Run("User Not Found", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM accounts").
			WithArgs("unknown-id", "packer").
			WillReturnError(sql.ErrNoRows)

		user, err := repo.GetUser(ctx, "unknown-id", "packer")
		assert.Nil(t, user)
		assert.ErrorContains(t, err, constants.ErrUserNotFound)
	})

	t.Run("DB Error", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM accounts").
			WithArgs("acc-xyz", "admin").
			WillReturnError(errors.New("db crash"))

		user, err := repo.GetUser(ctx, "acc-xyz", "admin")
		assert.Nil(t, user)
		assert.ErrorContains(t, err, constants.DBQueryFailed)
	})

	t.Run("Empty loginID", func(t *testing.T) {
		repo := repository.NewUserRepository(db)
		_, err := repo.GetUser(context.Background(), "", "admin")
		assert.Error(t, err)
		assert.Equal(t, constants.CredentialMissing, err.Error())
	})

	t.Run("Empty accountType", func(t *testing.T) {
		repo := repository.NewUserRepository(db)
		_, err := repo.GetUser(context.Background(), "emp123", "")
		assert.Error(t, err)
		assert.Equal(t, constants.CredentialMissing, err.Error())
	})

	// t.Run("Permissions_NULL", func(t *testing.T) {
	// 	db, mock, err := sqlmock.New()
	// 	require.NoError(t, err)
	// 	defer db.Close()

	// 	repo := repository.NewUserRepository(db)

	// 	user := models.User{
	// 		AccountID:   "user-4",
	// 		EmployeeID:  "emp004",
	// 		AccountType: "packer",
	// 		Name:        "Ghost User",
	// 		MobileNo:    "0000000004",
	// 		Password:    "$2a$10$CJdLlcTv3tBK6eqjqi5t6e/u1vTpYQkQQDqf1M/g3GMNMw9FORftm",
	// 		Role:        "packer",
	// 		Permissions: nil, // NULL permissions
	// 		Status:      "active",
	// 	}

	// 	// insertUser(t, mock, user)

	// 	// Now prepare the expected SELECT query
	// 	rows := sqlmock.NewRows([]string{
	// 		"account_id", "employee_id", "account_type", "name", "mobile_no", "password", "role", "permissions", "status",
	// 	}).AddRow(
	// 		user.AccountID,
	// 		user.EmployeeID,
	// 		user.AccountType,
	// 		user.Name,
	// 		user.MobileNo,
	// 		user.Password,
	// 		user.Role,
	// 		nil, // NULL permissions from DB
	// 		user.Status,
	// 	)

	// 	mock.ExpectQuery("SELECT (.+) FROM accounts").
	// 		WithArgs(user.AccountID, user.AccountType).
	// 		WillReturnRows(rows)

	// 	gotUser, err := repo.GetUser(context.Background(), user.AccountID, user.AccountType)
	// 	require.NoError(t, err)
	// 	require.Nil(t, gotUser.Permissions)
	// })

	// t.Run("LoginID with SQL Injection pattern", func(t *testing.T) {
	// 	repo := repository.NewUserRepository(db)
	// 	_, err := repo.GetUser(context.Background(), "' OR 1=1 --", "admin")
	// 	assert.Error(t, err)
	// 	assert.Contains(t, err.Error(), constants.ErrUserNotFound)
	// })

	t.Run("Special_characters_in_loginID", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := repository.NewUserRepository(db)

		loginID := "🧙‍♂️🔥✨"
		accountType := "admin"

		mock.ExpectQuery("SELECT (.+) FROM accounts").
			WithArgs(loginID, accountType).
			WillReturnRows(sqlmock.NewRows([]string{
				"account_id", "employee_id", "account_type", "name", "mobile_no", "password", "role", "permissions", "status",
			})) // no rows returned

		user, err := repo.GetUser(context.Background(), loginID, accountType)
		require.Error(t, err)
		require.Contains(t, err.Error(), "user not found")
		require.Nil(t, user)
	})

}

func TestVerifyPassword(t *testing.T) {
	repo := repository.NewUserRepository(nil)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("secure-password"), bcrypt.DefaultCost)

	t.Run("Correct Password", func(t *testing.T) {
		result := repo.VerifyPassword(string(hashedPassword), "secure-password")
		assert.True(t, result)
	})

	t.Run("Incorrect Password", func(t *testing.T) {
		result := repo.VerifyPassword(string(hashedPassword), "wrong-password")
		assert.False(t, result)
	})
}
