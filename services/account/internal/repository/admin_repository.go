package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

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

	GetAllFranchises(ctx context.Context, page int32, limit int32) ([]model.FranchiseResponse, error)
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
	var (
		method = constants.Methods.CreateFranchise
		table  = constants.DB.Table_Franchise
	)

	if err := dbutils.CheckDBConn(ar.db, method); err != nil {
		return nil, err
	}

	// Start transaction
	tx, err := ar.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf(constants.FailedToBeginTransaction, err)
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p) // re-throw panic after rollback
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	// Insert into franchise table
	opts := &dbutils.QueryBuilderOptions{
		Returning: []string{id, business_name, logo_url, sub_domain, theme_settings, status, created_at},
	}
	opts.Whilelist.Schemas = []string{schema}
	opts.Whilelist.Tables = []string{table}
	opts.Whilelist.Columns = []string{id, business_name, logo_url, sub_domain, theme_settings, status, created_at, updated_at, deleted_at}

	query, err := dbutils.BuildInsertQuery(method, schema, table, []string{business_name, logo_url, sub_domain, theme_settings, status}, opts)
	if err != nil {
		return nil, err
	}

	var franchise model.FranchiseResponse
	if err = dbutils.ExecuteAndScanRowTx(ctx, method, tx, query,
		[]any{id, franchiseInput.BusinessName, franchiseInput.LogoURL, franchiseInput.SubDomain, franchiseInput.ThemeSettings, franchiseInput.Status},
		&franchise.ID, &franchise.BusinessName, &franchise.LogoURL, &franchise.SubDomain, &franchise.ThemeSettings,
		&franchise.Status, &franchise.CreatedAt,
	); err != nil {
		return nil, err
	}

	// Insert into franchise owner table
	opts.Returning = []string{id, created_at}
	opts.Whilelist.Tables = []string{constants.DB.Table_Owner}
	opts.Whilelist.Columns = append(opts.Whilelist.Columns, franchise_id, name, gender, dob, mobile_no, email, address, aadhar_no, is_verified)

	query, err = dbutils.BuildInsertQuery(method, schema, constants.DB.Table_Owner, []string{
		franchise_id, name, gender, dob, mobile_no, email, address, aadhar_no, is_verified,
	}, opts)
	if err != nil {
		return nil, err
	}

	var ownerID string
	if err = dbutils.ExecuteAndScanRowTx(ctx, method, tx, query,
		[]any{franchise.ID, f_owner.Name, f_owner.Gender, f_owner.Dob, f_owner.MobileNo, f_owner.Email, f_owner.Address, f_owner.AadharNo, f_owner.IsVerified},
		&ownerID,
	); err != nil {
		return nil, err
	}

	return &franchise, nil
}

