package opencode

import "errors"

var (
	ErrNotImplemented = errors.New("not implemented")
	ErrConfigNotFound = errors.New("config not found")
)
