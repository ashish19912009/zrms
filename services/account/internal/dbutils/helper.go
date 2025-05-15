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

type tagCacheKey struct {
	Type reflect.Type
	Tag  string
}

var (
	scanTagCache   = make(map[string]map[reflect.Type]map[string][]int) // tag -> struct type -> tag name -> field index path
	scanTagMu      sync.RWMutex
	tagCache       sync.Map
	structTagCache sync.Map
)

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

func ExecuteAndScanRowTx(ctx context.Context, methodName string, tx *sql.Tx, query string, args []any, destStruct any, returningCols ...string) error {
	row := tx.QueryRowContext(ctx, query, args...)
	if row.Err() != nil {
		logCtx := logger.BaseLogContext("layer", constants.Repository, "method", methodName)
		logger.Error(constants.DBQueryFailed, row.Err(), logCtx)
		return row.Err()
	}

	scanTargets, err := StructScanDestByTag(destStruct, returningCols, "json") // or "db" depending on your tag convention
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

func StructScanDestByTag(target any, fields []string, tag string) ([]any, error) {
	t := reflect.TypeOf(target)
	v := reflect.ValueOf(target)

	if t.Kind() != reflect.Ptr || v.IsNil() {
		return nil, errors.New("target must be a non-nil pointer to a struct")
	}

	t = t.Elem()
	v = v.Elem()

	if t.Kind() != reflect.Struct {
		return nil, errors.New("target must point to a struct")
	}

	// Cache lookup
	scanTagMu.RLock()
	typeCache, ok := scanTagCache[tag][t]
	scanTagMu.RUnlock()

	if !ok {
		// Not cached yet — build and cache
		typeCache = make(map[string][]int)
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			tagVal := field.Tag.Get(tag)
			if tagVal != "" {
				typeCache[tagVal] = field.Index
			}
		}
		scanTagMu.Lock()
		if scanTagCache[tag] == nil {
			scanTagCache[tag] = make(map[reflect.Type]map[string][]int)
		}
		scanTagCache[tag][t] = typeCache
		scanTagMu.Unlock()
	}

	var dests []any
	for _, name := range fields {
		if idxPath, found := typeCache[name]; found {
			fieldVal := v.FieldByIndex(idxPath)
			dests = append(dests, fieldVal.Addr().Interface())
		} else {
			return nil, fmt.Errorf("field with tag %q not found in struct %s", name, t.Name())
		}
	}

	return dests, nil
}

func MapStructFieldsByTag(columns []string, inputStruct any, tag string) (map[string]any, error) {
	if inputStruct == nil {
		return nil, errors.New("inputStruct cannot be nil")
	}

	v := reflect.ValueOf(inputStruct)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return nil, errors.New("inputStruct must be a non-nil pointer to a struct")
	}
	v = v.Elem()
	t := v.Type()

	if t.Kind() != reflect.Struct {
		return nil, errors.New("inputStruct must point to a struct")
	}

	// Cache key based on type + tag
	cacheKey := tagCacheKey{Type: t, Tag: tag}

	tagMapInterface, ok := structTagCache.Load(cacheKey)
	var tagMap map[string]int
	if ok {
		tagMap = tagMapInterface.(map[string]int)
	} else {
		tagMap = make(map[string]int)
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			tagValue := field.Tag.Get(tag)
			if tagValue != "" {
				tagMap[tagValue] = i
			}
		}
		structTagCache.Store(cacheKey, tagMap)
	}

	result := make(map[string]any)
	for _, col := range columns {
		fieldIdx, exists := tagMap[col]
		if !exists {
			return nil, fmt.Errorf("field with tag %q not found in struct", col)
		}
		result[col] = v.Field(fieldIdx).Interface()
	}

	return result, nil
}
