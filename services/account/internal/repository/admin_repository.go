package repository

import (
	"context"
	"database/sql"
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

var outlet_schema = constants.DB.Schema_Outlet

// CFTC - Common Franchise Table Column
var CFTC = []string{
	constants.T_Fran.UUID,
	constants.T_Fran.Status,
	constants.T_Fran.BusinessName,
	constants.T_Fran.LogoUrl,
	constants.T_Fran.Subdomain,
	constants.T_Fran.ThemeSettings,
	constants.T_Fran.FranchiseOwnerID,
}

// CFTC - Common Owner Table Column
var COTC = []string{
	constants.T_Onr.UUID,
	constants.T_Onr.Name,
	constants.T_Onr.Gender,
	constants.T_Onr.DOB,
	constants.T_Onr.MobileNo,
	constants.T_Onr.Email,
	constants.T_Onr.Address,
	constants.T_Onr.AadharNo,
	constants.T_Onr.IsVerified,
	constants.T_Onr.Status,
}

type AdminRepository interface {
	CreateNewOwner(ctx context.Context, owner *model.FranchiseOwner) (*model.AddResponse, error)
	UpdateNewOwner(ctx context.Context, owner *model.FranchiseOwner) (*model.UpdateResponse, error)

	CreateFranchise(ctx context.Context, franchise *model.Franchise) (*model.AddResponse, error)
	UpdateFranchise(ctx context.Context, franchise *model.Franchise) (*model.UpdateResponse, error)
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

	opts := &dbutils.QueryBuilderOptions{
		Returning: []string{constants.T_Onr.UUID, constants.T_Onr.CreatedAt},
		Whitelist: struct {
			Schemas []string
			Tables  []string
			Columns []string
		}{
			Schemas: []string{outlet_schema},
			Tables:  []string{table},
			Columns: append(COTC, constants.T_Onr.CreatedAt),
		},
	}

	// Build query using join-aware builder
	query, err := dbutils.BuildInsertQuery(
		method,
		outlet_schema,
		table,
		COTC,
		opts,
	)
	if err != nil {
		return nil, err
	}
	// Prepare and execute the query
	values, err := dbutils.MapValuesDirect(owner, COTC, "json")
	if err != nil {
		return nil, err
	}

	repsonse := &model.AddResponse{}
	err = dbutils.ExecuteAndScanRow(ctx, method, ar.db, query, values,
		repsonse,
		opts.Returning...,
	)
	if err != nil {
		return nil, err
	}
	return repsonse, nil
}

func (ar *admin_repository) UpdateNewOwner(ctx context.Context, owner *model.FranchiseOwner) (*model.UpdateResponse, error) {
	var (
		method = constants.Methods.CreateOwner
		table  = constants.DB.Table_Owner
	)
	if err := dbutils.CheckDBConn(ar.db, method); err != nil {
		return nil, err
	}

	conditions := map[string]any{
		constants.T_Onr.UUID: owner.ID,
	}

	// Whitelist for safe updating
	opts := &dbutils.QueryBuilderOptions{
		Returning: []string{constants.T_Onr.UUID, constants.T_Onr.UpdatedAt},
		Whitelist: struct {
			Schemas []string
			Tables  []string
			Columns []string
		}{
			Schemas: []string{outlet_schema},
			Tables:  []string{table},
			Columns: append(COTC, constants.T_Onr.UUID, constants.T_Onr.UpdatedAt),
		},
	}

	query, args, err := dbutils.BuildUpdateQuery(method, outlet_schema, table, COTC, conditions, opts)
	if err != nil {
		return nil, err
	}
	values, err := dbutils.MapValuesDirect(owner, COTC, "json")
	if err != nil {
		logger.Error("failed to map update values", err, nil)
		return nil, err
	}
	copy(args[:len(COTC)], values)

	var updated model.UpdateResponse
	err = dbutils.ExecuteAndScanRow(ctx, method, ar.db, query, values,
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
	columns := []string{
		constants.T_Fran.UUID,
		constants.T_Fran.Status,
		constants.T_Fran.BusinessName,
		constants.T_Fran.LogoUrl,
		constants.T_Fran.Subdomain,
		constants.T_Fran.ThemeSettings,
		constants.T_Fran.FranchiseOwnerID,
	}
	// Insert into franchise table
	opts := &dbutils.QueryBuilderOptions{
		Returning: []string{constants.T_Fran.UUID, constants.T_Fran.CreatedAt},
		Whitelist: struct {
			Schemas []string
			Tables  []string
			Columns []string
		}{
			Schemas: []string{outlet_schema},
			Tables:  []string{table},
			Columns: append(columns, constants.T_Fran.CreatedAt),
		},
	}

	query, err := dbutils.BuildInsertQuery(method, outlet_schema, table, columns, opts)
	if err != nil {
		return nil, err
	}

	values, err := dbutils.MapValuesDirect(fInput, columns, "json")
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

func (ar *admin_repository) UpdateFranchise(ctx context.Context, franchise *model.Franchise) (*model.UpdateResponse, error) {
	var method = constants.Methods.UpdateFranchise
	var table = constants.DB.Table_Franchise
	if err := dbutils.CheckDBConn(ar.db, method); err != nil {
		return nil, err
	}

	columns := []string{
		constants.T_Fran.UUID,
		constants.T_Fran.BusinessName,
		constants.T_Fran.LogoUrl,
		constants.T_Fran.Subdomain,
		constants.T_Fran.ThemeSettings,
		constants.T_Fran.Status,
		constants.T_Fran.FranchiseOwnerID,
	}

	opts := &dbutils.QueryBuilderOptions{
		Returning: []string{constants.T_Fran.UUID, constants.T_Fran.UpdatedAt},
	}
	opts.Whitelist.Schemas = []string{outlet_schema}
	opts.Whitelist.Tables = []string{table}
	opts.Whitelist.Columns = append(columns, constants.T_Fran.UpdatedAt)

	// Define the condition map for WHERE clause
	conditions := map[string]any{
		constants.T_Fran.UUID: franchise.ID,
	}

	query, args, err := dbutils.BuildUpdateQuery(
		method,
		outlet_schema,
		table,
		columns,
		conditions,
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
	columns := []string{constants.T_Fran.Status, constants.T_Fran.UpdatedAt}
	condition := map[string]any{constants.T_Fran.UUID: id} // Use map for conditions

	// Set up whitelist and returning options
	opts := &dbutils.QueryBuilderOptions{
		Returning: []string{constants.T_Fran.UUID, constants.T_Fran.UpdatedAt},
	}
	opts.Whitelist.Schemas = []string{outlet_schema}
	opts.Whitelist.Tables = []string{table}
	opts.Whitelist.Columns = append(columns, constants.T_Fran.UUID)

	// Generate the update query using the BuildUpdateQuery helper function
	query, args, err := dbutils.BuildUpdateQuery(
		method,
		outlet_schema,
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
	var updated_franchise *model.UpdateResponse
	if err := dbutils.ExecuteAndScanRow(ctx, method, ar.db, query, args, &updated_franchise, opts.Returning...); err != nil {
		return nil, err
	}
	return updated_franchise, nil
}

func (ar *admin_repository) DeleteFranchise(ctx context.Context, id string) (*model.DeletedResponse, error) {
	var method = constants.Methods.DeleteFranchise
	var table = constants.DB.Table_Franchise
	if err := dbutils.CheckDBConn(ar.db, method); err != nil {
		return nil, err
	}

	// Columns to update
	columns := []string{constants.T_Fran.DeletedAt}

	// WHERE condition
	condition := map[string]any{constants.T_Fran.UUID: id}

	// Options: whitelist and returning (optional, here we skip RETURNING)
	opts := &dbutils.QueryBuilderOptions{
		Returning: []string{constants.T_Fran.UUID, constants.T_Fran.DeletedAt},
	}
	opts.Whitelist.Schemas = []string{outlet_schema}
	opts.Whitelist.Tables = []string{table}
	opts.Whitelist.Columns = append(columns, constants.T_Fran.UUID)

	// Build the update query
	query, args, err := dbutils.BuildUpdateQuery(
		method,
		outlet_schema,
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
	var deleted *model.DeletedResponse
	if err := dbutils.ExecuteAndScanRow(ctx, constants.Methods.DeleteFranchise, ar.db, query, args, &deleted, opts.Returning...); err != nil {
		return nil, err
	}
	return &model.DeletedResponse{
		ID:        id,
		DeletedAt: time.Now(),
	}, nil
}

func (ar *admin_repository) GetAllFranchises(ctx context.Context, page int32, limit int32) ([]model.FranchiseResponse, error) {
	var method = constants.Methods.GetAllFranchises
	var table = constants.DB.Table_Franchise

	if err := dbutils.CheckDBConn(ar.db, method); err != nil {
		return nil, err
	}

	offset := (page - 1) * limit
	columns := append(CFTC, constants.T_Fran.CreatedAt, constants.T_Fran.UpdatedAt)

	// WHERE conditions
	conditions := map[string]any{
		"deleted_at__null": true, // soft delete filter
	}

	// Build query options with whitelist
	opts := &dbutils.QueryBuilderOptions{
		Whitelist: struct {
			Schemas []string
			Tables  []string
			Columns []string
		}{
			Schemas: []string{outlet_schema},
			Tables:  []string{table},
			Columns: append(columns, constants.T_Fran.DeletedAt),
		},
		OrderBy: []string{"created_at DESC"},
		Limit:   limit,
		Offset:  offset,
	}

	// Build SELECT query using helper
	query, args, err := dbutils.BuildSelectQuery(
		method,
		outlet_schema,
		table,
		columns,
		conditions,
		opts,
	)
	if err != nil {
		return nil, err
	}

	// Execute the query
	rows, err := ar.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Scan result rows
	var franchises []model.FranchiseResponse
	for rows.Next() {
		var f model.FranchiseResponse
		f.ThemeSettings = make(map[string]interface{})
		if err := dbutils.ExecuteAndScanRow(ctx, method, ar.db, query, args, &f, columns...); err != nil {
			logger.Error("Failed to scan GetAllFranchises row", err, nil)
			return nil, err
		}
		franchises = append(franchises, f)
	}

	return franchises, nil
}
