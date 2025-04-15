package dbutils

import (
	"fmt"
	"strings"
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

func BuildSelectQuery(schema, table string, columns []string, conditions map[string]any, opts *QueryBuilderOptions) (string, []any, error) {
	if len(columns) == 0 {
		return "", nil, fmt.Errorf("no columns provided for select")
	}

	if opts != nil {
		if !contains(opts.Whilelist.Schemas, schema) {
			return "", nil, fmt.Errorf("unauthorized schema: %s", schema)
		}
		if !contains(opts.Whilelist.Tables, table) {
			return "", nil, fmt.Errorf("unauthorized table: %s", table)
		}
		for _, col := range columns {
			if !contains(opts.Whilelist.Columns, col) {
				return "", nil, fmt.Errorf("unauthorized column: %s", col)
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
				return "", nil, fmt.Errorf("unauthorized condition column: %s", key)
			}
			whereClauses = append(whereClauses, fmt.Sprintf(`"%s" = $%d`, key, argIndex))
			args = append(args, val)
			argIndex++
		}
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	return query, args, nil
}

// BuildInsertQuery constructs a safe INSERT query with optional RETURNING fields.
func BuildInsertQuery(schema, table string, columns []string, opts *QueryBuilderOptions) (string, error) {
	if len(columns) == 0 {
		return "", fmt.Errorf("no Columns provided")
	}
	// Validate schema, table, and columns against the whitelist if provided
	if opts != nil {
		if !contains(opts.Whilelist.Schemas, schema) {
			return "", fmt.Errorf("unauthorized schema: %s", schema)
		}
		if !contains(opts.Whilelist.Tables, table) {
			return "", fmt.Errorf("unauthorized table: %s", table)
		}
		for _, col := range columns {
			if !contains(opts.Whilelist.Columns, col) {
				return "", fmt.Errorf("unauthorized colum: %s", col)
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

	// Handle RETURNING fields
	if opts != nil && len(opts.Returning) > 0 {
		returninfCols := make([]string, len(opts.Returning))
		for i, col := range opts.Returning {
			returninfCols[i] = fmt.Sprintf(`"%s"`, col)
		}
		query += fmt.Sprintf(" RETURING %s", strings.Join(returninfCols, ", "))
	}
	return query, nil
}

// BuildUpdateQuery constructs a safe UPDATE query with optional RETURNING fields.
func BuildUpdateQuery(schema, table string, columns []string, condition string, opts *QueryBuilderOptions) (string, error) {
	if len(columns) == 0 {
		return "", fmt.Errorf("no Columns provided")
	}
	if opts != nil {
		if !contains(opts.Whilelist.Schemas, schema) {
			return "", fmt.Errorf("unauthorized schema: %s", schema)
		}
		if !contains(opts.Whilelist.Tables, table) {
			return "", fmt.Errorf("unauthorized table: %s", table)
		}
		for _, col := range columns {
			if !contains(opts.Whilelist.Columns, col) {
				return "", fmt.Errorf("unauthorized column: %s", col)
			}
		}
	}

	setClauses := make([]string, len(columns))
	for i, col := range columns {
		setClauses[i] = fmt.Sprintf(`"%s" = $%d`, col, i+1)
	}

	query := fmt.Sprintf(`UPDATE "%s"."%s" SET %s WHERE %s`, schema, table, strings.Join(setClauses, ", "), condition)

	if opts != nil && len(opts.Returning) > 0 {
		returningCols := make([]string, len(opts.Returning))
		for i, col := range opts.Returning {
			returningCols[i] = fmt.Sprintf(`"%s"`, col)
		}
		query += fmt.Sprintf(" RETURNING %s", strings.Join(returningCols, ", "))
	}

	return query, nil
}

// BuildDeleteQuery constructs a safe DELETE query with optional RETURNING fields.
func BuildDeleteQuery(schema, table, condition string, opts *QueryBuilderOptions) (string, error) {
	if opts != nil {
		if !contains(opts.Whilelist.Schemas, schema) {
			return "", fmt.Errorf("unauthorized schema: %s", schema)
		}
		if !contains(opts.Whilelist.Tables, table) {
			return "", fmt.Errorf("unauthorized table: %s", table)
		}
	}

	query := fmt.Sprintf(`DELETE FROM "%s"."%s" WHERE %s`, schema, table, condition)

	if opts != nil && len(opts.Returning) > 0 {
		returningCols := make([]string, len(opts.Returning))
		for i, col := range opts.Returning {
			if !contains(opts.Whilelist.Columns, col) {
				return "", fmt.Errorf("unauthorized returning column: %s", col)
			}
			returningCols[i] = fmt.Sprintf(`"%s"`, col)
		}
		query += fmt.Sprintf(" RETURNING %s", strings.Join(returningCols, ", "))
	}

	return query, nil
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
