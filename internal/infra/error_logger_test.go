package infra_test

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kgatilin/darwinflow-pub/internal/infra"
)

func TestNewErrorLogger(t *testing.T) {
	// Create temporary directory structure
	tmpDir := t.TempDir()
	logsDir := filepath.Join(tmpDir, "logs")
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		t.Fatalf("Failed to create logs directory: %v", err)
	}

	dbPath := filepath.Join(logsDir, "events.db")

	// Create error logger
	logger, err := infra.NewErrorLogger(dbPath)
	if err != nil {
		t.Fatalf("NewErrorLogger failed: %v", err)
	}

	// Verify error log path
	expectedPath := filepath.Join(tmpDir, "errors.log")
	if logger.GetLogPath() != expectedPath {
		t.Errorf("Expected log path %s, got %s", expectedPath, logger.GetLogPath())
	}
}

func TestErrorLogger_LogError(t *testing.T) {
	// Create temporary directory structure
	tmpDir := t.TempDir()
	logsDir := filepath.Join(tmpDir, "logs")
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		t.Fatalf("Failed to create logs directory: %v", err)
	}

	dbPath := filepath.Join(logsDir, "events.db")

	// Create error logger
	logger, err := infra.NewErrorLogger(dbPath)
	if err != nil {
		t.Fatalf("NewErrorLogger failed: %v", err)
	}

	// Log an error
	testErr := errors.New("test error message")
	context := map[string]interface{}{
		"session_id":     "test-session-123",
		"model":          "claude-sonnet-4",
		"duration_ms":    1500,
		"exit_code":      1,
		"tokens_estimate": 10000,
		"bool_value":     true,
		"float_value":    3.14,
	}

	logger.LogError("LLM_EXECUTION_FAILED", context, testErr)

	// Read the log file
	logPath := logger.GetLogPath()
	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("Failed to read error log: %v", err)
	}

	logContent := string(content)

	// Verify log content contains expected fields
	expectedFields := []string{
		"Category: LLM_EXECUTION_FAILED",
		"session_id: test-session-123",
		"model: claude-sonnet-4",
		"duration_ms: 1500",
		"exit_code: 1",
		"tokens_estimate: 10000",
		"bool_value: true",
		"float_value: 3.14",
		"Error: test error message",
		"========================================",
	}

	for _, field := range expectedFields {
		if !strings.Contains(logContent, field) {
			t.Errorf("Log content missing expected field: %s\nLog content:\n%s", field, logContent)
		}
	}
}

func TestErrorLogger_LogMultipleErrors(t *testing.T) {
	// Create temporary directory structure
	tmpDir := t.TempDir()
	logsDir := filepath.Join(tmpDir, "logs")
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		t.Fatalf("Failed to create logs directory: %v", err)
	}

	dbPath := filepath.Join(logsDir, "events.db")

	// Create error logger
	logger, err := infra.NewErrorLogger(dbPath)
	if err != nil {
		t.Fatalf("NewErrorLogger failed: %v", err)
	}

	// Log multiple errors
	logger.LogError("ERROR_1", map[string]interface{}{"field": "value1"}, errors.New("error 1"))
	logger.LogError("ERROR_2", map[string]interface{}{"field": "value2"}, errors.New("error 2"))
	logger.LogError("ERROR_3", map[string]interface{}{"field": "value3"}, errors.New("error 3"))

	// Read the log file
	content, err := os.ReadFile(logger.GetLogPath())
	if err != nil {
		t.Fatalf("Failed to read error log: %v", err)
	}

	logContent := string(content)

	// Verify all errors are logged
	if !strings.Contains(logContent, "ERROR_1") {
		t.Error("Log missing ERROR_1")
	}
	if !strings.Contains(logContent, "ERROR_2") {
		t.Error("Log missing ERROR_2")
	}
	if !strings.Contains(logContent, "ERROR_3") {
		t.Error("Log missing ERROR_3")
	}
	if !strings.Contains(logContent, "error 1") {
		t.Error("Log missing error 1")
	}
	if !strings.Contains(logContent, "error 2") {
		t.Error("Log missing error 2")
	}
	if !strings.Contains(logContent, "error 3") {
		t.Error("Log missing error 3")
	}
}

func TestErrorLogger_LogErrorNilError(t *testing.T) {
	// Create temporary directory structure
	tmpDir := t.TempDir()
	logsDir := filepath.Join(tmpDir, "logs")
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		t.Fatalf("Failed to create logs directory: %v", err)
	}

	dbPath := filepath.Join(logsDir, "events.db")

	// Create error logger
	logger, err := infra.NewErrorLogger(dbPath)
	if err != nil {
		t.Fatalf("NewErrorLogger failed: %v", err)
	}

	// Log with nil error (should not crash)
	logger.LogError("TEST_CATEGORY", map[string]interface{}{"field": "value"}, nil)

	// Read the log file
	content, err := os.ReadFile(logger.GetLogPath())
	if err != nil {
		t.Fatalf("Failed to read error log: %v", err)
	}

	logContent := string(content)

	// Verify category is logged
	if !strings.Contains(logContent, "TEST_CATEGORY") {
		t.Error("Log missing category")
	}
}

func TestErrorLogger_LogErrorEmptyContext(t *testing.T) {
	// Create temporary directory structure
	tmpDir := t.TempDir()
	logsDir := filepath.Join(tmpDir, "logs")
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		t.Fatalf("Failed to create logs directory: %v", err)
	}

	dbPath := filepath.Join(logsDir, "events.db")

	// Create error logger
	logger, err := infra.NewErrorLogger(dbPath)
	if err != nil {
		t.Fatalf("NewErrorLogger failed: %v", err)
	}

	// Log with empty context
	logger.LogError("TEST_CATEGORY", map[string]interface{}{}, errors.New("test error"))

	// Read the log file
	content, err := os.ReadFile(logger.GetLogPath())
	if err != nil {
		t.Fatalf("Failed to read error log: %v", err)
	}

	logContent := string(content)

	// Verify category and error are logged
	if !strings.Contains(logContent, "TEST_CATEGORY") {
		t.Error("Log missing category")
	}
	if !strings.Contains(logContent, "test error") {
		t.Error("Log missing error message")
	}
}

func TestErrorLogger_ConcurrentWrites(t *testing.T) {
	// Create temporary directory structure
	tmpDir := t.TempDir()
	logsDir := filepath.Join(tmpDir, "logs")
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		t.Fatalf("Failed to create logs directory: %v", err)
	}

	dbPath := filepath.Join(logsDir, "events.db")

	// Create error logger
	logger, err := infra.NewErrorLogger(dbPath)
	if err != nil {
		t.Fatalf("NewErrorLogger failed: %v", err)
	}

	// Log errors concurrently
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(id int) {
			logger.LogError("CONCURRENT_TEST", map[string]interface{}{
				"id": id,
			}, errors.New("concurrent error"))
			done <- true
		}(i)
	}

	// Wait for all goroutines to finish
	for i := 0; i < 10; i++ {
		<-done
	}

	// Read the log file
	content, err := os.ReadFile(logger.GetLogPath())
	if err != nil {
		t.Fatalf("Failed to read error log: %v", err)
	}

	logContent := string(content)

	// Verify all errors are logged
	count := strings.Count(logContent, "CONCURRENT_TEST")
	if count != 10 {
		t.Errorf("Expected 10 concurrent errors logged, got %d", count)
	}
}
