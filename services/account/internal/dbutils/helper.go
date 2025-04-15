package dbutils

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ashish19912009/zrms/services/account/internal/constants"
	"github.com/ashish19912009/zrms/services/account/internal/logger"
)

func CheckDBConn(db *sql.DB, logCtx map[string]interface{}) error {
	if db == nil {
		logger.Error(constants.DBConnectionFailure, nil, logCtx)
		return errors.New(constants.DBConnectionNil)
	}
	return nil
}

func ExecuteAndScanRow(ctx context.Context, db *sql.DB, query string, args []any, dest ...any) error {
	row := db.QueryRowContext(ctx, query, args...)
	if row.Err() != nil {
		return row.Err()
	}
	if err := row.Scan(dest...); err != nil {
		return err
	}
	return nil
}
