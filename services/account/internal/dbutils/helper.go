package dbutils

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ashish19912009/zrms/services/account/internal/constants"
	"github.com/ashish19912009/zrms/services/account/internal/logger"
)

func CheckDBConn(db *sql.DB, methodName string) error {
	logCtx := logger.BaseLogContext("layer", constants.Repository, "method", methodName)
	if db == nil {
		logger.Error(constants.DBConnectionFailure, nil, logCtx)
		return errors.New(constants.DBConnectionNil)
	}
	return nil
}

func ExecuteAndScanRow(ctx context.Context, methodName string, db *sql.DB, query string, args []any, dest ...any) error {
	row := db.QueryRowContext(ctx, query, args...)
	if row.Err() != nil {
		logCtx := logger.BaseLogContext("layer", constants.Repository, "method", methodName)
		logger.Fatal(constants.DBQueryFailed, row.Err(), logCtx)
		return row.Err()
	}
	if err := row.Scan(dest...); err != nil {
		logCtx := logger.BaseLogContext("layer", constants.Repository, "method", methodName)
		logger.Fatal(constants.FailedToRetrv, row.Err(), logCtx)
		return err
	}
	return nil
}

func ExecuteAndScanRowTx(ctx context.Context, methodName string, tx *sql.Tx, query string, args []any, dest ...any) error {
	row := tx.QueryRowContext(ctx, query, args...)
	if row.Err() != nil {
		logCtx := logger.BaseLogContext("layer", constants.Repository, "method", methodName)
		logger.Fatal(constants.DBQueryFailed, row.Err(), logCtx)
		return row.Err()
	}
	if err := row.Scan(dest...); err != nil {
		logCtx := logger.BaseLogContext("layer", constants.Repository, "method", methodName)
		logger.Fatal(constants.FailedToRetrv, err, logCtx)
		return err
	}
	return nil
}
