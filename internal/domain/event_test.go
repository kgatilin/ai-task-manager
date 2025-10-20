package domain_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/kgatilin/darwinflow-pub/internal/domain"
)

func TestNewEvent(t *testing.T) {
	tests := []struct {
		name      string
		eventType domain.EventType
		sessionID string
		payload   interface{}
		content   string
	}{
		{
			name:      "creates chat started event",
			eventType: domain.ChatStarted,
			sessionID: "test-session-1",
			payload:   domain.ChatPayload{Message: "Hello", Context: "greeting"},
			content:   "Hello greeting",
		},
		{
			name:      "creates tool invoked event",
			eventType: domain.ToolInvoked,
			sessionID: "test-session-2",
			payload:   domain.ToolPayload{Tool: "Read", Parameters: map[string]string{"file": "test.go"}},
			content:   "Reading test.go",
		},
		{
			name:      "creates event with nil payload",
			eventType: domain.FileRead,
			sessionID: "test-session-3",
			payload:   nil,
			content:   "file read",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := domain.NewEvent(tt.eventType, tt.sessionID, tt.payload, tt.content)

			// Verify required fields are set
			if event.ID == "" {
				t.Error("Expected ID to be generated, got empty string")
			}
			if event.SessionID != tt.sessionID {
				t.Errorf("Expected SessionID = %q, got %q", tt.sessionID, event.SessionID)
			}
			if event.Type != tt.eventType {
				t.Errorf("Expected Type = %q, got %q", tt.eventType, event.Type)
			}
			if event.Content != tt.content {
				t.Errorf("Expected Content = %q, got %q", tt.content, event.Content)
			}

			// Verify timestamp is recent (within last second)
			if time.Since(event.Timestamp) > time.Second {
				t.Errorf("Expected recent timestamp, got %v", event.Timestamp)
			}
		})
	}
}

func TestEvent_MarshalPayload(t *testing.T) {
	tests := []struct {
		name       string
		payload    interface{}
		wantErr    bool
		validateFn func([]byte) error
	}{
		{
			name: "marshals chat payload",
			payload: domain.ChatPayload{
				Message: "test message",
				Context: "test context",
			},
			wantErr: false,
			validateFn: func(data []byte) error {
				var p domain.ChatPayload
				if err := json.Unmarshal(data, &p); err != nil {
					return err
				}
				if p.Message != "test message" {
					t.Errorf("Expected Message = %q, got %q", "test message", p.Message)
				}
				return nil
			},
		},
		{
			name: "marshals tool payload",
			payload: domain.ToolPayload{
				Tool:       "Bash",
				Parameters: map[string]string{"command": "ls"},
				DurationMs: 100,
			},
			wantErr: false,
			validateFn: func(data []byte) error {
				var p domain.ToolPayload
				if err := json.Unmarshal(data, &p); err != nil {
					return err
				}
				if p.Tool != "Bash" {
					t.Errorf("Expected Tool = %q, got %q", "Bash", p.Tool)
				}
				return nil
			},
		},
		{
			name: "marshals file payload",
			payload: domain.FilePayload{
				FilePath:   "/test/path.go",
				Changes:    "added function",
				DurationMs: 50,
			},
			wantErr: false,
			validateFn: func(data []byte) error {
				var p domain.FilePayload
				if err := json.Unmarshal(data, &p); err != nil {
					return err
				}
				if p.FilePath != "/test/path.go" {
					t.Errorf("Expected FilePath = %q, got %q", "/test/path.go", p.FilePath)
				}
				return nil
			},
		},
		{
			name: "marshals error payload",
			payload: domain.ErrorPayload{
				Error:      "test error",
				StackTrace: "line 1\nline 2",
				Context:    "during test",
			},
			wantErr: false,
			validateFn: func(data []byte) error {
				var p domain.ErrorPayload
				if err := json.Unmarshal(data, &p); err != nil {
					return err
				}
				if p.Error != "test error" {
					t.Errorf("Expected Error = %q, got %q", "test error", p.Error)
				}
				return nil
			},
		},
		{
			name:    "marshals nil payload",
			payload: nil,
			wantErr: false,
			validateFn: func(data []byte) error {
				if string(data) != "null" {
					t.Errorf("Expected JSON null, got %q", string(data))
				}
				return nil
			},
		},
		{
			name: "marshals complex nested payload",
			payload: map[string]interface{}{
				"nested": map[string]interface{}{
					"key": "value",
					"arr": []int{1, 2, 3},
				},
			},
			wantErr: false,
			validateFn: func(data []byte) error {
				var result map[string]interface{}
				if err := json.Unmarshal(data, &result); err != nil {
					return err
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := domain.NewEvent(domain.ChatStarted, "test-session", tt.payload, "content")

			data, err := event.MarshalPayload()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalPayload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.validateFn != nil {
				if err := tt.validateFn(data); err != nil {
					t.Errorf("Validation failed: %v", err)
				}
			}
		})
	}
}

func TestEventTypes(t *testing.T) {
	// Test that all event type constants are defined
	eventTypes := []domain.EventType{
		domain.ChatStarted,
		domain.ChatMessageUser,
		domain.ChatMessageAssistant,
		domain.ToolInvoked,
		domain.ToolResult,
		domain.FileRead,
		domain.FileWritten,
		domain.ContextChanged,
		domain.Error,
	}

	// Verify each type has a non-empty value
	for _, et := range eventTypes {
		if string(et) == "" {
			t.Errorf("Event type %v has empty string value", et)
		}
	}

	// Verify types are unique
	seen := make(map[domain.EventType]bool)
	for _, et := range eventTypes {
		if seen[et] {
			t.Errorf("Duplicate event type: %v", et)
		}
		seen[et] = true
	}
}

func TestChatPayload(t *testing.T) {
	payload := domain.ChatPayload{
		Message: "test message",
		Context: "test context",
	}

	// Verify JSON marshaling preserves fields
	data, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Failed to marshal ChatPayload: %v", err)
	}

	var unmarshaled domain.ChatPayload
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal ChatPayload: %v", err)
	}

	if unmarshaled.Message != payload.Message {
		t.Errorf("Expected Message = %q, got %q", payload.Message, unmarshaled.Message)
	}
	if unmarshaled.Context != payload.Context {
		t.Errorf("Expected Context = %q, got %q", payload.Context, unmarshaled.Context)
	}
}

