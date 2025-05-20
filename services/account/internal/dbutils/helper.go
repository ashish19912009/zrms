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
type tagCacheDirectKey struct {
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

func ExecuteAndScanRow(
	ctx context.Context,
	methodName string,
	db *sql.DB,
	query string,
	args []any,
	destStruct any,
	returningCols ...string,
) error {
	row := db.QueryRowContext(ctx, query, args...)
	if row.Err() != nil {
		logCtx := logger.BaseLogContext("layer", constants.Repository, "method", methodName)
		logger.Error(constants.DBQueryFailed, row.Err(), logCtx)
		return row.Err()
	}

	// Get scan targets (now includes JSON-aware handling)
	scanTargets, jsonFieldIndexes, err := StructScanDestWithJSON(destStruct, returningCols, "json")
	if err != nil {
		logCtx := logger.BaseLogContext("layer", constants.Repository, "method", methodName)
		logger.Error("failed to map struct for scanning", err, logCtx)
		return err
	}

	// Perform the initial scan
	if err := row.Scan(scanTargets...); err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		logCtx := logger.BaseLogContext("layer", constants.Repository, "method", methodName)
		logger.Error(constants.FailedToRetrv, err, logCtx)
		return err
	}

	// Post-process JSON fields
	if err := unmarshalJSONFields(destStruct, jsonFieldIndexes, scanTargets); err != nil {
		logCtx := logger.BaseLogContext("layer", constants.Repository, "method", methodName)
		logger.Error("failed to unmarshal JSON fields", err, logCtx)
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

// MapValues extracts values from a struct based on the given tag and returns them
// in the order of the provided columns.
// columns: Ordered list of fields you want.

// input: Any struct or pointer to struct.

// tag: e.g., "json", "db" — whatever tags your struct uses.
func MapValues(columns []string, input any, tag string) ([]any, error) {
	// Extract field map using your existing utility
	fieldMap, err := MapStructFieldsByTag(columns, input, tag)
	if err != nil {
		return nil, fmt.Errorf("failed to map struct fields: %w", err)
	}

	values := make([]any, len(columns))
	for i, col := range columns {
		val, ok := fieldMap[col]
		if !ok {
			return nil, fmt.Errorf("missing value for column: %s", col)
		}
		values[i] = val
	}

	return values, nil
}

// leaner version of MapValues that avoids creating a map[string]any internally:
func MapValuesDirect(input any, columns []string, tag string) ([]any, error) {
	if input == nil {
		return nil, errors.New("input cannot be nil")
	}

	v := reflect.ValueOf(input)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()
	if t.Kind() != reflect.Struct {
		return nil, errors.New("input must be a struct or pointer to struct")
	}

	cacheKey := tagCacheDirectKey{Type: t, Tag: tag}
	var tagToIndex map[string]int
	if cached, ok := structTagCache.Load(cacheKey); ok {
		tagToIndex = cached.(map[string]int)
	} else {
		tagToIndex = make(map[string]int)
		for i := 0; i < t.NumField(); i++ {
			tagVal := t.Field(i).Tag.Get(tag)
			if tagVal != "" {
				tagToIndex[tagVal] = i
			}
		}
		structTagCache.Store(cacheKey, tagToIndex)
	}

	values := make([]any, len(columns))
	for i, col := range columns {
		fieldIndex, ok := tagToIndex[col]
		if !ok {
			return nil, fmt.Errorf("tag '%s' not found in struct", col)
		}
		values[i] = v.Field(fieldIndex).Interface()
	}
	return values, nil
}

// StructScanDestWithJSON prepares scan targets and tracks JSON fields
func StructScanDestWithJSON(dest any, columns []string, tagName string) ([]any, map[int]reflect.StructField, error) {
	destValue := reflect.ValueOf(dest)
	if destValue.Kind() != reflect.Ptr || destValue.Elem().Kind() != reflect.Struct {
		return nil, nil, fmt.Errorf("destination must be a pointer to struct")
	}

	destElem := destValue.Elem()
	destType := destElem.Type()

	scanTargets := make([]any, len(columns))
	jsonFields := make(map[int]reflect.StructField)

	for i, col := range columns {
		fieldFound := false
		for j := 0; j < destType.NumField(); j++ {
			field := destType.Field(j)
			tag := field.Tag.Get(tagName)
			if tag == col {
				fieldValue := destElem.Field(j)
				// Check for JSON-compatible fields (map or *map)
				if isJSONField(field.Type) {
					var jsonData []byte
					scanTargets[i] = &jsonData
					jsonFields[i] = field
				} else {
					scanTargets[i] = fieldValue.Addr().Interface()
				}
				fieldFound = true
				break
			}
		}
		if !fieldFound {
			return nil, nil, fmt.Errorf("no matching struct field for column: %s", col)
		}
	}

	return scanTargets, jsonFields, nil
}

func isJSONField(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Map:
		return t.Key().Kind() == reflect.String
	case reflect.Ptr:
		return t.Elem().Kind() == reflect.Map && t.Elem().Key().Kind() == reflect.String
	default:
		return false
	}
}

func unmarshalJSONFields(dest any, jsonFieldIndexes map[int]reflect.StructField, scanTargets []any) error {
	destValue := reflect.ValueOf(dest).Elem()

	for i, field := range jsonFieldIndexes {
		jsonDataPtr, ok := scanTargets[i].(*[]byte)
		if !ok {
			return fmt.Errorf("expected *[]byte for JSON field %s", field.Name)
		}

		fieldValue := destValue.FieldByName(field.Name)
		if !fieldValue.IsValid() {
			return fmt.Errorf("no such field: %s", field.Name)
		}

		// Skip if NULL or empty
		if len(*jsonDataPtr) == 0 {
			continue
		}

		// Handle both map and *map
		switch {
		case fieldValue.Kind() == reflect.Map:
			if fieldValue.IsNil() {
				fieldValue.Set(reflect.MakeMap(field.Type))
			}
			target := fieldValue.Interface() // Non-pointer map
			if err := json.Unmarshal(*jsonDataPtr, &target); err != nil {
				return fmt.Errorf("failed to unmarshal JSON for field %s: %v", field.Name, err)
			}

		case fieldValue.Kind() == reflect.Ptr && fieldValue.Type().Elem().Kind() == reflect.Map:
			if fieldValue.IsNil() {
				fieldValue.Set(reflect.New(fieldValue.Type().Elem()))
				fieldValue.Elem().Set(reflect.MakeMap(fieldValue.Type().Elem()))
			}
			target := fieldValue.Interface() // Pointer to map
			if err := json.Unmarshal(*jsonDataPtr, target); err != nil {
				return fmt.Errorf("failed to unmarshal JSON for field %s: %v", field.Name, err)
			}

		default:
			return fmt.Errorf("field %s is not a map or *map", field.Name)
		}
	}
	return nil
}
