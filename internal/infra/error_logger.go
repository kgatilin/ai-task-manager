package infra

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// ErrorLogger provides structured error logging to a file for debugging
// It logs detailed error information to .darwinflow/errors.log
type ErrorLogger struct {
	logPath string
	mu      sync.Mutex
}

// NewErrorLogger creates a new error logger
// dbPath is expected to be .darwinflow/logs/events.db
// The error log will be created at .darwinflow/errors.log
func NewErrorLogger(dbPath string) (*ErrorLogger, error) {
	// Derive .darwinflow directory from dbPath
	// dbPath is typically .darwinflow/logs/events.db
	// We want .darwinflow/errors.log
	dbDir := filepath.Dir(dbPath)           // .darwinflow/logs
	darwinflowDir := filepath.Dir(dbDir)    // .darwinflow
	logPath := filepath.Join(darwinflowDir, "errors.log")

	return &ErrorLogger{
		logPath: logPath,
	}, nil
}

// LogError logs an error with structured context information
// category: Error category (e.g., "LLM_EXECUTION_FAILED", "ANALYSIS_QUERY_FAILED")
// context: Map of contextual information (session_id, model, duration, etc.)
// err: The error that occurred
func (el *ErrorLogger) LogError(category string, context map[string]interface{}, err error) {
	el.mu.Lock()
	defer el.mu.Unlock()

	// Open file in append mode, create if doesn't exist
	f, openErr := os.OpenFile(el.logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if openErr != nil {
		// Can't log the error, just return silently
		// We don't want to crash the application if error logging fails
		return
	}
	defer f.Close()

	// Write formatted error entry
	timestamp := time.Now().Format(time.RFC3339)
	fmt.Fprintf(f, "\n========================================\n")
	fmt.Fprintf(f, "Timestamp: %s\n", timestamp)
	fmt.Fprintf(f, "Category: %s\n", category)

	// Write context fields in sorted order for consistency
	for key, value := range context {
		// Handle different value types
		switch v := value.(type) {
		case string:
			fmt.Fprintf(f, "%s: %s\n", key, v)
		case int, int64, uint, uint64:
			fmt.Fprintf(f, "%s: %d\n", key, v)
		case float64:
			fmt.Fprintf(f, "%s: %.2f\n", key, v)
		case bool:
			fmt.Fprintf(f, "%s: %t\n", key, v)
		default:
			fmt.Fprintf(f, "%s: %v\n", key, v)
		}
	}

	if err != nil {
		fmt.Fprintf(f, "Error: %v\n", err)
	}
	fmt.Fprintf(f, "========================================\n")
}

// GetLogPath returns the path to the error log file
func (el *ErrorLogger) GetLogPath() string {
	return el.logPath
}
