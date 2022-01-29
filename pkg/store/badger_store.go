package store

import (
	"github.com/dgraph-io/badger"
)

// BadgerStore implement the Store interface using badger as storage driver
type BadgerStore struct {
	db *badger.DB
}

// Set adds a key-value pair to the badger database.
func (s BadgerStore) Set(k string, v []byte) error {
	if err := CheckKeyAndValue(k, v); err != nil {
		return err
	}
	err := s.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(k), v)
	})
	if err != nil {
		return err
	}
	return nil
}

// Get looks for key and returns corresponding value.
func (s BadgerStore) Get(k string) ([]byte, error) {
	if err := CheckKey(k); err != nil {
		return nil, err
	}
	var val []byte
	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(k))
		if err != nil {
			return err
		}
		// TODO: explore if possible to avoid copying
		val, err = item.ValueCopy(nil)
		return err
	})
	if err == badger.ErrKeyNotFound {
		return nil, ErrKeyNotFound
	} else if err != nil {
		return nil, err
	}
	return val, nil
}

// Delete deletes a key.
func (s BadgerStore) Delete(k string) error {
	if err := CheckKey(k); err != nil {
		return err
	}
	return s.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(k))
	})
}

// Close closes the store.
// It must be called to make sure that all pending updates being synced to disk.
func (s BadgerStore) Close() error {
	return s.db.Close()
}

// BadgerStoreOptions are the options for the BadgerStore.
type BadgerStoreOptions struct {
	// Directory for storing the DB files.
	Dir string
}

// BadgerStoreDefaultOptions is an Options object with default values.
var BadgerStoreDefaultOptions = BadgerStoreOptions{
	Dir: "BadgerDB",
}

// NewBadgerStore creates a new BadgerDB store.
func NewBadgerStore(options BadgerStoreOptions) (Store, error) {
	result := BadgerStore{}

	// Set default values
	if options.Dir == "" {
		options.Dir = BadgerStoreDefaultOptions.Dir
	}

	// Open DB
	db, err := badger.Open(badger.DefaultOptions(options.Dir))
	if err != nil {
		return result, err
	}
	result.db = db

	return result, nil
}
