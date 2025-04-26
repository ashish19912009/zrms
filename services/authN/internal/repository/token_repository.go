package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/ashish19912009/zrms/services/authN/internal/constants"
	"github.com/ashish19912009/zrms/services/authN/internal/logger"
	"github.com/ashish19912009/zrms/services/authN/internal/store"
)

type TokenRepository interface {
	StoreToken(ctx context.Context, keyName, accountID string, token string, expire time.Duration) error
	CheckToken(ctx context.Context, keyName, accountID string, token string) (bool, error)
	DeleteToken(ctx context.Context, keyName, accountID string) error
}

type tokenRepository struct {
	store store.InMemoryStore
}

func NewTokenRepository(s store.InMemoryStore) TokenRepository {
	return &tokenRepository{
		store: s,
	}
}

func (r *tokenRepository) StoreToken(ctx context.Context, keyName, accountID string, token string, expiry time.Duration) error {
	key := r.tokenKey(keyName, accountID)
	var err error

	if expiry > 0 {
		err = r.store.SetWithTTL(key, token, expiry)
	} else {
		err = r.store.Set(key, token)
	}

	if err != nil {
		logger.Error(constants.FailedToStoreRshToken, err, map[string]interface{}{
			"method":     constants.Methods.StoreToken,
			"account_id": accountID,
		})
		return err
	}

	logger.Info(constants.TokenStoredSuccessfully, map[string]interface{}{
		"method":     constants.Methods.StoreToken,
		"account_id": accountID,
	})
	return nil
}

func (r *tokenRepository) CheckToken(ctx context.Context, keyName, accountID string, token string) (bool, error) {
	key := r.tokenKey(keyName, accountID)
	val, err := r.store.Get(key)
	if err != nil {
		if err == store.ErrKeyNotFound {
			logger.Warn(constants.AuthRshTokenInvalid, map[string]interface{}{
				"method":     constants.Methods.CheckToken,
				"account_id": accountID,
			})
			return false, nil
		}
		logger.Error(constants.RedisOperationFailed, err, map[string]interface{}{
			"method":     constants.Methods.CheckToken,
			"account_id": accountID,
		})
		return false, err
	}

	strVal, ok := val.(string)
	if !ok {
		logger.Warn(constants.AuthRshTokenInvalid, map[string]interface{}{
			"method":     constants.Methods.CheckToken,
			"account_id": accountID,
		})
		return false, nil
	}

	return strVal == token, nil
}

func (r *tokenRepository) DeleteToken(ctx context.Context, keyName, accountID string) error {
	key := r.tokenKey(keyName, accountID)
	err := r.store.Delete(key)
	if err != nil {
		logger.Error(constants.FailedToDeleteRshToken, err, map[string]interface{}{
			"method":     constants.Methods.DeleteToken,
			"account_id": accountID,
		})
		return err
	}

	logger.Info(constants.TokenDeleteSuccessfully, map[string]interface{}{
		"method":     constants.Methods.DeleteToken,
		"account_id": accountID,
	})
	return nil
}

func (r *tokenRepository) tokenKey(keyName string, accountID string) string {
	return fmt.Sprintf("%s:%s", accountID, keyName)
}
