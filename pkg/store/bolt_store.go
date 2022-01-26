package store

import (
	"fmt"
	"github.com/fort-io/fort/pkg/fileutil"
	"github.com/golang/glog"
	bolt "go.etcd.io/bbolt"
)

var defaultBucket = []byte("default")

type BoltStore struct {
	db  *bolt.DB
	opt *bolt.Options
}

func NewBoltStore() Store {
	return &BoltStore{}
}

func (s *BoltStore) Init(datadir string) {
	var err error
	if err = fileutil.TouchDirAll(datadir); err != nil {
		glog.Fatal(err)
	}
	s.db, err = bolt.Open(datadir+"/data", 0666, s.opt)
	if err != nil {
		glog.Fatal(err)
	}
	err = s.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(defaultBucket)
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
	if err != nil {
		glog.Fatal(err)
	}
}

func (s *BoltStore) Get(key []byte) (val []byte, err error) {
	err = s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(defaultBucket)
		val = bucket.Get(key)
		return nil
	})
	return val, err
}

func (s *BoltStore) Put(key []byte, val []byte) (err error) {
	return s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(defaultBucket)
		return bucket.Put(key, val)
	})
}

func (s *BoltStore) Delete(key []byte) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(defaultBucket)
		return bucket.Delete(key)
	})
}

func (s *BoltStore) Close() error {
	return s.db.Close()
}
