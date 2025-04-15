package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ashish19912009/zrms/services/account/internal/constants"
	"github.com/ashish19912009/zrms/services/account/internal/dbutils"
	"github.com/ashish19912009/zrms/services/account/internal/logger"
	"github.com/ashish19912009/zrms/services/account/internal/model"
)

/**
go:generate mockery --name=Repository --output=internal/repository/mocks --case=underscore
mockery --name=Repository --dir=services/account/internal/repository --output=services/account/internal/repository/mocks --case=underscore
**/

const (
	layer  = constants.Layer
	method = constants.Method
)
const (
	id             = "id"
	franchise_id   = "franchise_id"
	business_name  = "business_name"
	logo_url       = "logo_url"
	sub_domain     = "sub_domain"
	theme_settings = "theme_settings"
	status         = "status"
	created_at     = "created_at"
	updated_at     = "updated_at"
	deleted_at     = "deleted_at"

	name        = "name"
	gender      = "gender"
	dob         = "dob"
	mobile_no   = "mobile_no"
	email       = "email"
	address     = "address"
	aadhar_no   = "aadhar_no"
	is_verified = "is_verified"
)

type AdminRepository interface {
	CreateFranchise(ctx context.Context, franchise *model.Franchise, f_owner *model.FranchiseOwner) (*model.FranchiseResponse, error)
	UpdateFranchise(ctx context.Context, id string, franchise *model.Franchise) (*model.FranchiseResponse, error)
	UpdateFranchiseStatus(ctx context.Context, id string, status string) (string, error)
	DeleteFranchise(ctx context.Context, id string) (bool, error)

	GetFranchiseByID(ctx context.Context, id string) (*model.FranchiseResponse, error)
	GetAllFranchises(ctx context.Context, page int32, limit int32) ([]model.FranchiseResponse, error)
	GetFranchiseOwner(ctx context.Context, id string) (*model.FranchiseOwnerResponse, error)
	GetFranchiseDocuments(ctx context.Context, id string) ([]model.FranchiseDocumentResponse, error)
	GetFranchiseAccounts(ctx context.Context, id string) ([]model.FranchiseAccountResponse, error)
}

type admin_repository struct {
	db *sql.DB
}

func NewAdminRepository(db *sql.DB) AdminRepository {
	return &admin_repository{
		db: db,
	}
}

func (ar *admin_repository) CreateFranchise(ctx context.Context, franchiseInput *model.Franchise, f_owner *model.FranchiseOwner) (*model.FranchiseResponse, error) {
	logCtx := logger.BaseLogContext(layer, constants.Repository, method, constants.Methods.CreateFranchise)
	if err := dbutils.CheckDBConn(ar.db, logCtx); err != nil {
		return nil, err
	}
	opts := &dbutils.QueryBuilderOptions{
		Returning: []string{id, business_name, logo_url, sub_domain, theme_settings, status, created_at},
	}
	opts.Whilelist.Schemas = []string{constants.DB.Schema_Outlet}
	opts.Whilelist.Tables = []string{constants.DB.Table_Franchise}
	opts.Whilelist.Columns = []string{id, business_name, logo_url, sub_domain, theme_settings, status, created_at, updated_at, deleted_at}
	query, err := dbutils.BuildInsertQuery(constants.DB.Schema_Outlet, constants.DB.Table_Franchise, []string{business_name, logo_url, sub_domain, theme_settings, status}, opts)
	if err != nil {
		logCtx := logger.BaseLogContext(layer, constants.Repository, method, constants.Methods.CreateFranchise)
		logger.Fatal(constants.BuildInsertQuery, err, logCtx)
		return nil, err
	}

	var franchise model.FranchiseResponse
	if err = dbutils.ExecuteAndScanRow(ctx, ar.db, query,
		[]any{id, franchiseInput.BusinessName, franchiseInput.LogoURL, franchiseInput.SubDomain, franchiseInput.ThemeSettings, franchiseInput.Status},
		&franchise.ID, &franchise.BusinessName, &franchise.LogoURL, &franchise.SubDomain, &franchise.ThemeSettings,
		&franchise.Status,
		&franchise.CreatedAt,
	); err != nil {
		logCtx := logger.BaseLogContext(layer, constants.Repository, method, constants.Methods.CreateFranchise)
		logger.Fatal(constants.DBConnectionFailure, err, logCtx)
		return nil, err
	}

	// Update Franchise owner details
	opts.Returning = []string{id, created_at}
	opts.Whilelist.Tables = []string{constants.DB.Table_Owner}
	opts.Whilelist.Columns = append(opts.Whilelist.Columns, franchise_id, name, gender, dob, mobile_no, email, address, aadhar_no)

	query, err = dbutils.BuildInsertQuery(constants.DB.Schema_Outlet, constants.DB.Table_Owner, []string{franchise_id, name, gender, dob, mobile_no, email, address, aadhar_no, is_verified}, opts)
	if err != nil {
		logCtx := logger.BaseLogContext(layer, constants.Repository, method, constants.Methods.CreateFranchise)
		logger.Fatal(constants.BuildInsertQuery, err, logCtx)
		return nil, err
	}

	var owner_id string

	if err = dbutils.ExecuteAndScanRow(ctx, ar.db, query,
		[]any{franchise.ID, f_owner.Name, f_owner.Gender, f_owner.Dob, f_owner.MobileNo, f_owner.Email, f_owner.Address, f_owner.AadharNo, f_owner.IsVerified},
		&owner_id,
	); err != nil {
		logCtx := logger.BaseLogContext(layer, constants.Repository, method, constants.Methods.CreateFranchise)
		logger.Fatal(constants.DBConnectionFailure, err, logCtx)
		return nil, err
	}
	return &franchise, nil
}

