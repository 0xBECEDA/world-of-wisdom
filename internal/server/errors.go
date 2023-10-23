package server

import "errors"

var (
	ErrEmptyMessage      = errors.New("empty message")
	ErrFailedToMarshal   = errors.New("error marshaling timestamp")
	ErrFailedToGetRand   = errors.New("error get rand from cache")
	ErrChallengeUnsolved = errors.New("challenge is not solved")
	ErrUnknownRequest    = errors.New("unknown request received")
	ErrFailedToUnmarshal = errors.New("failed to unmarshal message data")
)
