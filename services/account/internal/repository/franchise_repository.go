package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ashish19912009/zrms/services/account/internal/constants"
	"github.com/ashish19912009/zrms/services/account/internal/dbutils"
	"github.com/ashish19912009/zrms/services/account/internal/logger"
	"github.com/ashish19912009/zrms/services/account/internal/model"
	_ "github.com/lib/pq"
)

/**
go:generate mockery --name=Repository --output=services/account/internal/repository/mocks --case=underscore
mockery --name=Repository --dir=services/account/internal/repository --output=services/account/internal/repository/mocks --case=underscore
**/

// CFATC - Common Franchise Account Table Columns
var CFATC = []string{
	constants.Acc.FranchiseID,
	constants.Acc.EmpID,
	constants.Acc.LoginID,
	constants.Acc.Password,
	constants.Acc.AccountType,
	constants.Acc.Name,
	constants.Acc.MobileNo,
	constants.Acc.Email,
	constants.Acc.RoleID,
	constants.Acc.Status,
}

// CFDTC - Common Franchise Document Table Columns
var CFDTC = []string{
	constants.F_doc.FranchiseID,
	constants.F_doc.DocumentTypeID,
	constants.F_doc.DocumentURL,
	constants.F_doc.UploadedBy,
	constants.F_doc.Status,
	constants.F_doc.Remark,
	constants.F_doc.VerifiedAt,
}

// CFDTC - Common Franchise Address Table Columns
var CFAddrTC = []string{
	constants.F_addr.FranchiseID,
	constants.F_addr.AddressLine,
	constants.F_addr.City,
	constants.F_addr.State,
	constants.F_addr.Country,
	constants.F_addr.Pincode,
	constants.F_addr.Latitude,
	constants.F_addr.Longitude,
	constants.F_addr.IsVerified,
}

// CFRTC - Common Franchise Role Table Columns
var CFRTC = []string{
	constants.F_role.FranchiseID,
	constants.F_role.Name,
	constants.F_role.Description,
	constants.F_role.IsDefault,
}

var CRPTC = []string{
	constants.F_Role_Per.RoleID,
	constants.F_Role_Per.PermissionID,
}

