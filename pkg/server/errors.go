package server

import "errors"

var (
	ErrFailedToAddIndicator    = errors.New("error adding indicator")
	ErrEmptyMessage            = errors.New("empty message")
	ErrFailedToMarshal         = errors.New("error marshaling timestamp")
	ErrFailedToDecodeRand      = errors.New("error decode rand")
	ErrFailedToGetRand         = errors.New("error get rand from cache")
	ErrChallengeUnsolved       = errors.New("challenge is not solved")
	ErrUnknownRequest          = errors.New("unknown request received")
	ErrFailedToUnmarshal       = errors.New("failed to unmarshal")
	ErrFailedToRemoveIndicator = errors.New("failed to remove indicator")
)
