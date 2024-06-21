package cache

type Repository interface {
	GetIterator() Iterator

	Get(uuid []byte) ([]byte, error)

	Set(key []byte, value []byte, expireSeconds int) error

	Del(key []byte) (affected bool)

	EntryCount() int64
	HitCount() int64
	MissCount() int64
}
