package app_test

import (
	"testing"

	"github.com/kgatilin/darwinflow-pub/internal/app"
	"github.com/kgatilin/darwinflow-pub/internal/domain"
)

func TestEventMapper_MapEventType(t *testing.T) {
	mapper := &app.EventMapper{}

	tests := []struct {
		name     string
		input    string
		expected domain.EventType
	}{
		// Chat events
		{name: "chat.started", input: "chat.started", expected: domain.ChatStarted},
		{name: "chat.ended", input: "chat.ended", expected: domain.ChatStarted},
		{name: "chat.end", input: "chat.end", expected: domain.ChatStarted},

		// User messages
		{name: "chat.message.user", input: "chat.message.user", expected: domain.ChatMessageUser},
		{name: "user.message", input: "user.message", expected: domain.ChatMessageUser},

		// Assistant messages
		{name: "chat.message.assistant", input: "chat.message.assistant", expected: domain.ChatMessageAssistant},
		{name: "assistant.message", input: "assistant.message", expected: domain.ChatMessageAssistant},

		// Tool events
		{name: "tool.invoked", input: "tool.invoked", expected: domain.ToolInvoked},
		{name: "tool.invoke", input: "tool.invoke", expected: domain.ToolInvoked},
		{name: "tool.result", input: "tool.result", expected: domain.ToolResult},

		// File events
		{name: "file.read", input: "file.read", expected: domain.FileRead},
		{name: "file.written", input: "file.written", expected: domain.FileWritten},
		{name: "file.write", input: "file.write", expected: domain.FileWritten},

		// Context events
		{name: "context.changed", input: "context.changed", expected: domain.ContextChanged},
		{name: "context.change", input: "context.change", expected: domain.ContextChanged},

		// Error events
		{name: "error", input: "error", expected: domain.Error},

		// Case insensitive
		{name: "uppercase CHAT.STARTED", input: "CHAT.STARTED", expected: domain.ChatStarted},
		{name: "mixed case Chat.Started", input: "Chat.Started", expected: domain.ChatStarted},

		// Underscore normalization
		{name: "underscore chat_started", input: "chat_started", expected: domain.ChatStarted},
		{name: "underscore tool_invoked", input: "tool_invoked", expected: domain.ToolInvoked},
		{name: "underscore file_read", input: "file_read", expected: domain.FileRead},

		// Unknown event types (returns as-is, normalized)
		{name: "unknown event", input: "custom.event", expected: domain.EventType("custom.event")},
		{name: "unknown with underscore", input: "custom_event", expected: domain.EventType("custom.event")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapper.MapEventType(tt.input)
			if result != tt.expected {
				t.Errorf("MapEventType(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestEventMapper_MapEventType_Normalization(t *testing.T) {
	mapper := &app.EventMapper{}

	// Test that normalization converts underscores to dots and lowercases
	tests := []struct {
		input    string
		expected domain.EventType
	}{
		{"Chat_Started", domain.ChatStarted},
		{"TOOL_INVOKED", domain.ToolInvoked},
		{"File_Read", domain.FileRead},
		{"Tool_Result", domain.ToolResult},
	}

	for _, tt := range tests {
		result := mapper.MapEventType(tt.input)
		if result != tt.expected {
			t.Errorf("MapEventType(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestNewLoggerService(t *testing.T) {
	repo := &MockEventRepository{}
	transcriptParser := &MockTranscriptParser{}
	contextDetector := &MockContextDetector{context: "test-context"}
	normalizer := func(eventType, payload string) string {
		return eventType + ":" + payload
	}

	logger := app.NewLoggerService(repo, transcriptParser, contextDetector, normalizer)

	if logger == nil {
		t.Error("Expected non-nil LoggerService")
	}
}

// MockTranscriptParser for testing
type MockTranscriptParser struct {
	toolName      string
	toolParams    string
	userMessage   string
	assistantMsg  string
	extractError  error
}

func (m *MockTranscriptParser) ExtractLastToolUse(transcriptPath string, maxParamLength int) (string, string, error) {
	if m.extractError != nil {
		return "", "", m.extractError
	}
	return m.toolName, m.toolParams, nil
}

func (m *MockTranscriptParser) ExtractLastUserMessage(transcriptPath string) (string, error) {
	if m.extractError != nil {
		return "", m.extractError
	}
	return m.userMessage, nil
}

func (m *MockTranscriptParser) ExtractLastAssistantMessage(transcriptPath string) (string, error) {
	if m.extractError != nil {
		return "", m.extractError
	}
	return m.assistantMsg, nil
}

// MockContextDetector for testing
type MockContextDetector struct {
	context string
}

func (m *MockContextDetector) DetectContext() string {
	return m.context
}

func TestLoggerService_Creation(t *testing.T) {
	repo := &MockEventRepository{}
	parser := &MockTranscriptParser{}
	detector := &MockContextDetector{context: "/test/path"}
	normalizer := func(eventType, payload string) string {
		return payload
	}

	service := app.NewLoggerService(repo, parser, detector, normalizer)

	if service == nil {
		t.Fatal("LoggerService should not be nil")
	}

	// Verify service was created (we can't access private fields, but creation is enough)
}

func TestLoggerService_ContentNormalizer(t *testing.T) {
	// Test that content normalizer is used
	repo := &MockEventRepository{}
	parser := &MockTranscriptParser{}
	detector := &MockContextDetector{context: "ctx"}

	called := false

	normalizer := func(eventType, payload string) string {
		called = true
		return "normalized"
	}

	service := app.NewLoggerService(repo, parser, detector, normalizer)

	// We can't directly test this without calling LogEvent, which requires
	// more complex setup. This test verifies the service accepts the normalizer.
	if service == nil {
		t.Error("Service should be created with normalizer")
	}

	// Verify normalizer wasn't called during construction
	if called {
		t.Error("Normalizer should not be called during construction")
	}
}
