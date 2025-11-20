package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// PromptGetter is a function that returns the system prompt.
// This abstraction avoids import cycles.
type PromptGetter func(context.Context) string

// NewPromptCommand creates the prompt command for displaying the LLM system prompt.
func NewPromptCommand(getPrompt PromptGetter) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "prompt",
		Short: "Display LLM system prompt",
		Long: `Displays the system prompt that explains the task manager to LLMs.

This prompt contains comprehensive documentation about the task manager's entity
hierarchy (Roadmap → Track → Task → Iteration), standard workflows, best practices,
and integration with other systems. Use this when working with AI assistants or
documenting the task manager usage.

The prompt explains:
- Entity definitions and relationships
- Required vs optional entities
- Standard workflows for different scenarios
- Best practices and pitfalls to avoid
- Integration with AC and ADR systems
- Command reference and examples`,
		Example: `  # Display prompt to terminal
  tm prompt

  # Save prompt to file for documentation
  tm prompt --output task-manager-prompt.md

  # Save and pipe to other tools
  tm prompt --output prompt.md && cat prompt.md

  # Use with LLM (example with Claude CLI)
  tm prompt | claude --system-prompt -`,
		RunE: func(cmd *cobra.Command, args []string) error {
			outputFile, _ := cmd.Flags().GetString("output")

			// Get the system prompt
			prompt := getPrompt(cmd.Context())

			// If output file is specified, save to file
			if outputFile != "" {
				// Create directory if needed
				dir := filepath.Dir(outputFile)
				if dir != "" && dir != "." {
					if err := os.MkdirAll(dir, 0755); err != nil {
						return fmt.Errorf("failed to create output directory: %w", err)
					}
				}

				// Write prompt to file
				if err := os.WriteFile(outputFile, []byte(prompt), 0644); err != nil {
					return fmt.Errorf("failed to write prompt to file: %w", err)
				}

				fmt.Fprintf(cmd.OutOrStdout(), "System prompt saved to: %s\n", outputFile)
				return nil
			}

			// Otherwise, print prompt to stdout
			fmt.Fprint(cmd.OutOrStdout(), prompt)

			return nil
		},
	}

	// Define flags
	cmd.Flags().String("output", "", "Save prompt to specified file instead of displaying")

	return cmd
}
