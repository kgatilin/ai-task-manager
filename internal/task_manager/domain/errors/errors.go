// Package errors defines standard error types for task manager.
package errors

import "errors"

var (
	// ErrNotFound indicates that a requested resource was not found.
	ErrNotFound = errors.New("not found")

	// ErrInvalidArgument indicates that an argument provided to a function is invalid.
	ErrInvalidArgument = errors.New("invalid argument")

	// ErrAlreadyExists indicates that a resource already exists and cannot be created.
	ErrAlreadyExists = errors.New("already exists")

	// ErrInternal indicates an internal error occurred.
	ErrInternal = errors.New("internal error")
)
