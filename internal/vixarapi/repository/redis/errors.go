package session

import "errors"

// common errors for redis repository layer
var (
	ErrNotFound = errors.New("data not found")
	ErrNilData  = errors.New("got unexpected nil data")
)
