package repository

import (
	"context"
	"database/sql"

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

type Repository interface {
	GetFranchiseByID(ctx context.Context, id string) (*model.FranchiseResponse, error)
	GetFranchiseByBusinessName(ctx context.Context, b_name string) (*model.FranchiseResponse, error)
	GetFranchiseOwnerByID(ctx context.Context, id string) (*model.FranchiseOwnerResponse, error)
	CheckIfOwnerExistsByAadharID(ctx context.Context, id string) (*model.FranchiseOwnerResponse, error)

	CreateFranchiseAccount(ctx context.Context, account *model.FranchiseAccount) (*model.FranchiseAccountResponse, error)
	UpdateFranchiseAccount(ctx context.Context, id string, account *model.FranchiseAccount) (*model.FranchiseAccountResponse, error)
	GetFranchiseAccountByID(ctx context.Context, id string) (*model.FranchiseAccountResponse, error)
	GetAllFranchiseAccounts(ctx context.Context, id string) ([]model.FranchiseAccountResponse, error)

	AddFranchiseDocument(ctx context.Context, doc *model.FranchiseDocument) (*model.AddResponse, error)
	UpdateFranchiseDocument(ctx context.Context, id string, doc *model.FranchiseDocument) (*model.FranchiseDocumentResponse, error)
	GetAllFranchiseDocuments(ctx context.Context, id string) ([]model.FranchiseDocumentResponseComplete, error)

	AddFranchiseAddress(ctx context.Context, addr *model.FranchiseAddress) (*model.AddResponse, error)
	UpdateFranchiseAddress(ctx context.Context, id string, addr *model.FranchiseAddress) (*model.FranchiseAddressResponse, error)
	GetFranchiseAddressByID(ctx context.Context, id string) (*model.FranchiseAddressResponse, error)

	AddFranchiseRole(ctx context.Context, role *model.FranchiseRole) (*model.AddResponse, error)
	UpdateFranchiseRole(ctx context.Context, id string, role *model.FranchiseRole) (*model.FranchiseRoleResponse, error)
	GetAllFranchiseRoles(ctx context.Context, id string) ([]model.FranchiseRoleResponse, error)

	AddPermissionsToRole(ctx context.Context, pRole *model.RoleToPermissions) (*model.RoleToPermissions, error)
	UpdatePermissionsToRole(ctx context.Context, id string, pRole *model.RoleToPermissions) (*model.RoleToPermissions, error)
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

func (ar *repository) GetFranchiseByBusinessName(ctx context.Context, b_name string) (*model.FranchiseResponse, error) {
	var method = constants.Methods.GetFranchiseByBusinessName
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
		"business_name": b_name,
		"deleted_at":    nil, // Filter deleted records
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
			Columns: append(columns, "deleted_at"),
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

func (ar *repository) GetFranchiseOwnerByID(ctx context.Context, id string) (*model.FranchiseOwnerResponse, error) {
	var method = constants.Methods.GetFranchiseOwnerByID
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

func (ar *repository) CheckIfOwnerExistsByAadharID(ctx context.Context, id string) (*model.FranchiseOwnerResponse, error) {
	var method = constants.Methods.CheckIfOwnerExistsByAadharID
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
		"aadhar_no": id, // Use franchise_id for filtering
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

func (ar *repository) CreateFranchiseAccount(ctx context.Context, account *model.FranchiseAccount) (*model.FranchiseAccountResponse, error) {
	var method = constants.Methods.CreateFranchiseAccount
	var table = constants.DB.Table_Franchise_Accounts

	if err := dbutils.CheckDBConn(ar.db, method); err != nil {
		return nil, err
	}

	// Fields to insert
	data := map[string]any{
		"franchise_id": account.FranchiseID,
		"employee_id":  account.EmployeeID,
		"login_id":     account.LoginID,
		"password":     account.Password,
		"account_type": account.AccountType,
		"name":         account.Name,
		"mobile_no":    account.MobileNo,
		"email":        account.Email,
		"role_id":      account.RoleID,
		"status":       account.Status,
	}

	columns := make([]string, 0, len(data))
	for col := range data {
		columns = append(columns, col)
	}

	// Whitelist for safe inserting
	opts := &dbutils.QueryBuilderOptions{
		Returning: append(columns, "id", "created_at"),
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

	query, err := dbutils.BuildInsertQuery(method, schema, table, columns, opts)
	if err != nil {
		return nil, err
	}

	var inserted model.FranchiseAccountResponse
	err = dbutils.ExecuteAndScanRow(ctx, method, ar.db, query, []any{account.FranchiseID, account.EmployeeID, account.LoginID, account.Password, account.AccountType, account.Name, account.MobileNo, account.Email, account.RoleID, account.Status},
		&inserted.ID,
		&inserted.FranchiseID,
		&inserted.EmployeeID,
		&inserted.LoginID,
		&inserted.AccountType,
		&inserted.Name,
		&inserted.MobileNo,
		&inserted.Email,
		&inserted.RoleName, // assuming role name is returned
		&inserted.Status,
		&inserted.CreatedAt,
		&inserted.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &inserted, nil
}

func (ar *repository) UpdateFranchiseAccount(ctx context.Context, id string, account *model.FranchiseAccount) (*model.FranchiseAccountResponse, error) {
	var method = constants.Methods.UpdateFranchiseAccount
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

func (ar *repository) GetFranchiseAccountByID(ctx context.Context, id string) (*model.FranchiseAccountResponse, error) {
	var method = constants.Methods.GetFranchiseAccountByID
	var table = constants.DB.Table_Franchise_Accounts
	// Check DB connection health
	if err := dbutils.CheckDBConn(ar.db, method); err != nil {
		return nil, err
	}

	columns := []string{
		"id", "employee_id", "name", "email", "mobile_no", "status",
		"franchise_id", "account_type", "created_at", "updated_at", "deleted_at",
	}

	conditions := map[string]any{
		"id": id,
	}

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

	// Build SELECT query
	query, args, err := dbutils.BuildSelectQuery(method, schema, table, columns, conditions, opts)
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

func (ar *repository) GetAllFranchiseAccounts(ctx context.Context, id string) ([]model.FranchiseAccountResponse, error) {
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

func (ar *repository) AddFranchiseDocument(ctx context.Context, doc *model.FranchiseDocument) (*model.AddResponse, error) {
	var method = constants.Methods.AddFranchiseDocument
	var table = constants.DB.Table_Document_Types
	if err := dbutils.CheckDBConn(ar.db, method); err != nil {
		return nil, err
	}

	columns := []string{
		"franchise_id",
		"document_type_id",
		"document_url",
		"uploaded_by",
		"status",
		"remark",
	}
	// Whitelist options
	opts := &dbutils.QueryBuilderOptions{
		Returning: []string{"id"},
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
	// Scan result into response struct
	var account = &model.AddResponse{}
	err = dbutils.ExecuteAndScanRow(ctx, method, ar.db, query, []any{doc.FranchiseID, doc.DocumentTypeID, doc.DocumentURL, doc.UploadedBy, doc.Status, doc.Remark}, &account.ID)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (ar *repository) UpdateFranchiseDocument(ctx context.Context, id string, doc *model.FranchiseDocument) (*model.FranchiseDocumentResponse, error) {
	var method = constants.Methods.UpdateFranchiseDocument
	var table = constants.DB.Table_Franchise_documents

	if err := dbutils.CheckDBConn(ar.db, method); err != nil {
		return nil, err
	}

	// Fields to update
	data := map[string]any{
		"franchise_id":     doc.FranchiseID,
		"document_type_id": doc.DocumentTypeID,
		"document_url":     doc.DocumentURL,
		"uploaded_by":      doc.UploadedBy,
		"status":           doc.Status,
		"remark":           doc.Remark,
		"verified_at":      doc.VerifiedAt,
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
		Returning: append(columns, "uploaded_at"),
	}

	query, args, err := dbutils.BuildUpdateQuery(method, schema, table, columns, conditions, opts)
	if err != nil {
		return nil, err
	}

	var updated model.FranchiseDocumentResponse
	err = dbutils.ExecuteAndScanRow(ctx, method, ar.db, query, args,
		&updated.ID,
		&updated.FranchiseID,
		&updated.DocumentTypeID,
		&updated.DocumentURL,
		&updated.UploadedBy,
		&updated.Status,
		&updated.Remark,
		&updated.VerifiedAt,
	)
	if err != nil {
		return nil, err
	}

	return &updated, nil
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

	columns := []string{
		"franchise_id",
		"address_line",
		"city",
		"state",
		"country",
		"pincode",
		"latitude",
		"longitude",
		"is_verified",
	}
	// Whitelist options
	opts := &dbutils.QueryBuilderOptions{
		Returning: []string{"id"},
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
	// Scan result into response struct
	var address = &model.AddResponse{}
	err = dbutils.ExecuteAndScanRow(ctx, method, ar.db, query, []any{addr.FranchiseID, addr.AddressLine, addr.City, addr.State, addr.Country, addr.Pincode, addr.Latitude, addr.Longitude, addr.IsVerified}, &address.ID)
	if err != nil {
		return nil, err
	}

	return address, nil
}
func (ar *repository) UpdateFranchiseAddress(ctx context.Context, id string, addr *model.FranchiseAddress) (*model.FranchiseAddressResponse, error) {
	var method = constants.Methods.UpdateFranchiseAddress
	var table = constants.DB.Table_Franchise_documents

	if err := dbutils.CheckDBConn(ar.db, method); err != nil {
		return nil, err
	}

	data := map[string]any{
		"franchise_id": addr.FranchiseID,
		"address_line": addr.AddressLine,
		"city":         addr.City,
		"state":        addr.State,
		"country":      addr.Country,
		"pincode":      addr.Pincode,
		"latitude":     addr.Latitude,
		"longitude":    addr.Longitude,
		"is_verified":  addr.IsVerified,
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
		Returning: append(columns, "uploaded_at"),
	}

	query, args, err := dbutils.BuildUpdateQuery(method, schema, table, columns, conditions, opts)
	if err != nil {
		return nil, err
	}

	var updated model.FranchiseAddressResponse
	err = dbutils.ExecuteAndScanRow(ctx, method, ar.db, query, args,
		&updated.ID,
		&updated.AddressLine,
		&updated.City,
		&updated.State,
		&updated.Country,
		&updated.Pincode,
		&updated.Latitude,
		&updated.Longitude,
		&updated.IsVerified,
		&updated.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &updated, nil
}
func (ar *repository) GetFranchiseAddressByID(ctx context.Context, id string) (*model.FranchiseAddressResponse, error) {
	var method = constants.Methods.GetFranchiseAddressByID
	var table = constants.DB.Table_Franchise_addresses
	// Check DB connection health
	if err := dbutils.CheckDBConn(ar.db, method); err != nil {
		return nil, err
	}

	columns := []string{
		"id", "address_line", "city", "state", "country", "pincode",
		"latitude", "longitude", "is_verified", "created_at", "updated_at",
	}

	conditions := map[string]any{
		"id": id,
	}

	// Whitelist for security
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

	// Build SELECT query
	query, args, err := dbutils.BuildSelectQuery(method, schema, table, columns, conditions, opts)
	if err != nil {
		return nil, err
	}

	// Scan result into response struct
	var addr model.FranchiseAddressResponse
	err = dbutils.ExecuteAndScanRow(ctx, method, ar.db, query, args,
		&addr.ID,
		&addr.AddressLine,
		&addr.City,
		&addr.State,
		&addr.Country,
		&addr.Pincode,
		&addr.Latitude,
		&addr.Longitude,
		&addr.IsVerified,
		&addr.CreatedAt,
		&addr.UpdatedAt,
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

	columns := []string{
		"franchise_id",
		"name",
		"description",
		"is_default",
	}
	// Whitelist options
	opts := &dbutils.QueryBuilderOptions{
		Returning: []string{"id"},
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
	// Scan result into response struct
	var newRole = &model.AddResponse{}
	err = dbutils.ExecuteAndScanRow(ctx, method, ar.db, query, []any{role.FranchiseID, role.Name, role.Description, role.IsDefault}, &newRole.ID)
	if err != nil {
		return nil, err
	}

	return newRole, nil
}

func (ar *repository) UpdateFranchiseRole(ctx context.Context, id string, role *model.FranchiseRole) (*model.FranchiseRoleResponse, error) {
	var method = constants.Methods.UpdateFranchiseRole
	var table = constants.DB.Table_Roles

	if err := dbutils.CheckDBConn(ar.db, method); err != nil {
		return nil, err
	}

	data := map[string]any{
		"franchise_id": role.FranchiseID,
		"name":         role.Name,
		"description":  role.Description,
		"is_default":   role.IsDefault,
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
		Returning: append(columns, "updated_at"),
	}

	query, args, err := dbutils.BuildUpdateQuery(method, schema, table, columns, conditions, opts)
	if err != nil {
		return nil, err
	}

	var updated model.FranchiseRoleResponse
	err = dbutils.ExecuteAndScanRow(ctx, method, ar.db, query, args,
		&updated.ID,
		&updated.FranchiseID,
		&updated.Name,
		&updated.Description,
		&updated.IsDefault,
		&updated.CreatedAt,
		&updated.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &updated, nil
}
func (ar *repository) GetAllFranchiseRoles(ctx context.Context, id string) ([]model.FranchiseRoleResponse, error) {
	var method = constants.Methods.GetAllFranchiseRoles
	var table = constants.DB.Table_Roles
	if err := dbutils.CheckDBConn(ar.db, method); err != nil {
		return nil, err
	}

	// Define columns with alias prefixes
	columns := []string{
		"franchise_id",
		"name",
		"description",
		"is_default",
		"created_at",
		"updated_at",
	}

	// Conditions
	conditions := map[string]any{
		"franchise_id": id,
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
			Columns: columns,
		},
	}

	// Build query using join-aware builder
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

	columns := []string{
		"role_id",
		"permission_id",
	}
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
	// Scan result into response struct
	var newPRole = &model.RoleToPermissions{}
	err = dbutils.ExecuteAndScanRow(ctx, method, ar.db, query, []any{pRole.RoleID, pRole.PermissionID}, &newPRole.RoleID, &newPRole.PermissionID)
	if err != nil {
		return nil, err
	}
	return newPRole, nil
}
func (ar *repository) UpdatePermissionsToRole(ctx context.Context, id string, pRole *model.RoleToPermissions) (*model.RoleToPermissions, error) {
	var method = constants.Methods.UpdatePermissionsToRole
	var table = constants.DB.Table_Role_Permissions

	if err := dbutils.CheckDBConn(ar.db, method); err != nil {
		return nil, err
	}

	data := map[string]any{
		"role_id":       pRole.RoleID,
		"permission_id": pRole.PermissionID,
	}

	columns := make([]string, 0, len(data))
	for col := range data {
		columns = append(columns, col)
	}

	conditions := map[string]any{
		"role_id": pRole.RoleID,
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

	var updated model.RoleToPermissions
	err = dbutils.ExecuteAndScanRow(ctx, method, ar.db, query, args,
		&updated.RoleID,
		&updated.PermissionID,
	)
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

	query, args, err := dbutils.BuildSelectQuery(method, schema, table, columns, conditions, opts)
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
