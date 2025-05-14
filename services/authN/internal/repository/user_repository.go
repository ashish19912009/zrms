package repository

import (
	"context"
	"database/sql"

	"github.com/ashish19912009/zrms/services/authN/internal/constants"
	"github.com/ashish19912009/zrms/services/authN/internal/dbutils"
	"github.com/ashish19912009/zrms/services/authN/internal/logger"
	"github.com/ashish19912009/zrms/services/authN/internal/model"
)

var schema_outlet = constants.DB.Schema_Outlet
var schema_global = constants.DB.Schema_Global

type UserRepository interface {
	GetUser(ctx context.Context, loginID_accountID string, accountType string) (*model.User, error)
	GetFranchiseRolePermissions(ctx context.Context, franchiseID string) ([]*model.ResourceAction, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) GetUser(ctx context.Context, indentifier string, accountType string) (*model.User, error) {
	var method = constants.Methods.GetUser
	var table = constants.DB.Table_Franchise_Accounts
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	if err := dbutils.CheckDBConn(r.db, method); err != nil {
		return nil, err
	}

	// Define columns with alias prefixes
	columns := []string{
		"id",
		"franchise_id",
		"employee_id",
		"password_hash",
		"account_type",
		"name",
		"mobile_no",
		"email",
		"role_id",
		"status",
	}

	conditions := map[string]any{
		"login_id":     indentifier,
		"account_type": accountType,
	}

	opts := &dbutils.QueryBuilderOptions{
		Whilelist: struct {
			Schemas []string
			Tables  []string
			Columns []string
		}{
			Schemas: []string{schema_outlet},
			Tables:  []string{table},
			Columns: columns,
		},
	}

	// Use the BuildSelectQuery helper function to build the query
	query, args, err := dbutils.BuildSelectQuery(method, schema_outlet, table, columns, conditions, opts)
	if err != nil {
		return nil, err
	}
	logger.Info("qquery", map[string]interface{}{
		"qquery":     query,
		"args":       args,
		"conditions": conditions,
	})
	var user model.User

	if err := dbutils.ExecuteAndScanRow(ctx, method, r.db, query, args,
		&user.AccountID,
		&user.FranchiseID,
		&user.EmployeeID,
		&user.Password,
		&user.AccountType,
		&user.Name,
		&user.MobileNo,
		&user.Email,
		&user.RoleID,
		&user.Status); err != nil {
		logger.Error(constants.DBQueryError, err, nil)
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetFranchiseRolePermissions(ctx context.Context, franchiseID string) ([]*model.ResourceAction, error) {
	var method = constants.Methods.GetFranchiseRolePermissions
	var table = constants.DB.Table_Role
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	if err := dbutils.CheckDBConn(r.db, method); err != nil {
		return nil, err
	}

	// Define columns with alias prefixes
	columns := []string{
		"p.resource AS resource",
		"p.action AS action",
	}

	// Join clause to fetch role name
	joins := []dbutils.JoinClause{
		{
			Type:   "INNER",
			Schema: schema_outlet,
			Table:  constants.DB.Table_Role_Permissions,
			Alias:  "rp",
			On:     "r.id = rp.role_id",
		},
		{
			Type:   "INNER",
			Schema: schema_outlet,
			Table:  constants.DB.Table_Permissions,
			Alias:  "p",
			On:     "rp.permission_id = p.id",
		},
	}

	conditions := map[string]any{
		"r.franchise_id": franchiseID,
	}

	// Whitelist for security
	opts := &dbutils.QueryBuilderOptions{
		Whilelist: struct {
			Schemas []string
			Tables  []string
			Columns []string
		}{
			Schemas: []string{schema_outlet},
			Tables:  []string{table, constants.DB.Table_Role_Permissions, constants.DB.Table_Permissions},
			Columns: columns,
		},
	}

	// Build the join query
	query, args, err := dbutils.BuildJoinSelectQuery(method, schema_outlet, table, "r", columns, joins, conditions, opts)
	if err != nil {
		return nil, err
	}

	// Run the query
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		logger.Error(constants.DBQueryError, err, nil)
		return nil, err
	}
	defer rows.Close()

	// Process the result
	var allResourceMatrix []*model.ResourceAction
	for rows.Next() {
		var res *model.ResourceAction
		err := rows.Scan(
			&res.Resource,
			&res.Action,
		)
		if err != nil {
			logger.Error("Failed to scan row into account response", err, nil)
			return nil, err
		}
		allResourceMatrix = append(allResourceMatrix, res)
	}

	return allResourceMatrix, nil
}
