package repository

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"time"

	"github.com/ashish19912009/zrms/services/account/internal/constants"
	"github.com/ashish19912009/zrms/services/account/internal/dbutils"
	"github.com/ashish19912009/zrms/services/account/internal/logger"
	"github.com/ashish19912009/zrms/services/account/internal/model"
	_ "github.com/lib/pq"
)

/**
go:generate mockery --name=Repository --output=internal/repository/mocks --case=underscore
mockery --name=Repository --dir=services/account/internal/repository --output=services/account/internal/repository/mocks --case=underscore
**/

var schema = constants.DB.Schema_Outlet

const (
	uuid               = "id"
	franchise_id       = "franchise_id"
	business_name      = "business_name"
	logo_url           = "logo_url"
	subdomain          = "sub_domain"
	theme_settings     = "theme_settings"
	status             = "status"
	created_at         = "created_at"
	updated_at         = "updated_at"
	deleted_at         = "deleted_at"
	franchise_owner_id = "franchise_owner_id"
)

type AdminRepository interface {
	CreateNewOwner(ctx context.Context, owner *model.FranchiseOwner) (*model.AddResponse, error)
	UpdateNewOwner(ctx context.Context, id string, owner *model.FranchiseOwner) (*model.UpdateResponse, error)

	CreateFranchise(ctx context.Context, franchise *model.Franchise) (*model.AddResponse, error)
	UpdateFranchise(ctx context.Context, id string, franchise *model.Franchise) (*model.UpdateResponse, error)
	UpdateFranchiseStatus(ctx context.Context, id string, status string) (*model.UpdateResponse, error)
	DeleteFranchise(ctx context.Context, id string) (*model.DeletedResponse, error)

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

func (ar *admin_repository) CreateNewOwner(ctx context.Context, owner *model.FranchiseOwner) (*model.AddResponse, error) {
	var (
		method = constants.Methods.CreateOwner
		table  = constants.DB.Table_Owner
	)
	if err := dbutils.CheckDBConn(ar.db, method); err != nil {
		return nil, err
	}

	columns := []string{
		"id",
		"name",
		"gender",
		"dob",
		"mobile_no",
		"email",
		"address",
		"aadhar_no",
		"is_verified",
		"status",
		"created_at",
	}
	//sort.Strings(columns)
	// Whitelist options
	opts := &dbutils.QueryBuilderOptions{
		Returning: columns,
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

	values, err := dbutils.StructToValuesByTag(owner, columns, "json")
	if err != nil {
		return nil, err
	}

	// Build query using join-aware builder
	query, err := dbutils.BuildInsertQuery(
		method,
		schema,
		table,
		columns,
		opts,
	)
	if err != nil {
		return nil, err
	}

	// return account, nil
	var repsonse *model.AddResponse
	err = dbutils.ExecuteAndScanRow(ctx, method, ar.db, query, values,
		&repsonse,
	)
	if err != nil {
		return nil, err
	}
	return repsonse, nil
	//return owner.ToResponse(createdID, created_at), nil
}

func (ar *admin_repository) UpdateNewOwner(ctx context.Context, id string, owner *model.FranchiseOwner) (*model.UpdateResponse, error) {
	var (
		method = constants.Methods.CreateOwner
		table  = constants.DB.Table_Owner
	)
	if err := dbutils.CheckDBConn(ar.db, method); err != nil {
		return nil, err
	}
	// Fields to update
	data := map[string]any{
		"id":          id,
		"name":        owner.Name,
		"gender":      owner.Gender,
		"dob":         owner.Dob,
		"mobile_no":   owner.MobileNo,
		"email":       owner.Email,
		"address":     owner.Address,
		"aadhar_no":   owner.AadharNo,
		"is_verified": owner.IsVerified,
		"status":      owner.Status,
	}

	columns := make([]string, 0, len(data))
	for col := range data {
		columns = append(columns, col)
	}
	columns = append(columns, "updated_at")
	sort.Strings(columns)
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
		Returning: columns,
	}

	query, args, err := dbutils.BuildUpdateQuery(method, schema, table, columns, conditions, opts)
	if err != nil {
		return nil, err
	}

	var updated model.UpdateResponse
	err = dbutils.ExecuteAndScanRow(ctx, method, ar.db, query, args,
		&updated,
	)
	if err != nil {
		return nil, err
	}

	return &updated, nil
}

func (ar *admin_repository) CreateFranchise(ctx context.Context, fInput *model.Franchise) (*model.AddResponse, error) {
	var (
		method = constants.Methods.CreateFranchise
		table  = constants.DB.Table_Franchise
	)

	if err := dbutils.CheckDBConn(ar.db, method); err != nil {
		return nil, err
	}
	columns := []string{uuid, logo_url, theme_settings, subdomain, status, business_name, franchise_owner_id}
	// Insert into franchise table
	opts := &dbutils.QueryBuilderOptions{
		Returning: []string{uuid, created_at},
		Whilelist: struct {
			Schemas []string
			Tables  []string
			Columns []string
		}{
			Schemas: []string{schema},
			Tables:  []string{table},
			Columns: append(columns, created_at),
		},
	}

	values, err := dbutils.StructToValuesByTag(fInput, columns, "json")
	if err != nil {
		return nil, err
	}

	query, err := dbutils.BuildInsertQuery(method, schema, table, columns, opts)
	if err != nil {
		return nil, err
	}

	var franchise model.AddResponse
	if err = dbutils.ExecuteAndScanRow(ctx, method, ar.db, query,
		values,
		&franchise,
		opts.Returning...,
	); err != nil {
		return nil, err
	}

	return &franchise, nil
}

func (ar *admin_repository) UpdateFranchise(ctx context.Context, id string, franchise *model.Franchise) (*model.UpdateResponse, error) {
	var method = constants.Methods.UpdateFranchise
	var table = constants.DB.Table_Franchise
	if err := dbutils.CheckDBConn(ar.db, method); err != nil {
		return nil, err
	}

	columns := []string{
		business_name, subdomain, logo_url, theme_settings, status,
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

	var updated model.UpdateResponse
	err = dbutils.ExecuteAndScanRow(ctx, constants.Methods.UpdateFranchise, ar.db, query,
		args,
		&updated,
		opts.Returning...,
	)
	if err != nil {
		return nil, err
	}
	return &updated, nil
}

func (ar *admin_repository) UpdateFranchiseStatus(ctx context.Context, id string, status string) (*model.UpdateResponse, error) {
	var method = constants.Methods.UpdateFranchiseStatus
	var table = constants.DB.Table_Franchise
	// Check database connection
	if err := dbutils.CheckDBConn(ar.db, method); err != nil {
		return nil, err
	}

	// Prepare columns and condition
	columns := []string{"status", "updated_at"}
	condition := map[string]any{"id": id} // Use map for conditions

	// Set up whitelist and returning options
	opts := &dbutils.QueryBuilderOptions{
		Returning: append(columns, "id"),
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
		return nil, err
	}

	// Prepare actual values matching the order of columns + conditions
	args = []any{status, time.Now(), id} // Manually prepare args matching placeholders

	// Execute the query and scan returned ID
	var updated_franchise *model.FranchiseResponse
	if err := dbutils.ExecuteAndScanRow(ctx, method, ar.db, query, args, &updated_franchise); err != nil {
		return nil, err
	}
	if updated_franchise.ID == id && updated_franchise.Status == status {
		return &model.UpdateResponse{
			ID:        updated_franchise.ID,
			UpdatedAt: *updated_franchise.UpdatedAt,
		}, nil
	}
	return nil, nil
}

func (ar *admin_repository) DeleteFranchise(ctx context.Context, id string) (*model.DeletedResponse, error) {
	var method = constants.Methods.DeleteFranchise
	var table = constants.DB.Table_Franchise
	if err := dbutils.CheckDBConn(ar.db, method); err != nil {
		return nil, err
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
		return nil, err
	}

	// Provide actual args: deleted_at value and id
	args = []any{time.Now(), id}

	// Execute the query
	var deletedID string
	if err := dbutils.ExecuteAndScanRow(ctx, constants.Methods.DeleteFranchise, ar.db, query, args, &deletedID); err != nil {
		return nil, err
	}
	return &model.DeletedResponse{
		ID:        id,
		DeletedAt: time.Now(),
	}, nil
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
