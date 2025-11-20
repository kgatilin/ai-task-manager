package logger

// Level represents the logging level
type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

// Logger defines the interface for logging in the task manager.
// This is the local equivalent of errors.Logger to avoid dependency on errors.
type Logger interface {
	// Debug logs a debug-level message.
	// Debug messages are only shown when debug logging is enabled.
	Debug(msg string, keysAndValues ...interface{})

	// Info logs an info-level message.
	// Info messages are shown during normal operation.
	Info(msg string, keysAndValues ...interface{})

	// Warn logs a warning-level message.
	// Warnings indicate potential issues that don't prevent operation.
	Warn(msg string, keysAndValues ...interface{})

	// Error logs an error-level message.
	// Errors indicate failures that may affect functionality.
	Error(msg string, keysAndValues ...interface{})

	// SetLevel sets the minimum logging level.
	SetLevel(level Level)

	// GetLevel returns the current logging level.
	GetLevel() Level
}
