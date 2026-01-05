package npm

import "errors"

var (
	ErrNotImplemented  = errors.New("not implemented")
	ErrPackageNotFound = errors.New("package not found")
)
