package pow

import (
	"world-of-wisdom/internal/storage"
)

type Repository interface {
	Add(indicator uint64)
	Exists(indicator uint64) bool
	Delete(indicator uint64)
}

type RepositoryHashCash struct {
	storage storage.DB
}

func NewHashCashRepository(storage storage.DB) *RepositoryHashCash {
	return &RepositoryHashCash{
		storage: storage,
	}
}

func (repo *RepositoryHashCash) Add(indicator uint64) {
	repo.storage.Add(indicator)
}

func (repo *RepositoryHashCash) Exists(indicator uint64) bool {
	_, err := repo.storage.Get(indicator)
	if err != nil {
		return false
	}
	return true
}

func (repo *RepositoryHashCash) Delete(indicator uint64) {
	repo.storage.Delete(indicator)
}
