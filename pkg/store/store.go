package store

type Store interface {
	Init(datadir string)
	Close() error
	Get(key []byte) (val []byte, err error)
	Put(key []byte, val []byte) error
	Delete(key []byte) error
}
