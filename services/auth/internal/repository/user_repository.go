package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ashish19912009/services/auth/internal/constants"
	"github.com/ashish19912009/services/auth/internal/logger"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	AccountID   string
	EmployeeID  string
	AccountType string
	Name        string
	MobileNo    string
	Password    string
	Role        string
	Permissions []string
	Status      string
}

type UserRepository interface {
	GetUser(ctx context.Context, loginID_accountID string, accountType string) (*User, error)
	VerifyPassword(hashedPassword, password string) bool
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) GetUser(ctx context.Context, loginID_accountID string, accountType string) (*User, error) {
	if loginID_accountID == "" || accountType == "" {
		logger.Error(constants.CredentialMissing, nil, map[string]interface{}{
			"method":       constants.Methods.GetUser,
			"login_id":     loginID_accountID,
			"account_type": accountType,
		})
		return nil, fmt.Errorf(constants.CredentialMissing)
	}

	query := "SELECT account_id, employee_id, account_type, name, mobile_no, password, role, permissions, status FROM accounts WHERE (login_id = $1 || account_id = $1) AND account_type = $2 AND deleted_at IS NULL"
	row := r.db.QueryRowContext(ctx, query, loginID_accountID, accountType)

	var user User
	if err := row.Scan(&user.AccountID, &user.AccountType, &user.Name, &user.MobileNo, &user.Password, &user.Role, &user.Permissions, &user.Status); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Warn(constants.ErrUserNotFound, map[string]interface{}{
				"method":       constants.Methods.GetUser,
				"login_id":     loginID_accountID,
				"account_type": accountType,
			})
			return nil, errors.New(constants.ErrUserNotFound)
		}
		logger.Error(constants.DBQueryFailed, err, map[string]interface{}{
			"method":       constants.Methods.GetUser,
			"login_id":     loginID_accountID,
			"account_type": accountType,
		})
		return nil, errors.New(constants.DBQueryFailed)
	}

	logger.Info(constants.UserFetchedSuccessful, map[string]interface{}{
		"method":       constants.Methods.GetUser,
		"account_id":   user.AccountID,
		"account_type": user.AccountType,
	})
	return &user, nil
}

func (r *userRepository) VerifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		logger.Warn(constants.ErrInvalidPassword, map[string]interface{}{
			"method": constants.Methods.VerifyPassword,
		})
		return false
	}
	return true
}
