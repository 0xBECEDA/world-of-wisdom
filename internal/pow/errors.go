package pow

import "errors"

var (
	ErrIndicatorNotFound     = errors.New("indicator not found")
	ErrMaxIterationsExceeded = errors.New("max iterations exceeded")
)
