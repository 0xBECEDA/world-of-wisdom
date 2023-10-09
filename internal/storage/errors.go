package storage

import "errors"

var (
	ErrKeyNotFound = errors.New("key does not exist in storage")
)
