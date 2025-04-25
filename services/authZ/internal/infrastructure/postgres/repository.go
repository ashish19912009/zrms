package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/ashish19912009/zrms/services/authz/internal/domain"
	_ "github.com/lib/pq"
)

type PolicyRepository struct {
	db *sql.DB
}

func NewPolicyRepository(dsn string) (*PolicyRepository, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return &PolicyRepository{db: db}, nil
}

func (r *PolicyRepository) AddPolicy(ctx context.Context, policy *domain.Policy) error {
	query := `INSERT INTO policies (subject, action, resource) VALUES ($1, $2, $3)`
	_, err := r.db.ExecContext(ctx, query, policy.Subject, policy.Action, policy.Resource)
	return err
}

func (r *PolicyRepository) GetPoliciesForSubject(ctx context.Context, subject string) ([]*domain.Policy, error) {
	query := `SELECT subject, action, resource FROM policies WHERE subject = $1`
	rows, err := r.db.QueryContext(ctx, query, subject)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var policies []*domain.Policy
	for rows.Next() {
		var p domain.Policy
		if err := rows.Scan(&p.Subject, &p.Action, &p.Resource); err != nil {
			return nil, err
		}
		policies = append(policies, &p)
	}

	return policies, nil
}

// Other repository methods...
