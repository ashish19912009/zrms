package repository

import (
	"context"
	"database/sql"
	"fmt"
)

type Repository interface {
	CreateAccount(ctx context.Context, account *model.Account) (*model.Account, error)
	UpdateAccount(ctx context.Context, account *model.Account) (*model.Account, error)
	GetAccount(ctx context.Context, id string) (*model.Account, error)
	ListAccounts(ctx context.Context, skip uint64, take uint64) ([]*model.Account, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateAccount(ctx context.Context, account *model.Account) (*model.Account, error) {
	query := `INSERT INTO accounts (id, mobile_no, name, role, created_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, query, account.ID, account.MobileNo, account.Name, account.Role, account.createdAt)
	if err != nil {
		return nil, fmt.Errorf("Failed to create account %w", err)
	}
	return account, nil
}
