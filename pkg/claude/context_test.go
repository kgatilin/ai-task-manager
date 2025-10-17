package claude_test

import (
	"os"
	"testing"

	"github.com/kgatilin/darwinflow-pub/pkg/claude"
)

func TestDetectContext_FromEnv(t *testing.T) {
	// Set environment variable
	os.Setenv("DW_CONTEXT", "project/test-project")
	defer os.Unsetenv("DW_CONTEXT")

	context := claude.DetectContext()
	if context != "project/test-project" {
		t.Errorf("Expected 'project/test-project', got '%s'", context)
	}
}

func TestDetectContext_Default(t *testing.T) {
	// Unset environment variable
	os.Unsetenv("DW_CONTEXT")

	context := claude.DetectContext()
	// Should return something, default is "unknown" or derived from path
	if context == "" {
		t.Error("DetectContext returned empty string")
	}
}

func TestNormalizeContent(t *testing.T) {
	tests := []struct {
		name      string
		eventType string
		payload   string
		want      string
	}{
		{
			name:      "simple event",
			eventType: "chat.started",
			payload:   `{"session": "123"}`,
			want:      `chat.started: {"session": "123"}`,
		},
		{
			name:      "empty payload",
			eventType: "tool.invoked",
			payload:   "",
			want:      "tool.invoked:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := claude.NormalizeContent(tt.eventType, tt.payload)
			if got != tt.want {
				t.Errorf("NormalizeContent() = %q, want %q", got, tt.want)
			}
		})
	}
}
