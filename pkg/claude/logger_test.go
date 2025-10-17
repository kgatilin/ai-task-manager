package claude_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kgatilin/darwinflow-pub/internal/events"
	"github.com/kgatilin/darwinflow-pub/internal/hooks"
	"github.com/kgatilin/darwinflow-pub/pkg/claude"
)

func TestNewLogger(t *testing.T) {
	// Create temporary directory for test database
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	logger, err := claude.NewLogger(dbPath)
	if err != nil {
		t.Fatalf("NewLogger failed: %v", err)
	}
	defer logger.Close()

	// Verify database file was created
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Error("Database file was not created")
	}
}

func TestLogEvent(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	logger, err := claude.NewLogger(dbPath)
	if err != nil {
		t.Fatalf("NewLogger failed: %v", err)
	}
	defer logger.Close()

	// Test logging a chat event
	payload := events.ChatPayload{
		Message: "Test message",
		Context: "test/project",
	}

	err = logger.LogEvent(events.ChatMessageUser, payload)
	if err != nil {
		t.Errorf("LogEvent failed: %v", err)
	}
}

func TestLogFromHookInput(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	logger, err := claude.NewLogger(dbPath)
	if err != nil {
		t.Fatalf("NewLogger failed: %v", err)
	}
	defer logger.Close()

	hookInput := hooks.HookInput{
		SessionID:      "test-session-123",
		TranscriptPath: "",
		CWD:            "/test/path",
		HookEventName:  "SessionStart",
	}

	err = logger.LogFromHookInput(hookInput, events.ChatStarted, 30)
	if err != nil {
		t.Errorf("LogFromHookInput failed: %v", err)
	}
}

func TestInitializeLogging(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	// Change to temp directory for this test
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tmpDir)

	err := claude.InitializeLogging(dbPath)
	if err != nil {
		t.Fatalf("InitializeLogging failed: %v", err)
	}

	// Verify database was created
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Error("Database file was not created")
	}
}
