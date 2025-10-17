package claude

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kgatilin/darwinflow-pub/internal/events"
	"github.com/kgatilin/darwinflow-pub/internal/hooks"
	"github.com/kgatilin/darwinflow-pub/internal/storage"
)

const (
	DefaultDBPath = ".darwinflow/logs/events.db"
)

// Logger orchestrates event logging for Claude Code interactions
type Logger struct {
	store   storage.Store
	context string
}

// NewLogger creates a new logger instance
func NewLogger(dbPath string) (*Logger, error) {
	// Create database directory if needed
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	// Create SQLite store
	store, err := NewSQLiteStore(dbPath)
	if err != nil {
		return nil, err
	}

	// Initialize schema
	ctx := context.Background()
	if err := store.Init(ctx); err != nil {
		store.Close()
		return nil, err
	}

	return &Logger{
		store:   store,
		context: DetectContext(),
	}, nil
}

// LogEvent logs a Claude Code event
func (l *Logger) LogEvent(eventType events.Type, payload interface{}) error {
	// Create event with payload
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Generate normalized content
	content := NormalizeContent(string(eventType), string(payloadJSON))

	// Create event
	event, err := events.NewEvent(eventType, payload, content)
	if err != nil {
		return fmt.Errorf("failed to create event: %w", err)
	}

	// Convert to storage record
	record := storage.Record{
		ID:        event.ID,
		Timestamp: event.Timestamp,
		EventType: string(event.Event),
		Payload:   event.Payload,
		Content:   event.Content,
	}

	// Store event
	ctx := context.Background()
	if err := l.store.Store(ctx, record); err != nil {
		return fmt.Errorf("failed to store event: %w", err)
	}

	return nil
}

// LogFromHookInput logs an event from Claude Code hook input
func (l *Logger) LogFromHookInput(hookInput hooks.HookInput, eventType events.Type, maxParamLength int) error {
	// Create appropriate payload based on event type
	var payload interface{}

	switch eventType {
	case events.ChatStarted:
		payload = map[string]interface{}{
			"session_id": hookInput.SessionID,
			"context":    l.context,
			"cwd":        hookInput.CWD,
		}

	case events.ChatMessageUser:
		// Extract user message from transcript
		message := ""
		if hookInput.TranscriptPath != "" {
			if msg, err := ExtractLastUserMessage(hookInput.TranscriptPath); err == nil {
				message = msg
			}
		}

		payload = events.ChatPayload{
			Message: message,
			Context: l.context,
		}

	case events.ChatMessageAssistant:
		// Extract assistant message from transcript
		message := ""
		if hookInput.TranscriptPath != "" {
			if msg, err := ExtractLastAssistantMessage(hookInput.TranscriptPath); err == nil {
				message = msg
			}
		}

		payload = events.ChatPayload{
			Message: message,
			Context: l.context,
		}

	case events.ToolInvoked:
		// Extract tool name and parameters from transcript
		toolName := "unknown"
		params := ""

		if hookInput.TranscriptPath != "" {
			if tool, p, err := ExtractLastToolUse(hookInput.TranscriptPath, maxParamLength); err == nil {
				toolName = tool
				params = p
			}
		}

		payload = events.ToolPayload{
			Tool:       toolName,
			Parameters: params,
			Context:    l.context,
		}

	case events.ToolResult:
		// Extract tool name from transcript
		toolName := "unknown"
		if hookInput.TranscriptPath != "" {
			if tool, _, err := ExtractLastToolUse(hookInput.TranscriptPath, maxParamLength); err == nil {
				toolName = tool
			}
		}

		payload = events.ToolPayload{
			Tool:    toolName,
			Context: l.context,
		}

	default:
		// Generic payload
		payload = map[string]interface{}{
			"hook_event": hookInput.HookEventName,
			"context":    l.context,
			"cwd":        hookInput.CWD,
		}
	}

	return l.LogEvent(eventType, payload)
}

// Close closes the logger and underlying store
func (l *Logger) Close() error {
	return l.store.Close()
}

// InitializeLogging sets up the logging infrastructure
func InitializeLogging(dbPath string) error {
	// Create database directory
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	// Initialize database
	store, err := NewSQLiteStore(dbPath)
	if err != nil {
		return err
	}
	defer store.Close()

	ctx := context.Background()
	if err := store.Init(ctx); err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	// Add hooks to Claude Code settings
	settingsManager, err := NewSettingsManager()
	if err != nil {
		return fmt.Errorf("failed to create settings manager: %w", err)
	}

	if err := settingsManager.AddDarwinFlowHooks(); err != nil {
		return fmt.Errorf("failed to add hooks: %w", err)
	}

	return nil
}
