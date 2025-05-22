package dbutils

import (
	"fmt"
	"strings"
	"time"

	"github.com/ashish19912009/zrms/services/account/internal/constants"
	"github.com/ashish19912009/zrms/services/account/internal/logger"
)

type cachedQuery struct {
	sql      string
	args     []any
	lastUsed time.Time
}

// QueryBuilderOptions holds optional settings for building queries.
// QueryBuilderOptions holds optional settings for building queries.
type QueryBuilderOptions struct {
	Returning []string
	Whitelist struct {
		Schemas []string
		Tables  []string
		Columns []string
	}
	AllowReturningAll bool // If true, allows RETURNING *
	OrderBy           []string
	Limit             int32
	Offset            int32
}

type JoinClause struct {
	Type   string // e.g. "INNER", "LEFT", etc.
	Schema string
	Table  string
	Alias  string
	On     string // e.g. "a.franchise_id = f.id"
}

// Pre-defined list of condition suffixes (optimized for extractColumnName)
var conditionSuffixes = []string{
	"__eq", "__ne", "__lt", "__lte", "__gt", "__gte",
	"__in", "__nin", "__null", "__notnull",
	"__like", "__ilike", "__between", "__exists", "__notexists",
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

// Condition key suffix					Behavior					SQL Example
// __null								IS NULL						"deleted_at" IS NULL
// __notnull							IS NOT NULL					"approved_at" IS NOT NULL
// (default)							Equals = $N					"id" = $1
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
		if !contains(opts.Whitelist.Schemas, schema) {
			err := fmt.Errorf(constants.UnauthorizedSchema, schema)
			logger.Error(constants.BuildSelectQuery, err, logCtx)
			return "", nil, err
		}
		if !contains(opts.Whitelist.Tables, table) {
			err := fmt.Errorf(constants.UnauthorizedTable, table)
			logger.Error(constants.BuildSelectQuery, err, logCtx)
			return "", nil, err
		}
		for _, col := range columns {
			if !contains(opts.Whitelist.Columns, col) {
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
			if opts != nil && !contains(opts.Whitelist.Columns, stripConditionSuffix(key)) {
				err := fmt.Errorf(constants.UnauthorizedConditionColumn, key)
				logger.Error(constants.BuildSelectQuery, err, logCtx)
			}

			switch {
			case strings.HasSuffix(key, "__null"):
				col := stripConditionSuffix(key)
				whereClauses = append(whereClauses, fmt.Sprintf(`"%s" IS NULL`, col))

			case strings.HasSuffix(key, "__notnull"):
				col := stripConditionSuffix(key)
				whereClauses = append(whereClauses, fmt.Sprintf(`"%s" IS NOT NULL`, col))

			default:
				whereClauses = append(whereClauses, fmt.Sprintf(`"%s" = $%d`, key, argIndex))
				args = append(args, val)
				argIndex++
			}
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
		if !contains(opts.Whitelist.Schemas, schema) {
			err := fmt.Errorf(constants.UnauthorizedSchema, schema)
			logger.Error(constants.BuildInsertQuery, err, logCtx)
			return "", err
		}
		if !contains(opts.Whitelist.Tables, table) {
			err := fmt.Errorf(constants.UnauthorizedTable, table)
			logger.Error(constants.BuildInsertQuery, err, logCtx)
			return "", err
		}
		for _, col := range columns {
			if !contains(opts.Whitelist.Columns, col) {
				err := fmt.Errorf(constants.UnauthorizedCloumn, col)
				logger.Error(constants.BuildInsertQuery, err, logCtx)
				return "", err
			}
		}
		if len(opts.Returning) > 0 {
			for _, col := range opts.Returning {
				if !contains(opts.Whitelist.Columns, col) {
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
func BuildUpdateQuery(
	methodName, schema, table string,
	columns []string,
	conditions map[string]any,
	opts *QueryBuilderOptions,
) (string, []any, error) {
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
		if !contains(opts.Whitelist.Schemas, schema) {
			err := fmt.Errorf(constants.UnauthorizedSchema, schema)
			logger.Error(constants.BuildUpdateQuery, err, logCtx)
			return "", nil, err
		}
		if !contains(opts.Whitelist.Tables, table) {
			err := fmt.Errorf(constants.UnauthorizedTable, table)
			logger.Error(constants.BuildUpdateQuery, err, logCtx)
			return "", nil, err
		}
		for _, col := range columns {
			if !contains(opts.Whitelist.Columns, col) {
				err := fmt.Errorf(constants.UnauthorizedCloumn, col)
				logger.Error(constants.BuildUpdateQuery, err, logCtx)
				return "", nil, err
			}
		}
		for condCol := range conditions {
			baseCol := extractColumnName(condCol) // strip suffixes like __gte, etc
			if !contains(opts.Whitelist.Columns, baseCol) {
				err := fmt.Errorf(constants.UnauthorizedConditionColumn, condCol)
				logger.Error(constants.BuildUpdateQuery, err, logCtx)
				return "", nil, err
			}
		}
		for _, col := range opts.Returning {
			if !contains(opts.Whitelist.Columns, col) {
				err := fmt.Errorf(constants.UnauthorizedReturningColumn, col)
				logger.Error(constants.BuildUpdateQuery, err, logCtx)
				return "", nil, err
			}
		}
	}

	setClauses := make([]string, len(columns))
	args := make([]any, 0, len(columns)+len(conditions))

	// SET clause, placeholders start at $1
	for i, col := range columns {
		setClauses[i] = fmt.Sprintf(`"%s" = $%d`, col, i+1)
		args = append(args, nil) // placeholder; caller fills actual value later
	}

	query := fmt.Sprintf(`UPDATE "%s"."%s" SET %s`, schema, table, strings.Join(setClauses, ", "))

	// WHERE clause
	if len(conditions) > 0 {
		whereClauses := []string{}
		argIndex := len(columns) + 1
		for col, val := range conditions {
			clause, clauseArgs, err := buildConditionClause(col, val, &argIndex)
			if err != nil {
				logger.Error(constants.BuildUpdateQuery, err, logCtx)
				return "", nil, err
			}
			whereClauses = append(whereClauses, clause)
			args = append(args, clauseArgs...)
		}
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	// RETURNING clause
	if opts != nil && len(opts.Returning) > 0 {
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
		if !contains(opts.Whitelist.Schemas, schema) {
			err := fmt.Errorf("unauthorized schema: %s", schema)
			logger.Error(constants.BuildDeleteQuery, err, logCtx)
			return "", nil, err
		}
		if !contains(opts.Whitelist.Tables, table) {
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
		if opts != nil && !contains(opts.Whitelist.Columns, key) {
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
			if !contains(opts.Whitelist.Columns, col) {
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
		if !contains(opts.Whitelist.Schemas, mainSchema) {
			return "", nil, fmt.Errorf(constants.UnauthorizedSchema, mainSchema)
		}
		if !contains(opts.Whitelist.Tables, mainTable) {
			return "", nil, fmt.Errorf(constants.UnauthorizedTable, mainTable)
		}
		for _, col := range columns {
			if !contains(opts.Whitelist.Columns, col) {
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
			if !contains(opts.Whitelist.Tables, join.Table) || !contains(opts.Whitelist.Schemas, join.Schema) {
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

func stripConditionSuffix(key string) string {
	key = strings.TrimSuffix(key, "__null")
	key = strings.TrimSuffix(key, "__notnull")
	return key
}

// buildConditionClause handles operators like __eq, __in, etc.
func buildConditionClause(key string, value any, argIndex *int) (string, []any, error) {
	for _, suffix := range conditionSuffixes {
		if strings.HasSuffix(key, suffix) {
			col := strings.TrimSuffix(key, suffix)
			switch suffix {
			case "__eq":
				return fmt.Sprintf(`"%s" = $%d`, col, *argIndex), []any{value}, nil
			case "__in":
				vals, ok := value.([]any)
				if !ok {
					return "", nil, fmt.Errorf("value for IN must be a slice, got %T", value)
				}
				placeholders := make([]string, len(vals))
				for i := range vals {
					placeholders[i] = fmt.Sprintf("$%d", *argIndex+i)
				}
				*argIndex += len(vals)
				return fmt.Sprintf(`"%s" IN (%s)`, col, strings.Join(placeholders, ", ")), vals, nil
			case "__between":
				vals, ok := value.([]any)
				if !ok || len(vals) != 2 {
					return "", nil, fmt.Errorf("value for BETWEEN must be [start, end], got %v", value)
				}
				clause := fmt.Sprintf(`"%s" BETWEEN $%d AND $%d`, col, *argIndex, *argIndex+1)
				*argIndex += 2
				return clause, vals, nil
				// Add other cases (__lt, __like, etc.) similarly...
			}
		}
	}
	// Default: equality check
	return fmt.Sprintf(`"%s" = $%d`, key, *argIndex), []any{value}, nil
}

// func buildConditionClause(key string, value any, argIndex *int) (string, []any, error) {

// 	for _, suffix := range conditionSuffixes {
// 		if strings.HasSuffix(key, suffix) {
// 			col := strings.TrimSuffix(key, suffix)

// 			switch suffix {
// 			case "__eq":
// 				return fmt.Sprintf(`"%s" = $%d`, col, *argIndex), []any{value}, advance(argIndex)
// 			case "__ne":
// 				return fmt.Sprintf(`"%s" != $%d`, col, *argIndex), []any{value}, advance(argIndex)
// 			case "__lt":
// 				return fmt.Sprintf(`"%s" < $%d`, col, *argIndex), []any{value}, advance(argIndex)
// 			case "__lte":
// 				return fmt.Sprintf(`"%s" <= $%d`, col, *argIndex), []any{value}, advance(argIndex)
// 			case "__gt":
// 				return fmt.Sprintf(`"%s" > $%d`, col, *argIndex), []any{value}, advance(argIndex)
// 			case "__gte":
// 				return fmt.Sprintf(`"%s" >= $%d`, col, *argIndex), []any{value}, advance(argIndex)

// 			case "__in", "__nin":
// 				v := reflect.ValueOf(value)
// 				if v.Kind() != reflect.Slice {
// 					return "", nil, fmt.Errorf("value for '%s' must be a slice", key)
// 				}
// 				if v.Len() == 0 {
// 					return "", nil, fmt.Errorf("value for '%s' cannot be an empty slice", key)
// 				}

// 				placeholders := make([]string, v.Len())
// 				args := make([]any, v.Len())
// 				for i := 0; i < v.Len(); i++ {
// 					placeholders[i] = fmt.Sprintf("$%d", *argIndex)
// 					args[i] = v.Index(i).Interface()
// 					*argIndex++
// 				}

// 				operator := "IN"
// 				if suffix == "__nin" {
// 					operator = "NOT IN"
// 				}
// 				return fmt.Sprintf(`"%s" %s (%s)`, col, operator, strings.Join(placeholders, ", ")), args, nil

// 			case "__null":
// 				return fmt.Sprintf(`"%s" IS NULL`, col), nil, nil
// 			case "__notnull":
// 				return fmt.Sprintf(`"%s" IS NOT NULL`, col), nil, nil

// 			case "__like":
// 				return fmt.Sprintf(`"%s" LIKE $%d`, col, *argIndex), []any{value}, advance(argIndex)
// 			case "__ilike":
// 				return fmt.Sprintf(`"%s" ILIKE $%d`, col, *argIndex), []any{value}, advance(argIndex)

// 			case "__between":
// 				v := reflect.ValueOf(value)
// 				if v.Kind() != reflect.Slice || v.Len() != 2 {
// 					return "", nil, fmt.Errorf("value for '%s' must be a 2-element slice", key)
// 				}
// 				arg1 := fmt.Sprintf("$%d", *argIndex)
// 				val1 := v.Index(0).Interface()
// 				*argIndex++
// 				arg2 := fmt.Sprintf("$%d", *argIndex)
// 				val2 := v.Index(1).Interface()
// 				*argIndex++
// 				return fmt.Sprintf(`"%s" BETWEEN %s AND %s`, col, arg1, arg2), []any{val1, val2}, nil

// 			case "__exists", "__notexists":
// 				query, ok := value.(string)
// 				if !ok {
// 					return "", nil, fmt.Errorf("value for '%s' must be a subquery string", key)
// 				}
// 				prefix := "EXISTS"
// 				if suffix == "__notexists" {
// 					prefix = "NOT EXISTS"
// 				}
// 				return fmt.Sprintf(`%s (%s)`, prefix, query), nil, nil
// 			}
// 		}
// 	}

// 	// default to __eq
// 	clause := fmt.Sprintf(`"%s" = $%d`, key, *argIndex)
// 	args := []any{value}
// 	*argIndex++
// 	return clause, args, nil
// }

func advance(argIndex *int) error {
	*argIndex++
	return nil
}

func extractColumnName(key string) string {
	for _, suffix := range conditionSuffixes {
		if strings.HasSuffix(key, suffix) {
			return strings.TrimSuffix(key, suffix)
		}
	}
	return key
}