type Repository interface {
	GetFranchiseByID(ctx context.Context, id string) (*model.FranchiseResponse, error)
	GetFranchiseByBusinessName(ctx context.Context, b_name string) (*model.FranchiseResponse, error)
	GetFranchiseOwnerByID(ctx context.Context, id string) (*model.FranchiseOwnerResponse, error)
	CheckIfOwnerExistsByAadharID(ctx context.Context, id string) (bool, error)

	CreateFranchiseAccount(ctx context.Context, account *model.FranchiseAccount) (*model.AddResponse, error)
	UpdateFranchiseAccount(ctx context.Context, id string, account *model.FranchiseAccount) (*model.UpdateResponse, error)
	GetFranchiseAccountByID(ctx context.Context, id string) (*model.FranchiseAccountResponse, error)
	GetAllFranchiseAccounts(ctx context.Context, fran *model.GetFranchisesRequest) ([]model.FranchiseAccountResponse, error)

	AddFranchiseDocument(ctx context.Context, doc *model.FranchiseDocument) (*model.AddResponse, error)
	UpdateFranchiseDocument(ctx context.Context, id string, doc *model.FranchiseDocument) (*model.UpdateResponse, error)
	GetAllFranchiseDocuments(ctx context.Context, id string) ([]model.FranchiseDocumentResponseComplete, error)

	AddFranchiseAddress(ctx context.Context, addr *model.FranchiseAddress) (*model.AddResponse, error)
	UpdateFranchiseAddress(ctx context.Context, id string, addr *model.FranchiseAddress) (*model.UpdateResponse, error)
	GetFranchiseAddressByID(ctx context.Context, id string) (*model.FranchiseAddressResponse, error)

	AddFranchiseRole(ctx context.Context, role *model.FranchiseRole) (*model.AddResponse, error)
	UpdateFranchiseRole(ctx context.Context, id string, role *model.FranchiseRole) (*model.UpdateResponse, error)
	GetAllFranchiseRoles(ctx context.Context, id string) ([]model.FranchiseRoleResponse, error)

	AddPermissionsToRole(ctx context.Context, pRole *model.RoleToPermissions) (*model.RoleToPermissions, error)
	UpdatePermissionsToRole(ctx context.Context, pRole *model.RoleToPermissions) (*model.RoleToPermissions, error)
	GetAllPermissionsToRole(ctx context.Context, id string) ([]model.RoleToPermissionsComplete, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (ar *repository) GetFranchiseByID(ctx context.Context, id string) (*model.FranchiseResponse, error) {
	var method = constants.Methods.GetFranchiseByID
	var table = constants.DB.Table_Franchise
	if err := dbutils.CheckDBConn(ar.db, method); err != nil {
		return nil, err
	}

	// Define the columns you need to retrieve
	columns := append(CFTC, constants.T_Fran.CreatedAt, constants.T_Fran.UpdatedAt)

	// Define the condition map for WHERE clause
	conditions := map[string]any{
		"id":               id,
		"deleted_at__null": true, // Filter deleted records
	}

	// Prepare query options (if any), you can add Returning or whitelist logic here
	opts := &dbutils.QueryBuilderOptions{
		Returning: columns,
		Whitelist: struct {
			Schemas []string
			Tables  []string
			Columns []string
		}{
			Schemas: []string{outlet_schema},
			Tables:  []string{table},
			Columns: append(columns, constants.T_Fran.DeletedAt),
		},
	}

	// Use the BuildSelectQuery helper function to build the query
	query, args, err := dbutils.BuildSelectQuery(method, outlet_schema, table, columns, conditions, opts)
	if err != nil {
		return nil, err
	}
	// Execute query and scan the result into the model
	franchise := &model.FranchiseResponse{
		ThemeSettings: make(map[string]interface{}),
	}
	if err := dbutils.ExecuteAndScanRow(ctx, method, ar.db, query, args,
		franchise,
		opts.Returning...,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("franchise not found")
		}
		return nil, fmt.Errorf("scan failed: %w", err)
	}
	return franchise, nil
}

func (ar *repository) GetFranchiseByBusinessName(ctx context.Context, b_name string) (*model.FranchiseResponse, error) {
	var method = constants.Methods.GetFranchiseByBusinessName
	var table = constants.DB.Table_Franchise
	if err := dbutils.CheckDBConn(ar.db, method); err != nil {
		return nil, err
	}

	// Define the columns you need to retrieve
	columns := append(CFTC, constants.T_Fran.CreatedAt, constants.T_Fran.UpdatedAt)

	// Define the condition map for WHERE clause
	conditions := map[string]any{
		constants.T_Fran.BusinessName: b_name,
		"deleted_at__null":            true, // Filter deleted records
	}

	// Prepare query options (if any), you can add Returning or whitelist logic here
	opts := &dbutils.QueryBuilderOptions{
		Returning: columns,
		Whitelist: struct {
			Schemas []string
			Tables  []string
			Columns []string
		}{
			Schemas: []string{outlet_schema},
			Tables:  []string{table},
			Columns: append(columns, constants.T_Fran.DeletedAt),
		},
	}
	// Use the BuildSelectQuery helper function to build the query
	query, args, err := dbutils.BuildSelectQuery(method, outlet_schema, table, columns, conditions, opts)
	if err != nil {
		return nil, err
	}

	// Execute query and scan the result into the model
	franchise := &model.FranchiseResponse{
		ThemeSettings: make(map[string]interface{}),
	}
	if err := dbutils.ExecuteAndScanRow(ctx, method, ar.db, query, args,
		franchise,
		opts.Returning...,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return franchise, nil
		}
		return nil, err
	}
	return franchise, nil
}

