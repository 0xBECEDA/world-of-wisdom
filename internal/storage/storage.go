package storage

import (
	"log"
	"sync"
)

type DB interface {
	Add(key uint64)
	Get(key uint64) (uint64, error)
	Delete(key uint64) error
}

type Storage struct {
	memoryDB map[uint64]struct{}
	rw       *sync.RWMutex
}

func NewStorage() *Storage {
	return &Storage{
		memoryDB: make(map[uint64]struct{}),
		rw:       &sync.RWMutex{},
	}
}

func (r *Storage) Add(key uint64) {
	r.rw.Lock()
	defer r.rw.Unlock()

	r.memoryDB[key] = struct{}{}
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

func (r *Storage) Delete(key uint64) error {
	log.Printf("deleting key: %d", key)

	r.rw.Lock()
	defer r.rw.Unlock()

	if _, ok := r.memoryDB[key]; !ok {
		return ErrKeyNotFound
	}

	delete(r.memoryDB, key)
	return nil
}
