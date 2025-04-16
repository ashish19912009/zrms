package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ashish19912009/zrms/services/account/internal/constants"
	"github.com/ashish19912009/zrms/services/account/internal/dbutils"
	"github.com/ashish19912009/zrms/services/account/internal/logger"
	"github.com/ashish19912009/zrms/services/account/internal/model"
)

/**
go:generate mockery --name=Repository --output=internal/repository/mocks --case=underscore
mockery --name=Repository --dir=services/account/internal/repository --output=services/account/internal/repository/mocks --case=underscore
**/

type Repository interface {
	UpdateAccount(ctx context.Context, id string, account *model.FranchiseAccount) (*model.FranchiseAccountResponse, error)
	GetFranchiseByID(ctx context.Context, id string) (*model.FranchiseResponse, error)
	GetFranchiseOwner(ctx context.Context, id string) (*model.FranchiseOwnerResponse, error)
	GetAccountByID(ctx context.Context, id string) (*model.FranchiseAccountResponse, error)
	GetFranchiseDocuments(ctx context.Context, id string) ([]model.FranchiseDocumentResponse, error)
	GetFranchiseAccounts(ctx context.Context, id string) ([]model.FranchiseAccountResponse, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (ar *repository) UpdateAccount(ctx context.Context, id string, account *model.FranchiseAccount) (*model.FranchiseAccountResponse, error) {
	var method = constants.Methods.UpdateAccount
	var schema = constants.DB.Schema_Outlet
	var table = constants.DB.Table_Franchise_Accounts

	if err := dbutils.CheckDBConn(ar.db, method); err != nil {
		return nil, err
	}

	// Fields to update
	data := map[string]any{
		"name":       account.Name,
		"email":      account.Email,
		"mobile_no":  account.MobileNo,
		"role":       account.RoleID,
		"status":     account.Status,
		"updated_at": account.UpdatedAt,
	}

	columns := make([]string, 0, len(data))
	for col := range data {
		columns = append(columns, col)
	}

	conditions := map[string]any{
		"id": id,
	}

	// Whitelist for safe updating
	opts := &dbutils.QueryBuilderOptions{
		Whilelist: struct {
			Schemas []string
			Tables  []string
			Columns []string
		}{
			Schemas: []string{schema},
			Tables:  []string{table},
			Columns: columns,
		},
		Returning: []string{
			"id", "franchise_id", "employee_id", "account_type", "name", "mobile_no",
			"email", "role", "status", "created_at", "updated_at",
		},
	}

	query, args, err := dbutils.BuildUpdateQuery(method, schema, table, columns, conditions, opts)
	if err != nil {
		return nil, err
	}

	var updated model.FranchiseAccountResponse
	err = dbutils.ExecuteAndScanRow(ctx, method, ar.db, query, args,
		&updated.ID,
		&updated.FranchiseID,
		&updated.EmployeeID,
		&updated.AccountType,
		&updated.Name,
		&updated.MobileNo,
		&updated.Email,
		&updated.RoleName, // assuming role name is returned
		&updated.Status,
		&updated.CreatedAt,
		&updated.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &updated, nil
}

func (ar *repository) GetFranchiseByID(ctx context.Context, id string) (*model.FranchiseResponse, error) {
	var method = constants.Methods.GetFranchiseByID
	var schema = constants.DB.Schema_Outlet
	var table = constants.DB.Table_Franchise
	if err := dbutils.CheckDBConn(ar.db, method); err != nil {
		return nil, err
	}

	// Define the columns you need to retrieve
	columns := []string{
		"id", "business_name", "logo_url", "sub_domain", "theme_settings",
		"status", "created_at", "updated_at",
	}

	// Define the condition map for WHERE clause
	conditions := map[string]any{
		"id":         id,
		"deleted_at": nil, // Filter deleted records
	}

	// Prepare query options (if any), you can add Returning or whitelist logic here
	opts := &dbutils.QueryBuilderOptions{
		Whilelist: struct {
			Schemas []string
			Tables  []string
			Columns []string
		}{
			Schemas: []string{schema},
			Tables:  []string{table},
			Columns: columns,
		},
	}

	// Use the BuildSelectQuery helper function to build the query
	query, args, err := dbutils.BuildSelectQuery(method, schema, table, columns, conditions, opts)
	if err != nil {
		return nil, err
	}

	// Execute query and scan the result into the model
	var franchise model.FranchiseResponse
	if err := dbutils.ExecuteAndScanRow(ctx, method, ar.db, query, args,
		&franchise.ID, &franchise.BusinessName, &franchise.LogoURL, &franchise.SubDomain,
		&franchise.ThemeSettings, &franchise.Status, &franchise.CreatedAt, &franchise.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &franchise, nil
}

func (ar *repository) GetFranchiseOwner(ctx context.Context, id string) (*model.FranchiseOwnerResponse, error) {
	var method = constants.Methods.GetFranchiseOwner
	var schema = constants.DB.Schema_Outlet
	var table = constants.DB.Table_Owner
	if err := dbutils.CheckDBConn(ar.db, method); err != nil {
		return nil, err
	}

	// Define the columns you need to retrieve
	columns := []string{
		"id", "name", "gender", "dob", "mobile_no", "email",
		"address", "aadhar_no", "is_verified", "created_at",
	}

	// Define the conditions for WHERE clause
	conditions := map[string]any{
		"franchise_id": id, // Use franchise_id for filtering
	}

	// Prepare query options (if any), you can add Returning or whitelist logic here
	opts := &dbutils.QueryBuilderOptions{
		Whilelist: struct {
			Schemas []string
			Tables  []string
			Columns []string
		}{
			Schemas: []string{schema},
			Tables:  []string{table},
			Columns: columns,
		},
	}

	// Use BuildSelectQuery to build the query
	query, args, err := dbutils.BuildSelectQuery(method, schema, table, columns, conditions, opts)
	if err != nil {
		return nil, err
	}

	// Execute the query and scan the result into the model
	var owner model.FranchiseOwnerResponse
	err = dbutils.ExecuteAndScanRow(ctx, method, ar.db, query, args,
		&owner.ID, &owner.Name, &owner.Gender, &owner.Dob, &owner.MobileNo,
		&owner.Email, &owner.Address, &owner.AadharNo, &owner.IsVerified, &owner.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &owner, nil
}

func (ar *repository) GetFranchiseDocuments(ctx context.Context, id string) ([]model.FranchiseDocumentResponse, error) {
	var method = constants.Methods.GetFranchiseDocuments
	var schema = constants.DB.Schema_Outlet
	var table = constants.DB.Table_Document_Types
	if err := dbutils.CheckDBConn(ar.db, method); err != nil {
		return nil, err
	}

	// Define columns with alias prefixes
	columns := []string{
		"fd.id",
		"dt.name AS document_type",
		"fd.document_url",
		"fd.uploaded_at",
	}

	// Join with global.document_types using document_type_id
	joins := []dbutils.JoinClause{
		{
			Type:   "INNER",
			Schema: schema,
			Table:  table,
			Alias:  "dt",
			On:     "fd.document_type_id = dt.id",
		},
	}

	// Conditions
	conditions := map[string]any{
		"fd.franchise_id": id,
	}

	// Whitelist options
	opts := &dbutils.QueryBuilderOptions{
		Whilelist: struct {
			Schemas []string
			Tables  []string
			Columns []string
		}{
			Schemas: []string{schema},
			Tables:  []string{table},
			Columns: []string{"fd.id", "dt.name", "fd.document_url", "fd.uploaded_at"},
		},
	}

	// Build query using join-aware builder
	query, args, err := dbutils.BuildJoinSelectQuery(
		method,
		schema,
		table,
		"fd",
		columns,
		joins,
		conditions,
		opts,
	)
	if err != nil {
		return nil, err
	}

	// Execute the query
	rows, err := ar.db.QueryContext(ctx, query, args...)
	if err != nil {
		logger.Error(constants.DBQueryError, err, nil)
		return nil, err
	}
	defer rows.Close()

	// Process result
	var docs []model.FranchiseDocumentResponse
	for rows.Next() {
		var doc model.FranchiseDocumentResponse
		err := rows.Scan(&doc.ID, &doc.DocumentName, &doc.IsMandate, &doc.DocumentURL, &doc.UploadedBy)
		if err != nil {
			logger.Error("Failed to scan row", err, nil)
			return nil, err
		}
		docs = append(docs, doc)
	}
	return docs, nil
}

func (ar *repository) GetFranchiseAccounts(ctx context.Context, id string) ([]model.FranchiseAccountResponse, error) {
	var method = constants.Methods.GetFranchiseAccounts
	var schema = constants.DB.Schema_Outlet
	var table = constants.DB.Table_Franchise_Accounts

	if err := dbutils.CheckDBConn(ar.db, method); err != nil {
		return nil, err
	}

	// Define columns (alias used for join clarity)
	columns := []string{
		"fa.id",
		"fa.franchise_id",
		"fa.employee_id",
		"fa.account_type",
		"fa.name",
		"fa.mobile_no",
		"fa.email",
		"r.name AS role", // joining role name
		"fa.status",
		"fa.created_at",
	}

	// Join clause to fetch role name
	joins := []dbutils.JoinClause{
		{
			Type:   "LEFT",
			Schema: schema,
			Table:  constants.DB.Table_Roles,
			Alias:  "r",
			On:     "fa.role = r.id",
		},
	}

	// WHERE conditions
	conditions := map[string]any{
		"fa.franchise_id": id,
	}

	// Whitelist for security
	opts := &dbutils.QueryBuilderOptions{
		Whilelist: struct {
			Schemas []string
			Tables  []string
			Columns []string
		}{
			Schemas: []string{schema},
			Tables:  []string{table, constants.DB.Table_Roles},
			Columns: columns,
		},
	}

	// Build the join query
	query, args, err := dbutils.BuildJoinSelectQuery(method, schema, table, "fa", columns, joins, conditions, opts)
	if err != nil {
		return nil, err
	}

	// Run the query
	rows, err := ar.db.QueryContext(ctx, query, args...)
	if err != nil {
		logger.Error(constants.DBQueryError, err, nil)
		return nil, err
	}
	defer rows.Close()

	// Process the result
	var accounts []model.FranchiseAccountResponse
	for rows.Next() {
		var acc model.FranchiseAccountResponse
		err := rows.Scan(
			&acc.ID,
			&acc.FranchiseID,
			&acc.EmployeeID,
			&acc.AccountType,
			&acc.Name,
			&acc.MobileNo,
			&acc.Email,
			&acc.RoleName, // assumes RoleName is a string field in response model
			&acc.Status,
			&acc.CreatedAt,
		)
		if err != nil {
			logger.Error("Failed to scan row into account response", err, nil)
			return nil, err
		}
		accounts = append(accounts, acc)
	}

	return accounts, nil
}

func (r *repository) GetAccountByID(ctx context.Context, id string) (*model.FranchiseAccountResponse, error) {
	var method = constants.Methods.GetAccountByID
	var schema = constants.DB.Schema_Outlet
	var table = constants.DB.Table_Franchise_Accounts
	// Check DB connection health
	if err := dbutils.CheckDBConn(r.db, method); err != nil {
		return nil, err
	}

	columns := []string{
		"id", "employee_id", "name", "email", "mobile_no", "status",
		"franchise_id", "account_type", "created_at", "updated_at", "deleted_at",
	}

	conditions := map[string]any{
		"id": id,
	}

	opts := &dbutils.QueryBuilderOptions{}

	// Build SELECT query
	query, args, err := dbutils.BuildSelectQuery(method, schema, table, columns, conditions, opts)
	if err != nil {
		return nil, err
	}

	// Scan result into response struct
	var account model.FranchiseAccountResponse
	err = dbutils.ExecuteAndScanRow(ctx, method, r.db, query, args, &account)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Not found
		}
		return nil, err
	}

	return &account, nil
}