func (ar *admin_repository) UpdateFranchise(ctx context.Context, id string, franchise *model.Franchise) (*model.FranchiseResponse, error) {
	logCtx := logger.BaseLogContext(layer, constants.Repository, method, constants.Methods.UpdateFranchise)
	if err := dbutils.CheckDBConn(ar.db, logCtx); err != nil {
		return nil, err
	}

	opts := &dbutils.QueryBuilderOptions{
		Returning: []string{id, business_name, logo_url, sub_domain, theme_settings, status, updated_at},
	}
	opts.Whilelist.Schemas = []string{constants.DB.Schema_Outlet}
	opts.Whilelist.Tables = []string{constants.DB.Table_Franchise}
	opts.Whilelist.Columns = []string{business_name, logo_url, sub_domain, theme_settings, status}

	query, err := dbutils.BuildUpdateQuery(constants.DB.Schema_Outlet, constants.DB.Table_Franchise,
		[]string{business_name, logo_url, sub_domain, theme_settings, status}, id, opts)
	if err != nil {
		logger.Fatal(constants.BuildUpdateQuery, err, logCtx)
		return nil, err
	}

	var updated model.FranchiseResponse
	err = dbutils.ExecuteAndScanRow(ctx, ar.db, query,
		[]any{franchise.BusinessName, franchise.LogoURL, franchise.SubDomain, franchise.ThemeSettings, franchise.Status, id},
		&updated.ID, &updated.BusinessName, &updated.LogoURL, &updated.SubDomain, &updated.ThemeSettings, &updated.Status, &updated.UpdatedAt,
	)
	if err != nil {
		logger.Fatal(constants.DBConnectionFailure, err, logCtx)
		return nil, err
	}
	return &updated, nil
}

func (ar *admin_repository) UpdateFranchiseStatus(ctx context.Context, id string, status string) (string, error) {
	logCtx := logger.BaseLogContext(layer, constants.Repository, method, constants.Methods.UpdateFranchiseStatus)

	// Check database connection
	if err := dbutils.CheckDBConn(ar.db, logCtx); err != nil {
		logger.Error("Database connection check failed", err, logCtx)
		return "", err
	}

	// Prepare update columns and condition for the query
	columns := []string{"status", "updated_at"}
	condition := fmt.Sprintf("id = '%s' AND deleted_at IS NULL", id)

	// Define options for returning the updated ID
	opts := &dbutils.QueryBuilderOptions{
		Returning: []string{"id"},
	}

	// Generate the update query using the BuildUpdateQuery helper function
	query, err := dbutils.BuildUpdateQuery("outlet", "franchise", columns, condition, opts)
	if err != nil {
		logger.Error("Error building update query", err, logCtx)
		return "", err
	}

	// Execute the query and scan the returned ID
	var updatedID string
	args := []any{status, "NOW()"} // Values for the update query parameters
	if err := dbutils.ExecuteAndScanRow(ctx, ar.db, query, args, &updatedID); err != nil {
		logger.Error("Failed to execute update query for franchise status", err, logCtx)
		return "", err
	}

	// Return the updated ID
	return updatedID, nil
}

func (ar *admin_repository) DeleteFranchise(ctx context.Context, id string) (bool, error) {
	logCtx := logger.BaseLogContext(layer, constants.Repository, method, constants.Methods.DeleteFranchise)
	if err := dbutils.CheckDBConn(ar.db, logCtx); err != nil {
		return false, err
	}

	query := `UPDATE outlet.franchise SET deleted_at = NOW() WHERE id = $1`
	_, err := ar.db.ExecContext(ctx, query, id)
	if err != nil {
		logger.Fatal(constants.DBConnectionFailure, err, logCtx)
		return false, err
	}
	return true, nil
}

