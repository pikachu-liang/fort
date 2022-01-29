package store

import "errors"

type Store interface {
	Set(k string, v []byte) error
	Get(k string) (v []byte, err error)
	Delete(k string) error
	Close() error
}

var (
	ErrKeyNotFound = errors.New("key not found")
	ErrKeyInvalid  = errors.New("key can not be empty")
	ErrValInvalid  = errors.New("value can not be nil")
)
