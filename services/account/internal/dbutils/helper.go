package dbutils

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sync"

	"github.com/ashish19912009/zrms/services/account/internal/constants"
	"github.com/ashish19912009/zrms/services/account/internal/logger"
)

var tagCache sync.Map

func CheckDBConn(db *sql.DB, methodName string) error {
	logCtx := logger.BaseLogContext("layer", constants.Repository, "method", methodName)
	if db == nil {
		logger.Error(constants.DBConnectionFailure, nil, logCtx)
		return errors.New(constants.DBConnectionNil)
	}
	return nil
}

func ConvertStringMapToJson(themeSettingsMap map[string]interface{}) ([]byte, error) {
	themeSettingsJSON, err := json.Marshal(themeSettingsMap)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal theme settings: %w", err)
	}
	return themeSettingsJSON, nil
}

func ExecuteAndScanRow(ctx context.Context, methodName string, db *sql.DB, query string, args []any, destStruct any, returningCols ...string) error {
	row := db.QueryRowContext(ctx, query, args...)
	if row.Err() != nil {
		logCtx := logger.BaseLogContext("layer", constants.Repository, "method", methodName)
		logger.Error(constants.DBQueryFailed, row.Err(), logCtx)
		return row.Err()
	}

	scanTargets, err := StructScanDestByTag(destStruct, returningCols, "json") // or "db"
	if err != nil {
		logCtx := logger.BaseLogContext("layer", constants.Repository, "method", methodName)
		logger.Error("failed to map struct for scanning", err, logCtx)
		return err
	}

	if err := row.Scan(scanTargets...); err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		logCtx := logger.BaseLogContext("layer", constants.Repository, "method", methodName)
		logger.Error(constants.FailedToRetrv, err, logCtx)
		return err
	}

	return nil
}

func ExecuteAndScanRowTx(ctx context.Context, methodName string, tx *sql.Tx, query string, args []any, dest ...any) error {
	row := tx.QueryRowContext(ctx, query, args...)
	if row.Err() != nil {
		logCtx := logger.BaseLogContext("layer", constants.Repository, "method", methodName)
		logger.Error(constants.DBQueryFailed, row.Err(), logCtx)
		return row.Err()
	}
	if err := row.Scan(dest...); err != nil {
		logCtx := logger.BaseLogContext("layer", constants.Repository, "method", methodName)
		logger.Error(constants.FailedToRetrv, err, logCtx)
		return err
	}
	return nil
}

// StructToValuesByTag extracts values from a struct based on a list of field tags (e.g., json or db)
func StructToValuesByTag(input interface{}, dbColumns []string, tagKey string) ([]any, error) {
	v := reflect.ValueOf(input)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()

	cacheKey := t.String() + "::" + tagKey

	var tagToIndex map[string]int
	if cached, ok := tagCache.Load(cacheKey); ok {
		tagToIndex = cached.(map[string]int)
	} else {
		tagToIndex = make(map[string]int)
		for i := 0; i < t.NumField(); i++ {
			tag := t.Field(i).Tag.Get(tagKey)
			if tag != "" {
				tagToIndex[tag] = i
			}
		}
		tagCache.Store(cacheKey, tagToIndex)
	}

	values := make([]any, 0, len(dbColumns))
	for _, col := range dbColumns {
		index, ok := tagToIndex[col]
		if !ok {
			return nil, fmt.Errorf("column '%s' not found in struct via tag '%s'", col, tagKey)
		}
		values = append(values, v.Field(index).Interface())
	}

	return values, nil
}

func StructScanDestByTag(ptr interface{}, dbColumns []string, tagKey string) ([]any, error) {
	v := reflect.ValueOf(ptr)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected pointer to struct")
	}
	v = v.Elem()
	t := v.Type()

	tagToField := make(map[string]int)
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get(tagKey)
		if tag != "" {
			tagToField[tag] = i
		}
	}

	values := make([]any, len(dbColumns))
	for i, col := range dbColumns {
		index, ok := tagToField[col]
		if !ok {
			return nil, fmt.Errorf("column '%s' not found in struct for tag '%s'", col, tagKey)
		}
		field := v.Field(index)
		if !field.CanAddr() {
			return nil, fmt.Errorf("cannot take address of field '%s'", col)
		}
		values[i] = field.Addr().Interface()
	}

	return values, nil
}