func TestToolPayload(t *testing.T) {
	payload := domain.ToolPayload{
		Tool:       "Read",
		Parameters: map[string]interface{}{"file": "test.go", "offset": 10},
		Result:     "file contents",
		DurationMs: 42,
		Context:    "reading file",
	}

	// Verify JSON marshaling preserves fields
	data, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Failed to marshal ToolPayload: %v", err)
	}

	var unmarshaled domain.ToolPayload
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal ToolPayload: %v", err)
	}

	if unmarshaled.Tool != payload.Tool {
		t.Errorf("Expected Tool = %q, got %q", payload.Tool, unmarshaled.Tool)
	}
	if unmarshaled.DurationMs != payload.DurationMs {
		t.Errorf("Expected DurationMs = %d, got %d", payload.DurationMs, unmarshaled.DurationMs)
	}
}

func TestFilePayload(t *testing.T) {
	payload := domain.FilePayload{
		FilePath:   "/test/file.go",
		Changes:    "added function Foo",
		DurationMs: 123,
		Context:    "writing changes",
	}

	// Verify JSON marshaling preserves fields
	data, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Failed to marshal FilePayload: %v", err)
	}

	var unmarshaled domain.FilePayload
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal FilePayload: %v", err)
	}

	if unmarshaled.FilePath != payload.FilePath {
		t.Errorf("Expected FilePath = %q, got %q", payload.FilePath, unmarshaled.FilePath)
	}
	if unmarshaled.Changes != payload.Changes {
		t.Errorf("Expected Changes = %q, got %q", payload.Changes, unmarshaled.Changes)
	}
}

func TestErrorPayload(t *testing.T) {
	payload := domain.ErrorPayload{
		Error:      "file not found",
		StackTrace: "line 1\nline 2\nline 3",
		Context:    "during read operation",
	}

	// Verify JSON marshaling preserves fields
	data, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Failed to marshal ErrorPayload: %v", err)
	}

	var unmarshaled domain.ErrorPayload
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal ErrorPayload: %v", err)
	}

	if unmarshaled.Error != payload.Error {
		t.Errorf("Expected Error = %q, got %q", payload.Error, unmarshaled.Error)
	}
	if unmarshaled.StackTrace != payload.StackTrace {
		t.Errorf("Expected StackTrace = %q, got %q", payload.StackTrace, unmarshaled.StackTrace)
	}
}

func TestContextPayload(t *testing.T) {
	payload := domain.ContextPayload{
		Context:     "new context",
		Description: "context changed due to user action",
	}

	// Verify JSON marshaling preserves fields
	data, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Failed to marshal ContextPayload: %v", err)
	}

	var unmarshaled domain.ContextPayload
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal ContextPayload: %v", err)
	}

	if unmarshaled.Context != payload.Context {
		t.Errorf("Expected Context = %q, got %q", payload.Context, unmarshaled.Context)
	}
	if unmarshaled.Description != payload.Description {
		t.Errorf("Expected Description = %q, got %q", payload.Description, unmarshaled.Description)
	}
}
