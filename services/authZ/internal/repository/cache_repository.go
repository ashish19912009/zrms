package repository

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/ashish19912009/zrms/services/authZ/internal/constants"
	"github.com/ashish19912009/zrms/services/authZ/internal/logger"
	"github.com/ashish19912009/zrms/services/authZ/internal/store"
	"github.com/klauspost/compress/zstd"
	"google.golang.org/protobuf/proto"
)

type CacheRepository interface {
	StoreWithTTL(ctx context.Context, tenantPrefix, resourceActionPostfix string, msg proto.Message, expire time.Duration) error
	Store(ctx context.Context, tenantPrefix, resourceActionPostfix string, msg proto.Message) error
	Check(ctx context.Context, tenantPrefix, resourceActionPostfix string) (*bool, error)
	Get(ctx context.Context, tenantPrefix, resourceActionPostfix string, out proto.Message) error
	Delete(ctx context.Context, tenantPrefix, resourceActionPostfix string) error
}

type cacheRepository struct {
	store   store.InMemoryStore
	encoder *zstd.Encoder
	decoder *zstd.Decoder
}

func NewCacheRepository(s store.InMemoryStore) (CacheRepository, error) {
	encoder, err := zstd.NewWriter(nil)
	if err != nil {
		return nil, fmt.Errorf(constants.ZSTDEncodingFailed, err)
	}
	decoder, err := zstd.NewReader(nil)
	if err != nil {
		return nil, fmt.Errorf(constants.ZSTDDecodingFailed, err)
	}
	return &cacheRepository{
		store:   s,
		encoder: encoder,
		decoder: decoder,
	}, nil
}

func (r *cacheRepository) StoreWithTTL(ctx context.Context, tenantPrefix, resourceActionPostfix string, msg proto.Message, expiry time.Duration) error {
	method := constants.Methods.StoreWithTTL
	key := r.key(tenantPrefix, resourceActionPostfix)

	// Debug logging
	// logger.Debug("Storing cache entry", map[string]interface{}{
	// 	"key":     key,
	// 	"message": msg, // Use String() for proto messages
	// 	"expiry":  expiry,
	// })

	data, err := proto.Marshal(msg)
	if err != nil {
		logger.Error(constants.FailedToMarshal, err, map[string]interface{}{
			"method": method,
			"key":    key,
		})
		return fmt.Errorf("%s: %w", constants.FailedToMarshal, err)
	}

	compressed := r.encoder.EncodeAll(data, make([]byte, 0, len(data))) // Pre-allocate with capacity

	var storeErr error
	if expiry > 0 {
		storeErr = r.store.SetWithTTL(key, compressed, expiry)
	} else {
		storeErr = r.store.Set(key, compressed)
	}

	if storeErr != nil {
		logger.Error(constants.FailedToStoreCache, storeErr, map[string]interface{}{
			"key":    key,
			"method": method,
			"expiry": expiry,
		})
		return fmt.Errorf(constants.FailedToStoreCache, storeErr)
	}

	// Verify write by immediate read (for debugging)
	// Verification logic
	// if err := r.verifyCacheWrite(ctx, key, tenantPrefix, resourceActionPostfix, msg); err != nil {
	// 	logger.Warn(constants.VerficationFailed, map[string]interface{}{
	// 		"key": key,
	// 		"err": err.Error(),
	// 	})
	// 	return fmt.Errorf(constants.VerficationFailed, err)
	// }

	return nil
}

// func (r *cacheRepository) verifyCacheWrite(ctx context.Context, key, tenantPrefix, resourceActionPostfix string, original proto.Message) error {
// 	// Create a new instance of the same type as original message
// 	retrieved := proto.Clone(original)
// 	proto.Reset(retrieved)

// 	// Get from cache
// 	if err := r.Get(ctx, tenantPrefix, resourceActionPostfix, retrieved); err != nil {
// 		logger.Error(constants.VerficationFailed, err, map[string]interface{}{
// 			"key": key,
// 			"err": err.Error(),
// 		})
// 		return fmt.Errorf(constants.FailedToRetrieve, err)
// 	}