func (ar *admin_repository) UpdateFranchise(ctx context.Context, id string, franchise *model.Franchise) (*model.FranchiseResponse, error) {
	var method = constants.Methods.UpdateFranchise
	var table = constants.DB.Table_Franchise
	if err := dbutils.CheckDBConn(ar.db, method); err != nil {
		return nil, err
	}

	columns := []string{
		business_name, logo_url, sub_domain, theme_settings, status,
	}

	opts := &dbutils.QueryBuilderOptions{
		Returning: append(columns, "id", "updated_at"),
	}
	opts.Whilelist.Schemas = []string{schema}
	opts.Whilelist.Tables = []string{table}
	opts.Whilelist.Columns = append(columns, "id", "updated_at")

	query, args, err := dbutils.BuildUpdateQuery(
		method,
		schema,
		table,
		append(columns, "id"), map[string]any{"id": id},
		opts,
	)
	if err != nil {
		return nil, err
	}

	var updated model.FranchiseResponse
	err = dbutils.ExecuteAndScanRow(ctx, constants.Methods.UpdateFranchise, ar.db, query,
		args,
		&updated.ID, &updated.BusinessName, &updated.LogoURL, &updated.SubDomain, &updated.ThemeSettings, &updated.Status, &updated.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &updated, nil
}

func (ar *admin_repository) UpdateFranchiseStatus(ctx context.Context, id string, status string) (string, error) {
	var method = constants.Methods.UpdateFranchiseStatus
	var table = constants.DB.Table_Franchise
	// Check database connection
	if err := dbutils.CheckDBConn(ar.db, method); err != nil {
		return "", err
	}

	// Prepare columns and condition
	columns := []string{"status", "updated_at"}
	condition := map[string]any{"id": id} // Use map for conditions

	// Set up whitelist and returning options
	opts := &dbutils.QueryBuilderOptions{
		Returning: []string{"id"},
	}
	opts.Whilelist.Schemas = []string{schema}
	opts.Whilelist.Tables = []string{table}
	opts.Whilelist.Columns = append(columns, "id")

	// Generate the update query using the BuildUpdateQuery helper function
	query, args, err := dbutils.BuildUpdateQuery(
		method,
		schema,
		table,
		columns,
		condition,
		opts,
	)
	if err != nil {
		return "", err
	}

	// Prepare actual values matching the order of columns + conditions
	args = []any{status, time.Now(), id} // Manually prepare args matching placeholders

	// Execute the query and scan returned ID
	var updatedID string
	if err := dbutils.ExecuteAndScanRow(ctx, method, ar.db, query, args, &updatedID); err != nil {
		return "", err
	}

	return updatedID, nil
}

func (ar *admin_repository) DeleteFranchise(ctx context.Context, id string) (bool, error) {
	var method = constants.Methods.DeleteFranchise
	var table = constants.DB.Table_Franchise
	if err := dbutils.CheckDBConn(ar.db, method); err != nil {
		return false, err
	}

	// Columns to update
	columns := []string{"deleted_at"}

	// WHERE condition
	condition := map[string]any{"id": id}

	// Options: whitelist and returning (optional, here we skip RETURNING)
	opts := &dbutils.QueryBuilderOptions{}
	opts.Whilelist.Schemas = []string{schema}
	opts.Whilelist.Tables = []string{table}
	opts.Whilelist.Columns = append(columns, "id")

	// Build the update query
	query, args, err := dbutils.BuildUpdateQuery(
		method,
		schema,
		table,
		columns,
		condition,
		opts,
	)
	if err != nil {
		return false, err
	}

	// Provide actual args: deleted_at value and id
	args = []any{time.Now(), id}

	// Execute the query
	var deletedID string
	if err := dbutils.ExecuteAndScanRow(ctx, constants.Methods.DeleteFranchise, ar.db, query, args, &deletedID); err != nil {
		return false, err
	}
	return true, nil
}

func (ar *admin_repository) GetAllFranchises(ctx context.Context, page int32, limit int32) ([]model.FranchiseResponse, error) {
	// Step 1: DB Connection Check
	var method = constants.Methods.GetAllFranchises
	var table = constants.DB.Table_Franchise
	if err := dbutils.CheckDBConn(ar.db, method); err != nil {
		return nil, err
	}

	// Step 2: Calculate offset for pagination
	offset := (page - 1) * limit

	// Step 3: Define SELECT columns
	columns := []string{
		"id", "business_name", "logo_url", "sub_domain",
		"theme_settings", "status", "created_at",
	}

	// Step 4: WHERE conditions
	conditions := map[string]any{
		"deleted_at": nil, // soft delete filter
	}

	// Step 5: Build query options with whitelist
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

	// Step 6: Build SELECT query using helper
	query, args, err := dbutils.BuildSelectQuery(
		method,
		schema,
		table,
		columns,
		conditions,
		opts,
	)
	if err != nil {
		return nil, err
	}

	// Step 7: Append ORDER BY, LIMIT, OFFSET
	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
	args = append(args, limit, offset)

	// Step 8: Execute the query
	rows, err := ar.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Step 9: Scan result rows
	var franchises []model.FranchiseResponse
	for rows.Next() {
		var f model.FranchiseResponse
		if err := rows.Scan(
			&f.ID, &f.BusinessName, &f.LogoURL, &f.SubDomain,
			&f.ThemeSettings, &f.Status, &f.CreatedAt,
		); err != nil {
			logger.Error("Failed to scan GetAllFranchises row", err, nil)
			return nil, err
		}
		franchises = append(franchises, f)
	}

	return franchises, nil
}
