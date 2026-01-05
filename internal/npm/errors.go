package npm

import "errors"

var (
	// ErrNotImplemented is a placeholder for unimplemented functionality.
	ErrNotImplemented = errors.New("not implemented")
	// ErrPackageNotFound indicates the npm package was not found in the registry.
	ErrPackageNotFound = errors.New("package not found")
)
