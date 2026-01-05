package snapshot

import "errors"

var (
	// ErrNotImplemented is a placeholder for unimplemented functionality.
	ErrNotImplemented = errors.New("not implemented")
	// ErrSnapshotNotFound indicates no snapshot exists for a plugin.
	ErrSnapshotNotFound = errors.New("snapshot not found")
	// ErrInvalidSnapshotDir indicates the snapshot directory is invalid.
	ErrInvalidSnapshotDir = errors.New("invalid snapshot directory")
)