func (ar *admin_repository) GetFranchiseByID(ctx context.Context, id string) (*model.FranchiseResponse, error) {
	logCtx := logger.BaseLogContext(layer, constants.Repository, method, constants.Methods.GetFranchiseByID)

	// Check if the database connection is valid
	if err := dbutils.CheckDBConn(ar.db, logCtx); err != nil {
		logger.Error("Database connection check failed", err, logCtx)
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
			Schemas: []string{"outlet"},
			Tables:  []string{"franchise"},
			Columns: columns,
		},
	}

	// Use the BuildSelectQuery helper function to build the query
	query, args, err := dbutils.BuildSelectQuery("outlet", "franchise", columns, conditions, opts)
	if err != nil {
		logger.Error("Error building select query", err, logCtx)
		return nil, err
	}

	// Execute query and scan the result into the model
	var franchise model.FranchiseResponse
	if err := dbutils.ExecuteAndScanRow(ctx, ar.db, query, args,
		&franchise.ID, &franchise.BusinessName, &franchise.LogoURL, &franchise.SubDomain,
		&franchise.ThemeSettings, &franchise.Status, &franchise.CreatedAt, &franchise.UpdatedAt,
	); err != nil {
		logger.Error("Failed to execute query and scan result", err, logCtx)
		return nil, err
	}

	return &franchise, nil
}

func (ar *admin_repository) GetAllFranchises(ctx context.Context, page int32, limit int32) ([]model.FranchiseResponse, error) {
	logCtx := logger.BaseLogContext(layer, constants.Repository, method, constants.Methods.GetAllFranchises)

	// Check if the database connection is valid
	if err := dbutils.CheckDBConn(ar.db, logCtx); err != nil {
		logger.Error("Database connection check failed", err, logCtx)
		return nil, err
	}

	offset := (page - 1) * limit

	// Define the columns you need to retrieve
	columns := []string{
		"id", "business_name", "logo_url", "sub_domain", "theme_settings",
		"status", "created_at",
	}

	// Define the conditions for WHERE clause
	conditions := map[string]any{
		"deleted_at": nil, // Filter deleted records
	}

	// Prepare query options (if any), you can add Returning or whitelist logic here
	opts := &dbutils.QueryBuilderOptions{
		Whilelist: struct {
			Schemas []string
			Tables  []string
			Columns []string
		}{
			Schemas: []string{"outlet"},
			Tables:  []string{"franchise"},
			Columns: columns,
		},
	}

	// Use the BuildSelectQuery helper function to build the query
	query, args, err := dbutils.BuildSelectQuery("outlet", "franchise", columns, conditions, opts)
	if err != nil {
		logger.Error("Error building select query", err, logCtx)
		return nil, err
	}

	// Add pagination (LIMIT and OFFSET)
	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
	args = append(args, limit, offset)

	// Execute the query and fetch rows
	rows, err := ar.db.QueryContext(ctx, query, args...)
	if err != nil {
		logger.Error("Failed to execute query", err, logCtx)
		return nil, err
	}
	defer rows.Close()

	// Scan the result into the model
	var franchises []model.FranchiseResponse
	for rows.Next() {
		var f model.FranchiseResponse
		if err := rows.Scan(&f.ID, &f.BusinessName, &f.LogoURL, &f.SubDomain, &f.ThemeSettings, &f.Status, &f.CreatedAt); err != nil {
			logger.Error("Failed to scan row", err, logCtx)
			return nil, err
		}
		franchises = append(franchises, f)
	}

	return franchises, nil
}

func (ar *admin_repository) GetFranchiseOwner(ctx context.Context, id string) (*model.FranchiseOwnerResponse, error) {
	logCtx := logger.BaseLogContext(layer, constants.Repository, method, constants.Methods.GetFranchiseOwner)

	// Check if the database connection is valid
	if err := dbutils.CheckDBConn(ar.db, logCtx); err != nil {
		logger.Error("Database connection check failed", err, logCtx)
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
			Schemas: []string{"outlet"},
			Tables:  []string{"owner"},
			Columns: columns,
		},
	}

	// Use BuildSelectQuery to build the query
	query, args, err := dbutils.BuildSelectQuery("outlet", "owner", columns, conditions, opts)
	if err != nil {
		logger.Error("Error building select query", err, logCtx)
		return nil, err
	}

	// Execute the query and scan the result into the model
	var owner model.FranchiseOwnerResponse
	err = dbutils.ExecuteAndScanRow(ctx, ar.db, query, args,
		&owner.ID, &owner.Name, &owner.Gender, &owner.Dob, &owner.MobileNo,
		&owner.Email, &owner.Address, &owner.AadharNo, &owner.IsVerified, &owner.CreatedAt,
	)
	if err != nil {
		logger.Error("Failed to execute query for franchise owner", err, logCtx)
		return nil, err
	}

	return &owner, nil
}

