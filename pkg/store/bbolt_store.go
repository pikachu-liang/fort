package store

import (
	"go.etcd.io/bbolt"
)

// BBoltStore implement Store interface using bbolt as storage driver
type BBoltStore struct {
	db         *bbolt.DB
	bucketName string
}

// Set adds a key-value pair to the bbolt database.
func (s BBoltStore) Set(k string, v []byte) error {
	if err := CheckKeyAndValue(k, v); err != nil {
		return err
	}
	err := s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(s.bucketName))
		return b.Put([]byte(k), v)
	})
	if err != nil {
		return err
	}
	return nil
}

// Get looks for key and returns corresponding value.
func (s BBoltStore) Get(k string) ([]byte, error) {
	if err := CheckKey(k); err != nil {
		return nil, err
	}
	var val []byte
	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(s.bucketName))
		txData := b.Get([]byte(k))
		// TODO: explore if possible to avoid copying
		if txData != nil {
			val = make([]byte, len(txData))
			copy(val, txData)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if val == nil {
		return nil, ErrKeyNotFound
	}
	return val, nil
}

// Delete deletes a key.
func (s BBoltStore) Delete(k string) error {
	if err := CheckKey(k); err != nil {
		return err
	}
	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(s.bucketName))
		return b.Delete([]byte(k))
	})
}

// Close closes the store.
// It must be called to make sure that all pending updates being synced to disk.
func (s BBoltStore) Close() error {
	return s.db.Close()
}

// BBoltStoreOptions are the options for the BBoltStore.
type BBoltStoreOptions struct {
	// Bucket name for storing the key-value pairs.
	BucketName string
	// Path for storing the DB file.
	Path string
}

// BBoltStoreDefaultOptions is an BBoltStoreOptions object with default values.
var BBoltStoreDefaultOptions = BBoltStoreOptions{
	BucketName: "default",
	Path:       "bbolt_store.db",
}

// NewBBoltStore creates a new BBoltStore.
func NewBBoltStore(options BBoltStoreOptions) (Store, error) {
	result := BBoltStore{}
	// Set default values
	if options.BucketName == "" {
		options.BucketName = BBoltStoreDefaultOptions.BucketName
	}
	if options.Path == "" {
		options.Path = BBoltStoreDefaultOptions.Path
	}

	// Open DB
	db, err := bbolt.Open(options.Path, 0600, nil)
	if err != nil {
		return result, err
	}

	// Create a bucket if it doesn't exist yet.
	// In bbolt, key/value pairs are stored to bucket.
	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(options.BucketName))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return result, err
	}

	result.db = db
	result.bucketName = options.BucketName

	return result, nil
}
