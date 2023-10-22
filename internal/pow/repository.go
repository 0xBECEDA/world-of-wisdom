package pow

import (
	"world-of-wisdom/internal/storage"
)

type Repository interface {
	AddIndicator(indicator uint64)
	IndicatorExists(indicator uint64) bool
	DeleteIndicator(indicator uint64) error
}

type RepositoryHashCash struct {
	storage storage.DB
}

func NewHashCashRepository(storage storage.DB) *RepositoryHashCash {
	return &RepositoryHashCash{
		storage: storage,
	}
}

func (repo *RepositoryHashCash) AddIndicator(indicator uint64) {
	repo.storage.Add(indicator)
}

func (repo *RepositoryHashCash) IndicatorExists(indicator uint64) bool {
	_, err := repo.storage.Get(indicator)
	if err != nil {
		return false
	}
	return true
}

func (repo *RepositoryHashCash) DeleteIndicator(indicator uint64) error {
	return repo.storage.Delete(indicator)
}
