package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ashish19912009/zrms/services/account/internal/model"
)

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
	query := `INSERT INTO accounts (id, mobile_no, employee_id, name, role, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.ExecContext(ctx, query, account.ID, account.MobileNo, account.Name, account.Role, account.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("Failed to create account %w", err)
	}
	return account, nil
}

func (r *repository) UpdateAccount(ctx context.Context, account *model.Account) (*model.Account, error) {
	query := `UPDATE accounts SET mobile_no = $1, name = $2, role = $3 WHERE id = $4`
	_, err := r.db.ExecContext(ctx, query, account.MobileNo, account.Name, account.Role, account.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to update account: %w", err)
	}
	return account, nil
}

func (r *repository) GetAccountByID(ctx context.Context, id string) (*model.Account, error) {
	query := `SELECT id, mobile_no, name, role, created_at FROM accounts WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)

	var acc model.Account
	err := row.Scan(&acc.ID, &acc.MobileNo, &acc.Name, &acc.Role, &acc.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // not found
		}
		return nil, fmt.Errorf("failed to get account: %w", err)
	}
	return &acc, nil
}

func (r *repository) ListAccounts(ctx context.Context, skip, take uint64) ([]*model.Account, error) {
	query := `SELECT id, mobile_no, name, role, created_at FROM accounts ORDER BY created_at DESC OFFSET $1 LIMIT $2`
	rows, err := r.db.QueryContext(ctx, query, skip, take)
	if err != nil {
		return nil, fmt.Errorf("failed to list accounts: %w", err)
	}
	defer rows.Close()

	var accounts []*model.Account
	for rows.Next() {
		var acc model.Account
		err := rows.Scan(&acc.ID, &acc.MobileNo, &acc.Name, &acc.Role, &acc.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan account: %w", err)
		}
		accounts = append(accounts, &acc)
	}
	return accounts, nil
}
