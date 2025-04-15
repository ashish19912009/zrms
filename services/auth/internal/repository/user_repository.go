package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ashish19912009/services/auth/internal/constants"
	"github.com/ashish19912009/services/auth/internal/logger"
	"github.com/ashish19912009/services/auth/internal/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	GetUser(ctx context.Context, loginID_accountID string, accountType string) (*models.User, error)
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

func isUUID(s string) bool {
	_, err := uuid.Parse(s)
	return err == nil
}

func (r *userRepository) GetUser(ctx context.Context, indentifier string, accountType string) (*models.User, error) {
	if indentifier == "" || accountType == "" {
		logger.Error(constants.CredentialMissing, nil, map[string]interface{}{
			"method":       constants.Methods.GetUser,
			"login_id":     indentifier,
			"account_type": accountType,
		})
		return nil, fmt.Errorf(constants.CredentialMissing)
	}
	var query string
	if isUUID(indentifier) {
		query = "SELECT account_id, employee_id, account_type, name, mobile_no, password_hash, role, permissions, status FROM users.team_accounts WHERE account_id = $1 AND account_type = $2 AND deleted_at IS NULL"
	}
	query = "SELECT account_id, employee_id, account_type, name, mobile_no, password_hash, role, permissions, status FROM users.team_accounts WHERE login_id = $1 AND account_type = $2 AND deleted_at IS NULL"
	row := r.db.QueryRowContext(ctx, query, indentifier, accountType)

	var user models.User
	if err := row.Scan(
		&user.AccountID,
		&user.EmployeeID,
		&user.AccountType,
		&user.Name,
		&user.MobileNo,
		&user.Password,
		&user.Role,
		&user.Permissions,
		&user.Status); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Warn(constants.ErrUserNotFound, map[string]interface{}{
				"method":       constants.Methods.GetUser,
				"login_id":     indentifier,
				"account_type": accountType,
			})
			return nil, fmt.Errorf("%s: %w", constants.ErrUserNotFound, err)
		}
		logger.Error(constants.DBQueryFailed, err, map[string]interface{}{
			"method":       constants.Methods.GetUser,
			"login_id":     indentifier,
			"account_type": accountType,
		})
		return nil, fmt.Errorf("%s: %w", constants.DBQueryFailed, err)
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
