package cli_test

import (
	"context"
	"testing"

	"github.com/kgatilin/ai-task-manager/internal/task_manager/presentation/cli"
)

// TestGetSystemPrompt verifies that GetSystemPrompt returns the default prompt.
func TestGetSystemPrompt(t *testing.T) {
	ctx := context.Background()
	prompt := cli.GetSystemPrompt(ctx)

	// Should return the default prompt
	if prompt != cli.DefaultSystemPrompt {
		t.Error("GetSystemPrompt should return DefaultSystemPrompt")
	}

	// Should not be empty
	if prompt == "" {
		t.Error("GetSystemPrompt should not return empty string")
	}
}
