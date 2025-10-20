package pluginsdk

import "errors"

// Common errors that plugins can return
var (
	// ErrNotFound is returned when an entity is not found
	ErrNotFound = errors.New("entity not found")

	// ErrNotSupported is returned when an operation is not supported
	ErrNotSupported = errors.New("operation not supported")

	// ErrInvalidInput is returned when input validation fails
	ErrInvalidInput = errors.New("invalid input")

	// ErrReadOnly is returned when attempting to modify a read-only entity
	ErrReadOnly = errors.New("entity is read-only")
)
