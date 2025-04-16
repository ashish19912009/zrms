package dbutils

import (
	"fmt"
	"strings"

	"github.com/ashish19912009/zrms/services/account/internal/constants"
	"github.com/ashish19912009/zrms/services/account/internal/logger"
)

// QueryBuilderOptions holds optional settings for building queries.
type QueryBuilderOptions struct {
	Returning []string
	Whilelist struct {
		Schemas []string
		Tables  []string
		Columns []string
	}
}

type JoinClause struct {
	Type   string // e.g. "INNER", "LEFT", etc.
	Schema string
	Table  string
	Alias  string
	On     string // e.g. "a.franchise_id = f.id"
}

// BuildSelectQuery dynamically constructs a parameterized SQL SELECT query string.
// It ensures schema, table, column, and condition validations using a whitelist.
// Uses structured logging for errors and returns both the query and its arguments.
//
// Parameters:
// - methodName: Name of the calling repository method, used for log traceability.
// - schema:           Database schema name (must be whitelisted).
// - table:            Table name (must be whitelisted).
// - columns:          Columns to select (must be whitelisted).
// - conditions:       Key-value map for WHERE clause (columns must be whitelisted).
// - opts:             QueryBuilderOptions pointer containing whitelist configuration.
//
// Returns:
// - query string (with parameter placeholders), slice of arguments for the placeholders, and error if any.
//
// Example:
//
//	query, args, err := BuildSelectQuery(
//		"GetActiveAccounts",
//		"public",
//		"accounts",
//		[]string{"id", "name"},
//		map[string]any{"status": "active"},
//		opts,
//	)
func BuildSelectQuery(methodName, schema, table string, columns []string, conditions map[string]any, opts *QueryBuilderOptions) (string, []any, error) {

	logCtx := logger.BaseLogContext(
		"layer", constants.Repository,
		"method", methodName,
		"schema", schema,
		"table", table,
		"columns", strings.Join(columns, ", "),
		"conditions", parseConditions(conditions),
	)
	if len(columns) == 0 {
		err := fmt.Errorf(constants.NoColumProvided)
		logger.Error(constants.BuildSelectQuery, err, logCtx)
		return "", nil, err
	}

	if opts != nil {
		if !contains(opts.Whilelist.Schemas, schema) {
			err := fmt.Errorf(constants.UnauthorizedSchema, schema)
			logger.Error(constants.BuildSelectQuery, err, logCtx)
			return "", nil, err
		}
		if !contains(opts.Whilelist.Tables, table) {
			err := fmt.Errorf(constants.UnauthorizedTable, table)
			logger.Error(constants.BuildSelectQuery, err, logCtx)
			return "", nil, err
		}
		for _, col := range columns {
			if !contains(opts.Whilelist.Columns, col) {
				err := fmt.Errorf(constants.UnauthorizedCloumn, col)
				logger.Error(constants.BuildSelectQuery, err, logCtx)
				return "", nil, err
			}
		}
	}

	quotedCols := make([]string, len(columns))
	for i, col := range columns {
		quotedCols[i] = fmt.Sprintf(`"%s"`, col)
	}

	query := fmt.Sprintf(`SELECT %s FROM "%s"."%s"`, strings.Join(quotedCols, ", "), schema, table)

	args := []any{}
	if len(conditions) > 0 {
		whereClauses := []string{}
		argIndex := 1
		for key, val := range conditions {
			if opts != nil && !contains(opts.Whilelist.Columns, key) {
				err := fmt.Errorf(constants.UnauthorizedConditionColumn, key)
				logger.Error(constants.BuildSelectQuery, err, logCtx)
			}
			whereClauses = append(whereClauses, fmt.Sprintf(`"%s" = $%d`, key, argIndex))
			args = append(args, val)
			argIndex++
		}
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	return query, args, nil
}

// BuildInsertQuery dynamically constructs a parameterized SQL INSERT query string.
// It ensures security through a whitelist mechanism and logs structured errors for any violations.
//
// Parameters:
// - methodName: Name of the repository method calling this function, used in structured logging.
// - schema:           Database schema name (must be whitelisted).
// - table:            Table name (must be whitelisted).
// - columns:          List of column names to insert values into (must be whitelisted).
// - opts:             QueryBuilderOptions pointer containing whitelist and optional RETURNING clause.
//
// Returns:
// - A parameterized SQL query string (e.g., `INSERT INTO ... (...) VALUES (...) RETURNING ...`).
// - Error if input validation fails.
//
// Example:
//
//		query, err := BuildInsertQuery(
//	     "CreateUser",
//			"public",
//			"users",
//			[]string{"name", "email"},
//			opts,
//		)
func BuildInsertQuery(methodName, schema, table string, columns []string, opts *QueryBuilderOptions) (string, error) {
	logCtx := logger.BaseLogContext(
		"layer", constants.Repository,
		"method", methodName,
		"schema", schema,
		"table", table,
		"columns", strings.Join(columns, ", "),
	)
	if len(columns) == 0 {
		err := fmt.Errorf(constants.NoColumProvided)
		logger.Error(constants.BuildUpdateQuery, err, logCtx)
		return "", nil
	}

	if opts != nil {
		if !contains(opts.Whilelist.Schemas, schema) {
			err := fmt.Errorf(constants.UnauthorizedSchema, schema)
			logger.Error(constants.BuildInsertQuery, err, logCtx)
			return "", err
		}
		if !contains(opts.Whilelist.Tables, table) {
			err := fmt.Errorf(constants.UnauthorizedTable, table)
			logger.Error(constants.BuildInsertQuery, err, logCtx)
			return "", err
		}
		for _, col := range columns {
			if !contains(opts.Whilelist.Columns, col) {
				err := fmt.Errorf(constants.UnauthorizedCloumn, col)
				logger.Error(constants.BuildInsertQuery, err, logCtx)
				return "", err
			}
		}
		if len(opts.Returning) > 0 {
			for _, col := range opts.Returning {
				if !contains(opts.Whilelist.Columns, col) {
					err := fmt.Errorf(constants.UnauthorizedReturningColumn, col)
					logger.Error(constants.BuildInsertQuery, err, logCtx)
					return "", err
				}
			}
		}
	}

	quotedCols := make([]string, len(columns))
	args := make([]string, len(columns))

	for i, col := range columns {
		quotedCols[i] = fmt.Sprintf(`"%s"`, col)
		args[i] = fmt.Sprintf("$%d", i+1)
	}

	query := fmt.Sprintf(`INSERT INTO "%s"."%s" (%s) VALUES (%s)`, schema, table, strings.Join(quotedCols, ", "), strings.Join(args, ", "))

	if opts != nil && len(opts.Returning) > 0 {
		returningCols := make([]string, len(opts.Returning))
		for i, col := range opts.Returning {
			returningCols[i] = fmt.Sprintf(`"%s"`, col)
		}
		query += fmt.Sprintf(" RETURNING %s", strings.Join(returningCols, ", "))
	}

	return query, nil
}

// BuildUpdateQuery dynamically constructs a parameterized SQL UPDATE query string.
// It ensures security through a whitelist mechanism and logs detailed errors if any validation fails.
//
// Parameters:
// - methodName: Name of the repository method calling this function, used in structured logging.
// - schema:        Database schema name (must be whitelisted).
// - table:         Table name (must be whitelisted).
// - columns:       List of column names to be updated (must be whitelisted).
// - conditions:    Map of condition column names and their corresponding values (WHERE clause, keys must be whitelisted).
// - opts:          QueryBuilderOptions pointer for whitelist validation and RETURNING columns.
//
// Returns:
// - A parameterized SQL query string (e.g., `UPDATE ... SET ... WHERE ... RETURNING ...`).
// - Slice of arguments corresponding to placeholders in the SQL query.
// - Error if validation fails or required input is missing.
//
// Example:
//
//		query, args, err := BuildUpdateQuery(
//	     "updateFranchise"
//			"public",
//			"users",
//			[]string{"name", "email"},
//			map[string]any{"id": 123},
//			opts,
//			"UpdateUserByID",
//		)
func BuildUpdateQuery(methodName, schema, table string, columns []string, conditions map[string]any, opts *QueryBuilderOptions) (string, []any, error) {
	// Initial log context with base details
	logCtx := logger.BaseLogContext(
		"layer", constants.Repository,
		"method", methodName,
		"schema", schema,
		"table", table,
		"columns", strings.Join(columns, ", "),
	)
	if len(columns) == 0 {
		err := fmt.Errorf(constants.NoColumProvided)
		logger.Error(constants.BuildUpdateQuery, err, logCtx)
		return "", nil, err
	}

	// Whitelist validations
	if opts != nil {
		if !contains(opts.Whilelist.Schemas, schema) {
			err := fmt.Errorf(constants.UnauthorizedSchema, schema)
			logger.Error(constants.BuildUpdateQuery, err, logCtx)
			return "", nil, err
		}
		if !contains(opts.Whilelist.Tables, table) {
			err := fmt.Errorf(constants.UnauthorizedTable, table)
			logger.Error(constants.BuildUpdateQuery, err, logCtx)
			return "", nil, err
		}
		for _, col := range columns {
			if !contains(opts.Whilelist.Columns, col) {
				err := fmt.Errorf(constants.UnauthorizedCloumn, col)
				logger.Error(constants.BuildUpdateQuery, err, logCtx)
				return "", nil, err
			}
		}
		for condCol := range conditions {
			if !contains(opts.Whilelist.Columns, condCol) {
				err := fmt.Errorf(constants.UnauthorizedConditionColumn, condCol)
				logger.Error(constants.BuildUpdateQuery, err, logCtx)
				return "", nil, err
			}
		}
		for _, col := range opts.Returning {
			if !contains(opts.Whilelist.Columns, col) {
				err := fmt.Errorf(constants.UnauthorizedReturningColumn, col)
				logger.Error(constants.BuildUpdateQuery, err, logCtx)
				return "", nil, err
			}
		}
	}

	setClauses := make([]string, len(columns))
	args := make([]any, 0, len(columns)+len(conditions))

	// SET columns
	for i, col := range columns {
		setClauses[i] = fmt.Sprintf(`"%s" = $%d`, col, i+1)
		args = append(args, nil) // placeholder, to be set by caller
	}

	query := fmt.Sprintf(`UPDATE "%s"."%s" SET %s`, schema, table, strings.Join(setClauses, ", "))

	// WHERE conditions
	if len(conditions) > 0 {
		whereClauses := []string{}
		argIndex := len(columns) + 1
		for key, val := range conditions {
			whereClauses = append(whereClauses, fmt.Sprintf(`"%s" = $%d`, key, argIndex))
			args = append(args, val)
			argIndex++
		}
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	// RETURNING clause
	if len(opts.Returning) > 0 {
		returningCols := make([]string, len(opts.Returning))
		for i, col := range opts.Returning {
			returningCols[i] = fmt.Sprintf(`"%s"`, col)
		}
		query += fmt.Sprintf(" RETURNING %s", strings.Join(returningCols, ", "))
	}

	return query, args, nil
}

// BuildDeleteQuery constructs a parameterized SQL DELETE query string with added security.
// It validates schema, table, and column names against a whitelist and parameterizes the WHERE condition
// to prevent SQL injection.
//
// Parameters:
// - schema:           Database schema name (must be whitelisted).
// - table:            Table name (must be whitelisted).
// - conditions:       A map of key-value pairs for WHERE condition (must be sanitized).
// - opts:             QueryBuilderOptions pointer containing whitelist configuration.
//
// Returns:
// - query string (parameterized) and error if any.
func BuildDeleteQuery(schema, table string, conditions map[string]any, opts *QueryBuilderOptions, methodName string) (string, []any, error) {
	logCtx := logger.BaseLogContext(
		"layer", constants.Repository,
		"method", methodName,
		"schema", schema,
		"table", table,
		"conditions", parseConditions(conditions),
	)

	// Validate schema and table against whitelist
	if opts != nil {
		if !contains(opts.Whilelist.Schemas, schema) {
			err := fmt.Errorf("unauthorized schema: %s", schema)
			logger.Error(constants.BuildDeleteQuery, err, logCtx)
			return "", nil, err
		}
		if !contains(opts.Whilelist.Tables, table) {
			err := fmt.Errorf("unauthorized table: %s", table)
			logger.Error(constants.BuildDeleteQuery, err, logCtx)
			return "", nil, err
		}
	}

	// Build the WHERE clause from conditions
	var args []any
	var whereClauses []string
	argIndex := 1
	for key, val := range conditions {
		// Validate condition columns against whitelist
		if opts != nil && !contains(opts.Whilelist.Columns, key) {
			err := fmt.Errorf("unauthorized condition column: %s", key)
			logger.Error(constants.BuildDeleteQuery, err, logCtx)
			return "", nil, err
		}
		whereClauses = append(whereClauses, fmt.Sprintf(`"%s" = $%d`, key, argIndex))
		args = append(args, val)
		argIndex++
	}

	// Form the DELETE query
	query := fmt.Sprintf(`DELETE FROM "%s"."%s" WHERE %s`, schema, table, strings.Join(whereClauses, " AND "))

	// Add RETURNING clause if specified
	if opts != nil && len(opts.Returning) > 0 {
		returningCols := make([]string, len(opts.Returning))
		for i, col := range opts.Returning {
			if !contains(opts.Whilelist.Columns, col) {
				err := fmt.Errorf("unauthorized returning column: %s", col)
				logger.Error(constants.BuildDeleteQuery, err, logCtx)
				return "", nil, err
			}
			returningCols[i] = fmt.Sprintf(`"%s"`, col)
		}
		query += fmt.Sprintf(" RETURNING %s", strings.Join(returningCols, ", "))
	}

	return query, args, nil
}

func BuildJoinSelectQuery(
	methodName string,
	mainSchema string,
	mainTable string,
	mainAlias string,
	columns []string,
	joins []JoinClause,
	conditions map[string]any,
	opts *QueryBuilderOptions,
) (string, []any, error) {

	logCtx := logger.BaseLogContext(
		"layer", constants.Repository,
		"method", methodName,
		"main_schema", mainSchema,
		"main_table", mainTable,
		"main_alias", mainAlias,
		"columns", strings.Join(columns, ", "),
		"conditions", parseConditions(conditions),
	)

	// Validation
	if len(columns) == 0 {
		err := fmt.Errorf(constants.NoColumProvided)
		logger.Error(constants.BuildSelectQuery, err, logCtx)
		return "", nil, err
	}
	if opts != nil {
		if !contains(opts.Whilelist.Schemas, mainSchema) {
			return "", nil, fmt.Errorf(constants.UnauthorizedSchema, mainSchema)
		}
		if !contains(opts.Whilelist.Tables, mainTable) {
			return "", nil, fmt.Errorf(constants.UnauthorizedTable, mainTable)
		}
		for _, col := range columns {
			if !contains(opts.Whilelist.Columns, col) {
				return "", nil, fmt.Errorf(constants.UnauthorizedCloumn, col)
			}
		}
	}

	// SELECT columns
	quotedCols := make([]string, len(columns))
	for i, col := range columns {
		quotedCols[i] = fmt.Sprintf(`%s`, col) // let caller prefix with alias like a.id
	}
	query := fmt.Sprintf(`SELECT %s FROM "%s"."%s" AS %s`, strings.Join(quotedCols, ", "), mainSchema, mainTable, mainAlias)

	// JOIN clauses
	for _, join := range joins {
		if opts != nil {
			if !contains(opts.Whilelist.Tables, join.Table) || !contains(opts.Whilelist.Schemas, join.Schema) {
				return "", nil, fmt.Errorf(constants.UnauthorizedJoinTable, join.Table)
			}
		}
		query += fmt.Sprintf(` %s JOIN "%s"."%s" AS %s ON %s`, join.Type, join.Schema, join.Table, join.Alias, join.On)
	}

	// WHERE conditions
	args := []any{}
	if len(conditions) > 0 {
		whereClauses := []string{}
		argIndex := 1
		for key, val := range conditions {
			whereClauses = append(whereClauses, fmt.Sprintf(`%s = $%d`, key, argIndex))
			args = append(args, val)
			argIndex++
		}
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	return query, args, nil
}

// contains checks if a value exists in a slice.
func contains(slice []string, value string) bool {
	for _, s := range slice {
		if s == value {
			return true
		}
	}
	return false
}

func parseConditions(conditions map[string]any) string {
	var parts []string
	for k, v := range conditions {
		parts = append(parts, fmt.Sprintf("%s=%v", k, v))
	}
	return strings.Join(parts, ", ")
}
