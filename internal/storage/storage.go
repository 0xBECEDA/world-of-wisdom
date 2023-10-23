package storage

import (
	"log"
	"sync"
	"time"
)

type DB interface {
	Add(key uint64)
	Get(key uint64) (uint64, error)
	Delete(key uint64)
	CleanUp()
}

type Storage struct {
	memoryDB map[uint64]time.Time
	keyTTL   time.Duration
	rw       *sync.RWMutex
}

func NewStorage(keyTTL time.Duration) *Storage {
	return &Storage{
		memoryDB: make(map[uint64]time.Time),
		rw:       &sync.RWMutex{},
		keyTTL:   keyTTL,
	}
}

func (r *Storage) Add(key uint64) {
	r.rw.Lock()
	defer r.rw.Unlock()

	r.memoryDB[key] = time.Now().Add(r.keyTTL)
	log.Printf("added key: %d", key)
}

func (r *Storage) Get(key uint64) (uint64, error) {
	log.Printf("getting key: %d", key)

	r.rw.RLock()
	defer r.rw.RUnlock()

	_, ok := r.memoryDB[key]
	if ok {
		return key, nil
	}

	return 0, ErrKeyNotFound
}

func (r *Storage) Delete(key uint64) {
	log.Printf("deleting key: %d", key)

	r.rw.Lock()
	defer r.rw.Unlock()

	delete(r.memoryDB, key)
}

// CleanUp cleans keys, which are expired, because challenge wasn't resolved in time.
func (r *Storage) CleanUp() {
	tick := time.NewTicker(r.keyTTL)

	for range tick.C {
		now := time.Now()
		log.Printf("clean up storage started at %v", now)

		r.rw.Lock()
		for key, ttl := range r.memoryDB {
			if ttl.Before(now) {
				delete(r.memoryDB, key)
				log.Printf("key %v expired and was deleted", key)
			}
		}
		r.rw.Unlock()
	}
}
