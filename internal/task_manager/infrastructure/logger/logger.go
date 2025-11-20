// Package logger provides logging interfaces and implementations for the task manager.
package logger

import (
	"fmt"
	"log"
	"os"

	"github.com/kgatilin/ai-task-manager/internal/task_manager/domain/logger"
)

// StandardLogger is a standard library-based implementation of Logger.
// It uses Go's standard log package for output and supports level-based filtering.
type StandardLogger struct {
	level  logger.Level
	debug  *log.Logger
	info   *log.Logger
	warn   *log.Logger
	errLog *log.Logger
}

// NewStandardLogger creates a new StandardLogger with the given level.
// All log output goes to stderr.
func NewStandardLogger(level logger.Level) *StandardLogger {
	return &StandardLogger{
		level:  level,
		debug:  log.New(os.Stderr, "[DEBUG] ", log.LstdFlags),
		info:   log.New(os.Stderr, "[INFO] ", log.LstdFlags),
		warn:   log.New(os.Stderr, "[WARN] ", log.LstdFlags),
		errLog: log.New(os.Stderr, "[ERROR] ", log.LstdFlags),
	}
}

// Debug logs a debug-level message if debug logging is enabled.
func (l *StandardLogger) Debug(msg string, keysAndValues ...interface{}) {
	if l.level <= logger.LevelDebug {
		l.debug.Println(l.formatMessage(msg, keysAndValues...))
	}
}

// Info logs an info-level message if info logging is enabled.
func (l *StandardLogger) Info(msg string, keysAndValues ...interface{}) {
	if l.level <= logger.LevelInfo {
		l.info.Println(l.formatMessage(msg, keysAndValues...))
	}
}

// Warn logs a warning-level message if warning logging is enabled.
func (l *StandardLogger) Warn(msg string, keysAndValues ...interface{}) {
	if l.level <= logger.LevelWarn {
		l.warn.Println(l.formatMessage(msg, keysAndValues...))
	}
}

// Error logs an error-level message if error logging is enabled.
func (l *StandardLogger) Error(msg string, keysAndValues ...interface{}) {
	if l.level <= logger.LevelError {
		l.errLog.Println(l.formatMessage(msg, keysAndValues...))
	}
}

// SetLevel sets the minimum logging level.
func (l *StandardLogger) SetLevel(level logger.Level) {
	l.level = level
}

// GetLevel returns the current logging level.
func (l *StandardLogger) GetLevel() logger.Level {
	return l.level
}

// formatMessage formats a message with optional key-value pairs.
// If keysAndValues is provided, it appends them to the message.
func (l *StandardLogger) formatMessage(msg string, keysAndValues ...interface{}) string {
	if len(keysAndValues) == 0 {
		return msg
	}

	// Format key-value pairs
	formatted := msg
	for i := 0; i < len(keysAndValues); i += 2 {
		if i+1 < len(keysAndValues) {
			formatted += fmt.Sprintf(" %v=%v", keysAndValues[i], keysAndValues[i+1])
		} else {
			formatted += fmt.Sprintf(" %v=<missing>", keysAndValues[i])
		}
	}
	return formatted
}
