package cache

import "errors"

var (
	ErrNotImplemented  = errors.New("not implemented")
	ErrUnsafePath      = errors.New("unsafe cache path")
	ErrInvalidCacheDir = errors.New("invalid cache directory")
)
