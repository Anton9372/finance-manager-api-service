package freecache

import (
	"finance-manager-api-service/pkg/cache"
	"github.com/coocood/freecache"
	"sync"
)

type repository struct {
	sync.Mutex
	cache *freecache.Cache
}

func NewCacheRepo(size int) cache.Repository {
	return &repository{cache: freecache.NewCache(size)}
}

func (r *repository) GetIterator() cache.Iterator {
	return &iterator{r.cache.NewIterator()}
}

func (r *repository) Get(uuid []byte) ([]byte, error) {
	r.Lock()
	defer r.Unlock()

	return r.cache.Get(uuid)
}

func (r *repository) Set(key []byte, value []byte, expireSeconds int) error {
	r.Lock()
	defer r.Unlock()

	return r.cache.Set(key, value, expireSeconds)
}

func (r *repository) Del(key []byte) (affected bool) {
	r.Lock()
	defer r.Unlock()

	return r.cache.Del(key)
}

func (r *repository) EntryCount() int64 {
	r.Lock()
	defer r.Unlock()

	return r.cache.EntryCount()
}

func (r *repository) HitCount() int64 {
	r.Lock()
	defer r.Unlock()

	return r.cache.HitCount()
}

func (r *repository) MissCount() int64 {
	r.Lock()
	defer r.Unlock()

	return r.cache.MissCount()
}