func (ar *admin_repository) GetFranchiseDocuments(ctx context.Context, id string) ([]model.FranchiseDocumentResponse, error) {
	logCtx := logger.BaseLogContext(layer, constants.Repository, method, constants.Methods.GetFranchiseDocuments)

	// Check database connection
	if err := dbutils.CheckDBConn(ar.db, logCtx); err != nil {
		logger.Error("Database connection check failed", err, logCtx)
		return nil, err
	}

	// Define the columns to retrieve
	columns := []string{"id", "document_type", "document_url", "uploaded_at"}

	// Define the conditions for WHERE clause
	conditions := map[string]any{"franchise_id": id}

	// Prepare query options (e.g., whitelist or other conditions if needed)
	opts := &dbutils.QueryBuilderOptions{
		Whilelist: struct {
			Schemas []string
			Tables  []string
			Columns []string
		}{
			Schemas: []string{"outlet"},
			Tables:  []string{"franchise_documents"},
			Columns: columns,
		},
	}

	// Use BuildSelectQuery to build the query
	query, args, err := dbutils.BuildSelectQuery("outlet", "franchise_documents", columns, conditions, opts)
	if err != nil {
		logger.Error("Failed to build SELECT query", err, logCtx, "franchise_id", id)
		return nil, err
	}

	// Execute the query and get rows
	rows, err := ar.db.QueryContext(ctx, query, args...)
	if err != nil {
		logger.Error("Query execution failed", err, logCtx, "query", query, "franchise_id", id)
		return nil, err
	}
	defer rows.Close()

	var docs []model.FranchiseDocumentResponse

	// Scan the rows into the response slice
	for rows.Next() {
		var doc model.FranchiseDocumentResponse
		err := rows.Scan(&doc.ID, &doc.DocumentType, &doc.DocumentURL, &doc.UploadedAt)
		if err != nil {
			logger.Error("Failed to scan row into document response", err, logCtx, "query", query, "franchise_id", id)
			return nil, err
		}
		docs = append(docs, doc)
	}

	// Log successful retrieval
	logger.Info("Successfully retrieved franchise documents", logCtx, "document_count", len(docs))

	return docs, nil
}

func (ar *admin_repository) GetFranchiseAccounts(ctx context.Context, id string) ([]model.FranchiseAccountResponse, error) {
	logCtx := logger.BaseLogContext(layer, constants.Repository, method, constants.Methods.GetFranchiseAccounts)

	// Check database connection
	if err := dbutils.CheckDBConn(ar.db, logCtx); err != nil {
		return nil, err
	}

	// Define the columns and conditions for the SELECT query
	// Define the columns to retrieve
	columns := []string{"id", "document_type", "document_url", "uploaded_at"}

	// Define the conditions for WHERE clause
	conditions := map[string]any{"franchise_id": id}

	// Prepare query options (e.g., whitelist or other conditions if needed)
	opts := &dbutils.QueryBuilderOptions{
		Whilelist: struct {
			Schemas []string
			Tables  []string
			Columns []string
		}{
			Schemas: []string{"outlet"},
			Tables:  []string{"franchise_documents"},
			Columns: columns,
		},
	}

	// Use BuildSelectQuery to build the query
	query, args, err := dbutils.BuildSelectQuery("outlet", "franchise_documents", columns, conditions, opts)
	if err != nil {
		logger.Error("Failed to build SELECT query", err, logCtx, "franchise_id", id)
		return nil, err
	}

	// Execute the query and get rows
	rows, err := ar.db.QueryContext(ctx, query, args...)
	if err != nil {
		logger.Error("Query execution failed", err, logCtx, "query", query, "franchise_id", id)
		return nil, err
	}
	defer rows.Close()

	var accounts []model.FranchiseAccountResponse

	// Scan the rows into the response slice
	for rows.Next() {
		var acc model.FranchiseAccountResponse
		err := rows.Scan(&acc.ID, &acc.EmployeeID, &acc.Name, &acc.MobileNo, &acc.AccountType, &acc.Status, &acc.CreatedAt)
		if err != nil {
			logger.Error("Failed to scan row into account response", err, logCtx, "query", query, "franchise_id", id)
			return nil, err
		}
		accounts = append(accounts, acc)
	}

	return accounts, nil
}
