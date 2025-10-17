package events

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Type represents the type of event captured from Claude Code
type Type string

const (
	// Chat events
	ChatStarted          Type = "chat.started"
	ChatMessageUser      Type = "chat.message.user"
	ChatMessageAssistant Type = "chat.message.assistant"

	// Tool events
	ToolInvoked Type = "tool.invoked"
	ToolResult  Type = "tool.result"

	// File events
	FileRead    Type = "file.read"
	FileWritten Type = "file.written"

	// Context events
	ContextChanged Type = "context.changed"

	// Error events
	Error Type = "error"
)

// Event represents a single logged interaction from Claude Code
type Event struct {
	ID        string          `json:"id"`
	Timestamp int64           `json:"timestamp"` // Unix timestamp in milliseconds
	Event     Type            `json:"event"`
	Payload   json.RawMessage `json:"payload"`
	Content   string          `json:"content"` // Normalized text for analysis
}

// NewEvent creates a new event with generated ID and current timestamp
func NewEvent(eventType Type, payload interface{}, content string) (*Event, error) {
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return &Event{
		ID:        uuid.New().String(),
		Timestamp: time.Now().UnixMilli(),
		Event:     eventType,
		Payload:   payloadJSON,
		Content:   content,
	}, nil
}

// Payload types for different events

// ChatPayload contains data for chat-related events
type ChatPayload struct {
	Message string `json:"message,omitempty"`
	Context string `json:"context,omitempty"`
}

// ToolPayload contains data for tool invocation and result events
type ToolPayload struct {
	Tool       string `json:"tool"`
	Parameters string `json:"parameters,omitempty"`
	Result     string `json:"result,omitempty"`
	DurationMs int64  `json:"duration_ms,omitempty"`
	Context    string `json:"context,omitempty"`
}

// FilePayload contains data for file access events
type FilePayload struct {
	FilePath   string `json:"file_path"`
	Changes    string `json:"changes,omitempty"`
	DurationMs int64  `json:"duration_ms,omitempty"`
	Context    string `json:"context,omitempty"`
}

// ContextPayload contains data for context change events
type ContextPayload struct {
	Context     string `json:"context"`
	Description string `json:"description,omitempty"`
}

// ErrorPayload contains data for error events
type ErrorPayload struct {
	Error      string `json:"error"`
	StackTrace string `json:"stack_trace,omitempty"`
	Context    string `json:"context,omitempty"`
}
