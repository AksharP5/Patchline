package cache

import "errors"

var (
	// ErrNotImplemented is a placeholder for unimplemented functionality.
	ErrNotImplemented = errors.New("not implemented")
	// ErrUnsafePath indicates a cache path fails safety checks.
	ErrUnsafePath = errors.New("unsafe cache path")
	// ErrInvalidCacheDir indicates a cache directory is missing or invalid.
	ErrInvalidCacheDir = errors.New("invalid cache directory")
)
