package dbutils_test

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ashish19912009/zrms/services/authZ/internal/constants"
	"github.com/ashish19912009/zrms/services/authZ/internal/dbutils"
	"github.com/stretchr/testify/assert"
)

func TestCheckDBConn(t *testing.T) {
	t.Run("returns nil when db is not nil", func(t *testing.T) {
		db, _, _ := sqlmock.New()
		defer db.Close()

		err := dbutils.CheckDBConn(db, "TestCheckDBConn")
		assert.NoError(t, err)
	})

	t.Run("returns error when db is nil", func(t *testing.T) {
		err := dbutils.CheckDBConn(nil, "TestCheckDBConn")
		assert.EqualError(t, err, constants.DBConnectionNil)
	})
}

func TestExecuteAndScanRow(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	ctx := context.Background()

	t.Run("successful scan", func(t *testing.T) {
		mock.ExpectQuery("SELECT id FROM users WHERE id = \\$1").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		var id int
		err := dbutils.ExecuteAndScanRow(ctx, "TestExecuteAndScanRow", db, "SELECT id FROM users WHERE id = $1", []any{1}, &id)
		assert.NoError(t, err)
		assert.Equal(t, 1, id)
	})

	t.Run("row error", func(t *testing.T) {
		mock.ExpectQuery("SELECT id FROM users WHERE id = \\$1").
			WithArgs(99).
			WillReturnError(errors.New("query error"))

		var id int
		err := dbutils.ExecuteAndScanRow(ctx, "TestExecuteAndScanRow", db, "SELECT id FROM users WHERE id = $1", []any{99}, &id)
		assert.Error(t, err)
	})
}

func TestExecuteAndScanRowTx(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	ctx := context.Background()

	tx, _ := db.Begin()

	t.Run("successful scan using tx", func(t *testing.T) {
		mock.ExpectQuery("SELECT name FROM users WHERE id = \\$1").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("Alice"))

		var name string
		err := dbutils.ExecuteAndScanRowTx(ctx, "TestExecuteAndScanRowTx", tx, "SELECT name FROM users WHERE id = $1", []any{1}, &name)
		assert.NoError(t, err)
		assert.Equal(t, "Alice", name)
	})
}

func TestBuildSelectQuery(t *testing.T) {
	opts := &dbutils.QueryBuilderOptions{}
	opts.Whilelist.Schemas = []string{"public"}
	opts.Whilelist.Tables = []string{"users"}
	opts.Whilelist.Columns = []string{"id", "name"}

	t.Run("valid select query", func(t *testing.T) {
		query, args, err := dbutils.BuildSelectQuery(
			"TestSelect",
			"public",
			"users",
			[]string{"id", "name"},
			map[string]any{"id": 10},
			opts,
		)
		assert.NoError(t, err)
		assert.Contains(t, query, `SELECT "id", "name" FROM "public"."users"`)
		assert.Contains(t, query, `WHERE "id" = $1`)
		assert.Equal(t, []any{10}, args)
	})

	t.Run("unauthorized schema", func(t *testing.T) {
		_, _, err := dbutils.BuildSelectQuery("TestSelect", "secret", "users", []string{"id"}, nil, opts)
		assert.Error(t, err)
	})
}

func TestBuildInsertQuery(t *testing.T) {
	opts := &dbutils.QueryBuilderOptions{}
	opts.Whilelist.Schemas = []string{"public"}
	opts.Whilelist.Tables = []string{"users"}
	opts.Whilelist.Columns = []string{"id", "name"}
	opts.Returning = []string{"id"}

	t.Run("valid insert query", func(t *testing.T) {
		query, err := dbutils.BuildInsertQuery("TestInsert", "public", "users", []string{"id", "name"}, opts)
		assert.NoError(t, err)
		assert.Contains(t, query, `INSERT INTO "public"."users" ("id", "name") VALUES ($1, $2) RETURNING "id"`)
	})

	t.Run("unauthorized column", func(t *testing.T) {
		_, err := dbutils.BuildInsertQuery("TestInsert", "public", "users", []string{"password"}, opts)
		assert.Error(t, err)
	})
}

func TestBuildUpdateQuery(t *testing.T) {
	opts := &dbutils.QueryBuilderOptions{}
	opts.Whilelist.Schemas = []string{"public"}
	opts.Whilelist.Tables = []string{"users"}
	opts.Whilelist.Columns = []string{"name", "id"}
	opts.Returning = []string{"id"}

	t.Run("valid update query", func(t *testing.T) {
		query, args, err := dbutils.BuildUpdateQuery(
			"TestUpdate",
			"public",
			"users",
			[]string{"name"},
			map[string]any{"id": 5},
			opts,
		)
		assert.NoError(t, err)
		assert.Contains(t, query, `UPDATE "public"."users" SET "name" = $1 WHERE "id" = $2`)
		assert.Equal(t, []any{"name", 5}[1], args[1])
	})

	t.Run("unauthorized update column", func(t *testing.T) {
		_, _, err := dbutils.BuildUpdateQuery("TestUpdate", "public", "users", []string{"email"}, nil, opts)
		assert.Error(t, err)
	})
}
