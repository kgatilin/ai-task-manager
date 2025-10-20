package pluginsdk

// Logger provides logging capabilities to plugins.
// This interface abstracts the logging implementation from plugins,
// allowing the app layer to inject different loggers.
type Logger interface {
	// Debug logs a debug message
	Debug(format string, args ...interface{})

	// Info logs an informational message
	Info(format string, args ...interface{})

	// Warn logs a warning message
	Warn(format string, args ...interface{})

	// Error logs an error message
	Error(format string, args ...interface{})
}
