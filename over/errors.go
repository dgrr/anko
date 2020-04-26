package over

import "errors"

var (
	ErrMethodNotImplemented = errors.New("method not implemented")
	ErrNoMatch              = errors.New("values doesn't match")
)