func (ar *repository) GetFranchiseOwnerByID(ctx context.Context, id string) (*model.FranchiseOwnerResponse, error) {
	var method = constants.Methods.GetFranchiseOwnerByID
	var table = constants.DB.Table_Owner
	if err := dbutils.CheckDBConn(ar.db, method); err != nil {
		return nil, err
	}

	// Define the columns you need to retrieve
	// columns := []string{
	// 	"id", "name", "gender", "dob", "mobile_no", "email",
	// 	"address", "aadhar_no", "is_verified", "created_at",
	// }
	columns := append(COTC, constants.T_Onr.CreatedAt, constants.T_Fran.UpdatedAt)

	// Define the conditions for WHERE clause
	conditions := map[string]any{
		constants.T_Onr.UUID: id,   // Use franchise_id for filtering
		"deleted_at__null":   true, // Filter deleted records
	}

	// Prepare query options (if any), you can add Returning or whitelist logic here
	opts := &dbutils.QueryBuilderOptions{
		Returning: columns,
		Whitelist: struct {
			Schemas []string
			Tables  []string
			Columns []string
		}{
			Schemas: []string{outlet_schema},
			Tables:  []string{table},
			Columns: append(columns, constants.T_Onr.DeletedAt),
		},
	}

	// Use BuildSelectQuery to build the query
	query, args, err := dbutils.BuildSelectQuery(method, outlet_schema, table, columns, conditions, opts)
	if err != nil {
		return nil, err
	}

	// Execute the query and scan the result into the model
	var owner model.FranchiseOwnerResponse
	err = dbutils.ExecuteAndScanRow(ctx, method, ar.db, query, args,
		&owner,
		opts.Returning...,
	)
	if err != nil {
		return nil, err
	}

	return &owner, nil
}

