package store

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgraph-io/badger/v3"
)

var (
	ErrBadgerConnection = errors.New("badger connection failed")
	ErrBadgerOperation  = errors.New("badger operation failed")
)

type BadgerStore struct {
	db      *badger.DB
	timeout time.Duration
}

func NewBadgerStore(config *BadgerConfig) (*BadgerStore, error) {
	if config == nil {
		return nil, errors.New("badger config cannot be nil")
	}

	opts := badger.DefaultOptions(config.Dir)
	opts.Logger = nil // Disable internal logging unless configured

	if !config.SyncWrites {
		opts.SyncWrites = false
	}

	db, err := badger.Open(opts)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrBadgerConnection, err)
	}

	return &BadgerStore{
		db:      db,
		timeout: 2 * time.Second, // Default timeout
	}, nil
}

func (b *BadgerStore) Set(key string, value interface{}) error {
	return b.SetWithTTL(key, value, 0)
}

func (b *BadgerStore) SetWithTTL(key string, value interface{}, ttl time.Duration) error {
	return b.db.Update(func(txn *badger.Txn) error {
		var val []byte
		switch v := value.(type) {
		case []byte:
			val = v
		case string:
			val = []byte(v)
		default:
			return fmt.Errorf("%w: unsupported value type", ErrBadgerOperation)
		}

		e := badger.NewEntry([]byte(key), val)
		if ttl > 0 {
			e.WithTTL(ttl)
		}
		return txn.SetEntry(e)
	})
}

func (b *BadgerStore) Get(key string) (interface{}, error) {
	var valCopy []byte
	err := b.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			valCopy = append([]byte{}, val...)
			return nil
		})
	})

	if err != nil {
		if errors.Is(err, badger.ErrKeyNotFound) {
			return nil, ErrKeyNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrBadgerOperation, err)
	}
	return valCopy, nil
}

func (b *BadgerStore) Delete(key string) error {
	return b.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	})
}

func (b *BadgerStore) Exists(key string) (bool, error) {
	err := b.db.View(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte(key))
		return err
	})

	if err == nil {
		return true, nil
	}
	if errors.Is(err, badger.ErrKeyNotFound) {
		return false, nil
	}
	return false, fmt.Errorf("%w: %v", ErrBadgerOperation, err)
}

func (b *BadgerStore) Keys(prefix string) ([]string, error) {
	var keys []string
	err := b.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.Prefix = []byte(prefix)
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			keys = append(keys, string(it.Item().Key()))
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrBadgerOperation, err)
	}
	return keys, nil
}

func (b *BadgerStore) FlushAll() error {
	return b.db.DropAll()
}

func (b *BadgerStore) Close() error {
	return b.db.Close()
}

// Interface compliance check
var _ InMemoryStore = (*BadgerStore)(nil)
