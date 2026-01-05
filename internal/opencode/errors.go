package opencode

import "errors"

var (
	// ErrNotImplemented is a placeholder for unimplemented functionality.
	ErrNotImplemented = errors.New("not implemented")
	// ErrConfigNotFound indicates a config file could not be found.
	ErrConfigNotFound = errors.New("config not found")
	// ErrPluginNotFound indicates the plugin entry was not found in config.
	ErrPluginNotFound = errors.New("plugin not found")
)
