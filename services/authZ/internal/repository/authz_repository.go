package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ashish19912009/zrms/services/authZ/internal/constants"
	"github.com/ashish19912009/zrms/services/authZ/internal/dbutils"
	"github.com/ashish19912009/zrms/services/authZ/internal/logger"
	_ "github.com/lib/pq"
)

type AuthZRepository interface {
	GetAccountRole(ctx context.Context, franchiseID, accountID string) (string, string, string, error)
	GetRolePermissions(ctx context.Context, role_id string) (map[string][]string, error)
	GetDirectPermissions(ctx context.Context, accountID string) (map[string]bool, error)
}

var schema_outlet = constants.DB.Schema_Outlet
var schema_global = constants.DB.Schema_Global

type authZRepo struct {
	db *sql.DB
}

func NewAuthZRepository(db *sql.DB) *authZRepo {
	return &authZRepo{
		db: db,
	}
}

// GetAccountRole fetches the role ID for a given account
func (r *authZRepo) GetAccountRole(ctx context.Context, franchiseID, accountID string) (string, string, string, error) {
	var method = constants.Methods.GetAccountRole
	var table = constants.DB.Table_Franchise_Accounts
	if err := dbutils.CheckDBConn(r.db, method); err != nil {
		return "", "", "", err
	}
	// Define columns with alias prefixes
	columns := []string{
		"id",
		"franchise_id",
		"role_id",
	}

	conditions := map[string]any{
		"franchise_id": franchiseID,
		"id":           accountID,
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
		return "", "", "", err
	}
	// Execute query and scan the result into the model
	var ID sql.NullString
	var FID sql.NullString
	var roleID sql.NullString
	if err := dbutils.ExecuteAndScanRow(ctx, method, r.db, query, args,
		&ID, &FID, &roleID); err != nil {
		return "", "", "", err
	}
	if roleID.Valid {
		return ID.String, FID.String, roleID.String, nil
	}
	return "", "", "", nil
	// query := `
	// 	SELECT ta.id,ta.franchise_id, ta.role_id
	// 	FROM outlet.team_accounts ta
	// 	WHERE ta.franchise_id = $1 AND ta.id = $2
	// `
}

// GetRolePermissions fetches all permissions for a role ID
func (r *authZRepo) GetRolePermissions(ctx context.Context, roleID string) (map[string][]string, error) {
	var method = constants.Methods.GetRolePermissions
	var table = constants.DB.Table_Permissions

	if err := dbutils.CheckDBConn(r.db, method); err != nil {
		return nil, err
	}
	// Define columns with alias prefixes
	columns := []string{
		"p.resource",
		"p.action",
	}
	// Join with global.document_types using document_type_id
	joins := []dbutils.JoinClause{
		{
			Type:   "INNER",
			Schema: schema_outlet,
			Table:  table,
			Alias:  "p",
			On:     "rp.permission_id = p.id",
		},
	}
	// Conditions
	conditions := map[string]any{
		"rp.role_id": roleID,
	}
	opts := &dbutils.QueryBuilderOptions{
		Whilelist: struct {
			Schemas []string
			Tables  []string
			Columns []string
		}{
			Schemas: []string{schema_outlet},
			Tables:  []string{table, constants.DB.Table_Role_Permissions},
			Columns: columns,
		},
	}
	// Build query using join-aware builder
	query, args, err := dbutils.BuildJoinSelectQuery(
		method,
		schema_outlet,
		constants.DB.Table_Role_Permissions,
		"rp",
		columns,
		joins,
		conditions,
		opts,
	)
	if err != nil {
		return nil, err
	}
	// Execute the query
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		logger.Error(constants.DBQueryError, err, nil)
		return nil, err
	}
	defer rows.Close()

	permissions := make(map[string][]string)
	for rows.Next() {
		var resource, action string
		if err := rows.Scan(&resource, &action); err != nil {
			return nil, fmt.Errorf("error scanning role permissions: %w", err)
		}
		permissions[resource] = append(permissions[resource], action)
	}
	return permissions, nil
	// query := `
	// 	SELECT p.resource, p.action
	// 	FROM outlet.role_permissions rp
	// 	INNER JOIN outlet.permissions p ON p.id = rp.permission_id
	// 	WHERE rp.role_id = $1
	// `

	// rows, err := r.db.QueryContext(ctx, query, roleID)
	// if err != nil {
	// 	return nil, fmt.Errorf("error fetching role permissions: %w", err)
	// }
	// defer rows.Close()
}

// GetDirectPermissions fetches all direct permissions for an account
func (r *authZRepo) GetDirectPermissions(ctx context.Context, accountID string) (map[string]bool, error) {
	var (
		method = constants.Methods.GetDirectPermissions
		table  = constants.DB.Table_Permissions
	)
	if err := dbutils.CheckDBConn(r.db, method); err != nil {
		return nil, err
	}

	// Define the columns you need to retrieve
	columns := []string{
		"p.key", "dp.is_granted",
	}

	// Join with global.document_types using document_type_id
	joins := []dbutils.JoinClause{
		{
			Type:   "INNER",
			Schema: schema_outlet,
			Table:  table,
			Alias:  "p",
			On:     "dp.permission_id = p.id",
		},
	}
	// Conditions
	conditions := map[string]any{
		"dp.account_id": accountID,
	}
	// Whitelist options
	opts := &dbutils.QueryBuilderOptions{
		Whilelist: struct {
			Schemas []string
			Tables  []string
			Columns []string
		}{
			Schemas: []string{schema_outlet},
			Tables:  []string{table, constants.DB.Table_Direct_Permissions},
			Columns: columns,
		},
	}
	// Build query using join-aware builder
	query, args, err := dbutils.BuildJoinSelectQuery(
		method,
		schema_outlet,
		constants.DB.Table_Direct_Permissions,
		"dp",
		columns,
		joins,
		conditions,
		opts,
	)
	if err != nil {
		return nil, err
	}
	// Execute the query
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		logger.Error(constants.DBQueryError, err, nil)
		return nil, err
	}
	defer rows.Close()

	directPerms := make(map[string]bool)
	for rows.Next() {
		var key string
		var isGranted bool
		if err := rows.Scan(&key, &isGranted); err != nil {
			return nil, fmt.Errorf("error scanning direct permission: %w", err)
		}
		directPerms[key] = isGranted
	}
	return directPerms, nil
}
