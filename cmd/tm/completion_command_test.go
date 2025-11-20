package main

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestCompletionCommand(t *testing.T) {
	cmd := NewCompletionCommand()

	if cmd == nil {
		t.Fatal("NewCompletionCommand() returned nil")
	}

	if cmd.Use != "completion [bash|zsh|fish|powershell]" {
		t.Errorf("Use = %v, want 'completion [bash|zsh|fish|powershell]'", cmd.Use)
	}

	// Test valid args
	expectedArgs := []string{"bash", "zsh", "fish", "powershell"}
	if len(cmd.ValidArgs) != len(expectedArgs) {
		t.Errorf("ValidArgs length = %v, want %v", len(cmd.ValidArgs), len(expectedArgs))
	}
}

func TestCompletionCommand_WithRoot(t *testing.T) {
	// Create a root command and add completion as subcommand
	rootCmd := &cobra.Command{Use: "tm"}
	completionCmd := NewCompletionCommand()
	rootCmd.AddCommand(completionCmd)

	// Test bash completion (basic check that it doesn't error)
	rootCmd.SetArgs([]string{"completion", "bash"})

	// We can't easily test the actual completion output without a full setup,
	// but we can verify the command structure is correct
	if completionCmd.Parent() != rootCmd {
		t.Error("completion command not properly attached to root")
	}
}
