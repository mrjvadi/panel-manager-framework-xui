package core

import "errors"

var (
	ErrNotImplemented  = errors.New("not implemented")
	ErrExtNotSupported = errors.New("extension not supported for this panel")
	ErrCircuitOpen     = errors.New("circuit breaker is open")
)
