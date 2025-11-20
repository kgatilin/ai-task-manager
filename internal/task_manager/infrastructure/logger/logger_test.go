package logger_test

import (
	"bytes"
	"io"
	"log"
	"os"
	"testing"

	"github.com/kgatilin/ai-task-manager/internal/task_manager/domain/logger"
	infralogger "github.com/kgatilin/ai-task-manager/internal/task_manager/infrastructure/logger"
	"github.com/stretchr/testify/assert"
)

func TestNewStandardLogger(t *testing.T) {
	l := infralogger.NewStandardLogger(logger.LevelInfo)
	assert.NotNil(t, l)
	assert.Equal(t, logger.LevelInfo, l.GetLevel())
}

func TestLoggerDebugLevel(t *testing.T) {
	// Capture stderr
	old := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	l := infralogger.NewStandardLogger(logger.LevelDebug)
	l.Debug("debug message", "key", "value")
	l.Info("info message")
	l.Warn("warn message")
	l.Error("error message")

	w.Close()
	os.Stderr = old

	output, _ := io.ReadAll(r)
	outputStr := string(output)

	// All messages should be logged at Debug level
	assert.Contains(t, outputStr, "debug message")
	assert.Contains(t, outputStr, "info message")
	assert.Contains(t, outputStr, "warn message")
	assert.Contains(t, outputStr, "error message")
	assert.Contains(t, outputStr, "[DEBUG]")
}

func TestLoggerInfoLevel(t *testing.T) {
	old := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	l := infralogger.NewStandardLogger(logger.LevelInfo)
	l.Debug("debug message")
	l.Info("info message")
	l.Warn("warn message")
	l.Error("error message")

	w.Close()
	os.Stderr = old

	output, _ := io.ReadAll(r)
	outputStr := string(output)

	// Debug should not be logged, others should
	assert.NotContains(t, outputStr, "debug message")
	assert.Contains(t, outputStr, "info message")
	assert.Contains(t, outputStr, "warn message")
	assert.Contains(t, outputStr, "error message")
}

func TestLoggerWarnLevel(t *testing.T) {
	old := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	l := infralogger.NewStandardLogger(logger.LevelWarn)
	l.Debug("debug message")
	l.Info("info message")
	l.Warn("warn message")
	l.Error("error message")

	w.Close()
	os.Stderr = old

	output, _ := io.ReadAll(r)
	outputStr := string(output)

	// Debug and info should not be logged, warn and error should
	assert.NotContains(t, outputStr, "debug message")
	assert.NotContains(t, outputStr, "info message")
	assert.Contains(t, outputStr, "warn message")
	assert.Contains(t, outputStr, "error message")
}

func TestLoggerErrorLevel(t *testing.T) {
	old := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	l := infralogger.NewStandardLogger(logger.LevelError)
	l.Debug("debug message")
	l.Info("info message")
	l.Warn("warn message")
	l.Error("error message")

	w.Close()
	os.Stderr = old

	output, _ := io.ReadAll(r)
	outputStr := string(output)

	// Only error should be logged
	assert.NotContains(t, outputStr, "debug message")
	assert.NotContains(t, outputStr, "info message")
	assert.NotContains(t, outputStr, "warn message")
	assert.Contains(t, outputStr, "error message")
}

func TestSetLevel(t *testing.T) {
	l := infralogger.NewStandardLogger(logger.LevelError)
	assert.Equal(t, logger.LevelError, l.GetLevel())

	l.SetLevel(logger.LevelInfo)
	assert.Equal(t, logger.LevelInfo, l.GetLevel())

	l.SetLevel(logger.LevelDebug)
	assert.Equal(t, logger.LevelDebug, l.GetLevel())
}

func TestLoggerWithKeyValues(t *testing.T) {
	old := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	l := infralogger.NewStandardLogger(logger.LevelInfo)
	l.Info("operation", "name", "create", "status", "success")

	w.Close()
	os.Stderr = old

	output, _ := io.ReadAll(r)
	outputStr := string(output)

	assert.Contains(t, outputStr, "operation")
	assert.Contains(t, outputStr, "name=create")
	assert.Contains(t, outputStr, "status=success")
}

func TestLoggerOddNumberOfKeyValues(t *testing.T) {
	old := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	l := infralogger.NewStandardLogger(logger.LevelInfo)
	l.Info("message", "key1", "value1", "key2") // odd number

	w.Close()
	os.Stderr = old

	output, _ := io.ReadAll(r)
	outputStr := string(output)

	assert.Contains(t, outputStr, "key1=value1")
	assert.Contains(t, outputStr, "key2=<missing>")
}

func TestLoggerInterface(t *testing.T) {
	// Verify that StandardLogger implements Logger interface
	var _ logger.Logger = (*infralogger.StandardLogger)(nil)
}

func TestLoggerWithCustomStdErr(t *testing.T) {
	// Create a buffer to capture output
	var buf bytes.Buffer

	// Temporarily replace stderr
	old := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	l := infralogger.NewStandardLogger(logger.LevelInfo)
	l.Info("test message")

	w.Close()
	os.Stderr = old

	output, _ := io.ReadAll(r)
	outputStr := string(output)

	assert.Contains(t, outputStr, "test message")
	assert.NotContains(t, buf.String(), "test message") // buf should still be empty
}

func BenchmarkLoggerDebug(b *testing.B) {
	// Silence output
	old := log.Writer()
	log.SetOutput(io.Discard)
	defer log.SetOutput(old)

	l := infralogger.NewStandardLogger(logger.LevelDebug)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Debug("debug message", "key", "value")
	}
}

func BenchmarkLoggerInfo(b *testing.B) {
	old := log.Writer()
	log.SetOutput(io.Discard)
	defer log.SetOutput(old)

	l := infralogger.NewStandardLogger(logger.LevelInfo)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Info("info message")
	}
}
