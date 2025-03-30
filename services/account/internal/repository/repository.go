package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/ashish19912009/zrms/services/account/internal/model"
)

/**
go:generate mockery --name=Repository --output=internal/repository/mocks --case=underscore
mockery --name=Repository --dir=services/account/internal/repository --output=services/account/internal/repository/mocks --case=underscore
**/

type Repository interface {
	CreateAccount(ctx context.Context, account *model.Account) (*model.Account, error)
	UpdateAccount(ctx context.Context, account *model.Account) (*model.Account, error)
	GetAccountByID(ctx context.Context, id string) (*model.Account, error)
	ListAccounts(ctx context.Context, skip uint64, take uint64) ([]*model.Account, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateAccount(ctx context.Context, account *model.Account) (*model.Account, error) {
	if r.db == nil {
		log.Fatal("DB is nil")
	}
	if strings.TrimSpace(account.MobileNo) == "" || strings.TrimSpace(account.Name) == "" {
		return nil, fmt.Errorf("mobile number and name cannot be empty")
	}
	query := `INSERT INTO users.accounts (id, mobile_no, employee_id, name, role, status, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.ExecContext(ctx, query, account.ID, account.MobileNo, account.EmployeeID, account.Name, account.Role, account.Status, account.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}
	return account, nil
}

func (r *repository) UpdateAccount(ctx context.Context, account *model.Account) (*model.Account, error) {
	if strings.TrimSpace(account.ID) == "" {
		return nil, fmt.Errorf("account ID cannot be empty")
	}

	query := `UPDATE users.accounts SET mobile_no = $1, name = $2, role = $3, status = $4, employee_id = $5 WHERE id = $6`
	res, err := r.db.ExecContext(ctx, query, account.MobileNo, account.Name, account.Role, account.Status, account.EmployeeID, account.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to update account: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get affected rows: %w", err)
	}

	if affected == 0 {
		return nil, fmt.Errorf("no account updated")
	}

	return account, nil
}

func (r *repository) GetAccountByID(ctx context.Context, id string) (*model.Account, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("account ID cannot be empty")
	}

	query := `SELECT id, mobile_no, name, role, status, employee_id, created_at FROM users.accounts WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)

	var acc model.Account
	err := row.Scan(&acc.ID, &acc.MobileNo, &acc.Name, &acc.Role, &acc.Status, &acc.EmployeeID, &acc.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("sql: no rows in result set")
		}
		return nil, fmt.Errorf("failed to get account: %w", err)
	}
	return &acc, nil
}

func (r *repository) ListAccounts(ctx context.Context, skip, take uint64) ([]*model.Account, error) {
	if take == 0 {
		take = 100
	}

	query := `SELECT id, mobile_no, name, role, status, employee_id, created_at FROM users.accounts ORDER BY created_at DESC OFFSET $1 LIMIT $2`
	rows, err := r.db.QueryContext(ctx, query, skip, take)
	if err != nil {
		return nil, fmt.Errorf("failed to list accounts: %w", err)
	}
	defer rows.Close()

	var accounts []*model.Account
	for rows.Next() {
		var acc model.Account
		err := rows.Scan(&acc.ID, &acc.MobileNo, &acc.Name, &acc.Role, &acc.Status, &acc.EmployeeID, &acc.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan account: %w", err)
		}
		accounts = append(accounts, &acc)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return accounts, nil
}