// 	// Compare the marshaled forms
// 	originalData, _ := proto.Marshal(original)
// 	retrievedData, _ := proto.Marshal(retrieved)

// 	if !bytes.Equal(originalData, retrievedData) {
// 		return errors.New(constants.DoesnotMatch)
// 	}
// 	// Log successful verification with actual data
// 	// logger.Debug("Cache write verification successful", map[string]interface{}{
// 	// 	"key":     key,
// 	// 	"stored":  retrieved, // Now using proper String() method
// 	// 	"service": "auth-service",
// 	// })

// 	return nil
// }

func (r *cacheRepository) Store(ctx context.Context, tenantPrefix, resourceActionPostfix string, msg proto.Message) error {
	method := constants.Methods.Store
	key := r.key(tenantPrefix, resourceActionPostfix)
	data, err := proto.Marshal(msg)
	if err != nil {
		logger.Error(constants.FailedToMarshal, err, map[string]interface{}{
			"method": method,
			"key":    key,
		})
		return fmt.Errorf(constants.FailedToMarshal, err)
	}
	err = r.store.Set(key, data)

	if err != nil {
		logger.Error(constants.FailedToStoreDecision, err, map[string]interface{}{
			"method": method,
			"key":    key,
		})
		return err
	}
	return nil
}

func (r *cacheRepository) Check(ctx context.Context, tenantPrefix, resourceActionPostfix string) (*bool, error) {
	method := constants.Methods.Check
	key := r.key(tenantPrefix, resourceActionPostfix)

	exist, err := r.store.Exists(key)

	if err != nil {
		logger.Error(constants.FailedToCheckDecision, err, map[string]interface{}{
			"method": method,
			"key":    key,
		})
		return nil, err
	}
	return &exist, nil
}

func (r *cacheRepository) Get(ctx context.Context, tenantPrefix, resourceActionPostfix string, out proto.Message) error {
	method := constants.Methods.Get
	key := r.key(tenantPrefix, resourceActionPostfix)

	raw, err := r.store.Get(key)
	if err != nil {
		if err == store.ErrKeyNotFound {
			return store.ErrKeyNotFound
		}
		logger.Error(constants.RedisOperationFailed, err, map[string]interface{}{
			"method": method,
			"key":    key,
		})
		return err
	}

	// Handle both []byte and string types for backward compatibility
	var compressed []byte
	switch v := raw.(type) {
	case []byte:
		compressed = v
	case string:
		compressed = []byte(v)
	default:
		logger.Warn(constants.InvalidCacheValue, map[string]interface{}{
			"method": method,
			"key":    key,
			"type":   reflect.TypeOf(raw),
		})
		return fmt.Errorf(constants.InvalidCacheValueType, raw)
	}

	decompressed, err := r.decoder.DecodeAll(compressed, nil)
	if err != nil {
		logger.Error(constants.DecompressionFailed, err, map[string]interface{}{
			"method": method,
			"key":    key,
		})
		return fmt.Errorf("%s: %w", constants.DecompressionFailed, err)
	}

	if err := proto.Unmarshal(decompressed, out); err != nil {
		logger.Error(constants.FailedToUnmarshal, err, map[string]interface{}{
			"method": method,
			"key":    key,
		})
		return fmt.Errorf("%s: %w", constants.FailedToUnmarshal, err)
	}

	return nil
}

func (r *cacheRepository) Delete(ctx context.Context, tenantPrefix, resourceActionPostfix string) error {
	key := r.key(tenantPrefix, resourceActionPostfix)
	err := r.store.Delete(key)
	if err != nil {
		logger.Error(constants.FailedToDeleteRshToken, err, map[string]interface{}{
			"method": constants.Methods.DeleteToken,
			"key":    key,
		})
		return err
	}
	return nil
}

func (r *cacheRepository) key(resourceActionPostfix string, tenantPrefix string) string {
	return fmt.Sprintf("%s:%s", tenantPrefix, resourceActionPostfix)
}