func (ar *repository) CheckIfOwnerExistsByAadharID(ctx context.Context, aadharNo string) (bool, error) {
	var method = constants.Methods.CheckIfOwnerExistsByAadharID
	var table = constants.DB.Table_Owner
	if err := dbutils.CheckDBConn(ar.db, method); err != nil {
		return false, err
	}

	// Define the columns you need to retrieve
	columns := append(COTC, constants.T_Onr.CreatedAt, constants.T_Onr.UpdatedAt)

	// Define the conditions for WHERE clause
	conditions := map[string]any{
		constants.T_Onr.AadharNo: aadharNo, // Use franchise_id for filtering
	}

	// Prepare query options (if any), you can add Returning or whitelist logic here
	opts := &dbutils.QueryBuilderOptions{
		Returning: columns,
		Whitelist: struct {
			Schemas []string
			Tables  []string
			Columns []string
		}{
			Schemas: []string{outlet_schema},
			Tables:  []string{table},
			Columns: columns,
		},
	}

	// Use BuildSelectQuery to build the query
	query, args, err := dbutils.BuildSelectQuery(method, outlet_schema, table, columns, conditions, opts)
	if err != nil {
		return false, err
	}
	// Execute the query and scan the result into the model
	owner := &model.FranchiseOwnerResponse{}
	err = dbutils.ExecuteAndScanRow(ctx, method, ar.db, query, args,
		owner,
		opts.Returning...,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (ar *repository) CreateFranchiseAccount(ctx context.Context, account *model.FranchiseAccount) (*model.AddResponse, error) {
	var method = constants.Methods.CreateFranchiseAccount
	var table = constants.DB.Table_Franchise_Accounts

	if err := dbutils.CheckDBConn(ar.db, method); err != nil {
		return nil, err
	}

	// Whitelist for safe inserting
	opts := &dbutils.QueryBuilderOptions{
		Returning: []string{constants.Acc.UUID, constants.Acc.CreatedAt},
		Whitelist: struct {
			Schemas []string
			Tables  []string
			Columns []string
		}{
			Schemas: []string{outlet_schema},
			Tables:  []string{table},
			Columns: append(CFATC, constants.Acc.CreatedAt),
		},
	}

	values, err := dbutils.MapValuesDirect(account, CFATC, "json")
	if err != nil {
		return nil, err
	}

	query, err := dbutils.BuildInsertQuery(method, outlet_schema, table, CFATC, opts)
	if err != nil {
		return nil, err
	}

	var newAccount *model.AddResponse
	err = dbutils.ExecuteAndScanRow(ctx, method, ar.db, query, values, &newAccount)
	if err != nil {
		return nil, err
	}

	return newAccount, nil
}

func (ar *repository) UpdateFranchiseAccount(ctx context.Context, id string, account *model.FranchiseAccount) (*model.UpdateResponse, error) {
	var method = constants.Methods.UpdateFranchiseAccount
	var table = constants.DB.Table_Franchise_Accounts

	if err := dbutils.CheckDBConn(ar.db, method); err != nil {
		return nil, err
	}

	conditions := map[string]any{
		constants.Acc.UUID: id,
	}

	// Whitelist for safe inserting
	opts := &dbutils.QueryBuilderOptions{
		Returning: []string{constants.Acc.UUID, constants.Acc.UpdatedAt},
		Whitelist: struct {
			Schemas []string
			Tables  []string
			Columns []string
		}{
			Schemas: []string{outlet_schema},
			Tables:  []string{table},
			Columns: append(CFATC, constants.Acc.CreatedAt, constants.Acc.UpdatedAt),
		},
	}

	query, args, err := dbutils.BuildUpdateQuery(method, outlet_schema, table, CFATC, conditions, opts)
	if err != nil {
		return nil, err
	}

	values, err := dbutils.MapValues(CFATC, account, "json")
	if err != nil {
		logger.Error("failed to map update values", err, nil)
		return nil, err
	}
	copy(args[:len(COTC)], values)

	var updated model.UpdateResponse
	err = dbutils.ExecuteAndScanRow(ctx, method, ar.db, query, values, &updated)
	if err != nil {
		return nil, err
	}

	return &updated, nil
}

func (ar *repository) GetFranchiseAccountByID(ctx context.Context, id string) (*model.FranchiseAccountResponse, error) {
	var method = constants.Methods.GetFranchiseAccountByID
	var table = constants.DB.Table_Franchise_Accounts
	// Check DB connection health
	if err := dbutils.CheckDBConn(ar.db, method); err != nil {
		return nil, err
	}

	conditions := map[string]any{
		constants.Acc.UUID: id,
	}
	allColumns := append(CFATC, constants.Acc.UUID, constants.Acc.CreatedAt, constants.Acc.UpdatedAt)
	opts := &dbutils.QueryBuilderOptions{
		Returning: allColumns,
		Whitelist: struct {
			Schemas []string
			Tables  []string
			Columns []string
		}{
			Schemas: []string{outlet_schema},
			Tables:  []string{table},
			Columns: allColumns,
		},
	}

	// Build SELECT query
	query, args, err := dbutils.BuildSelectQuery(method, outlet_schema, table, allColumns, conditions, opts)
	if err != nil {
		return nil, err
	}

	// Scan result into response struct
	var account model.FranchiseAccountResponse
	err = dbutils.ExecuteAndScanRow(ctx, method, ar.db, query, args, &account)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (ar *repository) GetAllFranchiseAccounts(ctx context.Context, fran *model.GetFranchisesRequest) ([]model.FranchiseAccountResponse, error) {
	var method = constants.Methods.GetAllFranchiseAccounts
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
			Schema: outlet_schema,
			Table:  constants.DB.Table_Roles,
			Alias:  "r",
			On:     "fa.role = r.id",
		},
	}

	// WHERE conditions
	conditions := map[string]any{
		"fa.franchise_id": fran.FranchiseID,
	}
	// Add query filter if provided
	if fran.GetPagination != nil && fran.GetPagination.Query != "" {
		conditions["search"] = map[string]string{
			"columns": "fa.name, fa.email, fa.mobile_no, fa.account_type, fa.role_id", // searchable columns
			"value":   fran.GetPagination.Query,
		}
	}
	// Whitelist for security
	opts := &dbutils.QueryBuilderOptions{
		Whitelist: struct {
			Schemas []string
			Tables  []string
			Columns []string
		}{
			Schemas: []string{outlet_schema},
			Tables:  []string{table, constants.DB.Table_Roles},
			Columns: columns,
		},
	}

	// Build the join query
	query, args, err := dbutils.BuildJoinSelectQuery(method, outlet_schema, table, "fa", columns, joins, conditions, opts)
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

func (ar *repository) AddFranchiseDocument(ctx context.Context, doc *model.FranchiseDocument) (*model.AddResponse, error) {
	var method = constants.Methods.AddFranchiseDocument
	var table = constants.DB.Table_Document_Types
	if err := dbutils.CheckDBConn(ar.db, method); err != nil {
		return nil, err
	}

	// Whitelist options
	opts := &dbutils.QueryBuilderOptions{
		Returning: []string{constants.F_doc.UUID, constants.F_doc.CreatedAt},
		Whitelist: struct {
			Schemas []string
			Tables  []string
			Columns []string
		}{
			Schemas: []string{outlet_schema},
			Tables:  []string{table},
			Columns: CFDTC,
		},
	}

	// Build query using join-aware builder
	query, err := dbutils.BuildInsertQuery(
		method,
		outlet_schema,
		table,
		CFDTC,
		opts,
	)
	if err != nil {
		return nil, err
	}
	// Scan result into response struct
	var account = &model.AddResponse{}
	err = dbutils.ExecuteAndScanRow(ctx, method, ar.db, query, []any{doc.FranchiseID, doc.DocumentTypeID, doc.DocumentURL, doc.UploadedBy, doc.Status, doc.Remark}, &account.ID)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (ar *repository) UpdateFranchiseDocument(ctx context.Context, id string, doc *model.FranchiseDocument) (*model.UpdateResponse, error) {
	var method = constants.Methods.UpdateFranchiseDocument
	var table = constants.DB.Table_Franchise_documents

	if err := dbutils.CheckDBConn(ar.db, method); err != nil {
		return nil, err
	}

	// Fields to update

	conditions := map[string]any{
		constants.F_doc.UUID: id,
	}
	allColumns := append(CFDTC, constants.F_doc.UUID, constants.F_doc.UpdatedAt)
	// Whitelist for safe updating
	opts := &dbutils.QueryBuilderOptions{
		Returning: []string{constants.F_doc.UUID, constants.F_doc.UpdatedAt},
		Whitelist: struct {
			Schemas []string
			Tables  []string
			Columns []string
		}{
			Schemas: []string{outlet_schema},
			Tables:  []string{table},
			Columns: allColumns,
		},
	}

	query, args, err := dbutils.BuildUpdateQuery(method, outlet_schema, table, CFDTC, conditions, opts)
	if err != nil {
		return nil, err
	}
	values, err := dbutils.MapValues(CFDTC, doc, "json")
	if err != nil {
		logger.Error("failed to map update values", err, nil)
		return nil, err
	}
	copy(args[:len(COTC)], values)

	var updatedDoc model.UpdateResponse
	err = dbutils.ExecuteAndScanRow(ctx, method, ar.db, query, values,
		&updatedDoc,
	)
	if err != nil {
		return nil, err
	}

	return &updatedDoc, nil
}

func (ar *repository) GetAllFranchiseDocuments(ctx context.Context, id string) ([]model.FranchiseDocumentResponseComplete, error) {
	var method = constants.Methods.GetAllFranchiseDocuments
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
			Schema: outlet_schema,
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
		Whitelist: struct {
			Schemas []string
			Tables  []string
			Columns []string
		}{
			Schemas: []string{outlet_schema},
			Tables:  []string{table},
			Columns: []string{"fd.id", "dt.name", "fd.document_url", "fd.uploaded_at"},
		},
	}

	// Build query using join-aware builder
	query, args, err := dbutils.BuildJoinSelectQuery(
		method,
		outlet_schema,
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
	var docs []model.FranchiseDocumentResponseComplete
	for rows.Next() {
		var doc model.FranchiseDocumentResponseComplete
		err := rows.Scan(&doc.ID, &doc.DocumentName, &doc.DocumentDescription, &doc.DocumentURL, &doc.UploadedBy, &doc.Status, &doc.Remark, &doc.IsMandate, &doc.VerifiedAt, &doc.UploadedAt)
		if err != nil {
			logger.Error("Failed to scan row", err, nil)
			return nil, err
		}
		docs = append(docs, doc)
	}
	return docs, nil
}

func (ar *repository) AddFranchiseAddress(ctx context.Context, addr *model.FranchiseAddress) (*model.AddResponse, error) {
	var method = constants.Methods.AddFranchiseAddress
	var table = constants.DB.Table_Franchise_addresses
	if err := dbutils.CheckDBConn(ar.db, method); err != nil {
		return nil, err
	}
	// Whitelist options
	opts := &dbutils.QueryBuilderOptions{
		Returning: []string{constants.F_addr.UUID, constants.F_addr.CreatedAt},
		Whitelist: struct {
			Schemas []string
			Tables  []string
			Columns []string
		}{
			Schemas: []string{outlet_schema},
			Tables:  []string{table},
			Columns: append(CFAddrTC, constants.F_addr.UUID, constants.F_addr.CreatedAt),
		},
	}

	// Build query using join-aware builder
	query, err := dbutils.BuildInsertQuery(
		method,
		outlet_schema,
		table,
		CFAddrTC,
		opts,
	)
	if err != nil {
		return nil, err
	}
	values, err := dbutils.MapValuesDirect(addr, CFAddrTC, "json")
	if err != nil {
		return nil, err
	}
	// Scan result into response struct
	var address = &model.AddResponse{}
	err = dbutils.ExecuteAndScanRow(ctx, method, ar.db, query, values, &address)
	if err != nil {
		return nil, err
	}

	return address, nil
}
func (ar *repository) UpdateFranchiseAddress(ctx context.Context, id string, addr *model.FranchiseAddress) (*model.UpdateResponse, error) {
	var method = constants.Methods.UpdateFranchiseAddress
	var table = constants.DB.Table_Franchise_documents

	if err := dbutils.CheckDBConn(ar.db, method); err != nil {
		return nil, err
	}

	conditions := map[string]any{
		constants.F_addr.UUID: id,
	}

	// Whitelist for safe updating
	opts := &dbutils.QueryBuilderOptions{
		Returning: []string{constants.F_addr.UUID, constants.F_addr.UpdatedAt},
		Whitelist: struct {
			Schemas []string
			Tables  []string
			Columns []string
		}{
			Schemas: []string{outlet_schema},
			Tables:  []string{table},
			Columns: append(CFAddrTC, constants.F_addr.UUID, constants.F_addr.UpdatedAt),
		},
	}

	query, args, err := dbutils.BuildUpdateQuery(method, outlet_schema, table, CFAddrTC, conditions, opts)
	if err != nil {
		return nil, err
	}

	values, err := dbutils.MapValuesDirect(addr, CFAddrTC, "json")
	if err != nil {
		return nil, err
	}

	copy(args[:len(CFAddrTC)], values)

	var updatedAddr model.UpdateResponse
	err = dbutils.ExecuteAndScanRow(ctx, method, ar.db, query, values,
		&updatedAddr,
	)
	if err != nil {
		return nil, err
	}
	return &updatedAddr, nil
}

func (ar *repository) GetFranchiseAddressByID(ctx context.Context, id string) (*model.FranchiseAddressResponse, error) {
	var method = constants.Methods.GetFranchiseAddressByID
	var table = constants.DB.Table_Franchise_addresses
	// Check DB connection health
	if err := dbutils.CheckDBConn(ar.db, method); err != nil {
		return nil, err
	}

	allColumns := append(CFAddrTC, constants.F_addr.UUID, constants.F_addr.CreatedAt, constants.F_addr.UpdatedAt)

	conditions := map[string]any{
		constants.F_addr.UUID: id,
	}

	// Whitelist for security
	opts := &dbutils.QueryBuilderOptions{
		Returning: allColumns,
		Whitelist: struct {
			Schemas []string
			Tables  []string
			Columns []string
		}{
			Schemas: []string{outlet_schema},
			Tables:  []string{table},
			Columns: allColumns,
		},
	}

	// Build SELECT query
	query, args, err := dbutils.BuildSelectQuery(method, outlet_schema, table, allColumns, conditions, opts)
	if err != nil {
		return nil, err
	}

	// Scan result into response struct
	var addr model.FranchiseAddressResponse
	err = dbutils.ExecuteAndScanRow(ctx, method, ar.db, query, args,
		&addr,
	)
	if err != nil {
		return nil, err
	}

	return &addr, nil
}

func (ar *repository) AddFranchiseRole(ctx context.Context, role *model.FranchiseRole) (*model.AddResponse, error) {
	var method = constants.Methods.AddFranchiseRole
	var table = constants.DB.Table_Roles
	if err := dbutils.CheckDBConn(ar.db, method); err != nil {
		return nil, err
	}

	// Whitelist options
	opts := &dbutils.QueryBuilderOptions{
		Returning: []string{constants.F_role.UUID, constants.F_role.CreatedAt},
		Whitelist: struct {
			Schemas []string
			Tables  []string
			Columns []string
		}{
			Schemas: []string{outlet_schema},
			Tables:  []string{table},
			Columns: append(CFRTC, constants.F_role.UUID, constants.F_role.CreatedAt),
		},
	}

	// Build query using join-aware builder
	query, err := dbutils.BuildInsertQuery(
		method,
		outlet_schema,
		table,
		CFRTC,
		opts,
	)
	if err != nil {
		return nil, err
	}
	values, err := dbutils.MapValuesDirect(role, COTC, "json")
	if err != nil {
		logger.Error("failed to map update values", err, nil)
		return nil, err
	}
	// Scan result into response struct
	var newRole = &model.AddResponse{}
	err = dbutils.ExecuteAndScanRow(ctx, method, ar.db, query, values, &newRole)
	if err != nil {
		return nil, err
	}

	return newRole, nil
}

func (ar *repository) UpdateFranchiseRole(ctx context.Context, id string, role *model.FranchiseRole) (*model.UpdateResponse, error) {
	var method = constants.Methods.UpdateFranchiseRole
	var table = constants.DB.Table_Roles

	if err := dbutils.CheckDBConn(ar.db, method); err != nil {
		return nil, err
	}
	allColumns := append(CFRTC, constants.F_role.UUID, constants.F_role.CreatedAt, constants.F_role.UpdatedAt)

	conditions := map[string]any{
		constants.F_role.UUID: id,
	}

	// Whitelist for safe updating
	opts := &dbutils.QueryBuilderOptions{
		Returning: []string{constants.F_role.UUID, constants.F_role.UpdatedAt},
		Whitelist: struct {
			Schemas []string
			Tables  []string
			Columns []string
		}{
			Schemas: []string{outlet_schema},
			Tables:  []string{table},
			Columns: allColumns,
		},
	}

	query, args, err := dbutils.BuildUpdateQuery(method, outlet_schema, table, CFRTC, conditions, opts)
	if err != nil {
		return nil, err
	}
	values, err := dbutils.MapValuesDirect(role, COTC, "json")
	if err != nil {
		logger.Error("failed to map update values", err, nil)
		return nil, err
	}
	copy(args[:len(COTC)], values)

	var updated = &model.UpdateResponse{}
	err = dbutils.ExecuteAndScanRow(ctx, method, ar.db, query, values,
		&updated,
	)
	if err != nil {
		return nil, err
	}
	return updated, nil
}
func (ar *repository) GetAllFranchiseRoles(ctx context.Context, id string) ([]model.FranchiseRoleResponse, error) {
	var method = constants.Methods.GetAllFranchiseRoles
	var table = constants.DB.Table_Roles
	if err := dbutils.CheckDBConn(ar.db, method); err != nil {
		return nil, err
	}

	// Conditions
	conditions := map[string]any{
		constants.F_role.FranchiseID: id,
	}

	allColumns := append(CFRTC, constants.F_role.UUID, constants.F_role.CreatedAt, constants.F_role.UpdatedAt)

	// Whitelist options
	opts := &dbutils.QueryBuilderOptions{
		Returning: allColumns,
		Whitelist: struct {
			Schemas []string
			Tables  []string
			Columns []string
		}{
			Schemas: []string{outlet_schema},
			Tables:  []string{table},
			Columns: allColumns,
		},
	}

	// Build query using join-aware builder
	query, args, err := dbutils.BuildSelectQuery(
		method,
		outlet_schema,
		table,
		allColumns,
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
	var roles []model.FranchiseRoleResponse
	for rows.Next() {
		var role model.FranchiseRoleResponse
		err := rows.Scan(&role.ID, &role.FranchiseID, &role.Name, &role.Description, &role.IsDefault, &role.CreatedAt, &role.UpdatedAt)
		if err != nil {
			logger.Error("Failed to scan row", err, nil)
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, nil
}

func (ar *repository) AddPermissionsToRole(ctx context.Context, pRole *model.RoleToPermissions) (*model.RoleToPermissions, error) {
	var method = constants.Methods.AddPermissionsToRole
	var table = constants.DB.Table_Role_Permissions
	if err := dbutils.CheckDBConn(ar.db, method); err != nil {
		return nil, err
	}

	// Whitelist options
	opts := &dbutils.QueryBuilderOptions{
		Returning: CRPTC,
		Whitelist: struct {
			Schemas []string
			Tables  []string
			Columns []string
		}{
			Schemas: []string{outlet_schema},
			Tables:  []string{table},
			Columns: CRPTC,
		},
	}

	query, err := dbutils.BuildInsertQuery(
		method,
		outlet_schema,
		table,
		CRPTC,
		opts,
	)
	if err != nil {
		return nil, err
	}

	values, err := dbutils.MapValuesDirect(pRole, COTC, "json")
	if err != nil {
		logger.Error("failed to map update values", err, nil)
		return nil, err
	}
	// Scan result into response struct
	var newPRole = &model.RoleToPermissions{}
	err = dbutils.ExecuteAndScanRow(ctx, method, ar.db, query, values, &newPRole)
	if err != nil {
		return nil, err
	}
	return newPRole, nil
}
func (ar *repository) UpdatePermissionsToRole(ctx context.Context, pRole *model.RoleToPermissions) (*model.RoleToPermissions, error) {
	var method = constants.Methods.UpdatePermissionsToRole
	var table = constants.DB.Table_Role_Permissions

	if err := dbutils.CheckDBConn(ar.db, method); err != nil {
		return nil, err
	}

	conditions := map[string]any{
		constants.F_Role_Per.RoleID: pRole.RoleID,
	}

	// Whitelist for safe updating
	opts := &dbutils.QueryBuilderOptions{
		Returning: CRPTC,
		Whitelist: struct {
			Schemas []string
			Tables  []string
			Columns []string
		}{
			Schemas: []string{outlet_schema},
			Tables:  []string{table},
			Columns: CRPTC,
		},
	}

	query, args, err := dbutils.BuildUpdateQuery(method, outlet_schema, table, CRPTC, conditions, opts)
	if err != nil {
		return nil, err
	}

	var updated model.RoleToPermissions
	err = dbutils.ExecuteAndScanRow(ctx, method, ar.db, query, args, &updated)
	if err != nil {
		return nil, err
	}
	return &updated, nil
}
func (ar *repository) GetAllPermissionsToRole(ctx context.Context, id string) ([]model.RoleToPermissionsComplete, error) {
	var method = constants.Methods.GetAllPermissionsToRole
	var table = constants.DB.Table_Roles

	if err := dbutils.CheckDBConn(ar.db, method); err != nil {
		return nil, err
	}

	columns := []string{"id", "franchise_id", "name", "description", "is_default", "p.key", "p.description", "created_at", "updated_at"}

	conditions := map[string]any{
		"franchise_id": id,
	}

	// Whitelist for safe updating
	opts := &dbutils.QueryBuilderOptions{
		Whitelist: struct {
			Schemas []string
			Tables  []string
			Columns []string
		}{
			Schemas: []string{outlet_schema},
			Tables:  []string{table},
			Columns: columns,
		},
		Returning: columns,
	}

	query, args, err := dbutils.BuildSelectQuery(method, outlet_schema, table, columns, conditions, opts)
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
	var roles []model.RoleToPermissionsComplete
	for rows.Next() {
		var role model.RoleToPermissionsComplete
		err := rows.Scan(
			&role.FranchiseID,
			&role.RoleName,
			&role.Role_Description,
			&role.Permission_Key,
			&role.Permission_Description,
			&role.IsDefault,
			&role.CreatedAt,
			&role.UpdatedAt,
		)
		if err != nil {
			logger.Error("Failed to scan row", err, nil)
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, nil
}
